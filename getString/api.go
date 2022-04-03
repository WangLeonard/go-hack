package getString

import (
	"errors"
	"os"

	binary "github.com/WangLeonard/go-hack/getString/internal"
)

var (
	ErrStringNotFound    = errors.New("string not found")
	ErrWrongBinaryFormat = errors.New("binary format error")
	ErrFileNotExist      = errors.New("file not exist")
)

func GetInitValByFullName(path, name string) (string, error) {
	var err error
	if path == "" {
		// Current binary
		path, err = os.Executable()
		if err != nil {
			return "", ErrWrongBinaryFormat
		}
	}

	fi, err := os.Open(path)
	if err != nil {
		return "", ErrFileNotExist
	}
	defer fi.Close()
	return binary.GetSymValue(fi, name)
}
