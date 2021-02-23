package scanner

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-courier/packagesx"
	. "github.com/onsi/gomega"
)

func TestScaner(t *testing.T) {
	cwd, _ := os.Getwd()
	p, _ := packagesx.Load(filepath.Join(cwd, "../__examples__"))

	g := NewScanner(p)

	t.Run("should scan int stringer enum", func(t *testing.T) {
		options, ok := g.Options(g.pkg.TypeName("Protocol"))
		NewWithT(t).Expect(ok).To(BeTrue())

		NewWithT(t).Expect(options.Len()).To(Equal(3))
		NewWithT(t).Expect(*options[0].Str).To(Equal("HTTP"))
		NewWithT(t).Expect(*options[2].Str).To(Equal("TCP"))
	})

	t.Run("should scan string enum", func(t *testing.T) {
		options, ok := g.Options(g.pkg.TypeName("PullPolicy"))
		NewWithT(t).Expect(ok).To(BeTrue())

		NewWithT(t).Expect(options.Len()).To(Equal(3))
		NewWithT(t).Expect(*options[0].Str).To(Equal("Always"))
		NewWithT(t).Expect(*options[2].Str).To(Equal("Never"))
	})
}
