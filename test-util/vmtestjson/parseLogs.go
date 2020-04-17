package vmtestjson

import (
	"errors"
	"fmt"

	oj "github.com/ElrondNetwork/elrond-vm-util/test-util/orderedjson"
)

func (p *Parser) processLogList(logsRaw oj.OJsonObject) ([]*LogEntry, error) {
	logList, isList := logsRaw.(*oj.OJsonList)
	if !isList {
		return nil, errors.New("unmarshalled logs list is not a list")
	}
	var logEntries []*LogEntry
	var err error
	for _, logRaw := range logList.AsList() {
		logMap, isMap := logRaw.(*oj.OJsonMap)
		if !isMap {
			return nil, errors.New("unmarshalled log entry is not a map")
		}
		logEntry := LogEntry{}
		for _, kvp := range logMap.OrderedKV {
			if kvp.Key == "address" {
				accountStr, err := p.parseString(kvp.Value)
				if err != nil {
					return nil, fmt.Errorf("unmarshalled log entry address is not a json string: %w", err)
				}
				logEntry.Address, err = p.parseAccountAddress(accountStr)
				if err != nil {
					return nil, err
				}
			}
			if kvp.Key == "identifier" {
				strVal, err := p.parseString(kvp.Value)
				if err != nil {
					return nil, fmt.Errorf("invalid log identifier: %w", err)
				}
				logEntry.Identifier, err = p.parseAnyValueAsByteArray(strVal)
				if err != nil {
					return nil, fmt.Errorf("invalid log identifier: %w", err)
				}
				if len(logEntry.Identifier) != 32 {
					return nil, fmt.Errorf("invalid log identifier - should be 32 bytes in length")
				}
			}
			if kvp.Key == "topics" {
				logEntry.Topics, err = p.parseByteArrayList(kvp.Value)
				if err != nil {
					return nil, fmt.Errorf("unmarshalled log entry topics is not big int list: %w", err)
				}
			}
			if kvp.Key == "data" {
				logEntry.Data, _, err = p.processAnyValueAsByteArray(kvp.Value)
				if err != nil {
					return nil, fmt.Errorf("cannot parse log entry data: %w", err)
				}

			}
		}
		logEntries = append(logEntries, &logEntry)
	}

	return logEntries, nil
}
