package vmtestjson

// Scenario is a json object representing a test scenario with steps.
type Scenario struct {
	Name     string
	Comment  string
	CheckGas bool
	Steps    []Step
}

type Step interface {
	StepTypeName() string
}

type SetStateStep struct {
	Accounts    []*Account
	BlockHashes [][]byte
}

type CheckStateStep struct {
	CheckAccounts []*Account
}

type TxStep struct {
	Tx             *Transaction
	ExpectedResult *TransactionResult
}

var _ Step = (*SetStateStep)(nil)
var _ Step = (*CheckStateStep)(nil)
var _ Step = (*TxStep)(nil)

const stepNameSetState = "setState"

func (*SetStateStep) StepTypeName() string {
	return stepNameSetState
}

const stepNameCheckState = "checkState"

func (*CheckStateStep) StepTypeName() string {
	return stepNameCheckState
}

const stepNameTx = "tx"

func (*TxStep) StepTypeName() string {
	return stepNameTx
}
