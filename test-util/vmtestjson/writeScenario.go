package vmtestjson

import (
	oj "github.com/ElrondNetwork/elrond-vm-util/test-util/orderedjson"
)

// ScenarioToJSONString converts a scenario object to its JSON representation.
func ScenarioToJSONString(scenario *Scenario) string {
	jobj := ScenarioToOrderedJSON(scenario)
	return oj.JSONString(jobj)
}

// ScenarioToOrderedJSON converts a scenario object to an ordered JSON object.
func ScenarioToOrderedJSON(scenario *Scenario) oj.OJsonObject {
	scenarioOJ := oj.NewMap()

	if len(scenario.Name) > 0 {
		scenarioOJ.Put("name", stringToOJ(scenario.Name))
	}

	if len(scenario.Comment) > 0 {
		scenarioOJ.Put("comment", stringToOJ(scenario.Comment))
	}

	if !scenario.CheckGas {
		ojFalse := oj.OJsonBool(false)
		scenarioOJ.Put("checkGas", &ojFalse)
	}

	var stepOJList []oj.OJsonObject

	for _, generalStep := range scenario.Steps {
		stepOJ := oj.NewMap()
		stepOJ.Put("step", stringToOJ(generalStep.StepTypeName()))
		switch step := generalStep.(type) {
		case *SetStateStep:
			if len(step.Comment) > 0 {
				stepOJ.Put("comment", stringToOJ(step.Comment))
			}
			stepOJ.Put("accounts", accountsToOJ(step.Accounts))
			if len(step.NewAddressMocks) > 0 {
				stepOJ.Put("newAddresses", newAddressMocksToOJ(step.NewAddressMocks))
			}
			if len(step.BlockHashes) > 0 {
				stepOJ.Put("blockhashes", blockHashesToOJ(step.BlockHashes))
			}
		case *CheckStateStep:
			if len(step.Comment) > 0 {
				stepOJ.Put("comment", stringToOJ(step.Comment))
			}
			stepOJ.Put("accounts", checkAccountsToOJ(step.CheckAccounts))
		case *TxStep:
			if len(step.TxIdent) > 0 {
				stepOJ.Put("txId", stringToOJ(step.TxIdent))
			}
			if len(step.Comment) > 0 {
				stepOJ.Put("comment", stringToOJ(step.Comment))
			}
			stepOJ.Put("tx", transactionToScenarioOJ(step.Tx))
			stepOJ.Put("expect", resultToOJ(step.ExpectedResult))
		}

		stepOJList = append(stepOJList, stepOJ)
	}

	stepsOJ := oj.OJsonList(stepOJList)
	scenarioOJ.Put("steps", &stepsOJ)

	return scenarioOJ
}

func transactionToScenarioOJ(tx *Transaction) oj.OJsonObject {
	transactionOJ := oj.NewMap()
	transactionOJ.Put("from", byteArrayToOJ(tx.From))
	transactionOJ.Put("to", byteArrayToOJ(tx.To))
	transactionOJ.Put("value", bigIntToOJ(tx.Value))
	transactionOJ.Put("function", stringToOJ(tx.Function))
	transactionOJ.Put("contractCode", byteArrayToOJ(tx.Code))

	var argList []oj.OJsonObject
	for _, arg := range tx.Arguments {
		argList = append(argList, byteArrayToOJ(arg))
	}
	argOJ := oj.OJsonList(argList)
	transactionOJ.Put("arguments", &argOJ)

	transactionOJ.Put("gasLimit", uint64ToOJ(tx.GasLimit))
	transactionOJ.Put("gasPrice", uint64ToOJ(tx.GasPrice))

	return transactionOJ
}

func newAddressMocksToOJ(newAddressMocks []*NewAddressMock) oj.OJsonObject {
	var namList []oj.OJsonObject
	for _, namEntry := range newAddressMocks {
		namOJ := oj.NewMap()
		namOJ.Put("creatorAddress", byteArrayToOJ(namEntry.CreatorAddress))
		namOJ.Put("creatorNonce", uint64ToOJ(namEntry.CreatorNonce))
		namOJ.Put("newAddress", byteArrayToOJ(namEntry.NewAddress))
		namList = append(namList, namOJ)
	}
	namOJList := oj.OJsonList(namList)
	return &namOJList
}
