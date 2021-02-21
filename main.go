package main

import (
	"fmt"
	"reflect"

	"github.com/WangLeonard/go-hack/hackCall"
)

var global = make(map[int64]interface{}, 10)

func init() {
	// Keep it alive.
	global[0] = Hello
	global[1] = HelloTwoInterfaces
	global[2] = HelloTwoStrings
}

//go:noinline
func Hello() {
	fmt.Println("HelloWorld")
}

//go:noinline
func HelloTwoInterfaces(name1 interface{}, name2 interface{}) {
	s1, ok := name1.(string)
	fmt.Println("Args1", s1, ok)

	s2, ok := name2.(string)
	fmt.Println("Args2", s2, ok)

	fmt.Println("Hello", s1, s2)
}

//go:noinline
func HelloTwoStrings(s1 string, s2 string) {
	fmt.Println("Hello", s1, s2)
}

type MyType struct {
	Name string
	Age  int64
}

func (t *MyType) Hello(args interface{}) {
	switch arg := args.(type) {
	case int64:
		fmt.Println("int64 arg:", arg)
		t.Age = arg
	case string:
		fmt.Println("string arg:", arg)
		t.Name = arg
	}
	fmt.Println("MyType Hello", t.Name, t.Age)
}

func main() {
	fmt.Println("Demo1--Call No Args Function")
	err := hackCall.CallFuncNoArgs("main.Hello")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Demo2--Call Two Interface Args Function")
	err = hackCall.CallFuncWithInterfaceArgs("main.HelloTwoInterfaces", "Leonard", "Wang")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Demo3--Call Two String Args Function")
	p, err := hackCall.GetFuncType("main.HelloTwoStrings")
	if err != nil {
		fmt.Println(err)
	} else {
		funcInstance := *(*func(string, string))(p)
		funcInstance("Leonard", "Wang")
	}

	fmt.Println("Demo4--Call Method By Func Pointer")
	var method interface{} = (*MyType).Hello
	t := &MyType{"LeonardWang", int64(18)}
	hackCall.MethodCall(method, t, "NewName")
	hackCall.MethodCall(method, t, int64(20))

	fmt.Println("Demo5--Call Method By Reflect")
	t2 := &MyType{"LeonardWang", int64(18)}
	val := reflect.ValueOf(t2)
	val.MethodByName("Hello").Call([]reflect.Value{reflect.ValueOf("NewName")})
	val.MethodByName("Hello").Call([]reflect.Value{reflect.ValueOf(int64(20))})

	fmt.Println("Demo5--Call Method By Name")
	t3 := &MyType{"LeonardWang", int64(18)}
	p2, err := hackCall.GetFuncPointer("main.(*MyType).Hello")
	if err != nil {
		fmt.Println(err)
	} else {
		hackCall.MethodCallByPtr(p2, t3, "NewName")
		hackCall.MethodCallByPtr(p2, t3, int64(20))
	}
}
