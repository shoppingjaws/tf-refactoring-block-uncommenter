# CLAUDE.md - Terraform Refactoring Block Uncommenter

このファイルは、Claude Code がこのプロジェクトで作業する際のガイドラインを定義します。

## 📋 プロジェクト概要

**Terraform Refactoring Block Uncommenter** は、Terraform のリファクタリングブロック（`moved`, `import`, `removed`）を実行後に自動的にコメントアウトする GitHub Actions composite action です。

- **言語**: Go 1.21+
- **配布形式**: GitHub Actions composite action
- **リポジトリ**: https://github.com/shoppingjaws/tf-refactoring-block-uncommenter

## 🏗️ ディレクトリ構造

```
.
├── action.yaml              # Composite action 定義
├── main.go                  # エントリーポイント
├── pkg/                     # 公開パッケージ（internal/ は使用不可）
│   ├── git/                # Git 操作
│   ├── parser/             # HCL パーサー
│   └── commenter/          # コメントアウト処理
├── test.tf                  # 動作確認用テストファイル
└── .github/workflows/       # ワークフロー定義
```

### ⚠️ 重要: `internal/` ディレクトリは使用禁止

このプロジェクトは composite action として配布されるため、Go の `internal/` パッケージ可視性制約により、外部から参照できません。**必ず `pkg/` ディレクトリを使用してください。**

## 🔧 開発ガイドライン

### 1. テスト実行

**必ず Docker を使用してテストを実行してください：**

```bash
# 全テスト実行
docker run --rm -v .:/app -w /app golang:1.25.4 go test -v ./pkg/...

# 特定パッケージのテスト
docker run --rm -v .:/app -w /app golang:1.25.4 go test -v ./pkg/parser

# ビルドテスト
docker run --rm -v .:/app -w /app golang:1.25.4 go build -o /tmp/uncommenter main.go
```

### 2. コード変更時の注意点

- **パッケージ追加**: 必ず `pkg/` 配下に作成
- **インポートパス**: `github.com/shoppingjaws/tf-refactoring-block-uncommenter/pkg/...`
- **テスト追加**: すべての新機能にテストを追加
- **エラーハンドリング**: 適切なエラーメッセージを提供

### 3. HCL パース処理

- ネストされたブロックに対応（括弧の深さを追跡）
- コメントアウト済みブロックはスキップ
- インデントを保持してコメントアウト

### 4. Git 操作

- **重要**: `git diff` に依存しない設計
- 常に全 `.tf` ファイルをスキャン
- GitHub Actions 環境でも動作することを確認

## 📦 GitHub Actions 固有の注意点

### Composite Action として

- `action.yaml` で定義されたステップは別リポジトリから実行される
- `internal/` パッケージにアクセスできない
- 環境変数とワークスペースパスに注意

### バージョン参照

```yaml
# 開発版・最新機能（推奨）
uses: shoppingjaws/tf-refactoring-block-uncommenter@main

# 安定版（将来リリース予定）
uses: shoppingjaws/tf-refactoring-block-uncommenter@v1
```

## 📝 コミット規約

Conventional Commits 形式を使用：

```
feat: 新機能追加
fix: バグ修正
test: テスト追加・修正
docs: ドキュメント更新
refactor: リファクタリング
chore: その他の変更
```

**例:**
```bash
git commit -m "feat: add support for nested lifecycle blocks"
git commit -m "fix: correct indentation preservation in commenter"
git commit -m "test: add test cases for import blocks"
```

## 🧪 テスト方針

### 1. 単体テスト

各パッケージに `*_test.go` を配置：

- `pkg/parser/hcl_test.go`: ブロック検出テスト
- `pkg/commenter/commenter_test.go`: コメントアウトテスト

### 2. 統合テスト

`test.tf` を使用した実際の動作確認：

1. `test.tf` にコメントアウトされていないブロックを追加
2. コミット＆プッシュ
3. GitHub Actions が自動実行
4. PR が作成されることを確認

## 🚀 リリースフロー

1. **開発**: `main` ブランチで開発
2. **テスト**: Docker でローカルテスト + GitHub Actions で統合テスト
3. **安定版**: 十分にテストしたら `v1.x.x` タグを作成

```bash
# 安定版タグ作成（将来）
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
git tag -a v1 -m "Release v1" -f
git push origin v1 -f
```

## 🔍 トラブルシューティング

### "use of internal package not allowed" エラー

→ `internal/` ではなく `pkg/` を使用してください

### テストが GitHub Actions で失敗

→ ローカルで Docker を使用してテストを実行してください

### ブロックが検出されない

→ `git diff` に依存せず、全ファイルをスキャンする設計になっています

## 📚 参考リソース

- [Terraform Refactoring](https://developer.hashicorp.com/terraform/language/modules/develop/refactoring)
- [GitHub Actions Composite Actions](https://docs.github.com/en/actions/creating-actions/creating-a-composite-action)
- [Go Project Layout](https://github.com/golang-standards/project-layout)
