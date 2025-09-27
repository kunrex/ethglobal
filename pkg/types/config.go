package types

import "math/big"

type TimeoutConfigs struct {
	GetSeconds int64
	SetMinutes int64
}

type ContractConfigurationJson struct {
	RPC     string
	ChainID int64
}

type ConfigurationJson struct {
	KeyStoreDirectory string
	LighthouseAPIKey  string
	ContractConfig    ContractConfigurationJson
	Timeout           TimeoutConfigs
}

type ContractConfiguration struct {
	RPC     string
	ChainID *big.Int
}

type Configuration struct {
	KeyStoreDirectory string
	LighthouseAPIKey  string
	ContractConfig    ContractConfiguration
	Timeout           TimeoutConfigs
}
