package generator

import (
	"fmt"
	"go/build"
	"go/types"
	"path"
	"path/filepath"

	"github.com/go-courier/codegen"
	"github.com/go-courier/loaderx"
	"github.com/sirupsen/logrus"
	"golang.org/x/tools/go/loader"

	"github.com/go-courier/enumeration"
)

func NewEnumGenerator(program *loader.Program, rootPkgInfo *loader.PackageInfo) *EnumGenerator {
	return &EnumGenerator{
		pkgInfo: rootPkgInfo,
		scanner: NewEnumScanner(program),
		enums:   map[string]*Enum{},
	}
}

type EnumGenerator struct {
	pkgInfo *loader.PackageInfo
	scanner *EnumScanner
	enums   map[string]*Enum
}

func (g *EnumGenerator) Scan(names ...string) {
	pkgInfo := loaderx.NewPackageInfo(g.pkgInfo)

	for _, name := range names {
		typeName := pkgInfo.TypeName(name)
		g.enums[name] = NewEnum(typeName, g.scanner.Enum(typeName))
	}
}

func (g *EnumGenerator) Output(cwd string) {
	for _, enum := range g.enums {
		p, _ := build.Import(enum.TypeName.Pkg().Path(), "", build.FindOnly)
		dir, _ := filepath.Rel(cwd, p.Dir)
		filename := codegen.GeneratedFileSuffix(path.Join(dir, codegen.LowerSnakeCase(enum.Name())+".go"))

		file := codegen.NewFile(enum.TypeName.Pkg().Name(), filename)
		enum.WriteToFile(file)

		if _, err := file.WriteFile(); err != nil {
			logrus.Printf("%s generated", file)
		}
	}
}

func NewEnum(typeName *types.TypeName, options []enumeration.EnumOption) *Enum {
	return &Enum{
		TypeName: typeName,
		Options:  options,
	}
}

type Enum struct {
	TypeName *types.TypeName
	Options  []enumeration.EnumOption
}

func (e *Enum) Name() string {
	return e.TypeName.Name()
}

func (e *Enum) ConstUnknown() codegen.Snippet {
	return codegen.Id(codegen.UpperSnakeCase(e.Name()) + "_UNKNOWN")
}

func (e *Enum) ConstValue(value string) codegen.Snippet {
	return codegen.Id(codegen.UpperSnakeCase(e.Name()) + "__" + value)
}

func (e *Enum) VarInvalidError() codegen.Snippet {
	return codegen.Id("Invalid" + e.Name())
}

func (e *Enum) WriteToFile(file *codegen.File) {
	e.WriteInit(file)
	e.WriteErrors(file)
	e.WriteLabelStringParser(file)
	e.WriteStringer(file)
	e.WriteStringParser(file)
	e.WriteLabeler(file)
	e.WriteTypeNameAndConstValues(file)
	e.TextMarshalerAndTextUnmarshaler(file)
	e.TextScanAndValuer(file)
}

func (e *Enum) WriteInit(file *codegen.File) {
	file.WriteBlock(
		codegen.Func().Named("init").Do(
			codegen.Sel(
				codegen.Id(file.Use("github.com/go-courier/enumeration", "DefaultEnumMap")),
				codegen.Call("Register", e.ConstUnknown()),
			),
		),
	)
}

func (e *Enum) WriteErrors(file *codegen.File) {
	file.WriteBlock(
		codegen.Expr(
			`var ? = ?("invalid ? type")`,
			e.VarInvalidError(),
			codegen.Id(file.Use("errors", "New")),
			codegen.Id(e.Name()),
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
		clauses = append(clauses, codegen.Clause(file.Val(option.Value)).Do(
			codegen.Return(
				e.ConstValue(option.Value),
				codegen.Nil,
			),
		))
	}

	file.WriteBlock(
		codegen.Func(codegen.Var(codegen.String, "s")).
			Named(fmt.Sprintf("Parse%sFromString", e.Name())).
			Return(codegen.Var(codegen.Type(e.Name())), codegen.Var(codegen.Error)).Do(
			codegen.Switch(codegen.Id("s")).When(
				clauses...,
			),
			codegen.Return(e.ConstUnknown(), e.VarInvalidError()),
		),
	)

	file.WriteRune('\n')
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
				e.ConstValue(option.Value),
				codegen.Nil,
			),
		))
	}

	file.WriteBlock(
		codegen.Func(codegen.Var(codegen.String, "s")).
			Named(fmt.Sprintf("Parse%sFromLabelString", e.Name())).
			Return(codegen.Var(codegen.Type(e.Name())), codegen.Var(codegen.Error)).Do(
			codegen.Switch(codegen.Id("s")).When(
				clauses...,
			),
			codegen.Return(e.ConstUnknown(), e.VarInvalidError()),
		),
	)

	file.WriteRune('\n')
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
		clauses = append(clauses, codegen.Clause(e.ConstValue(option.Value)).Do(
			codegen.Return(
				file.Val(option.Value),
			),
		))
	}

	file.WriteBlock(
		codegen.Func().
			MethodOf(codegen.Var(codegen.Type(e.Name()), "v")).
			Named("String").
			Return(codegen.Var(codegen.String)).Do(
			codegen.Switch(codegen.Id("v")).When(
				clauses...,
			),
			codegen.Return(file.Val("UNKNOWN")),
		),
	)

	file.WriteRune('\n')
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
		clauses = append(clauses, codegen.Clause(e.ConstValue(option.Value)).Do(
			codegen.Return(
				file.Val(option.Label),
			),
		))
	}

	file.WriteBlock(
		codegen.Func().
			MethodOf(codegen.Var(codegen.Type(e.Name()), "v")).
			Named("Label").
			Return(codegen.Var(codegen.String)).Do(
			codegen.Switch(codegen.Id("v")).When(
				clauses...,
			),
			codegen.Return(file.Val("UNKNOWN")),
		),
	)

	file.WriteRune('\n')
}

func (e *Enum) WriteTypeNameAndConstValues(file *codegen.File) {
	file.WriteBlock(
		codegen.Func().
			MethodOf(codegen.Var(codegen.Type(e.Name()))).
			Named("TypeName").
			Return(codegen.Var(codegen.String)).Do(
			codegen.Return(file.Val(e.Name())),
		),
	)

	file.WriteRune('\n')

	tpe := codegen.Slice(codegen.Type(file.Use("github.com/go-courier/enumeration", "Enum")))

	list := []interface{}{
		tpe,
	}
	holder := "?"

	for i, o := range e.Options {
		if i > 0 {
			holder += ",?"
		}
		list = append(list, e.ConstValue(o.Value))
	}

	file.WriteBlock(
		codegen.Func().
			MethodOf(codegen.Var(codegen.Type(e.Name()))).
			Named("ConstValues").
			Return(codegen.Var(tpe)).Do(
			codegen.Return(file.Expr(`?{`+holder+`}`, list...)),
		),
	)

	file.WriteRune('\n')
}

func (e *Enum) TextMarshalerAndTextUnmarshaler(file *codegen.File) {
	file.WriteBlock(
		codegen.Func().
			MethodOf(codegen.Var(codegen.Type(e.Name()), "v")).
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

	file.WriteRune('\n')

	file.WriteBlock(
		codegen.Func(codegen.Var(codegen.Slice(codegen.Byte), "data")).
			MethodOf(codegen.Var(codegen.Star(codegen.Type(e.Name())), "v")).
			Named("UnmarshalText").
			Return(codegen.Var(codegen.Error, "err")).Do(
			file.Expr("*v, err = Parse?FromString(string(?(data)))", codegen.Id(e.Name()), codegen.Id(file.Use("bytes", "ToUpper"))),
			codegen.Return(),
		),
	)

	file.WriteRune('\n')
}

func (e *Enum) TextScanAndValuer(file *codegen.File) {
	offsetExprs := file.Expr(`offset := 0
if o, ok := (interface{})(v).(?); ok {
	offset = o.Offset()
}`, codegen.Id(file.Use("github.com/go-courier/enumeration", "EnumDriverValueOffset")))

	file.WriteBlock(
		codegen.Func().
			MethodOf(codegen.Var(codegen.Type(e.Name()), "v")).
			Named("Value").
			Return(
				codegen.Var(codegen.Type(file.Use("database/sql/driver", "Value"))),
				codegen.Var(codegen.Error),
			).Do(
			offsetExprs,
			codegen.Return(file.Expr("int(v) + offset"), codegen.Nil),
		),
	)

	file.WriteRune('\n')

	file.WriteBlock(
		codegen.Func(codegen.Var(codegen.Interface(), "src")).
			MethodOf(codegen.Var(codegen.Star(codegen.Type(e.Name())), "v")).
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
return nil
`,
				codegen.Id(file.Use("github.com/go-courier/enumeration", "ScanEnum")),
				codegen.Id(e.Name()),
			),
		),
	)

	file.WriteRune('\n')
}
