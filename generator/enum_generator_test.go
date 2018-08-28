package generator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-courier/packagesx"
)

func TestGenerator(t *testing.T) {
	cwd, _ := os.Getwd()
	p, _ := packagesx.Load(filepath.Join(cwd, "../__examples__"))

	g := NewEnumGenerator(p)

	g.Scan("Protocol")
	g.Output(cwd)
}
