package vmtestjson

import (
	"bytes"
	"math/big"
)

// Account is a json object representing an account.
type Account struct {
	Address       []byte
	Nonce         *big.Int
	Balance       *big.Int
	Storage       []*StorageKeyValuePair
	Code          string
	OriginalCode  string
	AsyncCallData string
}

// StorageKeyValuePair is a json key value pair in the storage map.
type StorageKeyValuePair struct {
	Key   []byte
	Value []byte
}

// TransactionResult is a json object representing an expected transaction result.
type TransactionResult struct {
	Out        [][]byte
	Status     *big.Int
	Message    string
	CheckGas   bool
	Gas        uint64
	Refund     *big.Int
	IgnoreLogs bool
	LogHash    string
	Logs       []*LogEntry
}

// LogEntry is a json object representing an expected transaction result log entry.
type LogEntry struct {
	Address    []byte
	Identifier []byte
	Topics     [][]byte
	Data       []byte
}

// FindAccount searches an account list by address.
func FindAccount(accounts []*Account, address []byte) *Account {
	for _, acct := range accounts {
		if bytes.Equal(acct.Address, address) {
			return acct
		}
	}
	return nil
}
