package lighthouse

import (
	"ethglobal/pkg/types"
	"net/http"
)

func InitLightHouseClient(configuration types.Configuration) *types.LighthouseClient {
	return &types.LighthouseClient{
		ApiKey:      configuration.LighthouseKey,
		ApiKeyBytes: []byte(configuration.LighthouseKey),
		Client: &http.Client{
			Timeout: configuration.ConnectionTimeout,
		},
	}
}
