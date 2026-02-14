package purse

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func xdgStateHome() string {
	if v := os.Getenv("XDG_STATE_HOME"); v != "" {
		return v
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "state")
}

func writeMappingFile(dir string, mf MappingFile) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(mf, "", "  ")
	if err != nil {
		return err
	}

	data = append(data, '\n')

	return os.WriteFile(filepath.Join(dir, mf.Server+".json"), data, 0o644)
}

// WriteGlobal writes the mapping file to the global purse-first directory
// at $XDG_STATE_HOME/purse-first/{server}.json.
func WriteGlobal(mf MappingFile) error {
	dir := filepath.Join(xdgStateHome(), "purse-first")
	return writeMappingFile(dir, mf)
}

// WriteProject writes the mapping file to a project-local purse-first directory
// at {projectDir}/.purse-first/{server}.json.
func WriteProject(projectDir string, mf MappingFile) error {
	dir := filepath.Join(projectDir, ".purse-first")
	return writeMappingFile(dir, mf)
}
