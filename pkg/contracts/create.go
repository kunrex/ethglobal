package contracts

import (
	"git-server/pkg/abi"
	"git-server/pkg/types"

	"github.com/ethereum/go-ethereum/ethclient"
)

func CreateContract(wallet *types.AnonymousWallet, client *ethclient.Client) (*abi.Abi, error) {
	contract, err := abi.NewAbi(wallet.Address, client)
	return contract, err
}
