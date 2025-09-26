package api

import (
	"bytes"
	"fmt"
	"git-server/pkg/crypto"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

// LighthouseAPI handles Lighthouse (Filecoin) API operations
type LighthouseAPI struct {
	apiKey string
	client *http.Client
}

// NewLighthouseAPI creates a new Lighthouse API client
func NewLighthouseAPI(apiKey string) *LighthouseAPI {
	return &LighthouseAPI{
		apiKey: apiKey,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// UploadResponse represents the response from a file upload
type UploadResponse struct {
	Success bool   `json:"success"`
	CID     string `json:"cid"`
	Message string `json:"message"`
	Size    int64  `json:"size"`
}

// DownloadResponse represents the response from a file download
type DownloadResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Size    int64  `json:"size"`
}

// FileInfo represents file information
type FileInfo struct {
	CID      string `json:"cid"`
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
	Type     string `json:"type"`
}

// UploadFile uploads a file to Lighthouse (Filecoin)
func (lh *LighthouseAPI) UploadFile(file io.Reader, filename string, key []byte) (*UploadResponse, error) {
	// Create multipart form data
	var plainBuf bytes.Buffer
	writer := multipart.NewWriter(&plainBuf)

	// Add file field
	fileWriter, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %v", err)
	}

	_, err = io.Copy(fileWriter, file)
	if err != nil {
		return nil, fmt.Errorf("failed to copy file data: %v", err)
	}

	writer.Close()

	cipherText, err := crypto.Encrypt(key, plainBuf.Bytes())
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

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+lh.apiKey)

	// Send request
	resp, err := lh.client.Do(req)
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

// DownloadFile downloads a file from Lighthouse (Filecoin)
func (lh *LighthouseAPI) DownloadFile(cid string, key []byte) ([]byte, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.lighthouse.storage/api/v0/cat?arg=%s", cid), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+lh.apiKey)

	// Send request
	resp, err := lh.client.Do(req)
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

	plainBuf, err := crypto.Decrypt(key, cipherBuf)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt file: %v", err)
	}

	return plainBuf, nil
}

// GetFileInfo retrieves information about a file stored on Lighthouse
func (lh *LighthouseAPI) GetFileInfo(cid string, key []byte) (*FileInfo, error) {
	fileContent, err := lh.DownloadFile(cid, key)
	if err != nil {
		return nil, err
	}

	// Determine file type based on content
	contentType := http.DetectContentType(fileContent)

	fileInfo := &FileInfo{
		CID:      cid,
		Filename: fmt.Sprintf("file_%s", cid[:8]),
		Size:     int64(len(fileContent)),
		Type:     contentType,
	}

	return fileInfo, nil
}
