package getString

import (
	"testing"

	"github.com/WangLeonard/go-hack/getString/testdata/package1"
)

func TestGetByFullName(t *testing.T) {
	ver, err := GetInitValByFullName("", "runtime.buildVersion")
	if err != nil || ver == "" {
		t.Error("find buildVersion error")
	}

	gloVar, err := GetInitValByFullName("", "github.com/WangLeonard/go-hack/getString/testdata/package1.Var")
	if err != nil || gloVar != "const v1" {
		t.Error("global var value error")
	}

	package1.Var = "const v2"
	if package1.Var != "const v2" {
		t.Error("global var assignment failed")
	}
	gloVar, err = GetInitValByFullName("", "github.com/WangLeonard/go-hack/getString/testdata/package1.Var")
	if err != nil || gloVar != "const v1" {
		t.Error("not the initial value")
	}
}
