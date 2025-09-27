package types

import (
	"math/big"
	"time"
)

type Configuration struct {
	GetSeconds time.Duration
	SetMinutes time.Duration

	LighthouseKey string

	JsonRPC string
	Chain   *big.Int

	KeystoreDirectory string
}
