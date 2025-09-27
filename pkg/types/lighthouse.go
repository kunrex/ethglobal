package types

import (
	"bytes"
	"encoding/json"
	"ethglobal/pkg/utils"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

type LighthouseClient struct {
	ApiKey      string
	ApiKeyBytes []byte
	Client      *http.Client
}

func (lh *LighthouseClient) UploadFile(plainBuf []byte, encryptionKey []byte, commitHash string) (string, error) {
	cipherText, err := utils.Encrypt(encryptionKey, plainBuf)

	var cipherBuffer bytes.Buffer
	writer := multipart.NewWriter(&cipherBuffer)

	part, err := writer.CreateFormFile("file", fmt.Sprintf("%v.git", commitHash))
	if err != nil {
		return "", err
	}
	_, err = part.Write(cipherText)
	if err != nil {
		return "", err
	}

	defer func(writer *multipart.Writer) {
		_ = writer.Close()
	}(writer)

	req, err := http.NewRequest("POST", "https://upload.lighthouse.storage/api/v0/add", &cipherBuffer)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+lh.ApiKey)

	resp, err := lh.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %v", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("lighthouse API error: %s", string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	return strings.TrimSpace(string(body)), nil
}

func (lh *LighthouseClient) DownloadFile(cid string, encryptionKey []byte) ([]byte, error) {
	var result map[string]string
	err := json.Unmarshal([]byte(cid), &result)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://gateway.lighthouse.storage/ipfs/%s", result["Hash"]), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+lh.ApiKey)

	resp, err := lh.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %v", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("lighthouse API error: %s", string(body))
	}

	cipherBuf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read file content: %v", err)
	}

	plainBuf, err := utils.Decrypt(encryptionKey, cipherBuf)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt file: %v", err)
	}

	return plainBuf, nil
}
