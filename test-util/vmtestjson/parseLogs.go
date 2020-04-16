package vmtestjson

import (
	"errors"

	oj "github.com/ElrondNetwork/elrond-vm-util/test-util/orderedjson"
)

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
			if kvp.Key == "identifier" {
				strVal, valStrOk := parseString(kvp.Value)
				if !valStrOk {
					return nil, errors.New("invalid log identifier")
				}
				var identifierErr error
				logEntry.Identifier, identifierErr = parseAnyValueAsByteArray(strVal)
				if identifierErr != nil {
					return nil, errors.New("invalid log identifier")
				}
				if len(logEntry.Identifier) != 32 {
					return nil, errors.New("invalid log identifier - should be 32 bytes in length")
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
