# Kizuna Bot Go

Ruby版 Kizuna Bot ( https://github.com/nekonenene/kizuna_bot_image ) の、Go言語での実装です。  
Discord上で動作する多機能ボットで、天気予報、ニュース配信、グルメ検索、会話応答などの機能を提供します。

作成には [Claude Code](https://docs.anthropic.com/ja/docs/claude-code/overview) を大いに活用しました。  
（おかげで README も長すぎるので、メンテしにくい項は徐々に削っていきます）


## 機能

### 実装済みコマンド

- `/ping` - 応答時間テスト
- `/help` - 利用可能なコマンド一覧を表示
- `/weather` - 天気予報を取得
- `/news` - ランダムなニュースを配信
- `/dice [最大値]` - サイコロを振る（デフォルト6面）
- `/gourmet <地域> [キーワード]` - レストラン検索
- `/image`, `/img <検索ワード>` - 画像検索
- `/video`, `/youtube <検索ワード>` - YouTube動画検索
- `/vtuber [検索ワード]` - VTuber動画検索
- `/eng <テキスト>` - 英語翻訳
- `/jpn <テキスト>` - 日本語翻訳
- `/rank` - チャンネル内のユーザー発言数ランキング

### 実装済み応答機能

- 特定の文字列を含むメンション時の会話応答（天気、ニュース、翻訳、ランキングなど）
- あいさつなどのメンションに自動応答（人工無能）

### 未実装

- 日英の翻訳機能（昔は動いていたが、Google Apps Script の実行可能APIの仕様変更により外部から簡単に呼び出せなくなった）


## 開発に必要な環境

- Go 1.24以上
- Discord botのトークンとクライアントID
- 外部APIのキー（ `.env.sample` 参照）


## セットアップ

### 1. 依存関係のインストール

```bash
make deps
```

### 2. 環境設定

```bash
make init
```

`.env` ファイルが作成されるので、必要なAPIキーを設定する。

### 3. ビルドと実行

#### 開発環境での実行

```bash
make dev
```

#### 本番用ビルド（Linux サーバー用）

```bash
make build-linux
```

生成されたバイナリは `build/kizuna_bot-linux` に出力される。


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
2. Bot セクションで bot を作成し、トークンを取得
3. OAuth2 セクションで以下の権限を設定：
   - Scopes: `bot`
   - Bot Permissions: `Send Messages`, `Read Messages`, `Read Message History`
4. `make dev` や `make run` での実行時に出力される URL を使って bot をサーバーに招待


## デプロイ

### サーバーへのデプロイ

1. Linux用バイナリをビルドし、サーバーにアップロード
2. `.env` ファイルをサーバーにアップロードするか、環境変数として設定
3. サーバー上でバイナリを実行

### systemdサービスとしての実行（推奨）

`/etc/systemd/system/kizuna-bot.service` を作成：

```ini
[Unit]
Description=Kizuna Bot Go
After=network.target

[Service]
Type=simple
User=your-user-name
WorkingDirectory=/path/to/bot
ExecStart=/path/to/bot/kizuna_bot
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

そうして実行した場合、以下のコマンドでログを確認可能。

```bash
sudo journalctl -u kizuna-bot -f
```


## 参考

- [discordgo ライブラリ](https://github.com/bwmarrin/discordgo)
- [Discord Developer Portal](https://discord.com/developers/docs/)
