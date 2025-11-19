package commenter

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/shoppingjaws/tf-refactoring-block-uncommenter/pkg/parser"
)

func TestCommentOutBlocks(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		blocks   []parser.BlockPosition
		expected string
	}{
		{
			name: "comment out moved block",
			content: `moved {
  from = aws_instance.old
  to   = aws_instance.new
}`,
			blocks: []parser.BlockPosition{
				{StartLine: 1, EndLine: 4, BlockType: "moved"},
			},
			expected: `# moved {
  # from = aws_instance.old
  # to   = aws_instance.new
# }`,
		},
		{
			name: "comment out with indentation preserved",
			content: `  moved {
    from = aws_instance.old
    to   = aws_instance.new
  }`,
			blocks: []parser.BlockPosition{
				{StartLine: 1, EndLine: 4, BlockType: "moved"},
			},
			expected: `  # moved {
    # from = aws_instance.old
    # to   = aws_instance.new
  # }`,
		},
		{
			name: "skip already commented lines",
			content: `moved {
# from = aws_instance.old
  to   = aws_instance.new
}`,
			blocks: []parser.BlockPosition{
				{StartLine: 1, EndLine: 4, BlockType: "moved"},
			},
			expected: `# moved {
# from = aws_instance.old
  # to   = aws_instance.new
# }`,
		},
		{
			name: "comment out multiple blocks",
			content: `moved {
  from = aws_instance.old
  to   = aws_instance.new
}

resource "aws_instance" "example" {
  ami = "ami-12345678"
}

import {
  to = aws_s3_bucket.example
  id = "my-bucket"
}`,
			blocks: []parser.BlockPosition{
				{StartLine: 1, EndLine: 4, BlockType: "moved"},
				{StartLine: 10, EndLine: 13, BlockType: "import"},
			},
			expected: `# moved {
  # from = aws_instance.old
  # to   = aws_instance.new
# }

resource "aws_instance" "example" {
  ami = "ami-12345678"
}

# import {
  # to = aws_s3_bucket.example
  # id = "my-bucket"
# }`,
		},
		{
			name: "empty blocks list",
			content: `moved {
  from = aws_instance.old
  to   = aws_instance.new
}`,
			blocks: []parser.BlockPosition{},
			expected: `moved {
  from = aws_instance.old
  to   = aws_instance.new
}`,
		},
		{
			name: "comment out removed block with nested structure",
			content: `removed {
  from = aws_security_group.unused
  lifecycle {
    destroy = false
  }
}`,
			blocks: []parser.BlockPosition{
				{StartLine: 1, EndLine: 6, BlockType: "removed"},
			},
			expected: `# removed {
  # from = aws_security_group.unused
  # lifecycle {
    # destroy = false
  # }
# }`,
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
			if err := CommentOutBlocks(tmpFile, tt.blocks); err != nil {
				t.Fatalf("CommentOutBlocks() error = %v", err)
			}

			// Read the result
			result, err := os.ReadFile(tmpFile)
			if err != nil {
				t.Fatalf("failed to read result file: %v", err)
			}

			// Compare results
			resultStr := strings.TrimSpace(string(result))
			expectedStr := strings.TrimSpace(tt.expected)
			if resultStr != expectedStr {
				t.Errorf("CommentOutBlocks() result:\n%s\n\nexpected:\n%s", resultStr, expectedStr)
			}
		})
	}
}

func TestCommentOutBlocksFileOperations(t *testing.T) {
	t.Run("nonexistent file", func(t *testing.T) {
		blocks := []parser.BlockPosition{
			{StartLine: 1, EndLine: 4, BlockType: "moved"},
		}
		err := CommentOutBlocks("/nonexistent/file.tf", blocks)
		if err == nil {
			t.Error("expected error for nonexistent file, got nil")
		}
	})

	t.Run("empty file", func(t *testing.T) {
		tmpDir := t.TempDir()
		tmpFile := filepath.Join(tmpDir, "empty.tf")
		if err := os.WriteFile(tmpFile, []byte(""), 0644); err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}

		blocks := []parser.BlockPosition{
			{StartLine: 1, EndLine: 4, BlockType: "moved"},
		}
		err := CommentOutBlocks(tmpFile, blocks)
		if err != nil {
			t.Errorf("CommentOutBlocks() on empty file error = %v", err)
		}
	})
}
