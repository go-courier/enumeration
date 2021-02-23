package scanner

import (
	"go/ast"
	"go/constant"
	"go/types"
	"sort"
	"strconv"
	"strings"

	"github.com/go-courier/codegen"
	"github.com/go-courier/packagesx"
)

func NewScanner(pkg *packagesx.Package) *Scanner {
	return &Scanner{
		pkg: pkg,
	}
}

type Scanner struct {
	pkg     *packagesx.Package
	results map[*types.TypeName]Options
}

// When TypeName have multi const values, this TypeName should be a enum
func (s *Scanner) Options(typeName *types.TypeName) (Options, bool) {
	if typeName == nil {
		return nil, false
	}

	if enumOptions, ok := s.results[typeName]; ok {
		return enumOptions, ok
	}

	pkgInfo := s.pkg.Pkg(typeName.Pkg().Path())
	if pkgInfo == nil {
		return nil, false
	}

	typeNameString := typeName.Name()

	for ident, def := range pkgInfo.TypesInfo.Defs {
		c, ok := def.(*types.Const)
		if !ok {
			continue
		}

		if c.Type() != typeName.Type() {
			continue
		}

		name := c.Name()

		if strings.HasPrefix(name, "_") {
			continue
		}

		label := strings.TrimSpace(ident.Obj.Decl.(*ast.ValueSpec).Comment.Text())

		val := c.Val()

		switch val.Kind() {
		case constant.String:
			v, _ := strconv.Unquote(val.String())
			s.appendOption(typeName, NewStrOption(v, label))
		case constant.Int:
			// TYPE_NAME_UNKNOWN
			// TYPE_NAME__XXX
			if strings.HasPrefix(name, codegen.UpperSnakeCase(typeNameString)) {
				var values = strings.SplitN(name, "__", 2)
				if len(values) == 2 {
					intVal, _ := strconv.ParseInt(val.String(), 10, 64)
					s.appendOption(typeName, NewIntStringerOption(intVal, values[1], label))
				}
			} else {
				intVal, _ := strconv.ParseInt(val.String(), 10, 64)
				s.appendOption(typeName, NewIntOption(intVal, label))
			}
		case constant.Float:
			f, _ := strconv.ParseFloat(val.String(), 64)
			s.appendOption(typeName, NewFloatOption(f, label))
		case constant.Bool, constant.Complex:
			return nil, false
		}
	}

	return s.results[typeName], len(s.results[typeName]) > 0
}

func (s *Scanner) appendOption(typeName *types.TypeName, opt *Option) {
	if s.results == nil {
		s.results = map[*types.TypeName]Options{}
	}

	s.results[typeName] = append(s.results[typeName], *opt)

	sort.Sort(s.results[typeName])
}
