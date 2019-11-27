package vmtestjson

import (
	"bytes"
	"math/big"
)

// Test is a json object representing a test.
type Test struct {
	TestName    string
	CheckGas    bool
	Pre         []*Account
	Blocks      []*Block
	Network     string
	BlockHashes [][]byte
	PostState   []*Account
}

// Account is a json object representing an account.
type Account struct {
	Address      []byte
	Nonce        *big.Int
	Balance      *big.Int
	Storage      []*StorageKeyValuePair
	Code         string
	OriginalCode string
}

// StorageKeyValuePair is a json key value pair in the storage map.
type StorageKeyValuePair struct {
	Key   []byte
	Value []byte
}

// Block is a json object representing a block.
type Block struct {
	Results      []*TransactionResult
	Transactions []*Transaction
	BlockHeader  *BlockHeader
}

// BlockHeader is a json object representing the block header.
type BlockHeader struct {
	Beneficiary *big.Int // "coinbase"
	Difficulty  *big.Int
	Number      *big.Int
	GasLimit    *big.Int
	Timestamp   uint64
}

// TransactionResult is a json object representing an expected transaction result.
type TransactionResult struct {
	Out        [][]byte
	Status     *big.Int
	CheckGas   bool
	Gas        uint64
	Refund     *big.Int
	IgnoreLogs bool
	LogHash    string
	Logs       []*LogEntry
}

// LogEntry is a json object representing an expected transaction result log entry.
type LogEntry struct {
	Address []byte
	Topics  [][]byte
	Data    []byte
}

// Argument encodes an argument in a transaction.
// Can distinguish values written explicitly as poitive or negative (e.g. -0x01, +0xFF),
// in order to provide some additional context on how to interpret them in an actual test.
type Argument struct {
	value     *big.Int
	forceSign bool
}

// Transaction is a json object representing a transaction.
type Transaction struct {
	Nonce         uint64
	Value         *big.Int
	IsCreate      bool
	From          []byte
	To            []byte
	Function      string
	ContractCode  string
	AssembledCode string
	Arguments     []Argument
	GasPrice      uint64
	GasLimit      uint64
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
