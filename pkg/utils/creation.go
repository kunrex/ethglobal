package utils

import (
	"git-server/pkg/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func CreateContract(address common.Address, client *ethclient.Client) (*abi.Abi, error) {
	contract, err := abi.NewAbi(address, client)
	return contract, err
}
