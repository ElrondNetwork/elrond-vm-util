package vmtestjson

import (
	"errors"
	"math/big"

	twos "github.com/ElrondNetwork/big-int-util/twos-complement"
	oj "github.com/ElrondNetwork/elrond-vm-util/test-util/orderedjson"
)

type bigIntParseFormat int

const (
	bigIntSignedBytes   = iota
	bigIntUnsignedBytes = iota
)

func (p *Parser) processBigInt(obj oj.OJsonObject, format bigIntParseFormat) (*big.Int, error) {
	strVal, err := p.parseString(obj)
	if err != nil {
		return nil, err
	}
	return p.parseBigInt(strVal, format)
}

func (p *Parser) parseBigInt(strRaw string, format bigIntParseFormat) (*big.Int, error) {
	bytes, err := p.parseAnyValueAsByteArray(strRaw)
	if err != nil {
		return nil, err
	}
	switch format {
	case bigIntSignedBytes:
		return twos.FromBytes(bytes), nil
	case bigIntUnsignedBytes:
		return big.NewInt(0).SetBytes(bytes), nil
	default:
		return nil, errors.New("unknown format requested")
	}
}

func (p *Parser) parseUint64(obj oj.OJsonObject) (uint64, error) {
	bi, err := p.processBigInt(obj, bigIntUnsignedBytes)
	if err != nil {
		return 0, err
	}

	if !bi.IsUint64() {
		return 0, errors.New("value is not uint64")
	}

	return bi.Uint64(), nil
}

func (p *Parser) parseString(obj oj.OJsonObject) (string, error) {
	str, isStr := obj.(*oj.OJsonString)
	if !isStr {
		return "", errors.New("not a string value")
	}
	return str.Value, nil
}

func isStar(obj oj.OJsonObject) bool {
	str, isStr := obj.(*oj.OJsonString)
	if !isStr {
		return false
	}
	return str.Value == "*"
}
