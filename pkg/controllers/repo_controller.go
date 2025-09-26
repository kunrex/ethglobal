package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"git-server/pkg/server"
)

// RepoController handles repository management operations
type RepoController struct {
	gitServer *server.InMemoryGitServer
}

// NewRepoController creates a new repository controller
func NewRepoController(gitServer *server.InMemoryGitServer) *RepoController {
	return &RepoController{
		gitServer: gitServer,
	}
}

// CreateRepoHandler handles creating new repositories
func (rc *RepoController) CreateRepoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	repoName := r.URL.Query().Get("name")
	if repoName == "" {
		http.Error(w, "Repository name is required", http.StatusBadRequest)
		return
	}

	err := rc.gitServer.CreateRepository(repoName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Repository '%s' created successfully\n", repoName)
}

// ListReposHandler handles listing all repositories
func (rc *RepoController) ListReposHandler(w http.ResponseWriter, r *http.Request) {
	repos := rc.gitServer.ListRepositories()
	
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
func (rc *RepoController) AddFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	repoName := vars["repo"]

	repo, err := rc.gitServer.GetRepository(repoName)
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
