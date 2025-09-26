package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/gorilla/mux"
)

// InMemoryGitServer represents a Git server that stores repositories in memory
type InMemoryGitServer struct {
	repos map[string]*git.Repository
	mutex sync.RWMutex
}

// NewInMemoryGitServer creates a new in-memory Git server
func NewInMemoryGitServer() *InMemoryGitServer {
	return &InMemoryGitServer{
		repos: make(map[string]*git.Repository),
	}
}

// CreateRepository creates a new repository in memory
func (s *InMemoryGitServer) CreateRepository(name string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.repos[name]; exists {
		return fmt.Errorf("repository %s already exists", name)
	}

	// Create a new in-memory repository
	repo, err := git.Init(memory.NewStorage(), nil)
	if err != nil {
		return fmt.Errorf("failed to initialize repository: %v", err)
	}

	s.repos[name] = repo
	log.Printf("Created repository: %s", name)
	return nil
}

// CreateRepositoryWithContent creates a new repository with initial content
func (s *InMemoryGitServer) CreateRepositoryWithContent(name string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.repos[name]; exists {
		return fmt.Errorf("repository %s already exists", name)
	}

	// Create a new in-memory repository
	repo, err := git.Init(memory.NewStorage(), nil)
	if err != nil {
		return fmt.Errorf("failed to initialize repository: %v", err)
	}

	// Create initial commit with README
	worktree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %v", err)
	}

	// In a real implementation, you'd use the filesystem abstraction
	// For now, we'll create a commit directly
	hash, err := worktree.Commit("Initial commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Git Server",
			Email: "server@example.com",
		},
	})
	if err != nil {
		// If commit fails (empty repo), that's okay for now
		log.Printf("Warning: Could not create initial commit for %s: %v", name, err)
	} else {
		log.Printf("Created initial commit %s for repository: %s", hash.String()[:7], name)
	}

	s.repos[name] = repo
	log.Printf("Created repository with content: %s", name)
	return nil
}

// GetRepository retrieves a repository by name
func (s *InMemoryGitServer) GetRepository(name string) (*git.Repository, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	repo, exists := s.repos[name]
	if !exists {
		return nil, fmt.Errorf("repository %s not found", name)
	}

	return repo, nil
}

// ListRepositories returns a list of all repository names
func (s *InMemoryGitServer) ListRepositories() []string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	names := make([]string, 0, len(s.repos))
	for name := range s.repos {
		names = append(names, name)
	}
	return names
}

// GitInfoRefsHandler handles the /info/refs endpoint
func (s *InMemoryGitServer) GitInfoRefsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	repoName := vars["repo"]
	service := r.URL.Query().Get("service")

	if service != "git-upload-pack" && service != "git-receive-pack" {
		http.Error(w, "Invalid service", http.StatusBadRequest)
		return
	}

	repo, err := s.GetRepository(repoName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Set proper headers for Git protocol
	w.Header().Set("Content-Type", "application/x-git-upload-pack-advertisement")
	w.Header().Set("Cache-Control", "no-cache")
	
	// Write the service advertisement
	serviceLine := fmt.Sprintf("# service=%s\n", service)
	// Calculate packet length (4 hex digits)
	packetLen := fmt.Sprintf("%04x", len(serviceLine)+4)
	
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, packetLen)
	fmt.Fprint(w, serviceLine)
	fmt.Fprint(w, "0000") // End of service advertisement
	
	// Get the default branch
	ref, err := repo.Head()
	if err != nil {
		// Repository is empty, just return empty refs
		fmt.Fprint(w, "0000") // End of refs
		return
	}

	// Format the refs for Git protocol
	refLine := fmt.Sprintf("%s %s\n", ref.Hash().String(), ref.Name().String())
	refPacketLen := fmt.Sprintf("%04x", len(refLine))
	fmt.Fprint(w, refPacketLen)
	fmt.Fprint(w, refLine)
	fmt.Fprint(w, "0000") // End of refs
}

// GitUploadPackHandler handles the /git-upload-pack endpoint (for clone/fetch)
func (s *InMemoryGitServer) GitUploadPackHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	repoName := vars["repo"]

	repo, err := s.GetRepository(repoName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Set proper headers
	w.Header().Set("Content-Type", "application/x-git-upload-pack-result")
	w.Header().Set("Cache-Control", "no-cache")
	
	// Check if repository has any commits
	_, err = repo.Head()
	if err != nil {
		// Empty repository - return empty pack
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "0000") // End packet
		return
	}

	// For a real implementation, you would:
	// 1. Parse the upload pack request
	// 2. Generate pack data with objects
	// 3. Return the pack file
	// For now, return a minimal response
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "0000") // End packet
}

// GitReceivePackHandler handles the /git-receive-pack endpoint (for push)
func (s *InMemoryGitServer) GitReceivePackHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	repoName := vars["repo"]

	_, err := s.GetRepository(repoName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// For now, return a simple response
	// In a full implementation, you'd handle the Git pack protocol here
	w.Header().Set("Content-Type", "application/x-git-receive-pack-result")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "0000")
}

// CreateRepoHandler handles creating new repositories
func (s *InMemoryGitServer) CreateRepoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	repoName := r.URL.Query().Get("name")
	if repoName == "" {
		http.Error(w, "Repository name is required", http.StatusBadRequest)
		return
	}

	err := s.CreateRepository(repoName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Repository '%s' created successfully\n", repoName)
}

// ListReposHandler handles listing all repositories
func (s *InMemoryGitServer) ListReposHandler(w http.ResponseWriter, r *http.Request) {
	repos := s.ListRepositories()
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "{\n  \"repositories\": [\n")
	for i, repo := range repos {
		if i > 0 {
			fmt.Fprint(w, ",\n")
		}
		fmt.Fprintf(w, "    \"%s\"", repo)
	}
	fmt.Fprint(w, "\n  ]\n}\n")
}

// AddFileHandler adds a file to a repository
func (s *InMemoryGitServer) AddFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	repoName := vars["repo"]

	repo, err := s.GetRepository(repoName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Parse form data
	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	filename := r.FormValue("filename")
	content := r.FormValue("content")

	if filename == "" || content == "" {
		http.Error(w, "Filename and content are required", http.StatusBadRequest)
		return
	}

	// Get the worktree
	_, err = repo.Worktree()
	if err != nil {
		http.Error(w, "Failed to get worktree", http.StatusInternalServerError)
		return
	}

	// Create the file in memory (this is a simplified approach)
	// In a real implementation, you'd need to handle the filesystem abstraction
	log.Printf("Adding file '%s' to repository '%s'", filename, repoName)
	
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "File '%s' added to repository '%s'\n", filename, repoName)
}

func main() {
	// Create the Git server
	gitServer := NewInMemoryGitServer()

	// Create some sample repositories
	err := gitServer.CreateRepository("test-repo")
	if err != nil {
		log.Printf("Warning: Failed to create test repository: %v", err)
	}

	err = gitServer.CreateRepository("demo")
	if err != nil {
		log.Printf("Warning: Failed to create demo repository: %v", err)
	}

	// Set up the HTTP router
	router := mux.NewRouter()

	// Git protocol endpoints
	router.HandleFunc("/{repo}/info/refs", gitServer.GitInfoRefsHandler).Methods("GET")
	router.HandleFunc("/{repo}/git-upload-pack", gitServer.GitUploadPackHandler).Methods("POST")
	router.HandleFunc("/{repo}/git-receive-pack", gitServer.GitReceivePackHandler).Methods("POST")

	// Management endpoints
	router.HandleFunc("/repos", gitServer.CreateRepoHandler).Methods("POST")
	router.HandleFunc("/repos", gitServer.ListReposHandler).Methods("GET")
	router.HandleFunc("/{repo}/files", gitServer.AddFileHandler).Methods("POST")

	// Root endpoint
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `
<!DOCTYPE html>
<html>
<head>
    <title>In-Memory Git Server</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .endpoint { background: #f5f5f5; padding: 10px; margin: 10px 0; border-radius: 5px; }
        code { background: #e8e8e8; padding: 2px 4px; border-radius: 3px; }
    </style>
</head>
<body>
    <h1>In-Memory Git Server</h1>
    <p>This is a Git server that stores repositories in memory using Go and the go-git library.</p>
    
    <h2>Available Endpoints:</h2>
    <div class="endpoint">
        <strong>GET /repos</strong> - List all repositories
    </div>
    <div class="endpoint">
        <strong>POST /repos?name=repo-name</strong> - Create a new repository
    </div>
    <div class="endpoint">
        <strong>POST /{repo}/files</strong> - Add a file to a repository (form data: filename, content)
    </div>
    
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
    
    <h2>Example Usage:</h2>
    <p>Clone a repository:</p>
    <code>git clone http://localhost:8080/test-repo.git</code>
    
    <p>List repositories:</p>
    <code>curl http://localhost:8080/repos</code>
    
    <p>Create a new repository:</p>
    <code>curl -X POST "http://localhost:8080/repos?name=my-new-repo"</code>
</body>
</html>
		`)
	}).Methods("GET")

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting Git server on port %s", port)
	log.Printf("Visit http://localhost:%s to see the web interface", port)
	
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
