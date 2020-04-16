package vmtestjson

import (
	"errors"

	oj "github.com/ElrondNetwork/elrond-vm-util/test-util/orderedjson"
)

func parseAccountAddress(addrRaw string) ([]byte, error) {
	if len(addrRaw) == 0 {
		return []byte{}, errors.New("missing account address")
	}
	addrBytes, err := parseAnyValueAsByteArray(addrRaw)
	if err == nil && len(addrBytes) != 32 {
		return []byte{}, errors.New("account addressis not 32 bytes in length")
	}
	return addrBytes, err
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
				byteKey, keyOk := parseAnyValueAsByteArray(storageKvp.Key)
				if keyOk != nil {
					return nil, errors.New("invalid account storage key")
				}
				strVal, valStrOk := parseString(storageKvp.Value)
				if !valStrOk {
					return nil, errors.New("invalid account storage value")
				}
				byteVal, valOk := parseAnyValueAsByteArray(strVal)
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

func processAccountMap(acctMapRaw oj.OJsonObject) ([]*Account, error) {
	var accounts []*Account
	preMap, isPreMap := acctMapRaw.(*oj.OJsonMap)
	if !isPreMap {
		return nil, errors.New("unmarshalled account map object is not a map")
	}
	for _, acctKVP := range preMap.OrderedKV {
		if acctKVP.Key != "step" {
			acct, acctErr := processAccount(acctKVP.Value)
			if acctErr != nil {
				return nil, acctErr
			}
			acctAddr, hexErr := parseAccountAddress(acctKVP.Key)
			if hexErr != nil {
				return nil, hexErr
			}
			acct.Address = acctAddr
			accounts = append(accounts, acct)
		}
	}
	return accounts, nil
}
