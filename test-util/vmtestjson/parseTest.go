package vmtestjson

import (
	"errors"

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
