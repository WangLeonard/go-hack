package main

import (
	"fmt"
	"runtime"
	"strconv"
)

var pcs = make([]uintptr, 32)
var n int

//go:noinline
func f1() {
	//fmt.Println("f1")
	f2()
}

//go:noinline
func f2() {
	//fmt.Println("f2")
	n = runtime.Callers(0, pcs)
	for i := 0; i < n; i++ {
		fmt.Println(strconv.FormatUint(uint64(pcs[i]), 10))
	}
}

func main() {
	f1()
	//fmt.Println("n:", n)
	//for i := 0; i < n; i++ {
	//	fmt.Println(pcs[i])
	//}
	//frames := runtime.CallersFrames(pcs[:n])
	//for {
	//	frame, more := frames.Next()
	//	println(frame.File, frame.Line, frame.Function)
	//	if !more {
	//		break
	//	}
	//}
}
