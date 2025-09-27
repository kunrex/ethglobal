package controllers

import (
	"encoding/json"
	"errors"
	"ethglobal/pkg/types"
	"ethglobal/pkg/utils"
	"os"
)

type Controller struct {
	EncryptionKeyBytes []byte
	ActionContracts    *types.ContractActions
	Lighthouse         *types.LighthouseClient
}

func (c Controller) calculateMetaData(hash [32]byte, commitHash string) ([]byte, error) {
	metaData, exists, err := c.ActionContracts.GetProjectMetadata(hash)
	if err != nil {
		return nil, err
	}

	var versions []types.VersionMetaData
	if exists {
		err = json.Unmarshal(metaData, &versions)
		if err != nil {
			return nil, err
		}

		next := types.VersionMetaData{
			Version:    uint32(len(versions) + 1),
			CommitHash: commitHash,
		}
		versions = append(versions, next)
	} else {
		next := types.VersionMetaData{
			Version:    1,
			CommitHash: commitHash,
		}
		versions = append(versions, next)
	}

	marshalledMetaData, err := json.Marshal(versions)
	if err != nil {
		return nil, err
	}
	return marshalledMetaData, nil
}

func (c Controller) PushColdStorage(repository string, dotGitFile string, commitHash string) (string, error) {
	hash := utils.SHA256(repository)
	marshalledMetaData, err := c.calculateMetaData(hash, commitHash)
	if err != nil {
		return "", err
	}

	bytes, err := os.ReadFile(dotGitFile)
	if err != nil {
		return "", err
	}

	cid, err := c.Lighthouse.UploadFile(bytes, c.EncryptionKeyBytes, commitHash)
	if err != nil {
		return "", err
	}

	metaDataCid, err := c.Lighthouse.UploadFile(marshalledMetaData, c.EncryptionKeyBytes, commitHash+"_meta")
	if err != nil {
		return "", err
	}

	transactionId, err := c.ActionContracts.SetProject(hash, []byte(cid), []byte(metaDataCid))
	if err != nil {
		return "", err
	}

	return transactionId, nil
}

func (c Controller) RetrieveLatestMetaData(repository string) ([]byte, error) {
	hash := utils.SHA256(repository)
	metaDataCid, exists, err := c.ActionContracts.GetProjectMetadata(hash)
	if err != nil {
		return nil, err
	}

	if exists {
		metaData, err := c.Lighthouse.DownloadFile(string(metaDataCid), c.EncryptionKeyBytes)
		if err != nil {
			return nil, err
		}

		return metaData, nil
	} else {
		return nil, nil
	}
}

func (c Controller) RetrieveColdStorage(repository string, output string) ([]byte, error) {
	hash := utils.SHA256(repository)
	cid, metaDataCid, exists, err := c.ActionContracts.GetProject(hash)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, errors.New("failed to retrieve project code")
	}

	data, err := c.Lighthouse.DownloadFile(string(cid), c.EncryptionKeyBytes)
	if err != nil {
		return nil, err
	}

	metaData, err := c.Lighthouse.DownloadFile(string(metaDataCid), c.EncryptionKeyBytes)
	if err != nil {
		return nil, err
	}

	err = os.WriteFile(output, data, 0644)
	if err != nil {
		return nil, err
	}

	return metaData, nil
}
