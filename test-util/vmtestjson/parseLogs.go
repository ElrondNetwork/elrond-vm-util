package vmtestjson

import (
	"errors"

	oj "github.com/ElrondNetwork/elrond-vm-util/test-util/orderedjson"
)

func (p *Parser) processLogList(logsRaw oj.OJsonObject) ([]*LogEntry, error) {
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
				accountStr, strOk := p.parseString(kvp.Value)
				if !strOk {
					return nil, errors.New("unmarshalled log entry address is not a json string")
				}
				var err error
				logEntry.Address, err = p.parseAccountAddress(accountStr)
				if err != nil {
					return nil, err
				}
			}
			if kvp.Key == "identifier" {
				strVal, valStrOk := p.parseString(kvp.Value)
				if !valStrOk {
					return nil, errors.New("invalid log identifier")
				}
				var identifierErr error
				logEntry.Identifier, identifierErr = p.parseAnyValueAsByteArray(strVal)
				if identifierErr != nil {
					return nil, errors.New("invalid log identifier")
				}
				if len(logEntry.Identifier) != 32 {
					return nil, errors.New("invalid log identifier - should be 32 bytes in length")
				}
			}
			if kvp.Key == "topics" {
				var topicsOk bool
				logEntry.Topics, topicsOk = p.parseByteArrayList(kvp.Value)
				if !topicsOk {
					return nil, errors.New("unmarshalled log entry topics is not big int list")
				}
			}
			if kvp.Key == "data" {
				var dataOk bool
				dataAsInt, dataOk := p.processBigInt(kvp.Value)
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
