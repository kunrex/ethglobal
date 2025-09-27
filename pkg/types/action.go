package types

import (
	"context"
	"ethglobal/pkg/abi"
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

func (c *ContractActions) GetProject(repositoryIdentifier [32]byte) ([]byte, []byte, bool, error) {
	ctx, cancel := context.WithTimeout(c.RootContext, c.GetTimeout)
	defer cancel()

	cid, metaData, exists, err := c.Contract.GetProject(
		&bind.CallOpts{
			Context: ctx,
		},
		repositoryIdentifier,
	)

	if err != nil {
		return nil, nil, false, err
	}
	return cid, metaData, exists, nil
}

func (c *ContractActions) GetProjectMetadata(repositoryIdentifier [32]byte) ([]byte, bool, error) {
	ctx, cancel := context.WithTimeout(c.RootContext, c.GetTimeout)
	defer cancel()

	metaData, exists, err := c.Contract.GetMetaData(
		&bind.CallOpts{
			Context: ctx,
		},
		repositoryIdentifier,
	)

	if err != nil {
		return nil, false, err
	}
	return metaData, exists, nil
}

func (c *ContractActions) SetProject(repositoryIdentifier [32]byte, cid []byte, metaData []byte) (string, error) {
	err := c.Keystore.Unlock(c.Account, "")
	if err != nil {
		panic(err)
	}

	auth, err := bind.NewKeyStoreTransactorWithChainID(c.Keystore, c.Account, c.Chain)
	if err != nil {
		return "", nil
	}

	tx, err := c.Contract.SetProject(
		auth,
		repositoryIdentifier, cid, metaData)
	if err != nil {
		return "", err
	}

	return tx.Hash().Hex(), nil
}
