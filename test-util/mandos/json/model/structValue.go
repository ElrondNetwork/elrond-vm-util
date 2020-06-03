package mandosjsonmodel

import (
	"bytes"
	"math/big"
)

// JSONBytes stores the parsed byte slice value but also the original parsed string
type JSONBytes struct {
	Value    []byte
	Original string
}

// JSONCheckBytes holds a byte slice condition.
// Values are checked for equality.
// "*" allows all values.
type JSONCheckBytes struct {
	Value    []byte
	IsStar   bool
	Original string
}

// Check returns true if condition expressed in object holds for another value.
// Explicit values are interpreted as equals assertion.
func (jcbytes JSONCheckBytes) Check(other []byte) bool {
	if jcbytes.IsStar {
		return true
	}
	return bytes.Equal(jcbytes.Value, other)
}

// JSONBigInt stores the parsed big int value but also the original parsed string
type JSONBigInt struct {
	Value    *big.Int
	Original string
}

// JSONCheckBigInt holds a big int condition.
// Values are checked for equality.
// "*" allows all values.
type JSONCheckBigInt struct {
	Value    *big.Int
	IsStar   bool
	Original string
}

// Check returns true if condition expressed in object holds for another value.
// Explicit values are interpreted as equals assertion.
func (jcbi JSONCheckBigInt) Check(other *big.Int) bool {
	if jcbi.IsStar {
		return true
	}
	return jcbi.Value.Cmp(other) == 0
}

// JSONUint64 stores the parsed uint64 value but also the original parsed string
type JSONUint64 struct {
	Value    uint64
	Original string
}

// JSONCheckUint64 holds a uint64 condition.
// Values are checked for equality.
// "*" allows all values.
type JSONCheckUint64 struct {
	Value    uint64
	IsStar   bool
	Original string
}

// Check returns true if condition expressed in object holds for another value.
// Explicit values are interpreted as equals assertion.
func (jcu JSONCheckUint64) Check(other uint64) bool {
	if jcu.IsStar {
		return true
	}
	return jcu.Value == other
}
