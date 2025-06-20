package config

import (
	"os"
)

// Config はボットの動作に必要な全ての設定値を保持する構造体
type Config struct {
	// Discord関連の設定
	BotClientID string // DiscordのBot Client ID
	BotToken    string // DiscordのBotトークン（認証に使用）

	// 外部API接続用のキー
	RSS2JSONAPIKey       string // ニュース取得用のRSS2JSON APIキー
	RecruitAPIKey        string // グルメ検索用のリクルートAPIキー
	CustomSearchEngineID string // 画像検索用のGoogleカスタム検索エンジンID
	CustomSearchAPIKey   string // 画像検索用のGoogleカスタム検索APIキー
	YouTubeDataAPIKey    string // 動画検索用のYouTube Data APIキー

	// 各種APIのエンドポイント（接続先URL）
	LivedoorWeatherAPIHost string // ライブドア天気予報API
	RSS2JSONAPIHost        string // RSS2JSON API（ニュース取得用）
	HotPepperAPIHost       string // ホットペッパーAPI（グルメ検索用）
	CustomSearchAPIHost    string // Google カスタム検索API
	YouTubeDataAPIHost     string // YouTube Data API
	GoogleTranslateAPIHost string // Google翻訳API

	// アプリケーション定数
	TokyoCityID       int    // 天気予報で使用する東京の都市ID
	RankTotalCount    int    // ユーザーランキング機能で取得するメッセージ数
	HatenaHotentryRSS string // はてなホットエントリーのRSS URL
}

// NewConfig は環境変数から設定を読み込んで新しいConfigインスタンスを作成
func NewConfig() *Config {
	return &Config{
		// Discord関連の環境変数を取得
		BotClientID: os.Getenv("BOT_CLIENT_ID"),
		BotToken:    os.Getenv("BOT_TOKEN"),

		// 外部API用のキーを環境変数から取得
		RSS2JSONAPIKey:       os.Getenv("RSS2JSON_API_KEY"),
		RecruitAPIKey:        os.Getenv("RECRUIT_API_KEY"),
		CustomSearchEngineID: os.Getenv("CUSTOM_SEARCH_ENGINE_ID"),
		CustomSearchAPIKey:   os.Getenv("CUSTOM_SEARCH_API_KEY"),
		YouTubeDataAPIKey:    os.Getenv("YOUTUBE_DATA_API_KEY"),

		// 各APIのエンドポイントURL（固定値）
		LivedoorWeatherAPIHost: "https://weather.tsukumijima.net/api/forecast",
		RSS2JSONAPIHost:        "https://api.rss2json.com/v1/api.json",
		HotPepperAPIHost:       "https://webservice.recruit.co.jp/hotpepper/gourmet/v1",
		CustomSearchAPIHost:    "https://www.googleapis.com/customsearch/v1",
		YouTubeDataAPIHost:     "https://www.googleapis.com/youtube/v3",
		GoogleTranslateAPIHost: "https://script.google.com/macros/s/AKfycbzX3hgwpkCG-q-47nvu9CpeGXJ2uoQVbAngwNpbHjx6jCiOMXE/exec",

		// アプリケーションで使用する定数値
		TokyoCityID:       130010,                                     // ライブドア天気APIでの東京の都市コード
		RankTotalCount:    200,                                        // ユーザーランキングで過去何件のメッセージを集計するか
		HatenaHotentryRSS: "https://b.hatena.ne.jp/hotentry?mode=rss", // はてなホットエントリーのRSS配信URL
	}
}
