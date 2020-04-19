package vmtestjson

// Scenario is a json object representing a test scenario with steps.
type Scenario struct {
	Name     string
	Comment  string
	CheckGas bool
	Steps    []Step
}

// Step is the basic block of a scenario.
type Step interface {
	StepTypeName() string
}

// SetStateStep is a step where data is saved to the blockchain mock.
type SetStateStep struct {
	Accounts    []*Account
	BlockHashes [][]byte
}

// CheckStateStep is a step where the state of the blockchain mock is verified.
type CheckStateStep struct {
	CheckAccounts []*Account
}

// TxStep is a step where a transaction is executed.
type TxStep struct {
	Tx             *Transaction
	ExpectedResult *TransactionResult
}

var _ Step = (*SetStateStep)(nil)
var _ Step = (*CheckStateStep)(nil)
var _ Step = (*TxStep)(nil)

const stepNameSetState = "setState"

// StepTypeName type as string
func (*SetStateStep) StepTypeName() string {
	return stepNameSetState
}

const stepNameCheckState = "checkState"

// StepTypeName type as string
func (*CheckStateStep) StepTypeName() string {
	return stepNameCheckState
}

const stepNameTx = "tx"

// StepTypeName type as string
func (*TxStep) StepTypeName() string {
	return stepNameTx
}
