package eth

import (
	"github.com/ethereum/go-ethereum/ethclient"
)

type AnonymousWallet struct {
	Address    string `json:"address"`
	PrivateKey string `json:"privateKey"`
}

func (*AnonymousWallet) Connect(rpc string) (*ethclient.Client, error) {
	client, err := ethclient.Dial(rpc)
	return client, err
}
