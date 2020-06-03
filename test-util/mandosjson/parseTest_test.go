package mandosjson

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func loadExampleFile(path string) ([]byte, error) {
	// Open our jsonFile
	var jsonFile *os.File
	var err error
	jsonFile, err = os.Open(path)
	// if we os.Open returns an error then handle it
	if err != nil {
		return nil, err
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	return ioutil.ReadAll(jsonFile)
}

func TestParseTest(t *testing.T) {
	contents, err := loadExampleFile("example.test.json")
	require.Nil(t, err)

	p := Parser{
		FileResolver: NewDefaultFileResolver().ReplacePath(
			"smart-contract.wasm",
			"exampleFile.txt"),
	}

	_, parseErr := p.ParseTestFile(contents)
	require.Nil(t, parseErr)
}
