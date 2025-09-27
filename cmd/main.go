package cmd

import (
	"git-server/pkg/config"
	"git-server/pkg/controllers"
	"git-server/pkg/types"
	"net/http"
)

func initLighthouseClient(configuration types.Configuration) *types.LighthouseClient {
	return &types.LighthouseClient{
		ApiKey:      configuration.LighthouseKey,
		ApiKeyBytes: []byte(configuration.LighthouseKey),
		Client: &http.Client{
			Timeout:   configuration.ConnectionTimeout,
		},
	}
}

func main() {
	configuration := config.LoadConfig()

	actions, rootCtx, err := InitContractActions(&configuration)
	if err != nil {
		return
	}

	controller := controllers.Controller{
		ActionContracts: actions,
		Lighthouse:      initLighthouseClient(configuration),
	}
}
