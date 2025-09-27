package controllers

import (
	"errors"
	"git-server/pkg/types"
	"git-server/pkg/utils"
	"os"
)

type Controller struct {
	ActionContracts *types.ContractActions
	Lighthouse      *types.LighthouseClient
}

func (c Controller) PushColdStorage(repository string, dotGitFile string) (string, error) {
	hash := utils.SHA256(repository)
	bytes, err := os.ReadFile(dotGitFile)
	if err != nil {
		return "", err
	}

	file, err := c.Lighthouse.UploadFile(bytes, c.Lighthouse.ApiKeyBytes)
	if err != nil {
		return "", err
	}

	if !file.Success {
		return "", errors.New("failed to upload to file coin")
	}

	transactionId, err := c.ActionContracts.SetProjectCID(hash, []byte(file.CID))
	if err != nil {
		return "", err
	}

	return transactionId, nil
}

func (c Controller) RetrieveColdStorage(repository string, ) (string, error) {

}
