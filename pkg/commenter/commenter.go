package commenter

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/shoppingjaws/tf-refactoring-block-uncommenter/pkg/parser"
)

// CommentOutBlocks comments out the specified blocks in a file
func CommentOutBlocks(filepath string, blocks []parser.BlockPosition) error {
	if len(blocks) == 0 {
		return nil
	}

	// Read all lines from the file
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", filepath, err)
	}

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	file.Close()

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file %s: %w", filepath, err)
	}

	// Comment out the specified blocks
	for _, block := range blocks {
		for i := block.StartLine - 1; i < block.EndLine && i < len(lines); i++ {
			line := lines[i]
			// Only add # if the line is not already commented
			if !parser.IsLineCommented(line) {
				// Preserve indentation
				trimmed := strings.TrimLeft(line, " \t")
				indent := line[:len(line)-len(trimmed)]
				lines[i] = indent + "# " + trimmed
			}
		}
	}

	// Write back to the file
	output, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filepath, err)
	}
	defer output.Close()

	writer := bufio.NewWriter(output)
	for i, line := range lines {
		if i > 0 {
			writer.WriteString("\n")
		}
		writer.WriteString(line)
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to write to file %s: %w", filepath, err)
	}

	return nil
}
