package hackCall

import (
	"errors"
	"unsafe"
)

var funMaps = make(map[string]uintptr, 10)

func init() {
	mds := activeModules()

	for _, md := range mds {
		for _, fn := range md.ftab {
			fi := funcInfo{(*_func)(unsafe.Pointer(&md.pclntable[fn.funcoff])), md}
			name := funcname(fi)
			funMaps[name] = fn.entry
		}
	}
}

func CallFuncNoArgs(name string) error {
	if fn, ok := funMaps[name]; ok {
		funcType := unsafe.Pointer(&_func{entry: fn})
		funcInstance := *(*func())(unsafe.Pointer(&funcType))
		funcInstance()
		return nil
	}
	return errors.New("Function Not Found")
}

func CallFuncWithInterfaceArgs(name string, arg1 interface{}, arg2 interface{}) error {
	if fn, ok := funMaps[name]; ok {
		funcType := unsafe.Pointer(&_func{entry: fn})
		funcInstance := *(*func(interface{}, interface{}))(unsafe.Pointer(&funcType))
		funcInstance(arg1, arg2)
		return nil
	}
	return errors.New("Function Not Found")
}

func GetFuncType(name string) (unsafe.Pointer, error) {
	if fn, ok := funMaps[name]; ok {
		funcType := unsafe.Pointer(&_func{entry: fn})
		return unsafe.Pointer(&funcType), nil
	}
	return nil, errors.New("Function Not Found")
}

func GetFuncPointer(name string) (uintptr, error) {
	if fn, ok := funMaps[name]; ok {
		return fn, nil
	}
	return uintptr(0), errors.New("Function Not Found")
}

func MethodCall(fn interface{}, obj interface{}, args interface{})

func MethodCallByPtr(fn uintptr, obj interface{}, args interface{})
