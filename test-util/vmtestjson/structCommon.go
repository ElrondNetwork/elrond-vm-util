package vmtestjson

import (
	"bytes"
	"math/big"
)

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

// Account is a json object representing an account.
type Account struct {
	Address       JSONBytes
	Nonce         JSONBigInt
	Balance       JSONBigInt
	Storage       []*StorageKeyValuePair
	Code          JSONBytes
	OriginalCode  string
	AsyncCallData string
}

// StorageKeyValuePair is a json key value pair in the storage map.
type StorageKeyValuePair struct {
	Key   JSONBytes
	Value JSONBytes
}

// Transaction is a json object representing a transaction.
type Transaction struct {
	Nonce     uint64
	Value     JSONBigInt
	IsCreate  bool
	From      JSONBytes
	To        JSONBytes
	Function  string
	Code      JSONBytes
	Arguments []JSONBytes
	GasPrice  uint64
	GasLimit  uint64
}

// TransactionResult is a json object representing an expected transaction result.
type TransactionResult struct {
	Out        []JSONBytes
	Status     JSONBigInt
	Message    string
	CheckGas   bool
	Gas        uint64
	Refund     JSONBigInt
	IgnoreLogs bool
	LogHash    string
	Logs       []*LogEntry
}

// LogEntry is a json object representing an expected transaction result log entry.
type LogEntry struct {
	Address    JSONBytes
	Identifier JSONBytes
	Topics     []JSONBytes
	Data       JSONBytes
}

// FindAccount searches an account list by address.
func FindAccount(accounts []*Account, address JSONBytes) *Account {
	for _, acct := range accounts {
		if bytes.Equal(acct.Address.Value, address.Value) {
			return acct
		}
	}
	return nil
}
