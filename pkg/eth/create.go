package eth

import (
	"git-server/pkg/abi"
	"git-server/pkg/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

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

func CreateContract(wallet *types.AnonymousWallet, client *ethclient.Client) (*abi.Abi, error) {
	contract, err := abi.NewAbi(wallet.Address, client)
	return contract, err
}
