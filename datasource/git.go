package datasource

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/zupzup/blog-generator/config"
)

// GitDataSource is the git data source object
type GitDataSource struct{}

// Fetch creates the output folder, clears it and clones the repository there
func (ds *GitDataSource) Fetch(cfg *config.Config) ([]string, error) {
	from := cfg.Generator.Repo
	to := cfg.Generator.Tmp
	branch := cfg.Generator.Branch
	if branch == "" {
		branch = "master"
	}

	fmt.Printf("Fetching data from %s into %s...\n", from, to)
	if err := createFolderIfNotExist(to); err != nil {
		return nil, err
	}
	if err := clearFolder(to); err != nil {
		return nil, err
	}
	if err := cloneRepo(to, from, branch); err != nil {
		return nil, err
	}
	dirs, err := getContentFolders(to)
	if err != nil {
		return nil, err
	}
	fmt.Print("Fetching complete.\n")
	return dirs, nil
}

func createFolderIfNotExist(path string) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			if err = os.Mkdir(path, os.ModePerm); err != nil {
				return fmt.Errorf("error creating directory %s: %v", path, err)
			}
		} else {
			return fmt.Errorf("error accessing directory %s: %v", path, err)
		}
	}
	return nil
}

func clearFolder(path string) error {
	dir, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error accessing directory %s: %v", path, err)
	}
	defer dir.Close()
	names, err := dir.Readdirnames(-1)
	if err != nil {
		return fmt.Errorf("error reading directory %s: %v", path, err)
	}

	for _, name := range names {
		if err = os.RemoveAll(filepath.Join(path, name)); err != nil {
			return fmt.Errorf("error clearing file %s: %v", name, err)
		}
	}
	return nil
}

func cloneRepo(path, repositoryURL, branch string) error {
	cmdName := "git"
	initArgs := []string{"init", "."}
	cmd := exec.Command(cmdName, initArgs...)
	cmd.Dir = path
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error initializing git repository at %s: %v", path, err)
	}
	remoteArgs := []string{"remote", "add", "origin", repositoryURL}
	cmd = exec.Command(cmdName, remoteArgs...)
	cmd.Dir = path
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error setting remote %s: %v", repositoryURL, err)
	}
	pullArgs := []string{"pull", "origin", branch}
	cmd = exec.Command(cmdName, pullArgs...)
	cmd.Dir = path
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error pulling %s at %s: %v", branch, path, err)
	}
	return nil
}

func getContentFolders(path string) ([]string, error) {
	var result []string
	dir, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error accessing directory %s: %v", path, err)
	}
	defer dir.Close()
	files, err := dir.Readdir(-1)
	if err != nil {
		return nil, fmt.Errorf("error reading contents of directory %s: %v", path, err)
	}
	for _, file := range files {
		if file.IsDir() && file.Name()[0] != '.' {
			result = append(result, filepath.Join(path, file.Name()))
		}
	}
	return result, nil
}
