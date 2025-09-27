package cmd

import (
	"context"
	"errors"
	"git-server/pkg/abi"
	"git-server/pkg/types"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func initKeystoreWallet(ks *keystore.KeyStore) (*accounts.Account, error) {
	switch len(ks.Accounts()) {
	case 0:
		{
			//ui or tui
			password := "secret"
			account, err := ks.NewAccount(password)

			if err != nil {
				return nil, err
			}

			return &account, nil
		}
	case 1:
		return &ks.Accounts()[0], nil
	default:
		return nil, errors.New("too many accounts found")
	}
}

func initContract(address common.Address, client *ethclient.Client) (*abi.Abi, error) {
	contract, err := abi.NewAbi(address, client)
	return contract, err
}

// InitContractActions creates a usable context action that can be used to perform any action on the contract
func InitContractActions(configuration *types.Configuration) (*types.ContractActions, *context.Context, error) {
	ks := keystore.NewKeyStore(configuration.KeyStoreDirectory, keystore.StandardScryptN, keystore.StandardScryptP)
	accountPtr, err := initKeystoreWallet(ks)
	if err != nil {
		return nil, nil, err
	}

	account := *accountPtr

	ctx := context.Background()
	client, err := ethclient.DialContext(ctx, configuration.ContractConfig.RPC)
	if err != nil {
		return nil, nil, err
	}

	contract, err := initContract(account.Address, client)
	if err != nil {
		return nil, nil, err
	}

	return &types.ContractActions{
		Chain:    configuration.ContractConfig.ChainID,
		Contract: contract,
		Account:  account,
		Keystore: ks,
	}, &ctx, nil
}
