.PHONY: build clean run dev test deps

# バイナリファイル名
BINARY_NAME=kizuna_bot_go
BUILD_DIR=build

# ビルド時のフラグ（バイナリサイズを最小化）
LDFLAGS=-ldflags "-s -w"

# デフォルトターゲット
all: build

# 依存関係の整理・ダウンロード
deps:
	go mod tidy

# アプリケーションのビルド
build: deps
	mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .

# Linux用ビルド（サーバーデプロイ用）
build-linux: deps
	mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux .

# Windows用ビルド
build-windows: deps
	mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME).exe .

# macOS用ビルド
build-darwin: deps
	mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin .

# 全プラットフォーム用ビルド
build-all: build-linux build-windows build-darwin

# 開発モードでアプリケーションを実行
dev: build
	./$(BUILD_DIR)/$(BINARY_NAME)

# ビルドせずにアプリケーションを実行
run:
	go run .

# テストの実行
test:
	go test -v ./...

# ビルド成果物のクリーンアップ
clean:
	go clean
	rm -rf $(BUILD_DIR)

# アプリケーションのシステムインストール
install: build
	cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/

# サンプルから.envファイルを作成
init:
	cp .env.example .env
	@echo ".envファイルを作成しました。APIキーを設定してください。"

# ヘルプ表示
help:
	@echo "利用可能なコマンド:"
	@echo "  build        - アプリケーションのビルド"
	@echo "  build-linux  - Linux用ビルド（サーバーデプロイ用）"
	@echo "  build-all    - 全プラットフォーム用ビルド"
	@echo "  dev          - 開発モードで実行"
	@echo "  run          - ビルドせずに実行"
	@echo "  test         - テスト実行"
	@echo "  clean        - ビルド成果物のクリーンアップ"
	@echo "  deps         - 依存関係のインストール"
	@echo "  init         - .envファイルの作成"
	@echo "  install      - /usr/local/binにインストール"
	@echo "  help         - このヘルプを表示"
