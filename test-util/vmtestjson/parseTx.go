package vmtestjson

import (
	"errors"

	oj "github.com/ElrondNetwork/elrond-vm-util/test-util/orderedjson"
)

func (p *Parser) processBlockTransaction(blrRaw oj.OJsonObject) (*Transaction, error) {
	bltMap, isMap := blrRaw.(*oj.OJsonMap)
	if !isMap {
		return nil, errors.New("unmarshalled block transaction is not a map")
	}

	blt := Transaction{}
	var nonceOk, functionOk, valueOk, argumentsOk, contractCodeOk, gasPriceOk, gasLimitOk bool

	for _, kvp := range bltMap.OrderedKV {

		switch kvp.Key {
		case "nonce":
			blt.Nonce, nonceOk = p.parseUint64(kvp.Value)
			if !nonceOk {
				return nil, errors.New("invalid block transaction nonce")
			}
		case "from":
			fromStr, fromOk := p.parseString(kvp.Value)
			if !fromOk {
				return nil, errors.New("invalid block transaction from")
			}
			var fromErr error
			blt.From, fromErr = p.parseAccountAddress(fromStr)
			if fromErr != nil {
				return nil, fromErr
			}

		case "to":
			toStr, toOk := p.parseString(kvp.Value)
			if !toOk {
				return nil, errors.New("invalid block transaction to")
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
			blt.Function, functionOk = p.parseString(kvp.Value)
			if !functionOk {
				return nil, errors.New("invalid block transaction function")
			}
		case "value":
			blt.Value, valueOk = p.processBigInt(kvp.Value)
			if !valueOk {
				return nil, errors.New("invalid block transaction value")
			}
		case "arguments":
			blt.Arguments, argumentsOk = p.processArgumentList(kvp.Value)
			if !argumentsOk {
				return nil, errors.New("invalid block transaction arguments")
			}
		case "contractCode":
			blt.ContractCode, contractCodeOk = p.parseString(kvp.Value)
			if !contractCodeOk {
				return nil, errors.New("invalid block transaction contract code")
			}
		case "gasPrice":
			blt.GasPrice, gasPriceOk = p.parseUint64(kvp.Value)
			if !gasPriceOk {
				return nil, errors.New("invalid block transaction gasPrice")
			}
		case "gasLimit":
			blt.GasLimit, gasLimitOk = p.parseUint64(kvp.Value)
			if !gasLimitOk {
				return nil, errors.New("invalid block transaction gasLimit")
			}
		default:
			return nil, errors.New("unknown field in transaction")
		}
	}

	return &blt, nil
}

func (p *Parser) processArgumentList(obj interface{}) ([][]byte, bool) {
	listRaw, listOk := obj.(*oj.OJsonList)
	if !listOk {
		return nil, false
	}
	var result [][]byte
	for _, elemRaw := range listRaw.AsList() {
		strRaw, strOk := p.parseString(elemRaw)
		if !strOk {
			return nil, false
		}
		arg, argErr := p.parseAnyValueAsByteArray(strRaw)
		if argErr != nil {
			return nil, false
		}
		result = append(result, arg)
	}
	return result, true
}
