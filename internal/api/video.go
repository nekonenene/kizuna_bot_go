package api

import (
	"fmt"
	"math/rand"
)

// YouTubeSearchResponse はYouTube Data API検索結果のレスポンス構造体
type YouTubeSearchResponse struct {
	Items []struct {
		ID struct {
			VideoID string `json:"videoId"` // 動画のID
		} `json:"id"`
	} `json:"items"`
}

// GetVideoByQuery は検索クエリを使ってYouTube動画をランダムに取得
func (c *Client) GetVideoByQuery(query string) (string, error) {
	// YouTube Data API検索パラメータを構築
	params := map[string]string{
		"key":         c.config.YouTubeDataAPIKey,
		"part":        "id",
		"type":        "video",
		"maxResults":  "50",
		"order":       "date",  // 日付順でソート
		"regionCode":  "JP",    // 日本地域での検索
	}

	// クエリが空でない場合のみq パラメータを追加
	if query != "" {
		params["q"] = query
	}

	// リクエストURLを構築
	requestURL := c.buildURL(c.config.YouTubeDataAPIHost+"/search", params)

	var response YouTubeSearchResponse
	if err := c.makeGetRequest(requestURL, &response); err != nil {
		return "", fmt.Errorf("YouTube動画検索APIの呼び出しに失敗: %w", err)
	}

	// 検索結果がない場合
	if len(response.Items) == 0 {
		return "", fmt.Errorf("動画が見つかりませんでした")
	}

	// Ruby版と同様にランダムに1つの動画を選択
	selectedVideo := response.Items[rand.Intn(len(response.Items))]
	videoID := selectedVideo.ID.VideoID

	// YouTube視聴URLを構築
	return fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID), nil
}

// GetVideoByChannel は指定されたチャンネルIDから動画をランダムに取得
func (c *Client) GetVideoByChannel(channelID string) (string, error) {
	if channelID == "" {
		return "", fmt.Errorf("チャンネルIDが指定されていません")
	}

	// YouTube Data API検索パラメータを構築
	params := map[string]string{
		"key":         c.config.YouTubeDataAPIKey,
		"part":        "id",
		"type":        "video",
		"channelId":   channelID,
		"maxResults":  "50",
		"order":       "date", // 日付順でソート
	}

	// リクエストURLを構築
	requestURL := c.buildURL(c.config.YouTubeDataAPIHost+"/search", params)

	var response YouTubeSearchResponse
	if err := c.makeGetRequest(requestURL, &response); err != nil {
		return "", fmt.Errorf("YouTubeチャンネル動画検索APIの呼び出しに失敗: %w", err)
	}

	// 検索結果がない場合
	if len(response.Items) == 0 {
		return "", fmt.Errorf("チャンネルから動画が見つかりませんでした")
	}

	// Ruby版と同様にランダムに1つの動画を選択
	selectedVideo := response.Items[rand.Intn(len(response.Items))]
	videoID := selectedVideo.ID.VideoID

	// YouTube視聴URLを構築
	return fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID), nil
}

// GetVideoSearch は検索クエリに基づいて動画検索メッセージを生成
func (c *Client) GetVideoSearch(query string) (string, error) {
	// YouTube動画を検索
	videoURL, err := c.GetVideoByQuery(query)
	if err != nil {
		return "いい動画が見つけられなかったよ、ごめんね", nil
	}

	// Ruby版と同じメッセージ形式で返却
	var message string
	if query != "" {
		message = fmt.Sprintf("『%s』の", query)
	} else {
		message = "最近の"
	}
	message += "動画を探してきたよ！ ( ⁎ᵕᴗᵕ⁎ ) :heartbeat:\n"
	message += videoURL

	return message, nil
}