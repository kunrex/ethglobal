package contracts

import (
	"crypto/ecdsa"
	"git-server/pkg/abi"
	"git-server/pkg/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// creates an anonymous wallet for a project
func CreateWallet() (*types.AnonymousWallet, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}

	publicKey := privateKey.Public().(*ecdsa.PublicKey)

	return &types.AnonymousWallet{
		PrivateKey: privateKey,
		Address:    common.HexToAddress(crypto.PubkeyToAddress(*publicKey).Hex()),
	}, nil
}

// creates a contract used to call the GetProject and SetProject methods on the testnet
func CreateContract(wallet *types.AnonymousWallet, client *ethclient.Client) (*abi.Abi, error) {
	contract, err := abi.NewAbi(wallet.Address, client)
	return contract, err
}
