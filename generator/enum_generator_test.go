package generator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-courier/loaderx"
)

func TestGenerator(t *testing.T) {
	cwd, _ := os.Getwd()
	p, pkgInfo, _ := loaderx.LoadWithTests(filepath.Join(cwd, "../examples"))

	g := NewEnumGenerator(p, pkgInfo)

	g.Scan("Protocol")
	g.Output(cwd)
}
