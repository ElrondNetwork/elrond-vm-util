package vmtestjson

import (
	"errors"
	"math/big"

	oj "github.com/ElrondNetwork/elrond-vm-util/test-util/orderedjson"
)

func (p *Parser) processBlock(blockRaw oj.OJsonObject) (*Block, error) {
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
				blr, blrErr := p.processBlockResult(resRaw)
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
				tr, trErr := p.processBlockTransaction(trRaw)
				if trErr != nil {
					return nil, trErr
				}
				bl.Transactions = append(bl.Transactions, tr)
			}
		}

		if kvp.Key == "blockHeader" {
			blh, blhErr := p.processBlockHeader(kvp.Value)
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

func (p *Parser) processBlockResult(blrRaw oj.OJsonObject) (*TransactionResult, error) {
	blrMap, isMap := blrRaw.(*oj.OJsonMap)
	if !isMap {
		return nil, errors.New("unmarshalled block result is not a map")
	}

	blr := TransactionResult{}
	var outOk, statusOk, messageOk, gasOk, refundOk bool

	for _, kvp := range blrMap.OrderedKV {

		if kvp.Key == "out" {
			blr.Out, outOk = p.parseByteArrayList(kvp.Value)
			if !outOk {
				return nil, errors.New("invalid block result out")
			}
		}

		if kvp.Key == "status" {
			blr.Status, statusOk = p.processBigInt(kvp.Value)
			if !statusOk {
				return nil, errors.New("invalid block result status")
			}
		}

		if kvp.Key == "message" {
			blr.Message, messageOk = p.parseString(kvp.Value)
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
				blr.Gas, gasOk = p.parseUint64(kvp.Value)
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
				blr.LogHash, logHashOk = p.parseString(kvp.Value)
				if !logHashOk {
					var logListErr error
					blr.Logs, logListErr = p.processLogList(kvp.Value)
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
				blr.Refund, refundOk = p.processBigInt(kvp.Value)
				if !refundOk {
					return nil, errors.New("invalid block result refund")
				}
			}
		}
	}

	return &blr, nil
}

func (p *Parser) processBlockHeader(blhRaw interface{}) (*BlockHeader, error) {
	blhMap, isMap := blhRaw.(*oj.OJsonMap)
	if !isMap {
		return nil, errors.New("unmarshalled block header is not a map")
	}

	blh := BlockHeader{}
	var gasLimitOk, numberOk, difficultyOk, timestampOk, coinbaseOk bool

	for _, kvp := range blhMap.OrderedKV {

		if kvp.Key == "gasLimit" {
			blh.GasLimit, gasLimitOk = p.processBigInt(kvp.Value)
			if !gasLimitOk {
				return nil, errors.New("invalid block header gasLimit")
			}
		}

		if kvp.Key == "number" {
			blh.Number, numberOk = p.processBigInt(kvp.Value)
			if !numberOk {
				return nil, errors.New("invalid block header number")
			}
		}

		if kvp.Key == "difficulty" {
			blh.Difficulty, difficultyOk = p.processBigInt(kvp.Value)
			if !difficultyOk {
				return nil, errors.New("invalid block header difficulty")
			}
		}

		if kvp.Key == "timestamp" {
			blh.Timestamp, timestampOk = p.parseUint64(kvp.Value)
			if !timestampOk {
				return nil, errors.New("invalid block header timestamp")
			}
		}

		if kvp.Key == "coinbase" {
			blh.Beneficiary, coinbaseOk = p.processBigInt(kvp.Value)
			if !coinbaseOk {
				return nil, errors.New("invalid block header coinbase")
			}
		}
	}

	return &blh, nil
}

func (p *Parser) processBigInt(obj oj.OJsonObject) (*big.Int, bool) {
	str, isStr := obj.(*oj.OJsonString)
	if !isStr {
		return nil, false
	}
	if len(str.Value) == 0 {
		return nil, true // interpret "" as nil, so we can restore to empty string instead of 0
	}

	return p.parseBigInt(str.Value)
}

func (p *Parser) processStringList(obj interface{}) ([]string, bool) {
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

func (p *Parser) processBigIntList(obj interface{}) ([]*big.Int, bool) {
	listRaw, listOk := obj.(*oj.OJsonList)
	if !listOk {
		return nil, false
	}
	var result []*big.Int
	for _, elemRaw := range listRaw.AsList() {
		i, iOk := p.processBigInt(elemRaw)
		if !iOk {
			return nil, false
		}
		result = append(result, i)
	}
	return result, true
}

func (p *Parser) parseByteArrayList(obj interface{}) ([][]byte, bool) {
	listRaw, listOk := obj.(*oj.OJsonList)
	if !listOk {
		return nil, false
	}
	var result [][]byte
	for _, elemRaw := range listRaw.AsList() {
		str, strOk := p.parseString(elemRaw)
		if !strOk {
			return nil, false
		}
		ba, baErr := p.parseAnyValueAsByteArray(str)
		if baErr != nil {
			return nil, false
		}
		result = append(result, ba)
	}
	return result, true
}
