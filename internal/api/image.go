package api

import (
	"fmt"
	"math/rand"
	"strings"
)

// ImageSearchResponse はGoogle Custom Search APIの画像検索レスポンス構造体
type ImageSearchResponse struct {
	Items []struct {
		Link string `json:"link"` // 画像のURL
	} `json:"items"`
}

// GetImageSearch は指定されたクエリで画像を検索してランダムな結果を返す
func (c *Client) GetImageSearch(query string) (string, error) {
	// 検索ワードが空の場合
	if query == "" {
		return "検索ワードがないよ？ 『/image ねこ』みたいに書いてね！", nil
	}

	// Ruby版と同様にカンマを空白に置換
	query = strings.ReplaceAll(query, ",", " ")
	query = strings.ReplaceAll(query, "、", " ")

	// Google Custom Search API画像検索パラメータを構築
	params := map[string]string{
		"key":        c.config.CustomSearchAPIKey,
		"cx":         c.config.CustomSearchEngineID,
		"q":          query,
		"hl":         "ja",    // 言語設定（日本語）
		"searchType": "image", // 画像検索指定
		"num":        "10",    // Ruby版と同様に最大10件取得
	}

	// リクエストURLを構築
	requestURL := c.buildURL(c.config.CustomSearchAPIHost, params)

	var response ImageSearchResponse
	if err := c.makeGetRequest(requestURL, &response); err != nil {
		return "", fmt.Errorf("画像検索APIの呼び出しに失敗: %w", err)
	}

	// 検索結果がない場合
	if len(response.Items) == 0 {
		return "画像が見つからなかったよー", nil
	}

	// Ruby版と同様にランダムに1つの画像を選択
	selectedImage := response.Items[rand.Intn(len(response.Items))]
	imageLink := selectedImage.Link

	// Ruby版と同じメッセージ形式で返却
	message := "画像のお届けですよ〜 ヾﾉ｡ÒㅅÓ)ﾉｼ\n"
	message += imageLink

	return message, nil
}
