package main

import (
	"fmt"
	"os"

	"github.com/shoppingjaws/tf-refactoring-block-uncommenter/pkg/commenter"
	"github.com/shoppingjaws/tf-refactoring-block-uncommenter/pkg/git"
	"github.com/shoppingjaws/tf-refactoring-block-uncommenter/pkg/parser"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	fmt.Println("🔍 Scanning for Terraform files...")

	// Get all .tf files
	files, err := git.GetTerraformFiles()
	if err != nil {
		return fmt.Errorf("failed to get terraform files: %w", err)
	}

	fmt.Printf("📄 Found %d Terraform file(s)\n", len(files))

	totalBlocks := 0
	totalFiles := 0

	// Process each file
	for _, file := range files {
		blocks, err := parser.FindRefactoringBlocks(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "⚠️  Warning: failed to parse %s: %v\n", file, err)
			continue
		}

		if len(blocks) == 0 {
			continue
		}

		fmt.Printf("\n📝 Processing %s...\n", file)
		for _, block := range blocks {
			fmt.Printf("   Found %s block at lines %d-%d\n", block.BlockType, block.StartLine, block.EndLine)
		}

		// Comment out the blocks
		if err := commenter.CommentOutBlocks(file, blocks); err != nil {
			return fmt.Errorf("failed to comment out blocks in %s: %w", file, err)
		}

		totalBlocks += len(blocks)
		totalFiles++
	}

	if totalBlocks == 0 {
		fmt.Println("\nℹ️  No uncommented refactoring blocks found. Nothing to do.")
		os.Exit(1) // Exit with non-zero to indicate no changes
	}

	fmt.Printf("\n✨ Successfully commented out %d block(s) in %d file(s)\n", totalBlocks, totalFiles)
	return nil
}
