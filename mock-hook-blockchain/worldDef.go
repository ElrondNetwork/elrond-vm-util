package callbackblockchain

// BlockchainHookMock provides a mock representation of the blockchain to be used in VM tests.
type BlockchainHookMock struct {
	AcctMap                      AccountMap
	CurrentTimestamp             uint64
	Blockhashes                  [][]byte
	mockAddressGenerationEnabled bool
}

// NewMock creates a new mock instance
func NewMock() *BlockchainHookMock {
	return &BlockchainHookMock{
		AcctMap:                      NewAccountMap(),
		CurrentTimestamp:             0,
		Blockhashes:                  nil,
		mockAddressGenerationEnabled: false,
	}
}

// Clear resets all mock data between tests.
func (b *BlockchainHookMock) Clear() {
	b.AcctMap = NewAccountMap()
	b.Blockhashes = nil
}

// EnableMockAddressGeneration causes the mock to generate its own new addresses.
func (b *BlockchainHookMock) EnableMockAddressGeneration() {
	b.mockAddressGenerationEnabled = true
}
