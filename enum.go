package enumeration

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"github.com/stretchr/testify/require"
)

type Enum interface {
	TypeName() string
	ConstValues() []Enum
	Int() int
	String() string
	Label() string
}

// sql value of enum maybe have offset from value of enum in go
type EnumDriverValueOffset interface {
	Offset() int
}

type EnumInfo struct {
	TypeName string       `json:"typeName"`
	Options  []EnumOption `json:"options"`
}

type EnumOption struct {
	ConstValue int    `json:"constValue"`
	Value      string `json:"value"`
	Label      string `json:"label"`
}

var DefaultEnumMap = EnumMap{}

type EnumMap map[string]Enum

func (m EnumMap) Register(enum Enum) {
	registerEnumName := ""
	pc := make([]uintptr, 1) // at least 1 entry needed
	n := runtime.Callers(2, pc)
	for i := 0; i < n; i++ {
		f := runtime.FuncForPC(pc[i])
		callPath := strings.Split(f.Name(), `.init`)[0]
		splitStrings := strings.Split(callPath, `/`)
		orgName := splitStrings[1]
		repoName := splitStrings[2]
		packageName := splitStrings[len(splitStrings)-1]
		registerEnumName = fmt.Sprintf("%s|%s|%s|%s", orgName, repoName, packageName, enum.TypeName())
	}
	if registerEnumName == ""{
		registerEnumName = enum.TypeName()
	}
	if _, ok := m[registerEnumName]; ok {
		panic(fmt.Errorf("`%s` is already defined, please make enum name unqiue in one service", registerEnumName))
	}
	m[registerEnumName] = enum
}

func (m EnumMap) List() []EnumInfo {
	_ = require.Equal

	infoList := make([]EnumInfo, 0)

	for typeName, e := range m {
		options := make([]EnumOption, 0)

		for _, v := range e.ConstValues() {
			options = append(options, EnumOption{
				ConstValue: v.Int(),
				Value:      v.String(),
				Label:      v.Label(),
			})
		}

		infoList = append(infoList, EnumInfo{
			TypeName: typeName,
			Options:  options,
		})
	}

	return infoList
}

func ScanEnum(src interface{}, offset int) (int, error) {
	n, err := toInteger(src, offset)
	if err != nil {
		return offset, err
	}
	return n - offset, nil
}

func toInteger(src interface{}, defaultInteger int) (int, error) {
	switch v := src.(type) {
	case []byte:
		if len(v) > 0 {
			i, err := strconv.ParseInt(string(v), 10, 64)
			if err != nil {
				return defaultInteger, err
			}
			return int(i), err
		}
		return defaultInteger, nil
	case string:
		if v != "" {
			i, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return defaultInteger, err
			}
			return int(i), err
		}
		return defaultInteger, nil
	case int:
		return int(v), nil
	case int8:
		return int(v), nil
	case int16:
		return int(v), nil
	case int32:
		return int(v), nil
	case int64:
		return int(v), nil
	case uint:
		return int(v), nil
	case uint8:
		return int(v), nil
	case uint16:
		return int(v), nil
	case uint32:
		return int(v), nil
	case uint64:
		return int(v), nil
	case nil:
		return defaultInteger, nil
	default:
		return defaultInteger, nil
	}
}
