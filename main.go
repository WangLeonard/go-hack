package main

import (
	"fmt"

	"github.com/WangLeonard/go-hack/callFuncByName"
)

//go:noinline
func HelloWorld(args *callFuncByName.MyString) {
	fmt.Println("Hello " + args.Name + "!")
}

var global = make(map[int64]interface{}, 10)

func init() {
	// Keep it alive.
	global[0] = HelloWorld
	global[1] = HelloWorld2
}

type MyType struct {
	Name string
}

func (t *MyType) Hello() {
	fmt.Println("MyType Hello")
	//fmt.Println("MyType Hello" + t.Name) // it will throw?
}

//go:noinline
func HelloWorld2(args callFuncByName.Say) {
	args.Hello()
}

func main() {
	err := callFuncByName.HackCallFuncByNameWithStructArg("main.HelloWorld", &callFuncByName.MyString{"WWW"})
	fmt.Println(err)

	err = callFuncByName.HackCallFuncByNameWithInterfaceArg("main.HelloWorld2", &MyType{})
	fmt.Println(err)
}
