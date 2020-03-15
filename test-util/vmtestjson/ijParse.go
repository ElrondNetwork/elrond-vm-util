package vmtestjson

import (
	"errors"
	"math/big"

	oj "github.com/ElrondNetwork/elrond-vm-util/test-util/orderedjson"
)

// ParseTopLevel converts json string to object representation
func ParseTopLevel(jsonString []byte) ([]*Test, error) {

	jobj, err := oj.ParseOrderedJSON(jsonString)
	if err != nil {
		return nil, err
	}

	topMap, isMap := jobj.(*oj.OJsonMap)
	if !isMap {
		return nil, errors.New("unmarshalled test top level object is not a map")
	}

	var top []*Test
	for _, kvp := range topMap.OrderedKV {
		t, tErr := processTest(kvp.Value)
		if tErr != nil {
			return nil, tErr
		}
		t.TestName = kvp.Key
		top = append(top, t)
	}
	return top, nil
}

func processTest(testObj oj.OJsonObject) (*Test, error) {
	testMap, isTestMap := testObj.(*oj.OJsonMap)
	if !isTestMap {
		return nil, errors.New("unmarshalled test object is not a map")
	}
	test := Test{CheckGas: true}

	for _, kvp := range testMap.OrderedKV {
		if kvp.Key == "checkGas" {
			checkGasOJ, isBool := kvp.Value.(*oj.OJsonBool)
			if !isBool {
				return nil, errors.New("unmarshalled test checkGas flag is not boolean")
			}
			test.CheckGas = bool(*checkGasOJ)
		}

		if kvp.Key == "pre" {
			preMap, isPreMap := kvp.Value.(*oj.OJsonMap)
			if !isPreMap {
				return nil, errors.New("unmarshalled pre object is not a map")
			}
			for _, acctKVP := range preMap.OrderedKV {
				acct, acctErr := processAccount(acctKVP.Value)
				if acctErr != nil {
					return nil, acctErr
				}
				acctAddr, hexErr := parseAccountAddress(acctKVP.Key)
				if hexErr != nil {
					return nil, hexErr
				}
				acct.Address = acctAddr
				test.Pre = append(test.Pre, acct)
			}
		}

		if kvp.Key == "blocks" {
			blocksRaw, blocksOk := kvp.Value.(*oj.OJsonList)
			if !blocksOk {
				return nil, errors.New("unmarshalled blocks object is not a list")
			}
			for _, blRaw := range blocksRaw.AsList() {
				bl, blErr := processBlock(blRaw)
				if blErr != nil {
					return nil, blErr
				}
				test.Blocks = append(test.Blocks, bl)
			}
		}

		if kvp.Key == "network" {
			var networkOk bool
			test.Network, networkOk = parseString(kvp.Value)
			if !networkOk {
				return nil, errors.New("test network value not a string")
			}
		}

		if kvp.Key == "blockhashes" {
			var bhsOk bool
			test.BlockHashes, bhsOk = parseByteArrayList(kvp.Value)
			if !bhsOk {
				return nil, errors.New("unmarshalled blockHashes object is not a list")
			}
		}

		if kvp.Key == "postState" {
			postMap, isPostMap := kvp.Value.(*oj.OJsonMap)
			if !isPostMap {
				return nil, errors.New("unmarshalled postState object is not a map")
			}
			for _, acctKVP := range postMap.OrderedKV {
				acct, acctErr := processAccount(acctKVP.Value)
				if acctErr != nil {
					return nil, acctErr
				}
				acctAddr, hexErr := parseAccountAddress(acctKVP.Key)
				if hexErr != nil {
					return nil, hexErr
				}
				acct.Address = acctAddr
				test.PostState = append(test.PostState, acct)
			}
		}
	}

	return &test, nil
}

func processAccount(acctRaw oj.OJsonObject) (*Account, error) {
	acctMap, isMap := acctRaw.(*oj.OJsonMap)
	if !isMap {
		return nil, errors.New("unmarshalled account object is not a map")
	}

	acct := Account{}
	var nonceOk, balanceOk, codeOk, dataOk bool

	for _, kvp := range acctMap.OrderedKV {

		if kvp.Key == "nonce" {
			acct.Nonce, nonceOk = processBigInt(kvp.Value)
			if !nonceOk {
				return nil, errors.New("invalid account nonce")
			}
		}

		if kvp.Key == "balance" {
			acct.Balance, balanceOk = processBigInt(kvp.Value)
			if !balanceOk {
				return nil, errors.New("invalid account balance")
			}
		}

		if kvp.Key == "storage" {
			storageMap, storageOk := kvp.Value.(*oj.OJsonMap)
			if !storageOk {
				return nil, errors.New("invalid account storage")
			}
			for _, storageKvp := range storageMap.OrderedKV {
				byteKey, keyOk := parseByteArray(storageKvp.Key)
				if keyOk != nil {
					return nil, errors.New("invalid account storage key")
				}
				strVal, valStrOk := parseString(storageKvp.Value)
				if !valStrOk {
					return nil, errors.New("invalid account storage value")
				}
				byteVal, valOk := parseByteArray(strVal)
				if valOk != nil {
					return nil, errors.New("invalid account storage value")
				}
				stElem := StorageKeyValuePair{
					Key:   byteKey,
					Value: byteVal,
				}
				acct.Storage = append(acct.Storage, &stElem)
			}
		}

		if kvp.Key == "code" {
			acct.Code, codeOk = parseString(kvp.Value)
			if !codeOk {
				return nil, errors.New("invalid account code")
			}
			acct.OriginalCode = acct.Code
		}

		if kvp.Key == "asyncCallData" {
			acct.AsyncCallData, dataOk = parseString(kvp.Value)
			if !dataOk {
				return nil, errors.New("invalid asyncCallData string")
			}
		}
	}

	return &acct, nil
}

func processBlock(blockRaw oj.OJsonObject) (*Block, error) {
	blockMap, isMap := blockRaw.(*oj.OJsonMap)
	if !isMap {
		return nil, errors.New("unmarshalled block object is not a map")
	}
	bl := Block{}

	for _, kvp := range blockMap.OrderedKV {

		if kvp.Key == "results" {
			resultsRaw, resultsOk := kvp.Value.(*oj.OJsonList)
			if !resultsOk {
				return nil, errors.New("unmarshalled block results object is not a list")
			}
			for _, resRaw := range resultsRaw.AsList() {
				blr, blrErr := processBlockResult(resRaw)
				if blrErr != nil {
					return nil, blrErr
				}
				bl.Results = append(bl.Results, blr)
			}
		}

		if kvp.Key == "transactions" {
			transactionsRaw, transactionsOk := kvp.Value.(*oj.OJsonList)
			if !transactionsOk {
				return nil, errors.New("unmarshalled block transactions object is not a list")
			}
			for _, trRaw := range transactionsRaw.AsList() {
				tr, trErr := processBlockTransaction(trRaw)
				if trErr != nil {
					return nil, trErr
				}
				bl.Transactions = append(bl.Transactions, tr)
			}
		}

		if kvp.Key == "blockHeader" {
			blh, blhErr := processBlockHeader(kvp.Value)
			if blhErr != nil {
				return nil, blhErr
			}
			bl.BlockHeader = blh
		}
	}

	if len(bl.Results) != len(bl.Transactions) {
		return nil, errors.New("mismatched number of blocks and transactions")
	}

	return &bl, nil
}

func processBlockResult(blrRaw oj.OJsonObject) (*TransactionResult, error) {
	blrMap, isMap := blrRaw.(*oj.OJsonMap)
	if !isMap {
		return nil, errors.New("unmarshalled block result is not a map")
	}

	blr := TransactionResult{}
	var outOk, statusOk, messageOk, gasOk, refundOk bool

	for _, kvp := range blrMap.OrderedKV {

		if kvp.Key == "out" {
			blr.Out, outOk = parseByteArrayList(kvp.Value)
			if !outOk {
				return nil, errors.New("invalid block result out")
			}
		}

		if kvp.Key == "status" {
			blr.Status, statusOk = processBigInt(kvp.Value)
			if !statusOk {
				return nil, errors.New("invalid block result status")
			}
		}

		if kvp.Key == "message" {
			blr.Message, messageOk = parseString(kvp.Value)
			if !messageOk {
				return nil, errors.New("invalid block result message")
			}
		}

		if kvp.Key == "gas" {
			if isStar(kvp.Value) {
				blr.CheckGas = false
				blr.Gas = 0
			} else {
				blr.CheckGas = true
				blr.Gas, gasOk = parseUint64(kvp.Value)
				if !gasOk {
					return nil, errors.New("invalid block result gas")
				}
			}
		}

		if kvp.Key == "logs" {
			if isStar(kvp.Value) {
				blr.IgnoreLogs = true
			} else {
				blr.IgnoreLogs = false
				var logHashOk bool
				blr.LogHash, logHashOk = parseString(kvp.Value)
				if !logHashOk {
					var logListErr error
					blr.Logs, logListErr = processLogList(kvp.Value)
					if logListErr != nil {
						return nil, logListErr
					}
				}
			}
		}

		if kvp.Key == "refund" {
			if isStar(kvp.Value) {
				blr.Refund = nil
			} else {
				blr.Refund, refundOk = processBigInt(kvp.Value)
				if !refundOk {
					return nil, errors.New("invalid block result refund")
				}
			}
		}
	}

	return &blr, nil
}

func processLogList(logsRaw oj.OJsonObject) ([]*LogEntry, error) {
	logList, isList := logsRaw.(*oj.OJsonList)
	if !isList {
		return nil, errors.New("unmarshalled logs list is not a list")
	}
	var logEntries []*LogEntry
	for _, logRaw := range logList.AsList() {
		logMap, isMap := logRaw.(*oj.OJsonMap)
		if !isMap {
			return nil, errors.New("unmarshalled log entry is not a map")
		}
		logEntry := LogEntry{}
		for _, kvp := range logMap.OrderedKV {
			if kvp.Key == "address" {
				accountStr, strOk := parseString(kvp.Value)
				if !strOk {
					return nil, errors.New("unmarshalled log entry address is not a json string")
				}
				var err error
				logEntry.Address, err = parseAccountAddress(accountStr)
				if err != nil {
					return nil, err
				}
			}
			if kvp.Key == "topics" {
				var topicsOk bool
				logEntry.Topics, topicsOk = parseByteArrayList(kvp.Value)
				if !topicsOk {
					return nil, errors.New("unmarshalled log entry topics is not big int list")
				}
			}
			if kvp.Key == "data" {
				var dataOk bool
				dataAsInt, dataOk := processBigInt(kvp.Value)
				if !dataOk {
					return nil, errors.New("cannot parse log entry data")
				}
				logEntry.Data = dataAsInt.Bytes()

			}
		}
		logEntries = append(logEntries, &logEntry)
	}

	return logEntries, nil
}

func processBlockTransaction(blrRaw oj.OJsonObject) (*Transaction, error) {
	bltMap, isMap := blrRaw.(*oj.OJsonMap)
	if !isMap {
		return nil, errors.New("unmarshalled block transaction is not a map")
	}

	blt := Transaction{}
	var nonceOk, functionOk, valueOk, argumentsOk, contractCodeOk, gasPriceOk, gasLimitOk bool

	for _, kvp := range bltMap.OrderedKV {

		if kvp.Key == "nonce" {
			blt.Nonce, nonceOk = parseUint64(kvp.Value)
			if !nonceOk {
				return nil, errors.New("invalid block transaction nonce")
			}
		}

		if kvp.Key == "function" {
			blt.Function, functionOk = parseString(kvp.Value)
			if !functionOk {
				return nil, errors.New("invalid block transaction function")
			}
		}

		if kvp.Key == "contractCode" {
			blt.ContractCode, contractCodeOk = parseString(kvp.Value)
			if !contractCodeOk {
				return nil, errors.New("invalid block transaction contract code")
			}
		}

		if kvp.Key == "value" {
			blt.Value, valueOk = processBigInt(kvp.Value)
			if !valueOk {
				return nil, errors.New("invalid block transaction value")
			}
		}

		if kvp.Key == "to" {
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

		}

		if kvp.Key == "arguments" {
			blt.Arguments, argumentsOk = processArgumentList(kvp.Value)
			if !argumentsOk {
				return nil, errors.New("invalid block transaction arguments")
			}
		}

		if kvp.Key == "contractCode" {
			blt.ContractCode, contractCodeOk = parseString(kvp.Value)
			if !contractCodeOk {
				return nil, errors.New("invalid block transaction contractCode")
			}
		}

		if kvp.Key == "gasPrice" {
			blt.GasPrice, gasPriceOk = parseUint64(kvp.Value)
			if !gasPriceOk {
				return nil, errors.New("invalid block transaction gasPrice")
			}
		}

		if kvp.Key == "gasLimit" {
			blt.GasLimit, gasLimitOk = parseUint64(kvp.Value)
			if !gasLimitOk {
				return nil, errors.New("invalid block transaction gasLimit")
			}
		}

		if kvp.Key == "from" {
			fromStr, fromOk := parseString(kvp.Value)
			if !fromOk {
				return nil, errors.New("invalid block transaction from")
			}
			var fromErr error
			blt.From, fromErr = parseAccountAddress(fromStr)
			if fromErr != nil {
				return nil, fromErr
			}
		}
	}

	return &blt, nil
}

func processBlockHeader(blhRaw interface{}) (*BlockHeader, error) {
	blhMap, isMap := blhRaw.(*oj.OJsonMap)
	if !isMap {
		return nil, errors.New("unmarshalled block header is not a map")
	}

	blh := BlockHeader{}
	var gasLimitOk, numberOk, difficultyOk, timestampOk, coinbaseOk bool

	for _, kvp := range blhMap.OrderedKV {

		if kvp.Key == "gasLimit" {
			blh.GasLimit, gasLimitOk = processBigInt(kvp.Value)
			if !gasLimitOk {
				return nil, errors.New("invalid block header gasLimit")
			}
		}

		if kvp.Key == "number" {
			blh.Number, numberOk = processBigInt(kvp.Value)
			if !numberOk {
				return nil, errors.New("invalid block header number")
			}
		}

		if kvp.Key == "difficulty" {
			blh.Difficulty, difficultyOk = processBigInt(kvp.Value)
			if !difficultyOk {
				return nil, errors.New("invalid block header difficulty")
			}
		}

		if kvp.Key == "timestamp" {
			blh.Timestamp, timestampOk = parseUint64(kvp.Value)
			if !timestampOk {
				return nil, errors.New("invalid block header timestamp")
			}
		}

		if kvp.Key == "coinbase" {
			blh.Beneficiary, coinbaseOk = processBigInt(kvp.Value)
			if !coinbaseOk {
				return nil, errors.New("invalid block header coinbase")
			}
		}
	}

	return &blh, nil
}

func processBigInt(obj oj.OJsonObject) (*big.Int, bool) {
	str, isStr := obj.(*oj.OJsonString)
	if !isStr {
		return nil, false
	}
	if len(str.Value) == 0 {
		return nil, true // interpret "" as nil, so we can restore to empty string instead of 0
	}

	return parseBigInt(str.Value)
}

func processStringList(obj interface{}) ([]string, bool) {
	listRaw, listOk := obj.(*oj.OJsonList)
	if !listOk {
		return nil, false
	}
	var result []string
	for _, elemRaw := range listRaw.AsList() {
		str, strOk := elemRaw.(*oj.OJsonString)
		if !strOk {
			return nil, false
		}
		result = append(result, str.Value)
	}
	return result, true
}

func processBigIntList(obj interface{}) ([]*big.Int, bool) {
	listRaw, listOk := obj.(*oj.OJsonList)
	if !listOk {
		return nil, false
	}
	var result []*big.Int
	for _, elemRaw := range listRaw.AsList() {
		i, iOk := processBigInt(elemRaw)
		if !iOk {
			return nil, false
		}
		result = append(result, i)
	}
	return result, true
}

func parseByteArrayList(obj interface{}) ([][]byte, bool) {
	listRaw, listOk := obj.(*oj.OJsonList)
	if !listOk {
		return nil, false
	}
	var result [][]byte
	for _, elemRaw := range listRaw.AsList() {
		str, strOk := parseString(elemRaw)
		if !strOk {
			return nil, false
		}
		ba, baErr := parseByteArray(str)
		if baErr != nil {
			return nil, false
		}
		result = append(result, ba)
	}
	return result, true
}

func processArgumentList(obj interface{}) ([]Argument, bool) {
	listRaw, listOk := obj.(*oj.OJsonList)
	if !listOk {
		return nil, false
	}
	var result []Argument
	for _, elemRaw := range listRaw.AsList() {
		arg, argOk := processArgument(elemRaw)
		if !argOk {
			return nil, false
		}
		result = append(result, arg)
	}
	return result, true
}

func processArgument(obj oj.OJsonObject) (Argument, bool) {
	strRaw, strOk := parseString(obj)
	if !strOk {
		return Argument{}, false
	}

	// try to parse as big int
	// TODO: figure out how to only use byte representation, there are still some issues with IELE
	// all arguments get converted to 2's complement bytes
	bi, parseOk := processBigInt(obj)
	if !parseOk {
		return Argument{}, false
	}

	forceSign := len(strRaw) > 0 && (strRaw[0] == '-' || strRaw[0] == '+')
	return Argument{
		value:     bi,
		forceSign: forceSign,
	}, true
}
