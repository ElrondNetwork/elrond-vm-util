package vmtestjson

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestParseScenario(t *testing.T) {
	// Open our jsonFile
	var jsonFile *os.File
	var err error
	jsonFile, err = os.Open("example.scen.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		t.Error(err)
		return
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		t.Error(err)
		return
	}

	p := Parser{
		FileResolver: NewDefaultFileResolver().ReplacePath(
			"smart-contract.wasm",
			"exampleFile.txt"),
	}
	_, parseErr := p.ParseScenarioFile(byteValue)
	if parseErr != nil {
		t.Error(parseErr)
		return
	}
}
