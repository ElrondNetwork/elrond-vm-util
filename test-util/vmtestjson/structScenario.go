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

// NewAddressMock allows tests to specify what new addresses to generate
type NewAddressMock struct {
	CreatorAddress JSONBytes
	CreatorNonce   JSONUint64
	NewAddress     JSONBytes
}

// BlockInfo contains data for the block info hooks
type BlockInfo struct {
	BlockTimestamp JSONUint64
	BlockNonce     JSONUint64
	BlockRound     JSONUint64
	BlockEpoch     JSONUint64
}

// SetStateStep is a step where data is saved to the blockchain mock.
type SetStateStep struct {
	Comment           string
	Accounts          []*Account
	PreviousBlockInfo *BlockInfo
	CurrentBlockInfo  *BlockInfo
	BlockHashes       []JSONBytes
	NewAddressMocks   []*NewAddressMock
}

// CheckStateStep is a step where the state of the blockchain mock is verified.
type CheckStateStep struct {
	Comment       string
	CheckAccounts []*CheckAccount
}

// TxStep is a step where a transaction is executed.
type TxStep struct {
	TxIdent        string
	Comment        string
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

const stepNameScCall = "scCall"
const stepNameScDeploy = "scDeploy"
const stepNameTransfer = "transfer"

// StepTypeName type as string
func (t *TxStep) StepTypeName() string {
	switch t.Tx.Type {
	case ScCall:
		return stepNameScCall
	case ScDeploy:
		return stepNameScDeploy
	case Transfer:
		return stepNameTransfer
	default:
		panic("unknown TransactionType")
	}
}
