package examples

import (
	bytes "bytes"
	database_sql_driver "database/sql/driver"
	errors "errors"

	github_com_go_courier_enumeration "github.com/go-courier/enumeration"
)

var InvalidProtocol = errors.New("invalid Protocol type")

func ParseProtocolFromLabelString(s string) (Protocol, error) {
	switch s {
	case "":
		return PROTOCOL_UNKNOWN, nil
	case "http":
		return PROTOCOL__HTTP, nil
	case "https":
		return PROTOCOL__HTTPS, nil
	case "TCP":
		return PROTOCOL__TCP, nil
	}
	return PROTOCOL_UNKNOWN, InvalidProtocol
}

func (v Protocol) String() string {
	switch v {
	case PROTOCOL_UNKNOWN:
		return ""
	case PROTOCOL__HTTP:
		return "HTTP"
	case PROTOCOL__HTTPS:
		return "HTTPS"
	case PROTOCOL__TCP:
		return "TCP"
	}
	return "UNKNOWN"
}

func ParseProtocolFromString(s string) (Protocol, error) {
	switch s {
	case "":
		return PROTOCOL_UNKNOWN, nil
	case "HTTP":
		return PROTOCOL__HTTP, nil
	case "HTTPS":
		return PROTOCOL__HTTPS, nil
	case "TCP":
		return PROTOCOL__TCP, nil
	}
	return PROTOCOL_UNKNOWN, InvalidProtocol
}

func (v Protocol) Label() string {
	switch v {
	case PROTOCOL_UNKNOWN:
		return ""
	case PROTOCOL__HTTP:
		return "http"
	case PROTOCOL__HTTPS:
		return "https"
	case PROTOCOL__TCP:
		return "TCP"
	}
	return "UNKNOWN"
}

func (v Protocol) Int() int {
	return int(v)
}

func (Protocol) TypeName() string {
	return "github.com/go-courier/enumeration/__examples__.Protocol"
}

func (Protocol) ConstValues() []github_com_go_courier_enumeration.IntStringerEnum {
	return []github_com_go_courier_enumeration.IntStringerEnum{PROTOCOL__HTTP, PROTOCOL__HTTPS, PROTOCOL__TCP}
}

func (v Protocol) MarshalText() ([]byte, error) {
	str := v.String()
	if str == "UNKNOWN" {
		return nil, InvalidProtocol
	}
	return []byte(str), nil
}

func (v *Protocol) UnmarshalText(data []byte) (err error) {
	*v, err = ParseProtocolFromString(string(bytes.ToUpper(data)))
	return
}

func (v Protocol) Value() (database_sql_driver.Value, error) {
	offset := 0
	if o, ok := (interface{})(v).(github_com_go_courier_enumeration.DriverValueOffset); ok {
		offset = o.Offset()
	}
	return int64(v) + int64(offset), nil
}

func (v *Protocol) Scan(src interface{}) error {
	offset := 0
	if o, ok := (interface{})(v).(github_com_go_courier_enumeration.DriverValueOffset); ok {
		offset = o.Offset()
	}

	i, err := github_com_go_courier_enumeration.ScanIntEnumStringer(src, offset)
	if err != nil {
		return err
	}
	*v = Protocol(i)
	return nil
}
