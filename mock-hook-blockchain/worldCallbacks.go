package callbackblockchain

import (
	"bytes"
	"errors"
	"math/big"

	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

var zero = big.NewInt(0)

// AccountExists yields whether or not an account is considered to exist in the blockchain or not.
// Note: on Ethereum an account with Balance = 0 and Nonce = 0 is considered to not exist.
// In Elrond we can have existing accounts with zer obalance and nonce.
func (b *BlockchainHookMock) AccountExists(address []byte) (bool, error) {
	acct := b.AcctMap.GetAccount(address)
	if acct == nil {
		return false, nil
	}
	return acct.Exists, nil
}

// NewAddress adapts between K model and elrond function
func (b *BlockchainHookMock) NewAddress(creatorAddress []byte, creatorNonce uint64, _ []byte) ([]byte, error) {
	// explicit new address mocks
	for _, newAddressMock := range b.NewAddressMocks {
		if bytes.Equal(creatorAddress, newAddressMock.CreatorAddress) && creatorNonce == newAddressMock.CreatorNonce {
			return newAddressMock.NewAddress, nil
		}
	}

	// a simple mock algorithm
	if b.mockAddressGenerationEnabled {
		result := make([]byte, 32)
		result[10] = 0x11
		result[11] = 0x11
		result[12] = 0x11
		result[13] = 0x11
		copy(result[14:29], creatorAddress)

		result[29] = byte(creatorNonce)

		copy(result[30:], creatorAddress[30:])

		return result, nil
	}
	// empty byte array signals not implemented, fallback to default
	return []byte{}, nil
}

// GetBalance should retrieve the balance of an account
func (b *BlockchainHookMock) GetBalance(address []byte) (*big.Int, error) {
	acct := b.AcctMap.GetAccount(address)
	if acct == nil {
		return zero, nil
	}
	return acct.Balance, nil
}

// GetNonce should retrieve the nonce of an account
func (b *BlockchainHookMock) GetNonce(address []byte) (uint64, error) {
	acct := b.AcctMap.GetAccount(address)
	if acct == nil {
		return 0, nil
	}
	return acct.Nonce, nil
}

// GetStorageData yields the storage value for a certain account and index.
// Should return an empty byte array if the key is missing from the account storage
func (b *BlockchainHookMock) GetStorageData(accountAddress []byte, index []byte) (vmcommon.StorageData, error) {
	acct := b.AcctMap.GetAccount(accountAddress)
	if acct == nil {
		return vmcommon.StorageData{Data: []byte{}}, nil
	}
	return acct.StorageValue(string(index)), nil
}

// IsCodeEmpty should return whether of not an account is SC.
func (b *BlockchainHookMock) IsCodeEmpty(address []byte) (bool, error) {
	acct := b.AcctMap.GetAccount(address)
	if acct == nil {
		return true, nil
	}
	return len(acct.Code) == 0, nil
}

// GetCode should return the compiled and assembled SC code.
// Empty byte array if the account is a wallet.
func (b *BlockchainHookMock) GetCode(address []byte) ([]byte, error) {
	acct := b.AcctMap.GetAccount(address)
	if acct == nil {
		return []byte{}, nil
	}
	return acct.Code, nil
}

// GetBlockhash should return the hash of the nth previous blockchain.
// Offset specifies how many blocks we need to look back.
func (b *BlockchainHookMock) GetBlockhash(nonce uint64) ([]byte, error) {
	offsetInt32 := int(nonce)
	if offsetInt32 >= len(b.Blockhashes) {
		return nil, errors.New("blockhash offset exceeds the blockhashes slice")
	}
	return b.Blockhashes[offsetInt32], nil
}

// LastNonce returns the nonce from from the last committed block
func (b *BlockchainHookMock) LastNonce() uint64 {
	return 0
}

// LastRound returns the round from the last committed block
func (b *BlockchainHookMock) LastRound() uint64 {
	return 0
}

// LastTimeStamp returns the timeStamp from the last committed block
func (b *BlockchainHookMock) LastTimeStamp() uint64 {
	return 0
}

// LastRandomSeed returns the random seed from the last committed block
func (b *BlockchainHookMock) LastRandomSeed() []byte {
	return nil
}

// LastEpoch returns the epoch from the last committed block
func (b *BlockchainHookMock) LastEpoch() uint32 {
	return 0
}

// GetStateRootHash returns the state root hash from the last committed block
func (b *BlockchainHookMock) GetStateRootHash() []byte {
	return nil
}

// CurrentNonce returns the nonce from the current block
func (b *BlockchainHookMock) CurrentNonce() uint64 {
	return 0
}

// CurrentRound returns the round from the current block
func (b *BlockchainHookMock) CurrentRound() uint64 {
	return 0
}

// CurrentTimeStamp return the timestamp from the current block
func (b *BlockchainHookMock) CurrentTimeStamp() uint64 {
	return b.CurrentTimestamp
}

// CurrentRandomSeed returns the random seed from the current header
func (b *BlockchainHookMock) CurrentRandomSeed() []byte {
	return nil
}

// CurrentEpoch returns the current epoch
func (b *BlockchainHookMock) CurrentEpoch() uint32 {
	return 0
}

// ProcessBuiltInFunction -
func (b *BlockchainHookMock) ProcessBuiltInFunction(_ *vmcommon.ContractCallInput) (*big.Int, uint64, error) {
	return big.NewInt(0), 0, nil
}

// GetBuiltinFunctionNames -
func (b *BlockchainHookMock) GetBuiltinFunctionNames() vmcommon.FunctionNames {
	return make(vmcommon.FunctionNames)
}

// GetAllState -
func (b *BlockchainHookMock) GetAllState(_ []byte) (map[string][]byte, error) {
	return make(map[string][]byte), nil
}
