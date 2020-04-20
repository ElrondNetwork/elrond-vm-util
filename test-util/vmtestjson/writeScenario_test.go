package vmtestjson

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWriteScenario(t *testing.T) {
	contents, err := loadExampleFile("example.scen.json")
	require.Nil(t, err)

	p := Parser{
		FileResolver: NewDefaultFileResolver().ReplacePath(
			"smart-contract.wasm",
			"exampleFile.txt"),
	}

	scenario, parseErr := p.ParseScenarioFile(contents)
	require.Nil(t, parseErr)

	serialized := ScenarioToJSONString(scenario)

	// good for debugging:
	_ = ioutil.WriteFile("example_re.scen.json", []byte(serialized), 0644)

	require.Equal(t, contents, []byte(serialized))
}
