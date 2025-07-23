package repository

import (
	"fmt"
	"os"
	"os/exec"
)

type GitRepository struct{}

func NewGitRepository() *GitRepository {
	return &GitRepository{}
}

func (g *GitRepository) Clone(repoURL, branch string) (string, error) {
	tempDir, err := os.MkdirTemp("", "repo-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp dir: %w", err)
	}

	args := []string{"clone", "--depth=1", repoURL, tempDir}
	if branch != "" {
		args = []string{"clone", "--depth=1", "--branch", branch, repoURL, tempDir}
	}

	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		os.RemoveAll(tempDir)
		return "", fmt.Errorf("git clone failed: %w", err)
	}

	return tempDir, nil
}
