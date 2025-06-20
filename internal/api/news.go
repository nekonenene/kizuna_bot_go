package api

import (
	"fmt"
	"math/rand"
	"net/url"
)

// NewsResponse はRSS2JSON APIからのレスポンス構造体
type NewsResponse struct {
	Items []struct {
		Link string `json:"link"` // 記事のURL
	} `json:"items"`
}

// GetNews ははてなホットエントリーからランダムなニュースを取得
func (c *Client) GetNews() (string, error) {
	// RSSのURLをパーセントエンコードして安全なURL形式に変換
	encodedRSSURL := url.QueryEscape(c.config.HatenaHotentryRSS)

	// RSS2JSON APIのリクエストURLを構築
	// Ruby版と同様に最大50件取得（実際のRSSは30件程度）
	requestURL := fmt.Sprintf("%s?rss_url=%s&api_key=%s&count=%d",
		c.config.RSS2JSONAPIHost,
		encodedRSSURL,
		c.config.RSS2JSONAPIKey,
		50,
	)

	var response NewsResponse
	if err := c.makeGetRequest(requestURL, &response); err != nil {
		return "", fmt.Errorf("ニュース取得APIの呼び出しに失敗: %w", err)
	}

	// 取得した記事がない場合
	if len(response.Items) == 0 {
		return "ニュースが取得できませんでした。", nil
	}

	// Ruby版と同じようにランダムに1つの記事を選択
	randomLink := response.Items[rand.Intn(len(response.Items))].Link

	// Ruby版と同じメッセージ形式で返却
	return fmt.Sprintf("ニュースのお届けだよー！ ガシーン ヽ(•̀ω•́ )ゝ\n%s", randomLink), nil
}
