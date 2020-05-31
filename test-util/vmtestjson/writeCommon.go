package vmtestjson

import (
	"encoding/hex"
	"math/big"

	oj "github.com/ElrondNetwork/elrond-vm-util/test-util/orderedjson"
)

func accountsToOJ(accounts []*Account) oj.OJsonObject {
	acctsOJ := oj.NewMap()
	for _, account := range accounts {
		acctOJ := oj.NewMap()
		acctOJ.Put("nonce", uint64ToOJ(account.Nonce))
		acctOJ.Put("balance", bigIntToOJ(account.Balance))
		storageOJ := oj.NewMap()
		for _, st := range account.Storage {
			storageOJ.Put(byteArrayToString(st.Key), byteArrayToOJ(st.Value))
		}
		acctOJ.Put("storage", storageOJ)
		acctOJ.Put("code", byteArrayToOJ(account.Code))
		if len(account.AsyncCallData) > 0 {
			acctOJ.Put("asyncCallData", stringToOJ(account.AsyncCallData))
		}

		acctsOJ.Put(byteArrayToString(account.Address), acctOJ)
	}

	return acctsOJ
}

func checkAccountsToOJ(checkAccounts *CheckAccounts) oj.OJsonObject {
	acctsOJ := oj.NewMap()
	for _, checkAccount := range checkAccounts.Accounts {
		acctOJ := oj.NewMap()
		acctOJ.Put("nonce", checkUint64ToOJ(checkAccount.Nonce))
		acctOJ.Put("balance", checkBigIntToOJ(checkAccount.Balance))
		storageOJ := oj.NewMap()
		for _, st := range checkAccount.CheckStorage {
			storageOJ.Put(byteArrayToString(st.Key), byteArrayToOJ(st.Value))
		}
		if checkAccount.IgnoreStorage {
			acctOJ.Put("storage", stringToOJ("*"))
		} else {
			acctOJ.Put("storage", storageOJ)
		}
		acctOJ.Put("code", byteArrayToOJ(checkAccount.Code))
		if len(checkAccount.AsyncCallData) > 0 {
			acctOJ.Put("asyncCallData", stringToOJ(checkAccount.AsyncCallData))
		}

		acctsOJ.Put(byteArrayToString(checkAccount.Address), acctOJ)
	}

	if checkAccounts.OtherAccountsAllowed {
		acctsOJ.Put("+", stringToOJ(""))
	}

	return acctsOJ
}

func blockHashesToOJ(blockHashes []JSONBytes) oj.OJsonObject {
	var blockhashesList []oj.OJsonObject
	for _, blh := range blockHashes {
		blockhashesList = append(blockhashesList, byteArrayToOJ(blh))
	}
	blockhashesOJ := oj.OJsonList(blockhashesList)
	return &blockhashesOJ
}

func resultToOJ(res *TransactionResult) oj.OJsonObject {
	resultOJ := oj.NewMap()

	var outList []oj.OJsonObject
	for _, out := range res.Out {
		outList = append(outList, checkBytesToOJ(out))
	}
	outOJ := oj.OJsonList(outList)
	resultOJ.Put("out", &outOJ)

	resultOJ.Put("status", bigIntToOJ(res.Status))
	if len(res.Message) > 0 {
		resultOJ.Put("message", stringToOJ(res.Message))
	}
	if res.IgnoreLogs {
		resultOJ.Put("logs", stringToOJ("*"))
	} else {
		if len(res.LogHash) > 0 {
			resultOJ.Put("logs", stringToOJ(res.LogHash))
		} else {
			resultOJ.Put("logs", logsToOJ(res.Logs))
		}
	}
	resultOJ.Put("gas", checkUint64ToOJ(res.Gas))
	resultOJ.Put("refund", checkBigIntToOJ(res.Refund))

	return resultOJ
}

// LogToString returns a json representation of a log entry, we use it for debugging
func LogToString(logEntry *LogEntry) string {
	logOJ := logToOJ(logEntry)
	return oj.JSONString(logOJ)
}

func logToOJ(logEntry *LogEntry) oj.OJsonObject {
	logOJ := oj.NewMap()
	logOJ.Put("address", byteArrayToOJ(logEntry.Address))
	logOJ.Put("identifier", byteArrayToOJ(logEntry.Identifier))

	var topicsList []oj.OJsonObject
	for _, topic := range logEntry.Topics {
		topicsList = append(topicsList, byteArrayToOJ(topic))
	}
	topicsOJ := oj.OJsonList(topicsList)
	logOJ.Put("topics", &topicsOJ)

	logOJ.Put("data", byteArrayToOJ(logEntry.Data))

	return logOJ
}

func logsToOJ(logEntries []*LogEntry) oj.OJsonObject {
	var logList []oj.OJsonObject
	for _, logEntry := range logEntries {
		logOJ := logToOJ(logEntry)
		logList = append(logList, logOJ)
	}
	logOJList := oj.OJsonList(logList)
	return &logOJList
}

func intToString(i *big.Int) string {
	if i == nil {
		return ""
	}
	if i.Sign() == 0 {
		return "0x00"
	}

	isNegative := i.Sign() == -1
	str := i.Text(16)
	if isNegative {
		str = str[1:] // drop the minus in front
	}
	if len(str)%2 != 0 {
		str = "0" + str
	}
	str = "0x" + str
	if isNegative {
		str = "-" + str
	}
	return str
}

func bigIntToOJ(i JSONBigInt) oj.OJsonObject {
	return &oj.OJsonString{Value: i.Original}
}

func checkBigIntToOJ(i JSONCheckBigInt) oj.OJsonObject {
	return &oj.OJsonString{Value: i.Original}
}

func byteArrayToString(byteArray JSONBytes) string {
	if len(byteArray.Original) == 0 && len(byteArray.Value) > 0 {
		byteArray.Original = hex.EncodeToString(byteArray.Value)
	}
	return byteArray.Original
}

func byteArrayToOJ(byteArray JSONBytes) oj.OJsonObject {
	return &oj.OJsonString{Value: byteArrayToString(byteArray)}
}

func checkBytesToString(checkBytes JSONCheckBytes) string {
	if len(checkBytes.Original) == 0 && len(checkBytes.value) > 0 {
		checkBytes.Original = hex.EncodeToString(checkBytes.value)
	}
	return checkBytes.Original
}

func checkBytesToOJ(checkBytes JSONCheckBytes) oj.OJsonObject {
	return &oj.OJsonString{Value: checkBytesToString(checkBytes)}
}

func uint64ToOJ(i JSONUint64) oj.OJsonObject {
	return &oj.OJsonString{Value: i.Original}
}

func checkUint64ToOJ(i JSONCheckUint64) oj.OJsonObject {
	return &oj.OJsonString{Value: i.Original}
}

func stringToOJ(str string) oj.OJsonObject {
	return &oj.OJsonString{Value: str}
}
