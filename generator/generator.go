package generator

import (
	"fmt"
	"go/types"
	"log"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/go-courier/enumeration/scanner"

	"github.com/go-courier/codegen"
	"github.com/go-courier/packagesx"
	"golang.org/x/tools/go/packages"
)

func NewGenerator(pkg *packagesx.Package) *Generator {
	return &Generator{
		pkg:     pkg,
		scanner: scanner.NewScanner(pkg),
		enums:   map[*types.TypeName]*Enum{},
	}
}

type Generator struct {
	pkg     *packagesx.Package
	scanner *scanner.Scanner
	enums   map[*types.TypeName]*Enum
}

func (g *Generator) Scan(names ...string) {
	for _, name := range names {
		typeName := g.pkg.TypeName(name)

		if options, ok := g.scanner.Options(typeName); ok {
			// only generate int stringer enum
			if options[0].Int != nil && options[0].Str != nil {
				g.enums[typeName] = NewEnum(typeName.Pkg().Path()+"."+typeName.Name(), options)
			}
		}
	}
}

func getPkgDir(importPath string) string {
	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedName | packages.NeedFiles,
	}, importPath)
	if err != nil {
		panic(err)
	}
	if len(pkgs) == 0 {
		panic(fmt.Errorf("package `%s` not found", importPath))
	}
	return filepath.Dir(pkgs[0].GoFiles[0])
}

func (g *Generator) Output(cwd string) {
	for typeName, enum := range g.enums {
		dir, _ := filepath.Rel(cwd, getPkgDir(typeName.Pkg().Path()))
		filename := codegen.GeneratedFileSuffix(path.Join(dir, codegen.LowerSnakeCase(enum.Name)+".go"))

		file := codegen.NewFile(typeName.Pkg().Name(), filename)
		enum.WriteToFile(file)

		if _, err := file.WriteFile(); err != nil {
			log.Printf("%s generated", file)
		}
	}
}

func NewEnum(pkgTypeOrName string, options scanner.Options) *Enum {
	parts := strings.Split(pkgTypeOrName, ".")
	pkgPath, name := "", ""

	switch len(parts) {
	case 1:
		name = parts[0]
	default:
		pkgPath = strings.Join(parts[0:len(parts)-1], ".")
		name = parts[len(parts)-1]
	}

	return &Enum{
		PkgPath: pkgPath,
		Name:    name,
		Options: options,
	}
}

type Enum struct {
	PkgPath string
	Name    string
	Options scanner.Options
}

func (e *Enum) ConstUnknown() codegen.Snippet {
	return codegen.Id(codegen.UpperSnakeCase(e.Name) + "_UNKNOWN")
}

func (e *Enum) ConstName(value string) codegen.Snippet {
	return codegen.Id(codegen.UpperSnakeCase(e.Name) + "__" + value)
}

func (e *Enum) VarInvalidError() codegen.Snippet {
	return codegen.Id("Invalid" + e.Name)
}

func (e *Enum) WriteToFile(file *codegen.File) {
	e.WriteErrors(file)
	e.WriteLabelStringParser(file)
	e.WriteStringer(file)
	e.WriteStringParser(file)
	e.WriteLabeler(file)
	e.WriteInt(file)
	e.WriteTypeNameAndConstValues(file)
	e.TextMarshalerAndTextUnmarshaler(file)
	e.TextScanAndValuer(file)
}

func (e *Enum) WriteErrors(file *codegen.File) {
	file.WriteBlock(
		codegen.Expr(
			`var ? = ?("invalid ? type")`,
			e.VarInvalidError(),
			codegen.Id(file.Use("errors", "New")),
			codegen.Id(e.Name),
		),
	)
}

func (e *Enum) WriteStringParser(file *codegen.File) {
	clauses := []*codegen.SnippetClause{
		codegen.Clause(file.Val("")).Do(
			codegen.Return(
				e.ConstUnknown(),
				codegen.Nil,
			),
		),
	}

	for _, option := range e.Options {
		clauses = append(clauses, codegen.Clause(file.Val(*option.Str)).Do(
			codegen.Return(
				e.ConstName(*option.Str),
				codegen.Nil,
			),
		))
	}

	file.WriteBlock(
		codegen.Func(codegen.Var(codegen.String, "s")).
			Named(fmt.Sprintf("Parse%sFromString", e.Name)).
			Return(codegen.Var(codegen.Type(e.Name)), codegen.Var(codegen.Error)).Do(
			codegen.Switch(codegen.Id("s")).When(
				clauses...,
			),
			codegen.Return(e.ConstUnknown(), e.VarInvalidError()),
		),
	)

}

func (e *Enum) WriteLabelStringParser(file *codegen.File) {
	clauses := []*codegen.SnippetClause{
		codegen.Clause(file.Val("")).Do(
			codegen.Return(
				e.ConstUnknown(),
				codegen.Nil,
			),
		),
	}

	for _, option := range e.Options {
		clauses = append(clauses, codegen.Clause(file.Val(option.Label)).Do(
			codegen.Return(
				e.ConstName(*option.Str),
				codegen.Nil,
			),
		))
	}

	file.WriteBlock(
		codegen.Func(codegen.Var(codegen.String, "s")).
			Named(fmt.Sprintf("Parse%sFromLabelString", e.Name)).
			Return(codegen.Var(codegen.Type(e.Name)), codegen.Var(codegen.Error)).Do(
			codegen.Switch(codegen.Id("s")).When(
				clauses...,
			),
			codegen.Return(e.ConstUnknown(), e.VarInvalidError()),
		),
	)
}

func (e *Enum) WriteStringer(file *codegen.File) {
	clauses := []*codegen.SnippetClause{
		codegen.Clause(e.ConstUnknown()).Do(
			codegen.Return(
				file.Val(""),
			),
		),
	}

	for _, option := range e.Options {
		clauses = append(clauses, codegen.Clause(e.ConstName(*option.Str)).Do(
			codegen.Return(
				file.Val(*option.Str),
			),
		))
	}

	file.WriteBlock(
		codegen.Func().
			MethodOf(codegen.Var(codegen.Type(e.Name), "v")).
			Named("String").
			Return(codegen.Var(codegen.String)).Do(
			codegen.Switch(codegen.Id("v")).When(
				clauses...,
			),
			codegen.Return(file.Val("UNKNOWN")),
		),
	)

}

func (e *Enum) WriteInt(file *codegen.File) {
	file.WriteBlock(
		codegen.Func().
			MethodOf(codegen.Var(codegen.Type(e.Name), "v")).
			Named("Int").
			Return(codegen.Var(codegen.Int)).Do(
			codegen.Return(file.Expr("int(v)")),
		),
	)

}

func (e *Enum) WriteLabeler(file *codegen.File) {
	clauses := []*codegen.SnippetClause{
		codegen.Clause(e.ConstUnknown()).Do(
			codegen.Return(
				file.Val(""),
			),
		),
	}

	for _, option := range e.Options {
		clauses = append(clauses, codegen.Clause(e.ConstName(*option.Str)).Do(
			codegen.Return(
				file.Val(option.Label),
			),
		))
	}

	file.WriteBlock(
		codegen.Func().
			MethodOf(codegen.Var(codegen.Type(e.Name), "v")).
			Named("Label").
			Return(codegen.Var(codegen.String)).Do(
			codegen.Switch(codegen.Id("v")).When(
				clauses...,
			),
			codegen.Return(file.Val("UNKNOWN")),
		),
	)

}

func (e *Enum) WriteTypeNameAndConstValues(file *codegen.File) {
	file.WriteBlock(
		codegen.Func().
			MethodOf(codegen.Var(codegen.Type(e.Name))).
			Named("TypeName").
			Return(codegen.Var(codegen.String)).Do(
			codegen.Return(file.Val(e.PkgPath + "." + e.Name)),
		),
	)

	tpe := codegen.Slice(codegen.Type(file.Use("github.com/go-courier/enumeration", "IntStringerEnum")))

	list := []interface{}{
		tpe,
	}
	holder := "?"

	sort.Slice(e.Options, func(i, j int) bool {
		return *e.Options[i].Int < *e.Options[j].Int
	})

	for i, o := range e.Options {
		if i > 0 {
			holder += ",?"
		}
		list = append(list, e.ConstName(*o.Str))
	}

	file.WriteBlock(
		codegen.Func().
			MethodOf(codegen.Var(codegen.Type(e.Name))).
			Named("ConstValues").
			Return(codegen.Var(tpe)).Do(
			codegen.Return(file.Expr(`?{`+holder+`}`, list...)),
		),
	)

}

func (e *Enum) TextMarshalerAndTextUnmarshaler(file *codegen.File) {
	file.WriteBlock(
		codegen.Func().
			MethodOf(codegen.Var(codegen.Type(e.Name), "v")).
			Named("MarshalText").
			Return(
				codegen.Var(codegen.Slice(codegen.Byte)),
				codegen.Var(codegen.Error),
			).Do(
			file.Expr(`str := v.String()`),
			codegen.If(file.Expr(`str == ?`, "UNKNOWN")).Do(
				codegen.Return(codegen.Nil, e.VarInvalidError()),
			),
			codegen.Return(file.Expr(`[]byte(str)`), codegen.Nil),
		),
	)

	file.WriteBlock(
		codegen.Func(codegen.Var(codegen.Slice(codegen.Byte), "data")).
			MethodOf(codegen.Var(codegen.Star(codegen.Type(e.Name)), "v")).
			Named("UnmarshalText").
			Return(codegen.Var(codegen.Error, "err")).Do(
			file.Expr("*v, err = Parse?FromString(string(?(data)))", codegen.Id(e.Name), codegen.Id(file.Use("bytes", "ToUpper"))),
			codegen.Return(),
		),
	)

}

func (e *Enum) TextScanAndValuer(file *codegen.File) {
	offsetExprs := file.Expr(`offset := 0
if o, ok := (interface{})(v).(?); ok {
	offset = o.Offset()
}`, codegen.Id(file.Use("github.com/go-courier/enumeration", "DriverValueOffset")))

	file.WriteBlock(
		codegen.Func().
			MethodOf(codegen.Var(codegen.Type(e.Name), "v")).
			Named("Value").
			Return(
				codegen.Var(codegen.Type(file.Use("database/sql/driver", "Value"))),
				codegen.Var(codegen.Error),
			).Do(
			offsetExprs,
			codegen.Return(file.Expr("int64(v) + int64(offset)"), codegen.Nil),
		),
	)

	file.WriteBlock(
		codegen.Func(codegen.Var(codegen.Interface(), "src")).
			MethodOf(codegen.Var(codegen.Star(codegen.Type(e.Name)), "v")).
			Named("Scan").
			Return(codegen.Var(codegen.Error)).Do(
			offsetExprs,
			file.Expr(
				`
i, err := ?(src, offset)
if err != nil {
	return err
}
*v = ?(i)
return nil`,
				codegen.Id(file.Use("github.com/go-courier/enumeration", "ScanIntEnumStringer")),
				codegen.Id(e.Name),
			),
		),
	)

}
