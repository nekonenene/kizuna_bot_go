package bot

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// handleWeather sends weather information
func (b *KizunaBot) handleWeather(s *discordgo.Session, m *discordgo.MessageCreate) {
	message, err := b.apiClient.GetWeather()
	if err != nil {
		log.Printf("Error getting weather: %v", err)
		message = "天気情報の取得に失敗しました。しばらく時間をおいてからお試しください。"
	}
	s.ChannelMessageSend(m.ChannelID, message)
}

// handleNews sends news information
func (b *KizunaBot) handleNews(s *discordgo.Session, m *discordgo.MessageCreate) {
	message, err := b.apiClient.GetNews()
	if err != nil {
		log.Printf("Error getting news: %v", err)
		message = "ニュース取得に失敗しました。しばらく時間をおいてからお試しください。"
	}
	s.ChannelMessageSend(m.ChannelID, message)
}

// handleDice rolls a dice
func (b *KizunaBot) handleDice(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	max := 6 // Default to 6-sided dice

	if len(args) > 0 {
		if val, err := strconv.Atoi(args[0]); err == nil && val > 0 {
			max = val
		}
	}

	result := rand.Intn(max) + 1
	message := fmt.Sprintf("%d面サイコロを回したら、「%d」が出たよ！", max, result)
	s.ChannelMessageSend(m.ChannelID, message)
}

// handleGourmet searches for restaurants
func (b *KizunaBot) handleGourmet(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	address := ""
	keyword := ""

	if len(args) > 0 {
		address = args[0]
	}
	if len(args) > 1 {
		keyword = strings.Join(args[1:], " ")
	}

	message, err := b.apiClient.GetGourmet(address, keyword)
	if err != nil {
		log.Printf("Error getting gourmet info: %v", err)
		message = "グルメ検索に失敗しました。しばらく時間をおいてからお試しください。"
	}
	s.ChannelMessageSend(m.ChannelID, message)
}

// handleImage はGoogle Custom Search APIで画像を検索
func (b *KizunaBot) handleImage(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	query := strings.Join(args, " ")
	message, err := b.apiClient.GetImageSearch(query)
	if err != nil {
		log.Printf("画像検索エラー: %v", err)
		message = "画像検索に失敗しました。しばらく時間をおいてからお試しください。"
	}
	s.ChannelMessageSend(m.ChannelID, message)
}

// handleRank はチャンネル内のユーザーアクティビティランキングを表示
func (b *KizunaBot) handleRank(s *discordgo.Session, m *discordgo.MessageCreate) {
	message, err := b.apiClient.GetUserRanking(s, m.ChannelID)
	if err != nil {
		log.Printf("ランキング取得エラー: %v", err)
		message = "ランキングの取得に失敗しました。しばらく時間をおいてからお試しください。"
	}
	s.ChannelMessageSend(m.ChannelID, message)
}

// handleTranslate はテキストを指定された言語に翻訳
func (b *KizunaBot) handleTranslate(s *discordgo.Session, m *discordgo.MessageCreate, args []string, targetLang string) {
	text := strings.Join(args, " ")
	message, err := b.apiClient.GetTranslation(text, targetLang)
	if err != nil {
		log.Printf("翻訳エラー: %v", err)
		message = "翻訳に失敗しました。しばらく時間をおいてからお試しください。"
	}
	s.ChannelMessageSend(m.ChannelID, message)
}

// handleVideo はYouTubeから動画を検索
func (b *KizunaBot) handleVideo(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	query := strings.Join(args, " ")
	message, err := b.apiClient.GetVideoSearch(query)
	if err != nil {
		log.Printf("動画検索エラー: %v", err)
		message = "動画検索に失敗しました。しばらく時間をおいてからお試しください。"
	}
	s.ChannelMessageSend(m.ChannelID, message)
}

// handleVTuber はVTuber動画を検索
func (b *KizunaBot) handleVTuber(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	// Ruby版と同様に"VTuber "を前に付けて検索
	query := "VTuber " + strings.Join(args, " ")
	message, err := b.apiClient.GetVideoSearch(query)
	if err != nil {
		log.Printf("VTuber動画検索エラー: %v", err)
		message = "VTuber動画検索に失敗しました。しばらく時間をおいてからお試しください。"
	}
	s.ChannelMessageSend(m.ChannelID, message)
}

// getMunouMessage はメンション時の応答メッセージを生成
func (b *KizunaBot) getMunouMessage(message string, s *discordgo.Session, m *discordgo.MessageCreate) string {
	content := strings.ToLower(message)

	switch {
	case strings.Contains(content, "英語で"):
		// 「」で囲まれたテキストを英語に翻訳
		if result, err := b.apiClient.GetTranslationWithQuotes(m.Content, "en"); err == nil {
			return result
		}
		return "翻訳に失敗しました"
	case strings.Contains(content, "日本語で"):
		// 「」で囲まれたテキストを日本語に翻訳
		if result, err := b.apiClient.GetTranslationWithQuotes(m.Content, "ja"); err == nil {
			return result
		}
		return "翻訳に失敗しました"
	case strings.Contains(content, "天気"):
		// メンション応答での天気機能
		if weather, err := b.apiClient.GetWeather(); err == nil {
			return weather
		}
		return "天気情報の取得に失敗しました"
	case strings.Contains(content, "さいころ") || strings.Contains(content, "サイコロ"):
		result := rand.Intn(6) + 1
		return fmt.Sprintf("6面サイコロを回したら、「%d」が出たよ！", result)
	case strings.Contains(content, "ニュース"):
		// メンション応答でのニュース機能
		if news, err := b.apiClient.GetNews(); err == nil {
			return news
		}
		return "ニュース取得に失敗しました"
	case strings.Contains(content, "ランキング"):
		// メンション応答でのランキング機能
		if ranking, err := b.apiClient.GetUserRanking(s, m.ChannelID); err == nil {
			return ranking
		}
		return "ランキングの取得に失敗しました"
	case strings.Contains(content, "おなかすいた") || strings.Contains(content, "おなすき"):
		responses := []string{
			"栄養あるものをしっかり食べようね！",
			"ぐぐぅぅー",
			"実は /gurume コマンドは /gourmet や /grm と打っても使えるよ！",
		}
		return responses[rand.Intn(len(responses))]
	case strings.Contains(content, "おはよ"):
		responses := []string{
			"おはよう〜！",
			"きょうもがんばろうね！",
			"ねむい……顔、洗わなきゃ・・・・",
			"(n'∀')η ﾔｧｰｯﾎｫｰ!!",
		}
		return responses[rand.Intn(len(responses))]
	case strings.Contains(content, "こんにち"):
		return "こんにちは〜！"
	case strings.Contains(content, "おやすみ"):
		responses := []string{
			"おやすみ〜",
			"ｚｚｚ。。。。。",
			"また明日〜",
			"明日はもっといい日にしようね！",
			"もうこんな時間なんだね、おやすみなさい",
		}
		return responses[rand.Intn(len(responses))]
	case strings.Contains(content, "ねむい") || strings.Contains(content, "眠い"):
		responses := []string{
			"だよねわかる・・・・",
			"うとうと・・・・",
			"もう寝よう？",
			"眠って、楽になっちゃおうよ？",
		}
		return responses[rand.Intn(len(responses))]
	case strings.Contains(content, "元気？"):
		responses := []string{
			"うん！ ありがと！",
			"元気だよー！！",
		}
		return responses[rand.Intn(len(responses))]
	case strings.Contains(content, "かわい"):
		responses := []string{
			"えへへ :heartbeat:",
			"そ……そうかな…",
			"や……やっぱり？！ なんて…",
			"照れるよお・・・・",
		}
		return responses[rand.Intn(len(responses))]
	case strings.Contains(content, "大好き") || strings.Contains(content, "だいすき"):
		return "私もだよ！"
	case strings.Contains(content, "好き") || strings.Contains(content, "すき"):
		responses := []string{
			"いいねいいね！！ :sparkles: :sparkles:",
			"わたしもわたしも！ :white_flower:",
			":heartpulse: 大好きだよ！！ :heartpulse:",
		}
		return responses[rand.Intn(len(responses))]
	case strings.Contains(content, "愛して") || strings.Contains(content, "あいして"):
		responses := []string{
			"えっ………",
			"ちょっと気持ち悪い",
			"そういうのはちょっと………",
			"普通にキモいんですけど、、",
			"なに言ってるんですか？",
			"やめてください……体調悪くなりました",
			"ヒィッッ！！ 近付かないで！！！！",
		}
		return responses[rand.Intn(len(responses))]
	case strings.Contains(content, "ありがと"):
		responses := []string{
			"どういたしまして！",
			"いえいえ〜〜",
			"今後ともごひいきにー！",
		}
		return responses[rand.Intn(len(responses))]
	case strings.Contains(content, "がんば"):
		responses := []string{
			"いっしょにがんばろー！",
			"楽しい日になるといいね！",
			"わっしょい！ わっしょい！ └(ﾟ∀ﾟ└)",
			"ファイトオー！ :fire:",
		}
		return responses[rand.Intn(len(responses))]
	case strings.Contains(content, "くまくま"):
		responses := []string{
			"ざわ……ざわ……",
			"ʕ•̀ω•́ʔ  ʕ•̀ω•́ʔ  ʕ•̀ω•́ʔ  ʕ•̀ω•́ʔ",
			"ฅʕ•ᴥ•ʔฅ ʕ´•ᴥ•`ʔ",
			"(σ´･(ｪ)･)σﾖﾛｼｸﾏｰ!!",
			"つられクマー！！ ＞ ＜",
			"（´・(ェ)・｀） くま？",
			"いわもウェイ！！",
			"ざわわ、ざわわ、ざわわ",
		}
		return responses[rand.Intn(len(responses))]
	case strings.Contains(content, "疲れ") || strings.Contains(content, "つかれ"):
		responses := []string{
			"よしよし・・・ ( ,,´・ω・)ﾉ (´っω・｀｡)",
			"すこし休もうねー？ ヾ(´ー｀*)",
			"生きてるからまだ大丈夫だよ！",
		}
		return responses[rand.Intn(len(responses))]
	case strings.Contains(content, "つらい") || strings.Contains(content, "ちゅらい"):
		responses := []string{
			"わかる・・・",
			"5000兆円あげるから元気だして",
			"よしよし・・・ ( ,,´・ω・)ﾉ (´っω・｀｡)",
			"休んでもいいんだよ？ _(* v v)。",
		}
		return responses[rand.Intn(len(responses))]
	case strings.Contains(content, "死に"):
		responses := []string{
			"生きてーーーーーっっっっ！！！！ (> <!!!!",
			"へんじがない、ただのしかばねのようだ",
			"まだ死ぬには早いですよ！",
			"にゃにゃにゃにゃにゃ！",
		}
		return responses[rand.Intn(len(responses))]
	case strings.Contains(content, "にゃん") || strings.Contains(content, "にゃー"):
		responses := []string{
			"にゃ〜ん :cat2:",
			"わかるにゃ・・・・・",
			"みゃみゃ〜ん！ V(=^・ω・^=)v",
			"(」・ω・)」うー！(/・ω・)/にゃー！",
		}
		return responses[rand.Intn(len(responses))]
	case strings.Contains(content, "ひま") || strings.Contains(content, "ヒマ") || strings.Contains(content, "暇"):
		// Ruby版と同様に「ひま」でニュースを返す
		if news, err := b.apiClient.GetNews(); err == nil {
			return news
		}
		return "ニュース取得に失敗しました"
	case strings.Contains(content, "アニメ"):
		return "アニメといえばキルミーベイベーだよね！"
	case strings.HasSuffix(content, "！！") || strings.HasSuffix(content, "!!"):
		responses := []string{
			"そうだね！！！",
			"元気いっぱいだねー！！",
			"うん！！",
		}
		return responses[rand.Intn(len(responses))]
	case strings.Contains(content, "ゆーま") && (strings.HasSuffix(content, "？") || strings.HasSuffix(content, "?")):
		// Ruby版と同じ特定チャンネル（ゆーま）の動画検索
		if videoURL, err := b.apiClient.GetVideoByChannel("UC_9DxYZ_4Lhm9ujFvcHryNw"); err == nil {
			return fmt.Sprintf("ゆーまってこの人かな？！ (੭ु ›ω‹ )੭ु⁾⁾ %s", videoURL)
		}
		return "ゆーまの動画が見つからなかったよ"
	case strings.HasSuffix(content, "？") || strings.HasSuffix(content, "?"):
		responses := []string{
			"そうかも？",
			"わからぬ〜",
			"むずかしい質問だねー",
			"知らなーい",
			"そうなの？",
		}
		return responses[rand.Intn(len(responses))]
	case strings.Contains(content, "help"):
		// Return help message
		return "コマンドについては /help を使ってね！"
	default:
		responses := []string{
			"なるほど〜",
			"それそれ！！",
			"ニャンニャン (ﾉ*ФωФ) //",
			"そうなんだ〜",
			"うんうん！",
		}
		return responses[rand.Intn(len(responses))]
	}
}
