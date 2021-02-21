package hackCall

import _ "unsafe"

type nameOff int32
type typeOff int32
type textOff int32

// A ptabEntry is generated by the compiler for each exported function
// and global variable in the main package of a plugin. It is used to
// initialize the plugin module's symbol map.
type ptabEntry struct {
	name nameOff
	typ  typeOff
}

type functab struct {
	entry   uintptr
	funcoff uintptr
}

type textsect struct {
	vaddr    uintptr // prelinked section vaddr
	length   uintptr // section length
	baseaddr uintptr // relocated section address
}

type modulehash struct {
	modulename   string
	linktimehash string
	runtimehash  *string
}

type bitvector struct {
	n        int32 // # of bits
	bytedata *uint8
}

type moduledata struct {
	pclntable    []byte
	ftab         []functab
	filetab      []uint32
	findfunctab  uintptr
	minpc, maxpc uintptr

	text, etext           uintptr
	noptrdata, enoptrdata uintptr
	data, edata           uintptr
	bss, ebss             uintptr
	noptrbss, enoptrbss   uintptr
	end, gcdata, gcbss    uintptr
	types, etypes         uintptr

	textsectmap []textsect
	typelinks   []int32 // offsets from types
	itablinks   []uintptr

	ptab []ptabEntry

	pluginpath string
	pkghashes  []modulehash

	modulename   string
	modulehashes []modulehash

	hasmain uint8 // 1 if module contains the main function, 0 otherwise

	gcdatamask, gcbssmask bitvector

	typemap map[typeOff]uintptr // offset to *_rtype in previous module

	bad bool // module failed to load and should be ignored

	next uintptr
}

//go:linkname activeModules runtime.activeModules
func activeModules() []*moduledata

type funcID uint8

type _func struct {
	entry   uintptr // start pc
	nameoff int32   // function name

	args        int32  // in/out args size
	deferreturn uint32 // offset of start of a deferreturn call instruction from entry, if any.

	pcsp      int32
	pcfile    int32
	pcln      int32
	npcdata   int32
	funcID    funcID  // set for certain special runtime functions
	_         [2]int8 // unused
	nfuncdata uint8   // must be last
}

type funcInfo struct {
	*_func
	datap *moduledata
}

//go:linkname funcname runtime.funcname
func funcname(f funcInfo) string