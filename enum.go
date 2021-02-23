package enumeration

type IntStringerEnum interface {
	TypeName() string
	Int() int
	String() string
	Label() string
	ConstValues() []IntStringerEnum
}

// Deprecated use IntStringerEnum instead
type Enum = IntStringerEnum

// sql value of enum maybe have offset from value of enum in go
type DriverValueOffset interface {
	Offset() int
}

// Deprecated use DriverValueOffset instead
type EnumDriverValueOffset = DriverValueOffset
