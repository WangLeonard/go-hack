package autoCode

import (
	"fmt"
	"testing"
)

func TestAutoInjectionCode(t *testing.T) {
	err := AutoInjectionCode(
		"./testdata/router.go",
		"Routers",
		"Code generated by XXX Begin; DO NOT EDIT.",
		"Code generated by XXX End; DO NOT EDIT.",
		"router.InitLeonardWangRouter(Router)")
	fmt.Println(err)
}
