package utils

import (
	"encoding/json"
	"git-server/pkg/types"
	"os"
	"os/user"
	"strings"
)

var filepath string

func init() {
	filepath = expandPath("~/.ccg/ccg.json")
}

func expandPath(path string) string {
	if strings.HasPrefix(path, "~") {
		usr, err := user.Current()
		if err != nil {
			return path
		}
		return usr.HomeDir + path[1:]
	}
	return path
}

// Save projects to disk after changes
func SaveProjectsToFile(projects map[string]*types.AnonymousWallet) error {
	data, err := json.MarshalIndent(projects, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath, data, 0600)
}

// Load projects from disk on startup
func LoadProjectsFromFile() (map[string]*types.AnonymousWallet, error) {
	var projects map[string]*types.AnonymousWallet
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &projects)
	return projects, err
}
