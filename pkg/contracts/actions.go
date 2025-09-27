package contracts

import (
	"context"
	"errors"
	"git-server/pkg/abi"
	"git-server/pkg/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"math/big"
)

func ExistsWallet(wallet *types.AnonymousWallet, contract *abi.Abi) (bool, error) {
	_, exists, err := contract.GetProject(
		&bind.CallOpts{Context: context.Background()},
		wallet.Address,
	)

	if err != nil {
		return false, err
	}

	return exists, nil
}

func GetProjectCID(wallet *types.AnonymousWallet, contract *abi.Abi) (*[32]byte, error) {
	cid, exists, err := contract.GetProject(
		&bind.CallOpts{Context: context.Background()},
		wallet.Address,
	)

	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, errors.New("project not found")
	}

	return &cid, nil
}

func SetProjectCID(chain *big.Int, wallet *types.AnonymousWallet, contract *abi.Abi, cid []byte) (string, error) {
	auth, err := bind.NewKeyedTransactorWithChainID(wallet.PrivateKey, chain)
	if err != nil {
		return "", err
	}

	var cidBytes [32]byte
	copy(cidBytes[:], cid)

	tx, err := contract.SetProject(auth, cidBytes)
	if err != nil {
		return "", err
	}

	return tx.Hash().Hex(), nil
}
