package api

import (
	"fmt"
	"net/url"
)

// TranslateResponse は翻訳APIのレスポンス構造体
type TranslateResponse struct {
	TranslatedText string `json:"translatedText,omitempty"` // Google Cloud Translation API用
	// Ruby版のGoogle Apps Script用のフィールド（レスポンスが直接テキストの場合に対応）
}

// GetTranslation は指定されたテキストを翻訳
func (c *Client) GetTranslation(text, targetLang string) (string, error) {
	if text == "" {
		return "翻訳するテキストを入力してね！", nil
	}

	// デフォルトは英語
	if targetLang == "" {
		targetLang = "en"
	}

	// URLエンコードしてクエリパラメータを構築
	values := url.Values{}
	values.Set("text", text)
	values.Set("target", targetLang)
	requestURL := c.config.GoogleTranslateAPIHost + "?" + values.Encode()

	// まずは単純なHTTPリクエストを試行
	resp, err := c.httpClient.Get(requestURL)
	if err != nil {
		return "", fmt.Errorf("翻訳APIへのリクエストに失敗: %w", err)
	}
	defer resp.Body.Close()

	// ステータスコードをチェック
	if resp.StatusCode != 200 {
		// Ruby版で発生している認証エラーの可能性
		if resp.StatusCode == 401 || resp.StatusCode == 403 {
			return "申し訳ありません。翻訳機能は現在メンテナンス中です。しばらく時間をおいてからお試しください。", nil
		}
		return "", fmt.Errorf("翻訳API呼び出しでエラー（ステータス: %d）", resp.StatusCode)
	}

	// レスポンスボディを読み取り
	body := make([]byte, 1024*10) // 最大10KB
	n, err := resp.Body.Read(body)
	if err != nil && err.Error() != "EOF" {
		return "", fmt.Errorf("翻訳APIレスポンスの読み取りに失敗: %w", err)
	}

	translatedText := string(body[:n])

	// 空のレスポンスの場合
	if translatedText == "" {
		return "翻訳に失敗しました。テキストが翻訳できない形式の可能性があります。", nil
	}

	// Ruby版と同じメッセージ形式で返却
	message := "翻訳してみたよ！\n"
	message += fmt.Sprintf("「%s」 これでどうかな？ σ(．_．@)", translatedText)

	return message, nil
}

// GetTranslationWithQuotes は「」で囲まれたテキストを抽出して翻訳
func (c *Client) GetTranslationWithQuotes(content, targetLang string) (string, error) {
	// 「」で囲まれたテキストを抽出（正規表現を使用せずシンプルに）
	startIndex := -1
	endIndex := -1

	runes := []rune(content)
	for i, r := range runes {
		if r == '「' && startIndex == -1 {
			startIndex = i + 1 // 「の次の文字から開始
		} else if r == '」' && startIndex != -1 {
			endIndex = i
			break
		}
	}

	// 「」が見つからない場合
	if startIndex == -1 || endIndex == -1 {
		if targetLang == "en" {
			return "「」で囲ってくれると英語に翻訳するよ〜", nil
		} else {
			return "「」で囲ってくれると日本語に翻訳するよ〜", nil
		}
	}

	// 抽出したテキストを翻訳
	extractedText := string(runes[startIndex:endIndex])
	return c.GetTranslation(extractedText, targetLang)
}

