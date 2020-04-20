package vmtestjson

import (
	"errors"
	"fmt"

	oj "github.com/ElrondNetwork/elrond-vm-util/test-util/orderedjson"
)

func (p *Parser) processTxExpectedResult(blrRaw oj.OJsonObject) (*TransactionResult, error) {
	blrMap, isMap := blrRaw.(*oj.OJsonMap)
	if !isMap {
		return nil, errors.New("unmarshalled block result is not a map")
	}

	blr := TransactionResult{}
	var err error
	for _, kvp := range blrMap.OrderedKV {
		switch kvp.Key {
		case "out":
			blr.Out, err = p.parseByteArrayList(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid block result out: %w", err)
			}
		case "status":
			blr.Status, err = p.processBigInt(kvp.Value, bigIntSignedBytes)
			if err != nil {
				return nil, fmt.Errorf("invalid block result status: %w", err)
			}
		case "message":
			blr.Message, err = p.parseString(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid block result message: %w", err)
			}
		case "gas":
			if isStar(kvp.Value) {
				blr.CheckGas = false
				blr.Gas = 0
			} else {
				blr.CheckGas = true
				blr.Gas, err = p.parseUint64(kvp.Value)
				if err != nil {
					return nil, fmt.Errorf("invalid block result gas: %w", err)
				}
			}
		case "logs":
			if isStar(kvp.Value) {
				blr.IgnoreLogs = true
			} else {
				blr.IgnoreLogs = false
				blr.LogHash, err = p.parseString(kvp.Value)
				if err != nil {
					var logListErr error
					blr.Logs, logListErr = p.processLogList(kvp.Value)
					if logListErr != nil {
						return nil, logListErr
					}
				}
			}
		case "refund":
			if isStar(kvp.Value) {
				blr.Refund = JSONBigInt{Original: "*"}
			} else {
				blr.Refund, err = p.processBigInt(kvp.Value, bigIntUnsignedBytes)
				if err != nil {
					return nil, fmt.Errorf("invalid block result refund: %w", err)
				}
			}
		default:
			return nil, fmt.Errorf("unknown tx result field: %s", kvp.Key)
		}
	}

	return &blr, nil
}
