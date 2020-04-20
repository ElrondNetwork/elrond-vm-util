package vmtestjson

import (
	oj "github.com/ElrondNetwork/elrond-vm-util/test-util/orderedjson"
)

// TestToJSONString converts a test object to its JSON representation.
func TestToJSONString(testTopLevel []*Test) string {
	jobj := TestToOrderedJSON(testTopLevel)
	return oj.JSONString(jobj)
}

// TestToOrderedJSON converts a test object to an ordered JSON object.
func TestToOrderedJSON(testTopLevel []*Test) oj.OJsonObject {
	result := oj.NewMap()
	for _, test := range testTopLevel {
		result.Put(test.TestName, testToOJ(test))
	}

	return result
}

func testToOJ(test *Test) oj.OJsonObject {
	testOJ := oj.NewMap()

	if !test.CheckGas {
		ojFalse := oj.OJsonBool(false)
		testOJ.Put("checkGas", &ojFalse)
	}

	testOJ.Put("pre", accountsToOJ(test.Pre))

	var blockList []oj.OJsonObject
	for _, block := range test.Blocks {
		blockList = append(blockList, blockToOJ(block))
	}
	blocksOJ := oj.OJsonList(blockList)
	testOJ.Put("blocks", &blocksOJ)
	testOJ.Put("network", stringToOJ(test.Network))
	testOJ.Put("blockhashes", blockHashesToOJ(test.BlockHashes))
	testOJ.Put("postState", accountsToOJ(test.PostState))
	return testOJ
}

func transactionToTestOJ(tx *Transaction) oj.OJsonObject {
	transactionOJ := oj.NewMap()
	transactionOJ.Put("nonce", uint64ToOJ(tx.Nonce))
	transactionOJ.Put("function", stringToOJ(tx.Function))
	transactionOJ.Put("gasLimit", uint64ToOJ(tx.GasLimit))
	transactionOJ.Put("value", intToOJ(tx.Value))
	transactionOJ.Put("to", byteArrayToOJ(tx.To))

	var argList []oj.OJsonObject
	for _, arg := range tx.Arguments {
		argList = append(argList, byteArrayToOJ(arg))
	}
	argOJ := oj.OJsonList(argList)
	transactionOJ.Put("arguments", &argOJ)

	transactionOJ.Put("contractCode", byteArrayToOJ(tx.Code))
	transactionOJ.Put("gasPrice", uint64ToOJ(tx.GasPrice))
	transactionOJ.Put("from", byteArrayToOJ(tx.From))

	return transactionOJ
}

func blockToOJ(block *Block) oj.OJsonObject {
	blockOJ := oj.NewMap()

	var resultList []oj.OJsonObject
	for _, blr := range block.Results {
		resultList = append(resultList, resultToOJ(blr))
	}
	resultsOJ := oj.OJsonList(resultList)
	blockOJ.Put("results", &resultsOJ)

	var txList []oj.OJsonObject
	for _, tx := range block.Transactions {
		txList = append(txList, transactionToTestOJ(tx))
	}
	txsOJ := oj.OJsonList(txList)
	blockOJ.Put("transactions", &txsOJ)

	blockHeaderOJ := oj.NewMap()
	blockHeaderOJ.Put("gasLimit", intToOJ(block.BlockHeader.GasLimit))
	blockHeaderOJ.Put("number", intToOJ(block.BlockHeader.Number))
	blockHeaderOJ.Put("difficulty", intToOJ(block.BlockHeader.Difficulty))
	blockHeaderOJ.Put("timestamp", uint64ToOJ(block.BlockHeader.Timestamp))
	blockHeaderOJ.Put("coinbase", intToOJ(block.BlockHeader.Beneficiary))
	blockOJ.Put("blockHeader", blockHeaderOJ)

	return blockOJ
}
