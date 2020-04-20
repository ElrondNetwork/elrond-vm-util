package vmtestjson

import (
	"errors"
	"fmt"

	oj "github.com/ElrondNetwork/elrond-vm-util/test-util/orderedjson"
)

func (p *Parser) processTx(blrRaw oj.OJsonObject) (*Transaction, error) {
	bltMap, isMap := blrRaw.(*oj.OJsonMap)
	if !isMap {
		return nil, errors.New("unmarshalled block transaction is not a map")
	}

	blt := Transaction{}
	var err error
	for _, kvp := range bltMap.OrderedKV {

		switch kvp.Key {
		case "nonce":
			blt.Nonce, err = p.processUint64(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid block transaction nonce: %w", err)
			}
		case "from":
			fromStr, err := p.parseString(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid block transaction from: %w", err)
			}
			var fromErr error
			blt.From, fromErr = p.parseAccountAddress(fromStr)
			if fromErr != nil {
				return nil, fromErr
			}

		case "to":
			toStr, err := p.parseString(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid block transaction to: %w", err)
			}

			// note "to": "0x00" has to yield isCreate=false, even though it parses to 0, just like the 2 cases below
			blt.IsCreate = toStr == "" || toStr == "0x"

			if !blt.IsCreate {
				var toErr error
				blt.To, toErr = p.parseAccountAddress(toStr)
				if toErr != nil {
					return nil, toErr
				}
			}
		case "function":
			blt.Function, err = p.parseString(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid block transaction function: %w", err)
			}
		case "value":
			blt.Value, err = p.processBigInt(kvp.Value, bigIntUnsignedBytes)
			if err != nil {
				return nil, fmt.Errorf("invalid block transaction value: %w", err)
			}
		case "arguments":
			blt.Arguments, err = p.parseByteArrayList(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid block transaction arguments: %w", err)
			}
		case "contractCode":
			blt.Code, err = p.processAnyValueAsByteArray(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid block transaction contract code: %w", err)
			}
		case "gasPrice":
			blt.GasPrice, err = p.processUint64(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid block transaction gasPrice: %w", err)
			}
		case "gasLimit":
			blt.GasLimit, err = p.processUint64(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid block transaction gasLimit: %w", err)
			}
		default:
			return nil, fmt.Errorf("unknown field in transaction: %w", err)
		}
	}

	return &blt, nil
}
