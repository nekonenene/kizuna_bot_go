# Kizuna Bot Go

Ruby版Kizuna BotのGo実装版です。Discord上で動作する多機能ボットで、天気予報、ニュース配信、グルメ検索、会話応答などの機能を提供します。

## 機能

### 実装済みコマンド
- `/ping` - 応答時間テスト
- `/help` - 利用可能なコマンド一覧を表示
- `/weather` - 東京の天気予報を取得
- `/news` - はてなホットエントリーからランダムなニュースを配信
- `/dice [最大値]` - サイコロを振る（デフォルト6面）
- `/gourmet <地域> [キーワード]` - レストラン検索（ホットペッパーAPI使用）
- `/image`, `/img <検索ワード>` - 画像検索（Google Custom Search API使用）
- `/video`, `/youtube <検索ワード>` - YouTube動画検索
- `/vtuber [検索ワード]` - VTuber動画検索
- `/eng <テキスト>` - 英語翻訳
- `/jpn`, `/jap <テキスト>` - 日本語翻訳
- `/rank` - チャンネル内ユーザーアクティビティランキング

### 実装済み応答機能
- メンション時の会話応答（天気、ニュース、翻訳、ランキングなど）
- 「天気は？」などのパターンマッチング応答
- 挨拶や感情表現への自動応答
- 「ゆーま？」での特定チャンネル動画検索
- 「英語で「テキスト」」「日本語で「テキスト」」での翻訳機能

## 必要な環境

- Go 1.24以上
- Discord botのトークンとクライアントID
- 各種外部APIのキー（下記参照）

## セットアップ

### 1. リポジトリのクローン

```bash
git clone <repository-url>
cd kizuna_bot_go
```

### 2. 依存関係のインストール

```bash
make deps
```

### 3. 環境設定

```bash
make init
```

`.env`ファイルが作成されるので、必要なAPIキーを設定してください：

```env
BOT_CLIENT_ID=your_discord_bot_client_id
BOT_TOKEN=your_discord_bot_token

# ニュース機能に必要
RSS2JSON_API_KEY=your_rss2json_api_key

# グルメ検索に必要
RECRUIT_API_KEY=your_recruit_api_key

# 画像検索に必要（未実装）
CUSTOM_SEARCH_ENGINE_ID=your_custom_search_engine_id
CUSTOM_SEARCH_API_KEY=your_custom_search_api_key

# 動画検索に必要（未実装）
YOUTUBE_DATA_API_KEY=your_youtube_data_api_key
```

### 4. ビルドと実行

#### 開発環境での実行
```bash
make dev
```

#### 本番用ビルド（Linux サーバー用）
```bash
make build-linux
```

生成されたバイナリは `build/kizuna_bot_go-linux` に出力されます。

## 利用可能なMakeコマンド

```bash
make deps          # 依存関係のインストール
make build         # アプリケーションのビルド
make build-linux   # Linux用ビルド（サーバーデプロイ用）
make build-all     # 全プラットフォーム用ビルド
make dev           # 開発モードで実行
make run           # ビルドせずに実行
make test          # テスト実行
make clean         # ビルド成果物の削除
make init          # 環境設定ファイルの作成
make help          # ヘルプ表示
```

## Discord Botの設定

1. [Discord Developer Portal](https://discord.com/developers/applications)でアプリケーションを作成
2. Bot セクションでbotを作成し、トークンを取得
3. OAuth2 セクションで以下の権限を設定：
   - Scopes: `bot`
   - Bot Permissions: `Send Messages`, `Read Messages`, `Read Message History`
4. 生成されたURLでbotをサーバーに招待

## プロジェクト構造

```
kizuna_bot_go/
├── main.go                 # アプリケーションエントリーポイント
├── go.mod                  # Go モジュール定義
├── Makefile               # ビルドと開発用コマンド
├── .env.example           # 環境変数のテンプレート
├── internal/
│   ├── bot/               # Discord bot実装
│   │   ├── bot.go         # メインbot構造とセットアップ
│   │   └── handlers.go    # コマンドとメッセージハンドラー
│   ├── api/               # 外部API統合
│   │   ├── client.go      # HTTPクライアントラッパー
│   │   ├── weather.go     # 天気API統合
│   │   ├── news.go        # ニュースAPI統合
│   │   └── gourmet.go     # グルメ検索API
│   └── config/
│       └── config.go      # 設定管理
└── build/                 # ビルド出力ディレクトリ
```

## 使用している外部API

- **Livedoor Weather API**: 天気予報の取得
- **RSS2JSON API**: はてなホットエントリーの取得
- **ホットペッパーグルメAPI**: レストラン検索
- **Google Custom Search API**: 画像検索（未実装）
- **YouTube Data API**: 動画検索（未実装）
- **Google Translate API**: 翻訳機能（未実装）

## デプロイ

### サーバーへのデプロイ

1. Linux用バイナリをビルド：
   ```bash
   make build-linux
   ```

2. バイナリと `.env` ファイルをサーバーにアップロード

3. サーバーで実行：
   ```bash
   ./kizuna_bot_go-linux
   ```

### systemdサービスとしての実行（推奨）

`/etc/systemd/system/kizuna-bot.service` を作成：

```ini
[Unit]
Description=Kizuna Bot Go
After=network.target

[Service]
Type=simple
User=your-user
WorkingDirectory=/path/to/bot
ExecStart=/path/to/bot/kizuna_bot_go-linux
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

サービスの有効化と開始：
```bash
sudo systemctl enable kizuna-bot
sudo systemctl start kizuna-bot
```

## トラブルシューティング

### よくある問題

1. **APIキーが見つからない**: `.env` ファイルに必要な環境変数が設定されているか確認
2. **Discord権限エラー**: botに適切な権限（メッセージ送信、読み取り）が付与されているか確認
3. **API制限**: 一部のAPIには使用制限があります（例：画像検索は1日100回まで）

### ログの確認

アプリケーションは標準出力にログを出力します。systemdを使用している場合：

```bash
sudo journalctl -u kizuna-bot -f
```

## 開発への貢献

1. フォークしてください
2. フィーチャーブランチを作成してください (`git checkout -b feature/amazing-feature`)
3. 変更をコミットしてください (`git commit -m 'Add some amazing feature'`)
4. ブランチにプッシュしてください (`git push origin feature/amazing-feature`)
5. プルリクエストを開いてください

## ライセンス

このプロジェクトは元のRuby版に準拠します。

## 参考

- 元のRuby版: `kizuna_bot_image/` ディレクトリ
- [discordgo ライブラリ](https://github.com/bwmarrin/discordgo)
- [Discord Developer Portal](https://discord.com/developers/docs/)
