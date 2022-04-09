package offlinepc

import (
	"bytes"
	"debug/elf"
	"debug/gosym"
	"debug/macho"
	"errors"
	"os"
)

var (
	// errUnrecognizedFormat is returned when a given executable file doesn't
	// appear to be in a known format, or it breaks the rules of that format,
	// or when there are I/O errors reading the file.
	errUnrecognizedFormat = errors.New("unrecognized file format")
)

type PcInfoParser struct {
	table *gosym.Table
}

func (p *PcInfoParser) PCToLine(pc uint64) (file string, line int, fnName string) {
	file, line, fn := p.table.PCToLine(pc)
	return file, line, fn.Name
}

func NewPcInfoParser(path string) (*PcInfoParser, error) {
	fi, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fi.Close()
	// Read the first bytes of the file to identify the format, then delegate to
	// a format-specific function to load segment and section headers.
	ident := make([]byte, 16)
	if n, err := fi.ReadAt(ident, 0); n < len(ident) || err != nil {
		return nil, errUnrecognizedFormat
	}

	switch {
	case bytes.HasPrefix(ident, []byte("\x7FELF")):
		f, err := elf.NewFile(fi)
		if err != nil {
			return nil, errUnrecognizedFormat
		}
		defer f.Close()

		tab, err := parseelf(f)
		if err != nil {
			return nil, err
		}
		return &PcInfoParser{table: tab}, nil
	case bytes.HasPrefix(ident, []byte("\xFE\xED\xFA")) || bytes.HasPrefix(ident[1:], []byte("\xFA\xED\xFE")):
		f, err := macho.NewFile(fi)
		if err != nil {
			return nil, errUnrecognizedFormat
		}

		tab, err := parsemacho(f)
		if err != nil {
			return nil, err
		}
		return &PcInfoParser{table: tab}, nil
	default:
		return nil, errUnrecognizedFormat
	}
}

func parseelf(f *elf.File) (*gosym.Table, error) {
	s := f.Section(".gosymtab")
	if s == nil {
		return nil, errors.New("no .gosymtab section")
	}
	symdat, err := s.Data()
	if err != nil {
		return nil, err
	}
	pclndat, err := f.Section(".gopclntab").Data()
	if err != nil {
		return nil, err
	}

	pcln := gosym.NewLineTable(pclndat, f.Section(".text").Addr)
	tab, err := gosym.NewTable(symdat, pcln)
	if err != nil {
		return nil, err
	}
	return tab, nil
}

func parsemacho(f *macho.File) (*gosym.Table, error) {
	s := f.Section("__gosymtab")
	if s == nil {
		return nil, errors.New("no .gosymtab section")
	}
	symdat, err := s.Data()
	if err != nil {
		return nil, err
	}
	pclndat, err := f.Section("__gopclntab").Data()
	if err != nil {
		return nil, err
	}

	pcln := gosym.NewLineTable(pclndat, f.Section("__text").Addr)
	tab, err := gosym.NewTable(symdat, pcln)
	if err != nil {
		return nil, err
	}
	return tab, nil
}
