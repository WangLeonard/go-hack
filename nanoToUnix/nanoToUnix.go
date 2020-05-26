package nanoToUnix

import (
	"time"
	"unsafe"
)

type Time struct {
	// wall and ext encode the wall time seconds, wall time nanoseconds,
	// and optional monotonic clock reading in nanoseconds.
	//
	// From high to low bit position, wall encodes a 1-bit flag (hasMonotonic),
	// a 33-bit seconds field, and a 30-bit wall time nanoseconds field.
	// The nanoseconds field is in the range [0, 999999999].
	// If the hasMonotonic bit is 0, then the 33-bit field must be zero
	// and the full signed 64-bit wall seconds since Jan 1 year 1 is stored in ext.
	// If the hasMonotonic bit is 1, then the 33-bit field holds a 33-bit
	// unsigned wall seconds since Jan 1 year 1885, and ext holds a
	// signed 64-bit monotonic clock reading, nanoseconds since process start.
	wall uint64
	ext  int64

	// loc specifies the Location that should be used to
	// determine the minute, hour, month, day, and year
	// that correspond to this Time.
	// The nil location means UTC.
	// All UTC times are represented with loc==nil, never loc==&utcLoc.
	loc uintptr
}

//go:linkname nanotime runtime.nanotime
func nanotime() int64

//go:linkname startNano time.startNano
var startNano int64

var startTime = time.Now()
var nanoToUnix int64

func init() {
	startNanoTime := (*(*Time)(unsafe.Pointer(&startTime))).ext + startNano
	nanoToUnix = startTime.UnixNano() - startNanoTime
}

func CurrentUnix() int64 {
	return CurrentUnixNano() / 1e9
}

func CurrentUnixMill() int64 {
	return CurrentUnixNano() / 1e6
}

func CurrentUnixNano() int64 {
	return nanotime() + nanoToUnix
}

