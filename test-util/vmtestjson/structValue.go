package vmtestjson

import "math/big"

// JSONBytes stores the parsed byte slice value but also the original parsed string
type JSONBytes struct {
	Value    []byte
	Original string
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
	value    *big.Int
	isStar   bool
	Original string
}

// Check returns true if condition expressed in object holds for another value.
// Explicit values are interpreted as equals assertion.
func (jcbi JSONCheckBigInt) Check(other *big.Int) bool {
	if jcbi.isStar {
		return true
	}
	return jcbi.value.Cmp(other) == 0
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
	value    uint64
	isStar   bool
	Original string
}

// Check returns true if condition expressed in object holds for another value.
// Explicit values are interpreted as equals assertion.
func (jcu JSONCheckUint64) Check(other uint64) bool {
	if jcu.isStar {
		return true
	}
	return jcu.value == other
}
