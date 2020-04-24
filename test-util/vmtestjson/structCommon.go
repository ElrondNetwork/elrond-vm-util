package vmtestjson

import (
	"bytes"
)

// Account is a json object representing an account.
type Account struct {
	Address       JSONBytes
	Nonce         JSONUint64
	Balance       JSONBigInt
	Storage       []*StorageKeyValuePair
	Code          JSONBytes
	AsyncCallData string
}

// CheckAccount is a json object representing checks for an account.
type CheckAccount struct {
	Address       JSONBytes
	Nonce         JSONCheckUint64
	Balance       JSONCheckBigInt
	Storage       []*StorageKeyValuePair
	Code          JSONBytes
	AsyncCallData string
}

// StorageKeyValuePair is a json key value pair in the storage map.
type StorageKeyValuePair struct {
	Key   JSONBytes
	Value JSONBytes
}

// Transaction is a json object representing a transaction.
type Transaction struct {
	Nonce     JSONUint64
	Value     JSONBigInt
	IsCreate  bool
	From      JSONBytes
	To        JSONBytes
	Function  string
	Code      JSONBytes
	Arguments []JSONBytes
	GasPrice  JSONUint64
	GasLimit  JSONUint64
}

// TransactionResult is a json object representing an expected transaction result.
type TransactionResult struct {
	Out        []JSONBytes
	Status     JSONBigInt
	Message    string
	Gas        JSONCheckUint64
	Refund     JSONCheckBigInt
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
func FindAccount(accounts []*Account, address []byte) *Account {
	for _, acct := range accounts {
		if bytes.Equal(acct.Address.Value, address) {
			return acct
		}
	}
	return nil
}

// FindCheckAccount searches a check account list by address.
func FindCheckAccount(accounts []*CheckAccount, address []byte) *CheckAccount {
	for _, acct := range accounts {
		if bytes.Equal(acct.Address.Value, address) {
			return acct
		}
	}
	return nil
}
