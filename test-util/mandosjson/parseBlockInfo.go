package mandosjson

import (
	"errors"
	"fmt"

	oj "github.com/ElrondNetwork/elrond-vm-util/test-util/orderedjson"
)

func (p *Parser) processBlockInfo(blockInfoRaw oj.OJsonObject) (*BlockInfo, error) {
	blockMap, isMap := blockInfoRaw.(*oj.OJsonMap)
	if !isMap {
		return nil, errors.New("unmarshalled block info object is not a map")
	}
	blockInfo := &BlockInfo{}
	var err error

	for _, kvp := range blockMap.OrderedKV {
		switch kvp.Key {
		case "blockTimestamp":
			blockInfo.BlockTimestamp, err = p.processUint64(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("error parsing blockTimestamp: %w", err)
			}
		case "blockNonce":
			blockInfo.BlockNonce, err = p.processUint64(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("error parsing blockNonce: %w", err)
			}
		case "blockRound":
			blockInfo.BlockRound, err = p.processUint64(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("error parsing blockRound: %w", err)
			}
		case "blockEpoch":
			blockInfo.BlockEpoch, err = p.processUint64(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("error parsing blockEpoch: %w", err)
			}
		default:
			return nil, fmt.Errorf("unknown block info field: %s", kvp.Key)
		}
	}

	return blockInfo, nil
}
