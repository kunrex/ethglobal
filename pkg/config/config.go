package config

import (
	"git-server/pkg/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/joho/godotenv"
	"log"
	"math/big"
	"os"
)

func LoadConfig() *types.Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal("error loading .env file")
	}

	rpc := os.Getenv("RPC")
	contractAddress := os.Getenv("CONTRACT_ADDRESS")

	chainID, ok := new(big.Int).SetString(os.Getenv("CHAIN_ID"), 10)
	if !ok {
		log.Fatal("failed to parse chain id CHAIN_ID")
	}

	return &types.Config{
		RPC:             rpc,
		ChainID:         chainID,
		ContractAddress: common.HexToAddress(contractAddress),
	}
}
