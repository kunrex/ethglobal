package config

import (
	"ethglobal/pkg/types"
	"github.com/joho/godotenv"
	"log"
	"math/big"
	"os"
	"strconv"
	"time"
)

func readInt(variable string, address *int) {
	var temp int
	temp, err := strconv.Atoi(os.Getenv(variable))
	if err != nil {
		log.Fatalf("error converting %v to integer", variable)
	}

	*address = temp
}

func readString(variable string, address *string) {
	var temp string
	temp = os.Getenv(variable)
	*address = temp
}

func LoadConfig() types.Configuration {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	var configuration types.Configuration

	var seconds int
	var minutes int
	readInt("GET_SECONDS", &seconds)
	readInt("SET_MINUTES", &minutes)
	configuration.GetSeconds = time.Second * time.Duration(seconds)
	configuration.SetMinutes = time.Minute * time.Duration(minutes)

	readString("LIGHTHOUSE_KEY", &configuration.LighthouseKey)
	readInt("CONNECTION_TIMEOUT_SECONDS", &seconds)
	configuration.ConnectionTimeout = time.Second * time.Duration(seconds)

	readString("KEYSTORE_DIRECTORY", &configuration.KeystoreDirectory)

	var chain int
	readInt("CHAIN", &chain)
	configuration.Chain = big.NewInt(int64(chain))
	readString("CONTRACT_ADDRESS", &configuration.ContactAddress)

	readString("JSON_RPC", &configuration.JsonRPC)

	readString("ENCRYPTION_KEY", &configuration.EncryptionKey)

	return configuration
}
