package scanner

import (
	"fmt"
	"strconv"
)

func NewIntStringerOption(i int64, v string, label string) *Option {
	o := &Option{
		Int:   &i,
		Str:   &v,
		Label: label,
	}

	if label == "" {
		o.Label = v
	}

	return o
}

func NewStrOption(v string, label string) *Option {
	o := &Option{
		Str:   &v,
		Label: label,
	}

	if label == "" {
		o.Label = v
	}
	return o
}

func NewIntOption(v int64, label string) *Option {
	o := &Option{
		Int:   &v,
		Label: label,
	}

	if label == "" {
		o.Label = strconv.FormatInt(v, 64)
	}
	return o
}

func NewFloatOption(v float64, label string) *Option {
	o := &Option{
		Float: &v,
		Label: label,
	}

	if label == "" {
		o.Label = fmt.Sprintf("%v", v)
	}
	return o
}

type Option struct {
	Label string   `json:"label"`
	Str   *string  `json:"str,omitempty"`
	Int   *int64   `json:"int,omitempty"`
	Float *float64 `json:"float,omitempty"`
}

func (o Option) Value() interface{} {
	if o.Str != nil {
		return *o.Str
	}
	if o.Float != nil {
		return *o.Float
	}
	if o.Int != nil {
		return *o.Int
	}
	return nil
}

type Options []Option

func (o Options) Len() int {
	return len(o)
}

func (o Options) Values() []interface{} {
	values := make([]interface{}, len(o))

	for i, v := range o {
		values[i] = v.Value()
	}

	return values
}

func (o Options) Less(i, j int) bool {
	if o[i].Float != nil {
		return *o[i].Float < *o[j].Float
	}
	if o[i].Int != nil {
		return *o[i].Int < *o[j].Int
	}
	return *o[i].Str < *o[j].Str
}

func (o Options) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}
