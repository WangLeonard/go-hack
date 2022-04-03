package binary

import (
	"bytes"
	"debug/elf"
	"debug/macho"
	"encoding/binary"
	"errors"
	"io"
)

var (
	// errUnrecognizedFormat is returned when a given executable file doesn't
	// appear to be in a known format, or it breaks the rules of that format,
	// or when there are I/O errors reading the file.
	errUnrecognizedFormat = errors.New("unrecognized file format")

	// errNotGoExe is returned when a given executable file is valid but does
	// not contain Go build information.
	errNotGoExe = errors.New("not a Go executable")

	// The build info blob left by the linker is identified by
	// a 16-byte header, consisting of buildInfoMagic (14 bytes),
	// the binary's pointer size (1 byte),
	// and whether the binary is big endian (1 byte).
	buildInfoMagic = []byte("\xff Go buildinf:")

	errNotFoundSymbol = errors.New("symbol not found")
)

type exe interface {
	// ReadData reads and returns up to size bytes starting at virtual address addr.
	ReadData(addr, size uint64) ([]byte, error)

	// DataStart returns the virtual address of the segment or section that
	// should contain build information. This is either a specially named section
	// or the first writable non-zero data segment.
	DataStart() uint64

	ByteOrder() binary.ByteOrder
}

func GetSymValue(r io.ReaderAt, symName string) (symValue string, err error) {
	// Read the first bytes of the file to identify the format, then delegate to
	// a format-specific function to load segment and section headers.
	ident := make([]byte, 16)
	if n, err := r.ReadAt(ident, 0); n < len(ident) || err != nil {
		return "", errUnrecognizedFormat
	}

	var x exe
	var symAddr uint64
	switch {
	case bytes.HasPrefix(ident, []byte("\x7FELF")):
		f, err := elf.NewFile(r)
		if err != nil {
			return "", errUnrecognizedFormat
		}
		x = &elfExe{f}
		symbols, err := f.Symbols()
		if err != nil {
			return "", err
		}
		for _, v := range symbols {
			if v.Name == symName {
				symAddr = v.Value
				break
			}
		}
	case bytes.HasPrefix(ident, []byte("\xFE\xED\xFA")) || bytes.HasPrefix(ident[1:], []byte("\xFA\xED\xFE")):
		f, err := macho.NewFile(r)
		if err != nil {
			return "", errUnrecognizedFormat
		}
		x = &machoExe{f}
		for _, v := range f.Symtab.Syms {
			if v.Name == symName {
				symAddr = v.Value
				break
			}
		}
	default:
		return "", errUnrecognizedFormat
	}

	if symAddr == 0 {
		return "", errNotFoundSymbol
	}

	// Read the first 64kB of dataAddr to find the build info blob.
	// On some platforms, the blob will be in its own section, and DataStart
	// returns the address of that section. On others, it's somewhere in the
	// data segment; the linker puts it near the beginning.
	// See cmd/link/internal/ld.Link.buildinfo.
	dataAddr := x.DataStart()
	data, err := x.ReadData(dataAddr, 64*1024)
	if err != nil {
		return "", err
	}
	const (
		buildInfoAlign = 16
		buildInfoSize  = 32
	)
	for {
		i := bytes.Index(data, buildInfoMagic)
		if i < 0 || len(data)-i < buildInfoSize {
			return "", errNotGoExe
		}
		if i%buildInfoAlign == 0 && len(data)-i >= buildInfoSize {
			data = data[i:]
			break
		}
		data = data[(i+buildInfoAlign-1)&^buildInfoAlign:]
	}

	// Decode the blob.
	// The first 14 bytes are buildInfoMagic.
	// The next two bytes indicate pointer size in bytes (4 or 8) and endianness
	// (0 for little, 1 for big).
	// Two virtual addresses to Go strings follow that: runtime.buildVersion,
	// and runtime.modinfo.
	// On 32-bit platforms, the last 8 bytes are unused.
	// If the endianness has the 2 bit set, then the pointers are zero
	// and the 32-byte header is followed by varint-prefixed string data
	// for the two string values we care about.
	// But the byteOrder information cannot be obtained in this way.
	ptrSize := int(data[14])
	var bo = x.ByteOrder()

	var readPtr func([]byte) uint64
	if ptrSize == 4 {
		readPtr = func(b []byte) uint64 { return uint64(bo.Uint32(b)) }
	} else {
		readPtr = bo.Uint64
	}

	readPtr = func(b []byte) uint64 {
		_ = b[7] // bounds check hint to compiler; see golang.org/issue/14808
		return uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 |
			uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56
	}

	symValue = readString(x, ptrSize, readPtr, symAddr)
	return symValue, nil
}

// readString returns the string at address addr in the executable x.
func readString(x exe, ptrSize int, readPtr func([]byte) uint64, addr uint64) string {
	hdr, err := x.ReadData(addr, uint64(2*ptrSize))
	if err != nil || len(hdr) < 2*ptrSize {
		return ""
	}
	dataAddr := readPtr(hdr)
	dataLen := readPtr(hdr[ptrSize:])
	data, err := x.ReadData(dataAddr, dataLen)
	if err != nil || uint64(len(data)) < dataLen {
		return ""
	}
	return string(data)
}

// elfExe is the ELF implementation of the exe interface.
type elfExe struct {
	f *elf.File
}

func (x *elfExe) ReadData(addr, size uint64) ([]byte, error) {
	for _, prog := range x.f.Progs {
		if prog.Vaddr <= addr && addr <= prog.Vaddr+prog.Filesz-1 {
			n := prog.Vaddr + prog.Filesz - addr
			if n > size {
				n = size
			}
			data := make([]byte, n)
			_, err := prog.ReadAt(data, int64(addr-prog.Vaddr))
			if err != nil {
				return nil, err
			}
			return data, nil
		}
	}
	return nil, errUnrecognizedFormat
}

func (x *elfExe) DataStart() uint64 {
	for _, s := range x.f.Sections {
		if s.Name == ".go.buildinfo" {
			return s.Addr
		}
	}
	for _, p := range x.f.Progs {
		if p.Type == elf.PT_LOAD && p.Flags&(elf.PF_X|elf.PF_W) == elf.PF_W {
			return p.Vaddr
		}
	}
	return 0
}

func (x *elfExe) ByteOrder() binary.ByteOrder {
	return x.f.ByteOrder
}

// machoExe is the Mach-O (Apple macOS/iOS) implementation of the exe interface.
type machoExe struct {
	f *macho.File
}

func (x *machoExe) ReadData(addr, size uint64) ([]byte, error) {
	for _, load := range x.f.Loads {
		seg, ok := load.(*macho.Segment)
		if !ok {
			continue
		}
		if seg.Addr <= addr && addr <= seg.Addr+seg.Filesz-1 {
			if seg.Name == "__PAGEZERO" {
				continue
			}
			n := seg.Addr + seg.Filesz - addr
			if n > size {
				n = size
			}
			data := make([]byte, n)
			_, err := seg.ReadAt(data, int64(addr-seg.Addr))
			if err != nil {
				return nil, err
			}
			return data, nil
		}
	}
	return nil, errUnrecognizedFormat
}

func (x *machoExe) DataStart() uint64 {
	// Look for section named "__go_buildinfo".
	for _, sec := range x.f.Sections {
		if sec.Name == "__go_buildinfo" {
			return sec.Addr
		}
	}
	// Try the first non-empty writable segment.
	const RW = 3
	for _, load := range x.f.Loads {
		seg, ok := load.(*macho.Segment)
		if ok && seg.Addr != 0 && seg.Filesz != 0 && seg.Prot == RW && seg.Maxprot == RW {
			return seg.Addr
		}
	}
	return 0
}

func (x *machoExe) ByteOrder() binary.ByteOrder {
	return x.f.ByteOrder
}
