package eth

import (
	"crypto/ecdsa"
	"encoding/hex"
	"git-server/pkg/types/eth"
	"github.com/ethereum/go-ethereum/crypto"
)

func createWallet() (*eth.AnonymousWallet, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	privateKeyHex := hex.EncodeToString(privateKeyBytes)

	publicKey := privateKey.Public().(*ecdsa.PublicKey)
	address := crypto.PubkeyToAddress(*publicKey).Hex()

	return &eth.AnonymousWallet{
		Address:    address,
		PrivateKey: privateKeyHex,
	}, nil
}
