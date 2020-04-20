package vmtestjson

import (
	"bytes"
	"encoding/hex"
	"math/big"
)

// ResultEqual returns true if result bytes encode the same number.
func ResultEqual(expected JSONBytes, actual []byte) bool {
	if bytes.Equal(expected.Value, actual) {
		return true
	}

	return big.NewInt(0).SetBytes(expected.Value).Cmp(big.NewInt(0).SetBytes(actual)) == 0
}

// ResultAsString helps create nicer error messages.
func ResultAsString(result [][]byte) string {
	str := "["
	for i, res := range result {
		str += "0x" + hex.EncodeToString(res)
		if i < len(result)-1 {
			str += ", "
		}
	}
	return str + "]"
}

// JSONBytesValues extracts values from a slice of JSONBytes into a list
func JSONBytesValues(jbs []JSONBytes) [][]byte {
	result := make([][]byte, len(jbs))
	for i, jb := range jbs {
		result[i] = jb.Value
	}
	return result
}
