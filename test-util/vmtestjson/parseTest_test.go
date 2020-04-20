package vmtestjson

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func loadTestExample() ([]byte, error) {
	// Open our jsonFile
	var jsonFile *os.File
	var err error
	jsonFile, err = os.Open("example.test.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		return nil, err
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	return ioutil.ReadAll(jsonFile)
}

func TestParseTest(t *testing.T) {
	contents, err := loadTestExample()
	if err != nil {
		t.Error(err)
		return
	}

	p := Parser{
		FileResolver: NewDefaultFileResolver().ReplacePath(
			"smart-contract.wasm",
			"exampleFile.txt"),
	}

	_, parseErr := p.ParseTestFile(contents)
	require.Nil(t, parseErr)
}
