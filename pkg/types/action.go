package types

import (
	"context"
	"errors"
	"git-server/pkg/abi"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"math/big"
	"time"
)

type ContractActions struct {
	Chain      *big.Int
	GetTimeout time.Duration
	SetTimeout time.Duration

	Account  accounts.Account
	Keystore *keystore.KeyStore

	Contract    *abi.Abi
	RootContext context.Context
}

func (c *ContractActions) CheckProjectExists(repositoryIdentifier *[32]byte) (bool, error) {
	ctx, cancel := context.WithTimeout(c.RootContext, c.GetTimeout)

	defer cancel()

	_, exists, err := c.Contract.GetProject(
		&bind.CallOpts{
			Context: ctx,
		},
		*repositoryIdentifier,
	)

	if err != nil {
		return false, err
	}

	return exists, nil
}

func (c *ContractActions) GetProjectCID(repositoryIdentifier [32]byte) ([]byte, error) {
	ctx, cancel := context.WithTimeout(c.RootContext, c.GetTimeout)
	defer cancel()

	cid, exists, err := c.Contract.GetProject(
		&bind.CallOpts{
			Context: ctx,
		},
		repositoryIdentifier,
	)

	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, errors.New("project not found")
	}

	return cid, nil
}

func (c *ContractActions) SetProjectCID(repositoryIdentifier [32]byte, cid []byte) (string, error) {
	_, err := bind.NewKeyStoreTransactorWithChainID(c.Keystore, c.Account, c.Chain)
	if err != nil {
		return "", nil
	}

	tx, err := c.Contract.SetProject(
		&bind.TransactOpts{},
		repositoryIdentifier, cid)
	if err != nil {
		return "", err
	}

	return tx.Hash().Hex(), nil
}
