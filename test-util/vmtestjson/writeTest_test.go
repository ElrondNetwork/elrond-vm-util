package vmtestjson

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWriteTest(t *testing.T) {
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

	testTopLevel, parseErr := p.ParseTestFile(contents)
	require.Nil(t, parseErr)

	serialized := TestToJSONString(testTopLevel)

	// good for debugging:
	// _ = ioutil.WriteFile("example.re__.json", []byte(serialized), 0644)

	require.Equal(t, contents, []byte(serialized))
}
