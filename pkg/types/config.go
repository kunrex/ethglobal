package types

import (
	"math/big"
	"time"
)

type Configuration struct {
	GetSeconds time.Duration
	SetMinutes time.Duration

	LighthouseKey     string
	ConnectionTimeout time.Duration

	JsonRPC        string
	Chain          *big.Int
	ContactAddress string

	KeystoreDirectory string
}
