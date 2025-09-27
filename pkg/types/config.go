package types

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type Config struct {
	RPC             string
	ChainID         *big.Int
	ContractAddress common.Address
}
