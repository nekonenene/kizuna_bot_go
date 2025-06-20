package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"kizuna_bot_go/internal/config"
)

// Client は全ての外部API呼び出しを処理するHTTPクライアント
type Client struct {
	httpClient *http.Client   // HTTP通信用のクライアント
	config     *config.Config // 設定情報（APIキーやエンドポイントなど）
}

// NewClient は新しいAPIクライアントを作成
func NewClient(cfg *config.Config) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second, // タイムアウトを30秒に設定
		},
		config: cfg,
	}
}

// makeGetRequest は指定されたURLにGETリクエストを送信し、結果をJSONとして解析
func (c *Client) makeGetRequest(targetURL string, result interface{}) error {
	// HTTPリクエストを実行
	resp, err := c.httpClient.Get(targetURL)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// HTTPステータスコードをチェック
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	// レスポンスボディを読み取り
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// JSONを構造体にパース
	if err := json.Unmarshal(body, result); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return nil
}

// buildURL はベースURLにクエリパラメータを付加してURLを構築
func (c *Client) buildURL(baseURL string, params map[string]string) string {
	u, _ := url.Parse(baseURL)
	q := u.Query()
	// パラメータが空でない場合のみクエリに追加
	for key, value := range params {
		if value != "" {
			q.Set(key, value)
		}
	}
	u.RawQuery = q.Encode()
	return u.String()
}
