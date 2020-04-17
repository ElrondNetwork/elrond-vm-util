package vmtestjson

import (
	"errors"
	"fmt"

	oj "github.com/ElrondNetwork/elrond-vm-util/test-util/orderedjson"
)

func (p *Parser) parseAccountAddress(addrRaw string) ([]byte, error) {
	if len(addrRaw) == 0 {
		return []byte{}, errors.New("missing account address")
	}
	addrBytes, err := p.parseAnyValueAsByteArray(addrRaw)
	if err == nil && len(addrBytes) != 32 {
		return []byte{}, errors.New("account addressis not 32 bytes in length")
	}
	return addrBytes, err
}

func (p *Parser) processAccount(acctRaw oj.OJsonObject) (*Account, error) {
	acctMap, isMap := acctRaw.(*oj.OJsonMap)
	if !isMap {
		return nil, errors.New("unmarshalled account object is not a map")
	}

	acct := Account{}
	var err error

	for _, kvp := range acctMap.OrderedKV {

		if kvp.Key == "nonce" {
			acct.Nonce, err = p.processBigInt(kvp.Value, bigIntUnsignedBytes)
			if err != nil {
				return nil, errors.New("invalid account nonce")
			}
		}

		if kvp.Key == "balance" {
			acct.Balance, err = p.processBigInt(kvp.Value, bigIntUnsignedBytes)
			if err != nil {
				return nil, errors.New("invalid account balance")
			}
		}

		if kvp.Key == "storage" {
			storageMap, storageOk := kvp.Value.(*oj.OJsonMap)
			if !storageOk {
				return nil, errors.New("invalid account storage")
			}
			for _, storageKvp := range storageMap.OrderedKV {
				byteKey, err := p.parseAnyValueAsByteArray(storageKvp.Key)
				if err != nil {
					return nil, fmt.Errorf("invalid account storage key: %w", err)
				}
				if len(byteKey) != 32 {
					return nil, errors.New("storage keys must be 32 bytes in length")
				}
				byteVal, _, err := p.processAnyValueAsByteArray(storageKvp.Value)
				if err != nil {
					return nil, fmt.Errorf("invalid account storage value: %w", err)
				}
				stElem := StorageKeyValuePair{
					Key:   byteKey,
					Value: byteVal,
				}
				acct.Storage = append(acct.Storage, &stElem)
			}
		}

		if kvp.Key == "code" {
			acct.Code, acct.OriginalCode, err = p.processAnyValueAsByteArray(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid account code: %w", err)
			}
		}

		if kvp.Key == "asyncCallData" {
			acct.AsyncCallData, err = p.parseString(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid asyncCallData string: %w", err)
			}
		}
	}

	return &acct, nil
}

func (p *Parser) processAccountMap(acctMapRaw oj.OJsonObject) ([]*Account, error) {
	var accounts []*Account
	preMap, isPreMap := acctMapRaw.(*oj.OJsonMap)
	if !isPreMap {
		return nil, errors.New("unmarshalled account map object is not a map")
	}
	for _, acctKVP := range preMap.OrderedKV {
		if acctKVP.Key != "step" {
			acct, acctErr := p.processAccount(acctKVP.Value)
			if acctErr != nil {
				return nil, acctErr
			}
			acctAddr, hexErr := p.parseAccountAddress(acctKVP.Key)
			if hexErr != nil {
				return nil, hexErr
			}
			acct.Address = acctAddr
			accounts = append(accounts, acct)
		}
	}
	return accounts, nil
}
