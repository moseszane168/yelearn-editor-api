package base

import (
	"os"
	"path/filepath"
)

func GetAppPathAbs() (string, error) {
	return filepath.Abs(filepath.Dir(os.Args[0]))
}
