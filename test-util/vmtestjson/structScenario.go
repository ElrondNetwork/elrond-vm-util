package vmtestjson

// Scenario is a json object representing a test scenario with steps.
type Scenario struct {
	Steps []Step
}

type Step interface {
	StepTypeName() string
}

type InitStateStep struct {
	InitAccounts []*Account
}

type CheckStateStep struct {
	CheckAccounts []*Account
}

type TxStep struct {
	Tx       *Transaction
	Expected *TransactionResult
}

var _ Step = (*InitStateStep)(nil)
var _ Step = (*CheckStateStep)(nil)
var _ Step = (*TxStep)(nil)

func (*InitStateStep) StepTypeName() string {
	return "initState"
}

func (*CheckStateStep) StepTypeName() string {
	return "checkState"
}

func (*TxStep) StepTypeName() string {
	return "tx"
}
