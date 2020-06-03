package mandosjson

import (
	"errors"

	oj "github.com/ElrondNetwork/elrond-vm-util/test-util/orderedjson"
)

func (p *Parser) processStringList(obj interface{}) ([]string, error) {
	listRaw, listOk := obj.(*oj.OJsonList)
	if !listOk {
		return nil, errors.New("not a JSON list")
	}
	var result []string
	for _, elemRaw := range listRaw.AsList() {
		strVal, err := p.parseString(elemRaw)
		if err != nil {
			return nil, err
		}
		result = append(result, strVal)
	}
	return result, nil
}

func (p *Parser) parseByteArrayList(obj interface{}) ([]JSONBytes, error) {
	listRaw, listOk := obj.(*oj.OJsonList)
	if !listOk {
		return nil, errors.New("not a JSON list")
	}
	var result []JSONBytes
	for _, elemRaw := range listRaw.AsList() {
		ba, err := p.processAnyValueAsByteArray(elemRaw)
		if err != nil {
			return nil, err
		}
		result = append(result, ba)
	}
	return result, nil
}

func (p *Parser) parseCheckBytesList(obj interface{}) ([]JSONCheckBytes, error) {
	listRaw, listOk := obj.(*oj.OJsonList)
	if !listOk {
		return nil, errors.New("not a JSON list")
	}
	var result []JSONCheckBytes
	for _, elemRaw := range listRaw.AsList() {
		checkBytes, err := p.parseCheckBytes(elemRaw)
		if err != nil {
			return nil, err
		}
		result = append(result, checkBytes)
	}
	return result, nil
}
