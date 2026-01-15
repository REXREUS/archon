package utils

import (
	"fmt"
	"os/exec"
	"strings"
)

func GetStagedDiff() (string, error) {
	if !IsGitRepo() {
		return "", fmt.Errorf("this directory is not a git repository")
	}
	cmd := exec.Command("git", "diff", "--cached")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git diff --cached failed: %s (%v)", strings.TrimSpace(string(out)), err)
	}
	return string(out), nil
}

func GetUnstagedDiff() (string, error) {
	if !IsGitRepo() {
		return "", fmt.Errorf("this directory is not a git repository")
	}
	cmd := exec.Command("git", "diff")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git diff failed: %s (%v)", strings.TrimSpace(string(out)), err)
	}
	return string(out), nil
}

func IsGitRepo() bool {
	_, err := exec.LookPath("git")
	if err != nil {
		return false
	}
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	err = cmd.Run()
	return err == nil
}

func GetCommitMessage(diff string) (string, error) {
	// This is just a helper, main logic remains in AI
	if strings.TrimSpace(diff) == "" {
		return "", nil
	}
	return "", nil
}
