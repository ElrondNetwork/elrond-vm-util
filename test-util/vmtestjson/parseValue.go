package vmtestjson

import (
	"math/big"

	oj "github.com/ElrondNetwork/elrond-vm-util/test-util/orderedjson"
)

func (p *Parser) parseBigInt(strRaw string) (*big.Int, bool) {
	if len(strRaw) == 0 {
		return big.NewInt(0), true
	}

	negative := strRaw[0] == '-'

	if strRaw[0] == '-' || strRaw[0] == '+' {
		strRaw = strRaw[1:]
	}

	bytes, err := p.parseAnyValueAsByteArray(strRaw)
	if err != nil {
		return nil, false
	}
	bi := big.NewInt(0).SetBytes(bytes)
	if negative {
		bi = bi.Neg(bi)
	}

	return bi, true
}

func (p *Parser) parseUint64(obj oj.OJsonObject) (uint64, bool) {
	str, isStr := obj.(*oj.OJsonString)
	if !isStr {
		return uint64(0), false
	}
	if len(str.Value) == 0 {
		return uint64(0), true // interpret "" as nil, so we can restore to empty string instead of 0
	}

	bi, parseOk := p.processBigInt(obj)
	if !parseOk {
		return 0, false
	}

	if !bi.IsUint64() {
		return 0, false
	}

	return bi.Uint64(), true
}

func (p *Parser) parseString(obj oj.OJsonObject) (string, bool) {
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
