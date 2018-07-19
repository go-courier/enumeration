package enumeration_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/go-courier/enumeration"
	"github.com/go-courier/enumeration/__examples__"
)

func TestEnumMap(t *testing.T) {
	require.Equal(t, examples.PROTOCOL__HTTP.String(), "HTTP")

	list := enumeration.DefaultEnumMap.List()

	require.Equal(t, list, []enumeration.EnumInfo{
		{
			TypeName: "Protocol",
			Options: []enumeration.EnumOption{
				{
					Value:      "HTTP",
					Label:      "http",
					ConstValue: examples.PROTOCOL__HTTP.Int(),
				},
				{
					Value:      "HTTPS",
					Label:      "https",
					ConstValue: examples.PROTOCOL__HTTPS.Int(),
				},
				{
					Value:      "TCP",
					Label:      "TCP",
					ConstValue: examples.PROTOCOL__TCP.Int(),
				},
			},
		},
	})
}

func TestScanEnum(t *testing.T) {
	cases := []struct {
		offset int
		values []interface{}
		expect []int
	}{
		{
			-3,
			[]interface{}{
				nil,
				[]byte("-3"),
				"-2",
				int(-1),
				int8(0),
				int16(1),
				int32(2),
				int64(3),
				uint(4),
				uint8(5),
				uint16(6),
				uint32(7),
				uint64(8),
			},
			[]int{
				0,
				0,
				1,
				2,
				3,
				4,
				5,
				6,
				7,
				8,
				9,
				10,
				11,
			},
		},
	}

	for _, c := range cases {
		for i, v := range c.values {
			n, err := enumeration.ScanEnum(v, c.offset)
			require.NoError(t, err)
			require.Equal(t, c.expect[i], n)
		}
	}
}
