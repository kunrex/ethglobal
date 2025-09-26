package server

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/go-git/go-git/v5"
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
