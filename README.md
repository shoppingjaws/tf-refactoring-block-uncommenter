# Terraform Refactoring Block Uncommenter

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![GitHub release](https://img.shields.io/github/v/release/shoppingjaws/tf-refactoring-block-uncommenter?include_prereleases)](https://github.com/shoppingjaws/tf-refactoring-block-uncommenter/releases)
[![Go Version](https://img.shields.io/github/go-mod/go-version/shoppingjaws/tf-refactoring-block-uncommenter)](go.mod)

A GitHub Actions composite action that automatically comments out executed Terraform refactoring blocks (`moved`, `import`, `removed`) and creates a pull request.

## 🎯 Purpose

When refactoring Terraform code, you use special blocks like `moved`, `import`, and `removed` to safely migrate resources. After these blocks are executed:

- ✅ They should be commented out to prevent re-execution
- ✅ They should be kept as execution history
- ✅ This helps maintain clean Terraform state management

This action automates this tedious manual process by:
1. Detecting when `.tf` files are changed
2. Finding uncommented refactoring blocks
3. Commenting them out with `#`
4. Creating a pull request with the changes

## 📦 Features

- 🔍 Detects `moved`, `import`, and `removed` blocks in Terraform files
- 📝 Comments out blocks with proper indentation preservation
- 🚀 Automatically creates pull requests via GitHub CLI
- ⚙️ Fully configurable inputs
- 🎨 Written in Go for fast, reliable parsing

## 🚀 Usage

### Version Selection

Choose the version reference based on your needs:

| Reference | Use Case | Stability |
|-----------|----------|-----------|
| `@main` | Latest features, active development | 🔄 Frequently updated |
| `@v1` | Stable release, production use | ✅ Stable, tagged version |

**Recommendation:**
- 🚀 Use `@main` for: Development, testing new features, quick iterations
- 🛡️ Use `@v1` for: Production, CI/CD pipelines, stable environments

### Basic Usage

Create a workflow file (e.g., `.github/workflows/terraform-uncommenter.yaml`):

```yaml
name: Comment Out Terraform Refactoring Blocks

on:
  push:
    branches:
      - main
    paths:
      - '**.tf'

permissions:
  contents: write
  pull-requests: write

jobs:
  comment-out-blocks:
    runs-on: ubuntu-latest
    steps:
      # Use @main for latest features (recommended for development)
      - uses: shoppingjaws/tf-refactoring-block-uncommenter@main
        with:
          github-token: ${{ secrets.GITHUB_TOKEN }}

      # Or use @v1 for stable release (recommended for production)
      # - uses: shoppingjaws/tf-refactoring-block-uncommenter@v1
      #   with:
      #     github-token: ${{ secrets.GITHUB_TOKEN }}
```

### Advanced Usage

Customize the action with inputs:

```yaml
- uses: shoppingjaws/tf-refactoring-block-uncommenter@main
  with:
    github-token: ${{ secrets.GITHUB_TOKEN }}
    branch-name: 'auto/comment-tf-blocks'
    pr-title: 'chore(terraform): Comment out executed refactoring blocks'
    pr-body: |
      This PR comments out Terraform refactoring blocks that have been executed.

      Please review and merge.
    base-branch: 'main'
```

## 🔧 Inputs

| Input | Description | Required | Default |
|-------|-------------|----------|---------|
| `github-token` | GitHub token for creating PRs | No | `${{ github.token }}` |
| `branch-name` | Branch name for the PR | No | `chore/comment-out-tf-blocks` |
| `pr-title` | Title for the pull request | No | `chore: Comment out executed Terraform refactoring blocks` |
| `pr-body` | Body for the pull request | No | Auto-generated message |
| `base-branch` | Base branch for the PR | No | Repository default branch |

## 📋 Example Terraform Blocks

### Before (Uncommented)

```hcl
moved {
  from = aws_instance.old_name
  to   = aws_instance.new_name
}

import {
  to = aws_s3_bucket.example
  id = "my-bucket"
}

removed {
  from = aws_security_group.unused
  lifecycle {
    destroy = false
  }
}
```

### After (Commented Out)

```hcl
# moved {
#   from = aws_instance.old_name
#   to   = aws_instance.new_name
# }

# import {
#   to = aws_s3_bucket.example
#   id = "my-bucket"
# }

# removed {
#   from = aws_security_group.unused
#   lifecycle {
#     destroy = false
#   }
# }
```

## 🛠️ How It Works

1. **Detection**: Monitors push events for `.tf` file changes
2. **Scanning**: Scans all Terraform files for uncommented refactoring blocks
3. **Processing**: Comments out each block line-by-line with `#` while preserving indentation
4. **PR Creation**: Creates a pull request with the changes using GitHub CLI

## 📚 Requirements

- GitHub Actions runner (ubuntu-latest recommended)
- Go 1.21+ (automatically set up by the action)
- `contents: write` and `pull-requests: write` permissions

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🔗 Related

- [Terraform Refactoring Documentation](https://developer.hashicorp.com/terraform/language/modules/develop/refactoring)
- [Terraform Import](https://developer.hashicorp.com/terraform/language/import)
- [Terraform Removed](https://developer.hashicorp.com/terraform/language/resources/syntax#removing-resources)

## ⚠️ Important Notes

- This action only comments out blocks; it does not delete them
- Always review the generated PR before merging
- The action requires write permissions for contents and pull requests
- Blocks are detected based on simple pattern matching (not full HCL parsing)

## 💡 Tips

- Run this action after your Terraform apply succeeds
- Consider adding this as a separate job that depends on your Terraform workflow
- Customize the PR title and body to match your project's conventions
- Use protected branches to ensure PRs are reviewed before merging
