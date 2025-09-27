package eth

import (
	"context"
	"errors"
	"git-server/pkg/abi"
	"git-server/pkg/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"math/big"
)

type ContractActions struct {
	chain    *big.Int
	contract *abi.Abi
	wallet   *types.AnonymousWallet
}

func (c *ContractActions) Exists(repositoryIdentifier *[32]byte) (bool, error) {
	_, exists, err := c.contract.GetProject(
		&bind.CallOpts{Context: context.Background()},
		*repositoryIdentifier,
	)

	if err != nil {
		return false, err
	}

	return exists, nil
}

func (c *ContractActions) GetProjectCID(repositoryIdentifier [32]byte) (*[32]byte, error) {
	cid, exists, err := c.contract.GetProject(
		&bind.CallOpts{Context: context.Background()},
		repositoryIdentifier,
	)

	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, errors.New("project not found")
	}

	return &cid, nil
}

func (c *ContractActions) SetProjectCID(repositoryIdentifier [32]byte, cid *[32]byte) (string, error) {
	auth, err := bind.NewKeyedTransactorWithChainID(c.wallet.PrivateKey, c.chain)
	if err != nil {
		return "", err
	}

	tx, err := c.contract.SetProject(auth, repositoryIdentifier, *cid)
	if err != nil {
		return "", err
	}

	return tx.Hash().Hex(), nil
}
