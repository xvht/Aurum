package gitinfo

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// FindFile searches for a file starting from current directory up to root
func FindFile(filename string) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, filename)); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New("file not found")
		}
		dir = parent
	}
}

// ReadGitFile reads a file from the git repository root
func ReadGitFile(filename string) (string, error) {
	gitRoot, err := FindFile(".git")
	if err != nil {
		return "", errors.New("no git repository root found")
	}

	content, err := os.ReadFile(filepath.Join(gitRoot, filename))
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// GetCommit returns the latest commit hash
func GetCommit() (string, error) {
	content, err := ReadGitFile(".git/logs/HEAD")
	if err != nil {
		return "", err
	}

	lines := strings.Split(strings.TrimSpace(content), "\n")
	if len(lines) == 0 {
		return "", errors.New("no commit history found")
	}

	lastLine := lines[len(lines)-1]
	parts := strings.Split(lastLine, " ")
	if len(parts) < 2 {
		return "", errors.New("invalid commit log format")
	}

	return parts[1], nil
}

// GetBranch returns the current git branch
func GetBranch() (string, error) {
	if branch := os.Getenv("CF_PAGES_BRANCH"); branch != "" {
		return branch, nil
	}

	content, err := ReadGitFile(".git/HEAD")
	if err != nil {
		return "", err
	}

	branch := strings.TrimSpace(content)
	branch = strings.TrimPrefix(branch, "ref: refs/heads/")
	return branch, nil
}

// GetRemote returns the git remote repository name
func GetRemote() (string, error) {
	content, err := ReadGitFile(".git/config")
	if err != nil {
		return "", err
	}

	lines := strings.Split(content, "\n")
	var remote string
	for _, line := range lines {
		if strings.Contains(line, "url = ") {
			remote = strings.Split(line, "url = ")[1]
			break
		}
	}

	if remote == "" {
		return "", errors.New("could not parse remote")
	}

	if strings.HasPrefix(remote, "git@") {
		remote = strings.Split(remote, ":")[1]
	} else if strings.HasPrefix(remote, "http") {
		// Parse URL path
		parts := strings.Split(remote, "/")
		if len(parts) < 2 {
			return "", errors.New("invalid remote URL format")
		}
		remote = strings.Join(parts[len(parts)-2:], "/")
	}

	remote = strings.TrimSuffix(remote, ".git")
	return remote, nil
}
