package config

import (
	"encoding/json"
	"git-server/pkg/types"
	"log"
	"math/big"
	"os"
)

func LoadConfig() types.Configuration {
	data, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatalf("error encountered reading config: %v", err)
	}

	var configuration types.ConfigurationJson
	if err = json.Unmarshal(data, &configuration); err != nil {
		log.Fatalf("error encountered reading config: %v", err)
	}

	return types.Configuration{
		KeyStoreDirectory: configuration.KeyStoreDirectory,
		ContractConfig: types.ContractConfiguration{
			RPC:     configuration.ContractConfig.RPC,
			ChainID: big.NewInt(configuration.ContractConfig.ChainID),
		},
		Timeout:          configuration.Timeout,
		LighthouseAPIKey: configuration.LighthouseAPIKey,
	}
}
