package vmtestjson

import (
	"errors"

	oj "github.com/ElrondNetwork/elrond-vm-util/test-util/orderedjson"
)

func processBlockTransaction(blrRaw oj.OJsonObject) (*Transaction, error) {
	bltMap, isMap := blrRaw.(*oj.OJsonMap)
	if !isMap {
		return nil, errors.New("unmarshalled block transaction is not a map")
	}

	blt := Transaction{}
	var nonceOk, functionOk, valueOk, argumentsOk, contractCodeOk, gasPriceOk, gasLimitOk bool

	for _, kvp := range bltMap.OrderedKV {

		switch kvp.Key {
		case "nonce":
			blt.Nonce, nonceOk = parseUint64(kvp.Value)
			if !nonceOk {
				return nil, errors.New("invalid block transaction nonce")
			}
		case "from":
			fromStr, fromOk := parseString(kvp.Value)
			if !fromOk {
				return nil, errors.New("invalid block transaction from")
			}
			var fromErr error
			blt.From, fromErr = parseAccountAddress(fromStr)
			if fromErr != nil {
				return nil, fromErr
			}

		case "to":
			toStr, toOk := parseString(kvp.Value)
			if !toOk {
				return nil, errors.New("invalid block transaction to")
			}

			// note "to": "0x00" has to yield isCreate=false, even though it parses to 0, just like the 2 cases below
			blt.IsCreate = toStr == "" || toStr == "0x"

			if !blt.IsCreate {
				var toErr error
				blt.To, toErr = parseAccountAddress(toStr)
				if toErr != nil {
					return nil, toErr
				}
			}
		case "function":
			blt.Function, functionOk = parseString(kvp.Value)
			if !functionOk {
				return nil, errors.New("invalid block transaction function")
			}
		case "value":
			blt.Value, valueOk = processBigInt(kvp.Value)
			if !valueOk {
				return nil, errors.New("invalid block transaction value")
			}
		case "arguments":
			blt.Arguments, argumentsOk = processArgumentList(kvp.Value)
			if !argumentsOk {
				return nil, errors.New("invalid block transaction arguments")
			}
		case "contractCode":
			blt.ContractCode, contractCodeOk = parseString(kvp.Value)
			if !contractCodeOk {
				return nil, errors.New("invalid block transaction contract code")
			}
		case "gasPrice":
			blt.GasPrice, gasPriceOk = parseUint64(kvp.Value)
			if !gasPriceOk {
				return nil, errors.New("invalid block transaction gasPrice")
			}
		case "gasLimit":
			blt.GasLimit, gasLimitOk = parseUint64(kvp.Value)
			if !gasLimitOk {
				return nil, errors.New("invalid block transaction gasLimit")
			}
		default:
			return nil, errors.New("unknown field in transaction")
		}
	}

	return &blt, nil
}

func processArgumentList(obj interface{}) ([][]byte, bool) {
	listRaw, listOk := obj.(*oj.OJsonList)
	if !listOk {
		return nil, false
	}
	var result [][]byte
	for _, elemRaw := range listRaw.AsList() {
		strRaw, strOk := parseString(elemRaw)
		if !strOk {
			return nil, false
		}
		arg, argErr := parseAnyValueAsByteArray(strRaw)
		if argErr != nil {
			return nil, false
		}
		result = append(result, arg)
	}
	return result, true
}
