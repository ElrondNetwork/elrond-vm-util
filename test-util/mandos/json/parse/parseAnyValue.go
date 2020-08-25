package mandosjsonparse

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"

	twos "github.com/ElrondNetwork/big-int-util/twos-complement"
	mj "github.com/ElrondNetwork/elrond-vm-util/test-util/mandos/json/model"
	oj "github.com/ElrondNetwork/elrond-vm-util/test-util/orderedjson"
)

var strPrefixes = []string{"str:", "``", "''"}

const addrPrefix = "address:"
const filePrefix = "file:"
const keccak256Prefix = "keccak256:"

const u64Prefix = "u64:"
const u32Prefix = "u32:"
const u16Prefix = "u16:"
const u8Prefix = "u8:"
const i64Prefix = "i64:"
const i32Prefix = "i32:"
const i16Prefix = "i16:"
const i8Prefix = "i8:"

func (p *Parser) parseCheckBytes(obj oj.OJsonObject) (mj.JSONCheckBytes, error) {
	if IsStar(obj) {
		// "*" means any value, skip checking it
		return mj.JSONCheckBytes{
			Value:    nil,
			IsStar:   true,
			Original: "*"}, nil
	}

	jb, err := p.processAnyValueAsByteArray(obj)
	if err != nil {
		return mj.JSONCheckBytes{}, err
	}
	return mj.JSONCheckBytes{
		Value:    jb.Value,
		IsStar:   false,
		Original: jb.Original,
	}, nil
}

func (p *Parser) processAnyValueAsByteArray(obj oj.OJsonObject) (mj.JSONBytes, error) {
	strVal, err := p.parseString(obj)
	if err != nil {
		return mj.JSONBytes{}, err
	}
	result, err := p.parseAnyValueAsByteArray(strVal)
	return mj.JSONBytes{
		Value:    result,
		Original: strVal,
	}, err
}

func (p *Parser) parseAnyValueAsByteArray(strRaw string) ([]byte, error) {
	if len(strRaw) == 0 {
		return []byte{}, nil
	}

	// file contents
	// TODO: make this part of a proper parser
	if strings.HasPrefix(strRaw, filePrefix) {
		if p.FileResolver == nil {
			return []byte{}, errors.New("parser FileResolver not provided")
		}
		fileContents, err := p.FileResolver.ResolveFileValue(strRaw[len(filePrefix):])
		if err != nil {
			return []byte{}, err
		}
		return fileContents, nil
	}

	// keccak256
	// TODO: make this part of a proper parser
	if strings.HasPrefix(strRaw, keccak256Prefix) {
		arg, err := p.parseAnyValueAsByteArray(strRaw[len(keccak256Prefix):])
		if err != nil {
			return []byte{}, fmt.Errorf("cannot parse keccak256 argument: %w", err)
		}
		hash, err := keccak256(arg)
		if err != nil {
			return []byte{}, fmt.Errorf("error computing keccak256: %w", err)
		}
		return hash, nil
	}

	// concatenate values of different formats
	// TODO: make this part of a proper parser
	parts := strings.Split(strRaw, "|")
	if len(parts) > 1 {
		concat := make([]byte, 0)
		for _, part := range parts {
			eval, err := p.parseAnyValueAsByteArray(part)
			if err != nil {
				return []byte{}, err
			}
			concat = append(concat, eval...)
		}
		return concat, nil
	}

	if strRaw == "false" {
		return []byte{}, nil
	}

	if strRaw == "true" {
		return []byte{0x01}, nil
	}

	// allow ascii strings, for readability
	for _, strPrefix := range strPrefixes {
		if strings.HasPrefix(strRaw, strPrefix) {
			str := strRaw[len(strPrefix):]
			return []byte(str), nil
		}
	}

	// address
	if strings.HasPrefix(strRaw, addrPrefix) {
		addrName := strRaw[len(addrPrefix):]
		return address([]byte(addrName))
	}

	// fixed width numbers
	parsed, result, err := p.parseFixedWidthBasicTypes(strRaw)
	if err != nil {
		return nil, err
	}
	if parsed {
		return result, nil
	}

	// general numbers, arbitrary length
	return p.parseAnyNumberAsByteArray(strRaw, 0)
}

// targetWidth = 0 means minimum length that can contain the result
func (p *Parser) parseAnyNumberAsByteArray(strRaw string, targetWidth int) ([]byte, error) {
	// signed numbers
	if strRaw[0] == '-' || strRaw[0] == '+' {
		numberBytes, err := p.parseUnsignedNumberAsByteArray(strRaw[1:])
		if err != nil {
			return []byte{}, err
		}
		number := big.NewInt(0).SetBytes(numberBytes)
		if strRaw[0] == '-' {
			number = number.Neg(number)
		}
		if targetWidth == 0 {
			return twos.ToBytes(number), nil
		}

		return twos.ToBytesOfLength(number, targetWidth)
	}

	// unsigned numbers
	return p.parseUnsignedNumberAsByteArrayOfLength(strRaw, targetWidth)
}

func (p *Parser) parseUnsignedNumberAsByteArrayOfLength(strRaw string, targetWidth int) ([]byte, error) {
	numberBytes, err := p.parseUnsignedNumberAsByteArray(strRaw)
	if err != nil {
		return []byte{}, err
	}
	if targetWidth == 0 {
		return numberBytes, nil
	}

	if len(numberBytes) > targetWidth {
		return []byte{}, fmt.Errorf("representation of %s does not fit in %d bytes", strRaw, targetWidth)
	}
	return twos.CopyAlignRight(numberBytes, targetWidth), nil
}

func (p *Parser) parseUnsignedNumberAsByteArray(strRaw string) ([]byte, error) {
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

func (p *Parser) parseFixedWidthBasicTypes(strRaw string) (bool, []byte, error) {
	if strings.HasPrefix(strRaw, u64Prefix) {
		r, err := p.parseUnsignedNumberAsByteArrayOfLength(strRaw[len(u64Prefix):], 8)
		return true, r, err
	}
	if strings.HasPrefix(strRaw, u32Prefix) {
		r, err := p.parseUnsignedNumberAsByteArrayOfLength(strRaw[len(u32Prefix):], 4)
		return true, r, err
	}
	if strings.HasPrefix(strRaw, u16Prefix) {
		r, err := p.parseUnsignedNumberAsByteArrayOfLength(strRaw[len(u16Prefix):], 2)
		return true, r, err
	}
	if strings.HasPrefix(strRaw, u8Prefix) {
		r, err := p.parseUnsignedNumberAsByteArrayOfLength(strRaw[len(u8Prefix):], 1)
		return true, r, err
	}

	if strings.HasPrefix(strRaw, i64Prefix) {
		r, err := p.parseAnyNumberAsByteArray(strRaw[len(i64Prefix):], 8)
		return true, r, err
	}
	if strings.HasPrefix(strRaw, i32Prefix) {
		r, err := p.parseAnyNumberAsByteArray(strRaw[len(i32Prefix):], 4)
		return true, r, err
	}
	if strings.HasPrefix(strRaw, i16Prefix) {
		r, err := p.parseAnyNumberAsByteArray(strRaw[len(i16Prefix):], 2)
		return true, r, err
	}
	if strings.HasPrefix(strRaw, i8Prefix) {
		r, err := p.parseAnyNumberAsByteArray(strRaw[len(i8Prefix):], 1)
		return true, r, err
	}

	return false, []byte{}, nil
}
