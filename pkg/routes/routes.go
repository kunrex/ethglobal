package routes

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"git-server/pkg/controllers"
	"git-server/pkg/server"
)

// SetupRoutes configures all routes for the application
func SetupRoutes(gitServer *server.InMemoryGitServer) *mux.Router {
	router := mux.NewRouter()

	// Create controllers
	repoController := controllers.NewRepoController(gitServer)
	
	// Git protocol endpoints
	router.HandleFunc("/{repo}/info/refs", gitServer.GitInfoRefsHandler).Methods("GET")
	router.HandleFunc("/{repo}/git-upload-pack", gitServer.GitUploadPackHandler).Methods("POST")
	router.HandleFunc("/{repo}/git-receive-pack", gitServer.GitReceivePackHandler).Methods("POST")

	// Repository management endpoints
	router.HandleFunc("/repos", repoController.CreateRepoHandler).Methods("POST")
	router.HandleFunc("/repos", repoController.ListReposHandler).Methods("GET")
	router.HandleFunc("/{repo}/files", repoController.AddFileHandler).Methods("POST")

	// Lighthouse (Filecoin) API endpoints
	lighthouseAPIKey := os.Getenv("LIGHTHOUSE_API_KEY")
	if lighthouseAPIKey == "" {
		log.Println("Warning: LIGHTHOUSE_API_KEY environment variable not set. Lighthouse endpoints will not work.")
		log.Println("Set LIGHTHOUSE_API_KEY to enable Filecoin storage functionality.")
	} else {
		log.Println("Lighthouse API key found. Setting up Filecoin storage endpoints...")
		lighthouseController := controllers.NewLighthouseController(lighthouseAPIKey)
		
		// Lighthouse routes
		router.HandleFunc("/lighthouse/upload", lighthouseController.UploadHandler).Methods("POST")
		router.HandleFunc("/lighthouse/download", lighthouseController.DownloadHandler).Methods("GET")
		router.HandleFunc("/lighthouse/upload-text", lighthouseController.UploadTextHandler).Methods("POST")
		router.HandleFunc("/lighthouse/file-info", lighthouseController.GetFileInfoHandler).Methods("GET")
		router.HandleFunc("/lighthouse/uploads", lighthouseController.ListUploadsHandler).Methods("GET")
		router.HandleFunc("/lighthouse/help", lighthouseController.HelpHandler).Methods("GET")
	}

	// Root endpoint with web interface
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, getHomePageHTML())
	}).Methods("GET")

	return router
}

// getHomePageHTML returns the HTML content for the home page
func getHomePageHTML() string {
	return `
<!DOCTYPE html>
<html>
<head>
    <title>In-Memory Git Server</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .endpoint { background: #f5f5f5; padding: 10px; margin: 10px 0; border-radius: 5px; }
        code { background: #e8e8e8; padding: 2px 4px; border-radius: 3px; }
        .section { margin: 20px 0; }
    </style>
</head>
<body>
    <h1>In-Memory Git Server</h1>
    <p>This is a Git server that stores repositories in memory using Go and the go-git library.</p>
    
    <div class="section">
        <h2>Repository Management Endpoints:</h2>
        <div class="endpoint">
            <strong>GET /repos</strong> - List all repositories
        </div>
        <div class="endpoint">
            <strong>POST /repos?name=repo-name</strong> - Create a new repository
        </div>
        <div class="endpoint">
            <strong>POST /{repo}/files</strong> - Add a file to a repository (form data: filename, content)
        </div>
    </div>
    
    <div class="section">
        <h2>Git Protocol Endpoints:</h2>
        <div class="endpoint">
            <strong>GET /{repo}/info/refs?service=git-upload-pack</strong> - Get refs for clone/fetch
        </div>
        <div class="endpoint">
            <strong>POST /{repo}/git-upload-pack</strong> - Handle clone/fetch operations
        </div>
        <div class="endpoint">
            <strong>POST /{repo}/git-receive-pack</strong> - Handle push operations
        </div>
    </div>
    
    <div class="section">
        <h2>Lighthouse (Filecoin) Endpoints:</h2>
        <div class="endpoint">
            <strong>POST /lighthouse/upload</strong> - Upload a file to Filecoin network
        </div>
        <div class="endpoint">
            <strong>GET /lighthouse/download?cid=...</strong> - Download a file from Filecoin network
        </div>
        <div class="endpoint">
            <strong>POST /lighthouse/upload-text</strong> - Upload text content to Filecoin network
        </div>
        <div class="endpoint">
            <strong>GET /lighthouse/file-info?cid=...</strong> - Get file information
        </div>
        <div class="endpoint">
            <strong>GET /lighthouse/help</strong> - Lighthouse API help
        </div>
    </div>
    
    <div class="section">
        <h2>Example Usage:</h2>
        <p><strong>Git Operations:</strong></p>
        <code>git clone http://localhost:8080/test-repo</code><br>
        <code>curl http://localhost:8080/repos</code><br>
        <code>curl -X POST "http://localhost:8080/repos?name=my-new-repo"</code>
        
        <p><strong>Filecoin Operations:</strong></p>
        <code>curl -X POST -F 'file=@example.txt' http://localhost:8080/lighthouse/upload</code><br>
        <code>curl 'http://localhost:8080/lighthouse/download?cid=QmXXX...' -o downloaded_file</code><br>
        <code>curl -X POST -H 'Content-Type: application/json' -d '{"content":"Hello World","filename":"hello.txt"}' http://localhost:8080/lighthouse/upload-text</code>
    </div>
</body>
</html>
	`
}
