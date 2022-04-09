package offlinepc

import (
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestOfflineParsePC(t *testing.T) {
	dir := t.TempDir()
	fileName := "TestOfflineParsePC"
	out := filepath.Join(dir, fileName)
	cmd := exec.Command("go", "build", "-gcflags=-N -l", "-o="+out, "./testdata")
	_, err := cmd.CombinedOutput()
	if err != nil {
		t.Error(err)
	}
	cmd2 := exec.Command(out)
	got, err := cmd2.CombinedOutput()
	if err != nil {
		t.Error(err)
	}
	pcs := strings.Split(string(got), "\n")
	p, err := NewPcInfoParser(out)
	if err != nil {
		t.Error(err)
	}
	funcs := make([]string, 0, len(pcs))
	for _, pc := range pcs {
		if pc == "" {
			continue
		}
		pvValue, err := strconv.ParseUint(pc, 10, 64)
		if err != nil {
			t.Error(err)
		}
		_, _, fnName := p.PCToLine(pvValue)
		funcs = append(funcs, fnName)
	}
	expect := []string{"runtime.Callers", "main.f2", "main.f1", "main.main", "runtime.main", "runtime.goexit"}
	if len(funcs) != len(expect) {
		t.Error("func not equal")
	}
	for i := 0; i < len(funcs); i++ {
		if funcs[i] != expect[i] {
			t.Error("func not equal")
		}
	}
}
