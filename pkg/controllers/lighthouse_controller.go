package controllers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"git-server/pkg/api"
	"git-server/pkg/contracts"
	"git-server/pkg/main"
	"git-server/pkg/utils"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

// LighthouseController handles Lighthouse (Filecoin) API operations
type LighthouseController struct {
	lighthouseAPI *api.LighthouseAPI
}

// NewLighthouseController creates a new Lighthouse controller
func NewLighthouseController(apiKey string) *LighthouseController {
	return &LighthouseController{
		lighthouseAPI: api.NewLighthouseAPI(apiKey),
	}
}

func getKeyFromEnv() ([]byte, error) {
	keyStr := os.Getenv("CCG_SECRET_KEY")
	if keyStr == "" {
		return nil, fmt.Errorf("CCG_SECRET_KEY not set")
	}

	key, err := base64.StdEncoding.DecodeString(keyStr)
	if err != nil {
		return nil, fmt.Errorf("invalid base64 key: %w", err)
	}

	if l := len(key); l != 32 {
		return nil, fmt.Errorf("invalid AES key size: %d bytes", l)
	}

	return key, nil
}

// UploadHandler handles file uploads to Lighthouse (Filecoin)
func (lc *LighthouseController) UploadHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	repoName := vars["repo"]
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var err error
	projectWallet := main.ProjectWallets[repoName]
	if projectWallet == nil {
		projectWallet, err = contracts.CreateWallet() // Used ONLY for setting up the contract if needed
		if err != nil {
			log.Fatalf("Error creating master wallet: %v", err)
		}
		main.ProjectWallets[repoName] = projectWallet
		if err := utils.SaveProjectsToFile(main.ProjectWallets); err != nil {
			log.Fatalf("Error saving master wallet: %v", err)
		}
	}
	// Parse multipart form
	err = r.ParseMultipartForm(32 << 20) // 32 MB max
	if err != nil {
		http.Error(w, "Failed to parse multipart form/Size Exceeded", http.StatusBadRequest)
		return
	}

	// Get the file from form data
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "No file provided", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Get optional parameters
	fileName := r.FormValue("filename")
	if fileName == "" {
		fileName = header.Filename
	}

	key, err := getKeyFromEnv()

	if err != nil {
		http.Error(w, "Failed to get key from env", http.StatusBadRequest)
	}

	// Upload to Lighthouse
	log.Printf("Uploading file '%s' to Lighthouse...", fileName)
	uploadResponse, err := lc.lighthouseAPI.UploadFile(file, fileName, key)
	if err != nil {
		log.Printf("Failed to upload file: %v", err)
		response := api.UploadResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to upload file: %v", err),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	uploadResponse.Size = header.Size
	log.Printf("File uploaded successfully. CID: %s", uploadResponse.CID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(uploadResponse)
}

// DownloadHandler handles file downloads from Lighthouse (Filecoin)
func (lc *LighthouseController) DownloadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get CID from query parameters
	cid := r.URL.Query().Get("cid")
	if cid == "" {
		http.Error(w, "CID parameter is required", http.StatusBadRequest)
		return
	}

	// Optional filename parameter
	filename := r.URL.Query().Get("filename")

	log.Printf("Downloading file with CID: %s", cid)

	key, err := getKeyFromEnv()

	if err != nil {
		http.Error(w, "Failed to get key from env", http.StatusBadRequest)
	}

	// Download from Lighthouse
	fileContent, err := lc.lighthouseAPI.DownloadFile(cid, key)
	if err != nil {
		log.Printf("Failed to download file: %v", err)
		response := api.DownloadResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to download file: %v", err),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Set appropriate headers
	if filename == "" {
		filename = fmt.Sprintf("file_%s", cid[:8])
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(fileContent)))

	// Write file content to response
	_, err = w.Write(fileContent)
	if err != nil {
		log.Printf("Failed to write file content: %v", err)
		return
	}

	log.Printf("File downloaded successfully. Size: %d bytes", len(fileContent))
}

// UploadTextHandler handles text content uploads to Lighthouse
func (lc *LighthouseController) UploadTextHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse JSON body
	var request struct {
		Content  string `json:"content"`
		Filename string `json:"filename"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if request.Content == "" {
		http.Error(w, "Content is required", http.StatusBadRequest)
		return
	}

	if request.Filename == "" {
		request.Filename = "text_content.txt"
	}

	log.Printf("Uploading text content to Lighthouse...")

	// Create a reader from the text content
	reader := strings.NewReader(request.Content)

	key, err := getKeyFromEnv()

	if err != nil {
		http.Error(w, "Failed to get key from env", http.StatusBadRequest)
	}

	// Upload to Lighthouse
	uploadResponse, err := lc.lighthouseAPI.UploadFile(reader, request.Filename, key)
	if err != nil {
		log.Printf("Failed to upload text content: %v", err)
		response := api.UploadResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to upload text content: %v", err),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	uploadResponse.Size = int64(len(request.Content))
	log.Printf("Text content uploaded successfully. CID: %s", uploadResponse.CID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(uploadResponse)
}

// GetFileInfoHandler retrieves information about a file stored on Lighthouse
func (lc *LighthouseController) GetFileInfoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get CID from query parameters
	cid := r.URL.Query().Get("cid")
	if cid == "" {
		http.Error(w, "CID parameter is required", http.StatusBadRequest)
		return
	}

	log.Printf("Getting file info for CID: %s", cid)

	key, err := getKeyFromEnv()

	if err != nil {
		http.Error(w, "Failed to get key from env", http.StatusBadRequest)
	}

	// Get file info from Lighthouse
	fileInfo, err := lc.lighthouseAPI.GetFileInfo(cid, key)
	if err != nil {
		log.Printf("Failed to get file info: %v", err)
		response := map[string]interface{}{
			"success": false,
			"message": fmt.Sprintf("Failed to get file info: %v", err),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fileInfo)
}

// ListUploadsHandler lists recent uploads (placeholder)
func (lc *LighthouseController) ListUploadsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// For now, return a placeholder response
	// In a real implementation, you'd query Lighthouse for user's uploads
	response := map[string]interface{}{
		"success": true,
		"message": "List uploads functionality would require additional Lighthouse API integration",
		"uploads": []api.FileInfo{},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// HelpHandler provides help information for Lighthouse endpoints
func (lc *LighthouseController) HelpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	help := map[string]interface{}{
		"endpoints": map[string]string{
			"POST /lighthouse/upload":      "Upload a file to Filecoin network",
			"GET /lighthouse/download":     "Download a file from Filecoin network (requires ?cid=...)",
			"POST /lighthouse/upload-text": "Upload text content to Filecoin network",
			"GET /lighthouse/file-info":    "Get information about a file (requires ?cid=...)",
			"GET /lighthouse/uploads":      "List recent uploads",
			"GET /lighthouse/help":         "Show this help message",
		},
		"examples": map[string]string{
			"upload_file":   "curl -X POST -F 'file=@example.txt' http://localhost:8080/lighthouse/upload",
			"download_file": "curl 'http://localhost:8080/lighthouse/download?cid=QmXXX...' -o downloaded_file",
			"upload_text":   "curl -X POST -H 'Content-Type: application/json' -d '{\"content\":\"Hello World\",\"filename\":\"hello.txt\"}' http://localhost:8080/lighthouse/upload-text",
			"get_file_info": "curl 'http://localhost:8080/lighthouse/file-info?cid=QmXXX...'",
		},
	}
	json.NewEncoder(w).Encode(help)
}
