package types

import (
	"bytes"
	"encoding/json"
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

func (lh *LighthouseClient) UploadFile(plainBuf []byte, key []byte) (*UploadResponse, error) {
	//cipherText, err := utils.Encrypt(key, plainBuf)
	cipherText := plainBuf
	var cipherBuf bytes.Buffer
	cipherBuf.Write(cipherText)
	//if err != nil {
	//	return nil, fmt.Errorf("failed to encrypt file: %v", err)
	//}

	// Prepare multipart form
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// "file" is the field name from curl -F 'file=@...'
	part, err := writer.CreateFormFile("file", "unnamed.jpeg")
	if err != nil {
		panic(err)
	}
	_, err = part.Write(cipherText)
	if err != nil {
		panic(err)
	}

	writer.Close()

	// Create HTTP request to Lighthouse API
	req, err := http.NewRequest("POST", "https://upload.lighthouse.storage/api/v0/add", &buf)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
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
	var result map[string]string

	err := json.Unmarshal([]byte(cid), &result)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return nil, err
	}

	fmt.Println(result)

	req, err := http.NewRequest("GET", fmt.Sprintf("https://gateway.lighthouse.storage/ipfs/%s", result["Hash"]), nil)
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

	//plainBuf, err := utils.Decrypt(key, cipherBuf)
	plainBuf := cipherBuf
	/*if err != nil {
		return nil, fmt.Errorf("failed to decrypt file: %v", err)
	}*/

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
