package git

import (
	"fmt"
	"os/exec"
	"strings"
)

// HasTerraformChanges checks if there are any *.tf file changes in the current diff
func HasTerraformChanges() (bool, error) {
	cmd := exec.Command("git", "diff", "--name-only", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("failed to run git diff: %w", err)
	}

	files := strings.Split(string(output), "\n")
	for _, file := range files {
		if strings.HasSuffix(strings.TrimSpace(file), ".tf") {
			return true, nil
		}
	}

	return false, nil
}

// GetTerraformFiles returns all *.tf files in the repository
func GetTerraformFiles() ([]string, error) {
	cmd := exec.Command("git", "ls-files", "*.tf")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list terraform files: %w", err)
	}

	files := strings.Split(strings.TrimSpace(string(output)), "\n")
	var result []string
	for _, file := range files {
		if file != "" {
			result = append(result, file)
		}
	}

	return result, nil
}
