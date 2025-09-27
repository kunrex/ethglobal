package types

import (
	"bytes"
	"ethglobal/pkg/utils"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type LighthouseClient struct {
	ApiKey      string
	ApiKeyBytes []byte
	Client      *http.Client
}

func (lh *LighthouseClient) UploadFile(plainBuf []byte, key []byte) (*UploadResponse, error) {
	cipherText, err := utils.Encrypt(key, plainBuf)
	var cipherBuf bytes.Buffer
	cipherBuf.Write(cipherText)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt file: %v", err)
	}

	// Create HTTP request to Lighthouse API
	req, err := http.NewRequest("POST", "https://api.lighthouse.storage/api/v0/add", &cipherBuf)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Authorization", "Bearer "+lh.ApiKey)

	// Send request
	resp, err := lh.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("lighthouse API error: %s", string(body))
	}

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	// Parse CID from response (Lighthouse returns just the CID string)
	cid := strings.TrimSpace(string(body))

	return &UploadResponse{
		Success: true,
		CID:     cid,
		Message: "File uploaded successfully to Filecoin network",
	}, nil
}

func (lh *LighthouseClient) DownloadFile(cid string, key []byte) ([]byte, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.lighthouse.storage/api/v0/cat?arg=%s", cid), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+lh.ApiKey)

	// Send request
	resp, err := lh.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("lighthouse API error: %s", string(body))
	}

	// Read file content
	cipherBuf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read file content: %v", err)
	}

	plainBuf, err := utils.Decrypt(key, cipherBuf)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt file: %v", err)
	}

	return plainBuf, nil
}

func (lh *LighthouseClient) GetFileInfo(cid string, key []byte) (*FileInfo, error) {
	fileContent, err := lh.DownloadFile(cid, key)
	if err != nil {
		return nil, err
	}

	contentType := http.DetectContentType(fileContent)

	fileInfo := &FileInfo{
		CID:      cid,
		Filename: fmt.Sprintf("file_%s", cid[:8]),
		Size:     int64(len(fileContent)),
		Type:     contentType,
	}

	return fileInfo, nil
}
