package enumeration

import (
	"fmt"
	"strconv"
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
	typeName := enum.TypeName()
	if _, ok := m[typeName]; ok {
		panic(fmt.Errorf("`%s` is already defined, please make enum name unqiue in one service", typeName))
	}
	m[typeName] = enum
}

func (m EnumMap) List() []EnumInfo {
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
