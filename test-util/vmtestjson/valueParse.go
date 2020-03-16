package vmtestjson

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"

	oj "github.com/ElrondNetwork/elrond-vm-util/test-util/orderedjson"
)

func parseAccountAddress(addrRaw string) ([]byte, error) {
	if len(addrRaw) == 0 {
		return []byte{}, errors.New("missing account address")
	}
	return parseByteArray(addrRaw)
}

func parseByteArray(strRaw string) ([]byte, error) {
	if len(strRaw) == 0 {
		return []byte{}, nil
	}

	if strRaw == "false" {
		return []byte{0x00}, nil
	}

	if strRaw == "true" {
		return []byte{0x01}, nil
	}

	// hex, the usual representation
	if strings.HasPrefix(strRaw, "0x") || strings.HasPrefix(strRaw, "0X") {
		str := strRaw[2:]
		if len(str)%2 == 1 {
			str = "0" + str
		}
		return hex.DecodeString(str)
	}

	// allow ascii strings, for readability
	if strings.HasPrefix(strRaw, "``") {
		str := strRaw[2:]
		return []byte(str), nil
	}

	// default: parse as BigInt, base 10
	str := strings.ReplaceAll(strRaw, "_", "") // allow underscores, to group digits
	str = strings.ReplaceAll(str, ",", "")     // also allow commas to group digits
	result := new(big.Int)
	var parseOk bool
	result, parseOk = result.SetString(str, 10)
	if !parseOk {
		return []byte{}, fmt.Errorf("could not parse value: %s", strRaw)
	}

	return result.Bytes(), nil
}

func parseBigInt(strRaw string) (*big.Int, bool) {
	if len(strRaw) == 0 {
		return big.NewInt(0), true
	}

	negative := strRaw[0] == '-'

	if strRaw[0] == '-' || strRaw[0] == '+' {
		strRaw = strRaw[1:]
	}

	bytes, err := parseByteArray(strRaw)
	if err != nil {
		return nil, false
	}
	bi := big.NewInt(0).SetBytes(bytes)
	if negative {
		bi = bi.Neg(bi)
	}

	return bi, true
}

func parseUint64(obj oj.OJsonObject) (uint64, bool) {
	str, isStr := obj.(*oj.OJsonString)
	if !isStr {
		return uint64(0), false
	}
	if len(str.Value) == 0 {
		return uint64(0), true // interpret "" as nil, so we can restore to empty string instead of 0
	}

	bi, parseOk := processBigInt(obj)
	if !parseOk {
		return 0, false
	}

	if !bi.IsUint64() {
		return 0, false
	}

	return bi.Uint64(), true
}

func parseString(obj oj.OJsonObject) (string, bool) {
	str, isStr := obj.(*oj.OJsonString)
	if !isStr {
		return "", false
	}
	return str.Value, true
}

func isStar(obj oj.OJsonObject) bool {
	str, isStr := obj.(*oj.OJsonString)
	if !isStr {
		return false
	}
	return str.Value == "*"
}
