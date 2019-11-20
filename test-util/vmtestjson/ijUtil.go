package vmtestjson

import (
	"bytes"
	"encoding/hex"
	"math/big"

	twos "github.com/ElrondNetwork/big-int-util/twos-complement"
)

// ToBytes yields value as byte array. Can force sign if so specified in the test json.
func (arg Argument) ToBytes() []byte {
	if arg.forceSign {
		return twos.ToBytes(arg.value)
	}
	return arg.value.Bytes()
}

// ToBytesAlwaysForceSign yields value as byte array,
// always forcing the correct test bit, even if it means adding an extra byte.
func (arg Argument) ToBytesAlwaysForceSign() []byte {
	return twos.ToBytes(arg.value)
}

// ResultEqual returns true if result bytes encode the same number.
func ResultEqual(expected, actual []byte) bool {
	if bytes.Equal(expected, actual) {
		return true
	}

	return big.NewInt(0).SetBytes(expected).Cmp(big.NewInt(0).SetBytes(actual)) == 0
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
