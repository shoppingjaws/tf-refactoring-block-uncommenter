package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindRefactoringBlocks(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected []BlockPosition
	}{
		{
			name: "detect moved block",
			content: `moved {
  from = aws_instance.old
  to   = aws_instance.new
}`,
			expected: []BlockPosition{
				{StartLine: 1, EndLine: 4, BlockType: "moved"},
			},
		},
		{
			name: "detect import block",
			content: `import {
  to = aws_s3_bucket.example
  id = "my-bucket"
}`,
			expected: []BlockPosition{
				{StartLine: 1, EndLine: 4, BlockType: "import"},
			},
		},
		{
			name: "detect removed block",
			content: `removed {
  from = aws_security_group.unused
  lifecycle {
    destroy = false
  }
}`,
			expected: []BlockPosition{
				{StartLine: 1, EndLine: 6, BlockType: "removed"},
			},
		},
		{
			name: "detect multiple blocks",
			content: `moved {
  from = aws_instance.old
  to   = aws_instance.new
}

import {
  to = aws_s3_bucket.example
  id = "my-bucket"
}

removed {
  from = aws_security_group.unused
}`,
			expected: []BlockPosition{
				{StartLine: 1, EndLine: 4, BlockType: "moved"},
				{StartLine: 6, EndLine: 9, BlockType: "import"},
				{StartLine: 11, EndLine: 13, BlockType: "removed"},
			},
		},
		{
			name: "ignore commented blocks",
			content: `# moved {
#   from = aws_instance.old
#   to   = aws_instance.new
# }

moved {
  from = aws_instance.active
  to   = aws_instance.current
}`,
			expected: []BlockPosition{
				{StartLine: 6, EndLine: 9, BlockType: "moved"},
			},
		},
		{
			name:    "detect single-line block",
			content: `moved { from = aws_instance.old, to = aws_instance.new }`,
			expected: []BlockPosition{
				{StartLine: 1, EndLine: 1, BlockType: "moved"},
			},
		},
		{
			name: "mixed with other terraform code",
			content: `resource "aws_instance" "example" {
  ami           = "ami-12345678"
  instance_type = "t2.micro"
}

moved {
  from = aws_instance.old
  to   = aws_instance.new
}

variable "environment" {
  type = string
}`,
			expected: []BlockPosition{
				{StartLine: 6, EndLine: 9, BlockType: "moved"},
			},
		},
		{
			name: "no refactoring blocks",
			content: `resource "aws_instance" "example" {
  ami           = "ami-12345678"
  instance_type = "t2.micro"
}`,
			expected: []BlockPosition{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary file
			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, "test.tf")
			if err := os.WriteFile(tmpFile, []byte(tt.content), 0644); err != nil {
				t.Fatalf("failed to create temp file: %v", err)
			}

			// Run the function
			blocks, err := FindRefactoringBlocks(tmpFile)
			if err != nil {
				t.Fatalf("FindRefactoringBlocks() error = %v", err)
			}

			// Compare results
			if len(blocks) != len(tt.expected) {
				t.Errorf("expected %d blocks, got %d", len(tt.expected), len(blocks))
				return
			}

			for i, block := range blocks {
				if block.StartLine != tt.expected[i].StartLine ||
					block.EndLine != tt.expected[i].EndLine ||
					block.BlockType != tt.expected[i].BlockType {
					t.Errorf("block %d: expected %+v, got %+v", i, tt.expected[i], block)
				}
			}
		})
	}
}

func TestIsLineCommented(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		expected bool
	}{
		{"commented line", "# moved {", true},
		{"commented with spaces", "  # moved {", true},
		{"commented with tabs", "\t# moved {", true},
		{"not commented", "moved {", false},
		{"not commented with spaces", "  moved {", false},
		{"empty line", "", false},
		{"only spaces", "   ", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsLineCommented(tt.line)
			if result != tt.expected {
				t.Errorf("IsLineCommented(%q) = %v, expected %v", tt.line, result, tt.expected)
			}
		})
	}
}
