package vmtestjson

import (
	"errors"
	"fmt"

	oj "github.com/ElrondNetwork/elrond-vm-util/test-util/orderedjson"
)

func (p *Parser) processBlock(blockRaw oj.OJsonObject) (*Block, error) {
	blockMap, isMap := blockRaw.(*oj.OJsonMap)
	if !isMap {
		return nil, errors.New("unmarshalled block object is not a map")
	}
	bl := Block{}

	for _, kvp := range blockMap.OrderedKV {
		switch kvp.Key {
		case "results":
			resultsRaw, resultsOk := kvp.Value.(*oj.OJsonList)
			if !resultsOk {
				return nil, errors.New("unmarshalled block results object is not a list")
			}
			for _, resRaw := range resultsRaw.AsList() {
				blr, blrErr := p.processTxExpectedResult(resRaw)
				if blrErr != nil {
					return nil, blrErr
				}
				bl.Results = append(bl.Results, blr)
			}
		case "transactions":
			transactionsRaw, transactionsOk := kvp.Value.(*oj.OJsonList)
			if !transactionsOk {
				return nil, errors.New("unmarshalled block transactions object is not a list")
			}
			for _, trRaw := range transactionsRaw.AsList() {
				tr, trErr := p.processTx(trRaw)
				if trErr != nil {
					return nil, trErr
				}
				bl.Transactions = append(bl.Transactions, tr)
			}
		case "blockHeader":
			blh, blhErr := p.processBlockHeader(kvp.Value)
			if blhErr != nil {
				return nil, blhErr
			}
			bl.BlockHeader = blh
		default:
			return nil, fmt.Errorf("unknown block field: %s", kvp.Key)
		}
	}

	if len(bl.Results) != len(bl.Transactions) {
		return nil, errors.New("mismatched number of blocks and transactions")
	}

	return &bl, nil
}

func (p *Parser) processBlockHeader(blhRaw interface{}) (*BlockHeader, error) {
	blhMap, isMap := blhRaw.(*oj.OJsonMap)
	if !isMap {
		return nil, errors.New("unmarshalled block header is not a map")
	}

	blh := BlockHeader{}
	var err error

	for _, kvp := range blhMap.OrderedKV {
		switch kvp.Key {
		case "gasLimit":
			blh.GasLimit, err = p.processBigInt(kvp.Value, bigIntUnsignedBytes)
			if err != nil {
				return nil, fmt.Errorf("invalid block header gasLimit: %w", err)
			}
		case "number":
			blh.Number, err = p.processBigInt(kvp.Value, bigIntUnsignedBytes)
			if err != nil {
				return nil, fmt.Errorf("invalid block header number: %w", err)
			}
		case "difficulty":
			blh.Difficulty, err = p.processBigInt(kvp.Value, bigIntUnsignedBytes)
			if err != nil {
				return nil, fmt.Errorf("invalid block header difficulty: %w", err)
			}
		case "timestamp":
			blh.Timestamp, err = p.processUint64(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid block header timestamp: %w", err)
			}
		case "coinbase":
			blh.Beneficiary, err = p.processBigInt(kvp.Value, bigIntUnsignedBytes)
			if err != nil {
				return nil, fmt.Errorf("invalid block header coinbase: %w", err)
			}
		default:
			return nil, fmt.Errorf("unknown block header field: %s", kvp.Key)
		}
	}

	return &blh, nil
}
