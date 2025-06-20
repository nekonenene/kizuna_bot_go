package bot

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"kizuna_bot_go/internal/api"
	"kizuna_bot_go/internal/config"
)

// KizunaBot はDiscordボットのメイン構造体
// Discordサーバーとの通信、設定管理、外部API呼び出しの機能を持つ
type KizunaBot struct {
	session   *discordgo.Session // Discord APIとの通信セッション
	config    *config.Config     // ボットの設定情報
	apiClient *api.Client        // 外部API呼び出し用のクライアント
}

// NewKizunaBot は新しいKizunaBotインスタンスを作成
func NewKizunaBot() (*KizunaBot, error) {
	// 設定を環境変数から読み込み
	cfg := config.NewConfig()

	// DiscordのBotトークンが設定されていない場合はエラー
	if cfg.BotToken == "" {
		return nil, fmt.Errorf("BOT_TOKEN is required")
	}

	// Discord APIとの通信セッションを作成
	// "Bot "プレフィックスを付けることでBot認証を行う
	session, err := discordgo.New("Bot " + cfg.BotToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create Discord session: %w", err)
	}

	// KizunaBotインスタンスを作成し、必要な構成要素を設定
	bot := &KizunaBot{
		session:   session,
		config:    cfg,
		apiClient: api.NewClient(cfg),
	}

	// Discordからメッセージ内容を受信するためのIntent（権限）を設定
	// これにより、ボットがメッセージの内容を読み取れるようになる
	session.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages | discordgo.IntentsMessageContent

	// イベントハンドラーを登録
	session.AddHandler(bot.messageCreate) // メッセージが投稿された時の処理
	session.AddHandler(bot.ready)         // ボットがDiscordに接続完了した時の処理

	return bot, nil
}

// Start はDiscordサーバーへの接続を開始
func (b *KizunaBot) Start() error {
	err := b.session.Open()
	if err != nil {
		return fmt.Errorf("failed to open Discord session: %w", err)
	}
	return nil
}

// Close はDiscordとの接続を安全に切断
func (b *KizunaBot) Close() {
	b.session.Close()
}

// ready はボットがDiscordに正常に接続された時に呼ばれるイベントハンドラー
func (b *KizunaBot) ready(s *discordgo.Session, event *discordgo.Ready) {
	log.Printf("ボットが正常に起動しました！ ログイン名: %s", event.User.String())
	log.Printf("招待URL: https://discord.com/api/oauth2/authorize?client_id=%s&permissions=2048&scope=bot", b.config.BotClientID)
}

// messageCreate はDiscordでメッセージが投稿された時に呼ばれるイベントハンドラー
func (b *KizunaBot) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ボット自身が投稿したメッセージは無視（無限ループを防ぐため）
	if m.Author.ID == s.State.User.ID {
		return
	}

	// スラッシュ（/）で始まるコマンドの処理
	if strings.HasPrefix(m.Content, "/") {
		b.handleCommand(s, m)
		return
	}

	// ボットがメンション（@名前）された時の処理
	for _, user := range m.Mentions {
		if user.ID == s.State.User.ID {
			b.handleMention(s, m)
			return
		}
	}

	// 特定のキーワードを含むメッセージに対する自動応答
	b.handlePatternMatching(s, m)
}

// handleCommand はスラッシュコマンド（/で始まるコマンド）を解析して適切な処理関数を呼び出す
func (b *KizunaBot) handleCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	// メッセージの前後の空白を除去
	content := strings.TrimSpace(m.Content)
	// スペースで区切ってコマンドと引数に分離
	parts := strings.Fields(content)

	if len(parts) == 0 {
		return
	}

	// コマンド部分（最初の要素）を小文字に変換
	command := strings.ToLower(parts[0])
	// 引数部分（2番目以降の要素）を取得
	args := parts[1:]

	// コマンドに応じて対応する処理関数を呼び出し
	switch command {
	case "/ping":
		b.handlePing(s, m) // 応答速度テスト
	case "/help":
		b.handleHelp(s, m) // ヘルプメッセージ表示
	case "/weather":
		b.handleWeather(s, m) // 天気予報取得
	case "/news":
		b.handleNews(s, m) // ニュース記事取得
	case "/dice":
		b.handleDice(s, m, args) // サイコロ機能
	case "/gourmet", "/gurume", "/grm":
		b.handleGourmet(s, m, args) // グルメ検索（複数のエイリアス対応）
	case "/image", "/img":
		b.handleImage(s, m, args) // 画像検索
	case "/rank":
		b.handleRank(s, m) // ユーザーアクティビティランキング
	case "/eng":
		b.handleTranslate(s, m, args, "en") // 英語翻訳
	case "/jpn", "/jap":
		b.handleTranslate(s, m, args, "ja") // 日本語翻訳
	case "/video", "/youtube":
		b.handleVideo(s, m, args) // 動画検索
	case "/vtuber":
		b.handleVTuber(s, m, args) // VTuber動画検索
	}
}

// handleMention はボットがメンション（@で呼び出し）された時の処理
func (b *KizunaBot) handleMention(s *discordgo.Session, m *discordgo.MessageCreate) {
	// メッセージからメンション部分を除去して実際の内容を取得
	content := m.Content
	for _, user := range m.Mentions {
		// <@ユーザーID> または <@!ユーザーID> の形式のメンションを削除
		content = strings.ReplaceAll(content, fmt.Sprintf("<@%s>", user.ID), "")
		content = strings.ReplaceAll(content, fmt.Sprintf("<@!%s>", user.ID), "")
	}
	content = strings.TrimSpace(content)

	// メンション内容に応じた会話応答を生成
	reply := b.getMunouMessage(content, m)
	if reply != "" {
		s.ChannelMessageSend(m.ChannelID, reply)
	}
}

// handlePatternMatching processes message patterns
func (b *KizunaBot) handlePatternMatching(s *discordgo.Session, m *discordgo.MessageCreate) {
	content := strings.ToLower(m.Content)

	// Only respond to specific patterns when not mentioned
	if strings.Contains(content, "天気は？") && len(m.Mentions) == 0 {
		b.handleWeather(s, m)
	}
}

// handlePing responds to ping command
func (b *KizunaBot) handlePing(s *discordgo.Session, m *discordgo.MessageCreate) {
	start := time.Now()
	msg, err := s.ChannelMessageSend(m.ChannelID, "Pong！")
	if err != nil {
		log.Printf("Error sending ping message: %v", err)
		return
	}

	duration := time.Since(start)
	editedContent := fmt.Sprintf("Pong！ 応答までに %.3f 秒かかったよ！", duration.Seconds())
	s.ChannelMessageEdit(m.ChannelID, msg.ID, editedContent)
}

// handleHelp shows help message
func (b *KizunaBot) handleHelp(s *discordgo.Session, m *discordgo.MessageCreate) {
	helpMessage := `/weather : 天気を教えるよ〜 :white_sun_small_cloud:
/news : 話題の記事をお届けしちゃうよ！ 暇な時はこれ！ :newspaper:
/gurume, /grm : お料理屋さんを探すよ、「/gurume 新宿 焼肉,個室,食べ放題」みたいに使ってね。カンマは「、」でもOK！ :fork_knife_plate:
/image, /img : いい写真を見つけてくるよ！ 1日100回までしか検索できないみたい… :art:
/dice : サイコロを回すよ。引数があると、それを最大値とするサイコロを回すよ :game_die:
/rank : 最近ヒマそうにしてる人を教えてあげるね :kiss_ww:
/eng : 英語でなんて言うのかがんばって翻訳するよ！ :capital_abcd:
/jpn : 日本語でどう言うのか考えるよ！ :flag_jp:
/video, /youtube : YouTubeから動画を探してくるよ！ 「/video ゲーム実況」みたいに使ってね :arrow_forward:
/vtuber : VTuberさんの動画を探してくるよ！ :dancer:
/ping : テスト用だよ
/help : これだよ`

	s.ChannelMessageSend(m.ChannelID, helpMessage)
}
