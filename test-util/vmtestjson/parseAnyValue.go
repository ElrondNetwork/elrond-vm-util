package vmtestjson

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	twos "github.com/ElrondNetwork/big-int-util/twos-complement"
)

func parseAnyValueAsByteArray(strRaw string) ([]byte, error) {
	if len(strRaw) == 0 {
		return []byte{}, nil
	}

	if strRaw == "false" {
		return []byte{}, nil
	}

	if strRaw == "true" {
		return []byte{0x01}, nil
	}

	// allow ascii strings, for readability
	if strings.HasPrefix(strRaw, "``") {
		str := strRaw[2:]
		return []byte(str), nil
	}

	// signed numbers
	if strRaw[0] == '-' || strRaw[0] == '+' {
		numberBytes, err := parseUnsignedNumberAsByteArray(strRaw[1:])
		if err != nil {
			return []byte{}, err
		}
		number := big.NewInt(0).SetBytes(numberBytes)
		if strRaw[0] == '-' {
			number = number.Neg(number)
		}
		return twos.ToBytes(number), nil
	}

	// unsigned numbers
	return parseUnsignedNumberAsByteArray(strRaw)
}

func parseUnsignedNumberAsByteArray(strRaw string) ([]byte, error) {
	str := strings.ReplaceAll(strRaw, "_", "") // allow underscores, to group digits
	str = strings.ReplaceAll(str, ",", "")     // also allow commas to group digits

	// hex, the usual representation
	if strings.HasPrefix(strRaw, "0x") || strings.HasPrefix(strRaw, "0X") {
		str := strRaw[2:]
		if len(str)%2 == 1 {
			str = "0" + str
		}
		return hex.DecodeString(str)
	}

	// binary representation
	if strings.HasPrefix(strRaw, "0b") || strings.HasPrefix(strRaw, "0B") {
		result := new(big.Int)
		var parseOk bool
		result, parseOk = result.SetString(str[2:], 2)
		if !parseOk {
			return []byte{}, fmt.Errorf("could not parse binary value: %s", strRaw)
		}

		return result.Bytes(), nil
	}

	// default: parse as BigInt, base 10
	result := new(big.Int)
	var parseOk bool
	result, parseOk = result.SetString(str, 10)
	if !parseOk {
		return []byte{}, fmt.Errorf("could not parse base 10 value: %s", strRaw)
	}

	return result.Bytes(), nil
}
