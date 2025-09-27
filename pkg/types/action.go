package types

import (
	"context"
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

func (c *ContractActions) GetProjectCID(repositoryIdentifier [32]byte) (bool, []byte, error) {
	ctx, cancel := context.WithTimeout(c.RootContext, c.GetTimeout)
	defer cancel()

	cid, exists, err := c.Contract.GetProject(
		&bind.CallOpts{
			Context: ctx,
		},
		repositoryIdentifier,
	)

	if err != nil {
		return false, nil, err
	}
	return exists, cid, nil
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
