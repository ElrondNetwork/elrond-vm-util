package callbackblockchain

// BlockchainHookMock provides a mock representation of the blockchain to be used in VM tests.
type BlockchainHookMock struct {
	AcctMap     AccountMap
	Blockhashes [][]byte
}

// NewMock creates a new mock instance
func NewMock() *BlockchainHookMock {
	return &BlockchainHookMock{
		AcctMap:     NewAccountMap(),
		Blockhashes: nil,
	}
}

// Clear resets all mock data between tests.
func (b *BlockchainHookMock) Clear() {
	b.AcctMap = NewAccountMap()
	b.Blockhashes = nil
}
