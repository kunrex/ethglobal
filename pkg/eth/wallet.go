package eth

import (
	"crypto/ecdsa"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/crypto"
)

type AnonymousWallet struct {
	Address    string `json:"address"`
	PrivateKey string `json:"privateKey"`
}

func createWallet() (*AnonymousWallet, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	privateKeyHex := hex.EncodeToString(privateKeyBytes)

	publicKey := privateKey.Public().(*ecdsa.PublicKey)
	address := crypto.PubkeyToAddress(*publicKey).Hex()

	return &AnonymousWallet{
		Address:    address,
		PrivateKey: privateKeyHex,
	}, nil
}

func linkWallet() {
	
}
