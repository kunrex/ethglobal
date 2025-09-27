package types

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type AnonymousWallet struct {
	Address    common.Address
	PrivateKey *ecdsa.PrivateKey
}

func (*AnonymousWallet) Connect(rpc string) (*ethclient.Client, error) {
	client, err := ethclient.Dial(rpc)
	return client, err
}
