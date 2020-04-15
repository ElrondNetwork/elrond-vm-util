package vmtestjson

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEmpty(t *testing.T) {
	result, err := parseAnyValueAsByteArray("")
	require.Nil(t, err)
	require.Equal(t, []byte{}, result)
}

func TestBool(t *testing.T) {
	result, err := parseAnyValueAsByteArray("true")
	require.Nil(t, err)
	require.Equal(t, []byte{0x01}, result)

	result, err = parseAnyValueAsByteArray("false")
	require.Nil(t, err)
	require.Equal(t, []byte{}, result)
}

func TestString(t *testing.T) {
	result, err := parseAnyValueAsByteArray("``abcdefg")
	require.Nil(t, err)
	require.Equal(t, []byte("abcdefg"), result)

	result, err = parseAnyValueAsByteArray("``")
	require.Nil(t, err)
	require.Equal(t, []byte{}, result)

	result, err = parseAnyValueAsByteArray("```")
	require.Nil(t, err)
	require.Equal(t, []byte("`"), result)

	result, err = parseAnyValueAsByteArray("`` ")
	require.Nil(t, err)
	require.Equal(t, []byte(" "), result)
}

func TestUnsignedNumber(t *testing.T) {
	result, err := parseAnyValueAsByteArray("0x1234")
	require.Nil(t, err)
	require.Equal(t, []byte{0x12, 0x34}, result)

	result, err = parseAnyValueAsByteArray("0x")
	require.Nil(t, err)
	require.Equal(t, []byte{}, result)

	result, err = parseAnyValueAsByteArray("0")
	require.Nil(t, err)
	require.Equal(t, []byte{}, result)

	result, err = parseAnyValueAsByteArray("12")
	require.Nil(t, err)
	require.Equal(t, []byte{12}, result)

	result, err = parseAnyValueAsByteArray("256")
	require.Nil(t, err)
	require.Equal(t, []byte{0x01, 0x00}, result)

	result, err = parseAnyValueAsByteArray("0b1")
	require.Nil(t, err)
	require.Equal(t, []byte{0x01}, result)

	result, err = parseAnyValueAsByteArray("0b101")
	require.Nil(t, err)
	require.Equal(t, []byte{0x05}, result)
}

func TestSignedNumber(t *testing.T) {
	result, err := parseAnyValueAsByteArray("-1")
	require.Nil(t, err)
	require.Equal(t, []byte{0xff}, result)

	result, err = parseAnyValueAsByteArray("255")
	require.Nil(t, err)
	require.Equal(t, []byte{0xff}, result)

	result, err = parseAnyValueAsByteArray("+255")
	require.Nil(t, err)
	require.Equal(t, []byte{0x00, 0xff}, result)

	result, err = parseAnyValueAsByteArray("0xff")
	require.Nil(t, err)
	require.Equal(t, []byte{0xff}, result)

	result, err = parseAnyValueAsByteArray("+0xff")
	require.Nil(t, err)
	require.Equal(t, []byte{0x00, 0xff}, result)

	result, err = parseAnyValueAsByteArray("-256")
	require.Nil(t, err)
	require.Equal(t, []byte{0xff, 0x00}, result)

	result, err = parseAnyValueAsByteArray("-0b101")
	require.Nil(t, err)
	require.Equal(t, []byte{0xfb}, result)

}
