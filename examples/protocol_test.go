package examples

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProtocol(t *testing.T) {
	tt := assert.New(t)

	{
		var v Protocol
		err := v.Scan(-3)
		tt.NoError(err)
		tt.Equal(PROTOCOL__HTTP, v)
		dv, _ := v.Value()
		tt.Equal(-3, dv)
	}

	{
		var v Protocol
		err := v.Scan("-3")
		tt.NoError(err)
		tt.Equal(PROTOCOL__HTTP, v)
		dv, _ := v.Value()
		tt.Equal(-3, dv)
	}

	{
		var v Protocol
		err := v.Scan("")
		tt.NoError(err)
		tt.Equal(PROTOCOL_UNKNOWN, v)
	}

	{
		var v Protocol
		err := v.Scan(nil)
		tt.NoError(err)
		tt.Equal(PROTOCOL_UNKNOWN, v)
	}
}
