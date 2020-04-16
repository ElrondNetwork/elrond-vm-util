package vmtestjson

import "math/big"

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
