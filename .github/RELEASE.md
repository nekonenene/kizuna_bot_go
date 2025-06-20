# リリース手順

## 自動リリースの使用方法

このプロジェクトでは、Gitタグをpushすることで自動的にリリースが作成されます。

### 1. バージョンタグの作成

```bash
# バージョン番号を決定（例: v1.0.0）
git tag v1.0.0

# タグをリモートにpush
git push origin v1.0.0
```

### 2. 自動ビルド開始

タグをpushすると、GitHub ActionsがGoReleaserを使用して自動的に以下を実行します：

- 複数プラットフォーム向けのバイナリビルド
  - Linux (x64, ARM64)
  - Windows (x64)
  - macOS (Intel, Apple Silicon)
- バイナリの圧縮とアーカイブ作成
- GitHubリリースページの作成
- バイナリファイルのアップロード
- 変更履歴（Changelog）の自動生成

### 3. リリースの確認

- [Releases](https://github.com/nekonenene/kizuna_bot_go/releases) ページで新しいリリースを確認
- 各プラットフォーム向けのバイナリがダウンロード可能

## サポートプラットフォーム

| プラットフォーム | アーキテクチャ | ファイル名 |
|----------------|----------------|-----------|
| Linux | x64 | `kizuna_bot-Linux-x86_64.tar.gz` |
| Linux | ARM64 | `kizuna_bot-Linux-arm64.tar.gz` |
| Windows | x64 | `kizuna_bot-Windows-x86_64.zip` |
| macOS | Intel | `kizuna_bot-Darwin-x86_64.tar.gz` |
| macOS | Apple Silicon | `kizuna_bot-Darwin-arm64.tar.gz` |

## バージョニング規則

セマンティックバージョニング（SemVer）に従ってバージョン番号を決定：

- **メジャーバージョン** (`v2.0.0`): 互換性のない変更
- **マイナーバージョン** (`v1.1.0`): 後方互換性のある新機能
- **パッチバージョン** (`v1.0.1`): 後方互換性のあるバグ修正

### 例

```bash
# 新機能追加
git tag v1.1.0
git push origin v1.1.0

# バグ修正
git tag v1.0.1
git push origin v1.0.1

# 互換性のない変更
git tag v2.0.0
git push origin v2.0.0
```

## 手動ビルド（オプション）

GitHub Actionsを使わずに手動でビルドする場合：

```bash
# GoReleaserを使用したローカルビルド（リリースなし）
goreleaser build --snapshot --clean

# 全プラットフォーム向けビルド（従来の方法）
make build-all

# 特定プラットフォーム向けビルド
make build-linux
make build-windows
make build-darwin
```
