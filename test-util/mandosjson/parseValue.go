package mandosjson

import (
	"errors"
	"math/big"

	twos "github.com/ElrondNetwork/big-int-util/twos-complement"
	oj "github.com/ElrondNetwork/elrond-vm-util/test-util/orderedjson"
)

type bigIntParseFormat int

const (
	bigIntSignedBytes bigIntParseFormat = iota
	bigIntUnsignedBytes
)

func (p *Parser) processCheckBigInt(obj oj.OJsonObject, format bigIntParseFormat) (JSONCheckBigInt, error) {
	if isStar(obj) {
		// "*" means any value, skip checking it
		return JSONCheckBigInt{
			value:    nil,
			isStar:   true,
			Original: "*"}, nil
	}

	jbi, err := p.processBigInt(obj, format)
	if err != nil {
		return JSONCheckBigInt{}, err
	}
	return JSONCheckBigInt{
		value:    jbi.Value,
		isStar:   false,
		Original: jbi.Original,
	}, nil
}

func (p *Parser) processBigInt(obj oj.OJsonObject, format bigIntParseFormat) (JSONBigInt, error) {
	strVal, err := p.parseString(obj)
	if err != nil {
		return JSONBigInt{}, err
	}

	bi, err := p.parseBigInt(strVal, format)
	return JSONBigInt{
		Value:    bi,
		Original: strVal,
	}, err
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

func (p *Parser) processCheckUint64(obj oj.OJsonObject) (JSONCheckUint64, error) {
	if isStar(obj) {
		// "*" means any value, skip checking it
		return JSONCheckUint64{
			value:    0,
			isStar:   true,
			Original: "*"}, nil
	}

	ju, err := p.processUint64(obj)
	if err != nil {
		return JSONCheckUint64{}, err
	}
	return JSONCheckUint64{
		value:    ju.Value,
		isStar:   false,
		Original: ju.Original}, nil

}

func (p *Parser) processUint64(obj oj.OJsonObject) (JSONUint64, error) {
	bi, err := p.processBigInt(obj, bigIntUnsignedBytes)
	if err != nil {
		return JSONUint64{}, err
	}

	if bi.Value == nil || !bi.Value.IsUint64() {
		return JSONUint64{}, errors.New("value is not uint64")
	}

	return JSONUint64{
		Value:    bi.Value.Uint64(),
		Original: bi.Original}, nil
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
