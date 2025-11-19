// Package parser provides functionality to parse Terraform files and find refactoring blocks.
package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// BlockPosition represents the position of a refactoring block in a file
type BlockPosition struct {
	StartLine int
	EndLine   int
	BlockType string // "moved", "import", or "removed"
}

// FindRefactoringBlocks finds all uncommented moved/import/removed blocks in a Terraform file
func FindRefactoringBlocks(filepath string) ([]BlockPosition, error) {
	file, err := os.Open(filepath) //nolint:gosec // filepath is a required parameter
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filepath, err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("failed to close file %s: %w", filepath, closeErr)
		}
	}()

	var blocks []BlockPosition
	scanner := bufio.NewScanner(file)
	lineNum := 0
	var currentBlock *BlockPosition
	braceDepth := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Check for block start (moved, import, or removed)
		if currentBlock == nil {
			switch {
			case strings.HasPrefix(line, "moved {"):
				currentBlock = &BlockPosition{
					StartLine: lineNum,
					BlockType: "moved",
				}
				braceDepth = 1
			case strings.HasPrefix(line, "import {"):
				currentBlock = &BlockPosition{
					StartLine: lineNum,
					BlockType: "import",
				}
				braceDepth = 1
			case strings.HasPrefix(line, "removed {"):
				currentBlock = &BlockPosition{
					StartLine: lineNum,
					BlockType: "removed",
				}
				braceDepth = 1
			}

			// Check for single-line blocks
			if currentBlock != nil && strings.HasSuffix(line, "}") && strings.Count(line, "{") == strings.Count(line, "}") {
				currentBlock.EndLine = lineNum
				blocks = append(blocks, *currentBlock)
				currentBlock = nil
				braceDepth = 0
			}
		} else {
			// Inside a block, track brace depth to handle nested structures
			openBraces := strings.Count(line, "{")
			closeBraces := strings.Count(line, "}")
			braceDepth += openBraces - closeBraces

			// Block ends when brace depth returns to 0
			if braceDepth == 0 {
				currentBlock.EndLine = lineNum
				blocks = append(blocks, *currentBlock)
				currentBlock = nil
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file %s: %w", filepath, err)
	}

	return blocks, nil
}

// IsLineCommented checks if a line is already commented out
func IsLineCommented(line string) bool {
	trimmed := strings.TrimSpace(line)
	return strings.HasPrefix(trimmed, "#")
}
