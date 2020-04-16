package vmtestjson

import (
	"errors"

	oj "github.com/ElrondNetwork/elrond-vm-util/test-util/orderedjson"
)

// ParseTestFile converts json string to object representation
func (p *Parser) ParseTestFile(jsonString []byte) ([]*Test, error) {

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
		t, tErr := p.processTest(kvp.Value)
		if tErr != nil {
			return nil, tErr
		}
		t.TestName = kvp.Key
		top = append(top, t)
	}
	return top, nil
}

func (p *Parser) processTest(testObj oj.OJsonObject) (*Test, error) {
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
			var err error
			test.Pre, err = p.processAccountMap(kvp.Value)
			if err != nil {
				return nil, err
			}
		}

		if kvp.Key == "blocks" {
			blocksRaw, blocksOk := kvp.Value.(*oj.OJsonList)
			if !blocksOk {
				return nil, errors.New("unmarshalled blocks object is not a list")
			}
			for _, blRaw := range blocksRaw.AsList() {
				bl, blErr := p.processBlock(blRaw)
				if blErr != nil {
					return nil, blErr
				}
				test.Blocks = append(test.Blocks, bl)
			}
		}

		if kvp.Key == "network" {
			var networkOk bool
			test.Network, networkOk = p.parseString(kvp.Value)
			if !networkOk {
				return nil, errors.New("test network value not a string")
			}
		}

		if kvp.Key == "blockhashes" {
			var bhsOk bool
			test.BlockHashes, bhsOk = p.parseByteArrayList(kvp.Value)
			if !bhsOk {
				return nil, errors.New("unmarshalled blockHashes object is not a list")
			}
		}

		if kvp.Key == "postState" {
			var err error
			test.PostState, err = p.processAccountMap(kvp.Value)
			if err != nil {
				return nil, err
			}
		}
	}

	return &test, nil
}
