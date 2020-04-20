package vmtestjson

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseScenario(t *testing.T) {
	contents, err := loadExampleFile("example.scen.json")
	require.Nil(t, err)

	p := Parser{
		FileResolver: NewDefaultFileResolver().ReplacePath(
			"smart-contract.wasm",
			"exampleFile.txt"),
	}
	_, parseErr := p.ParseScenarioFile(contents)
	require.Nil(t, parseErr)
}
