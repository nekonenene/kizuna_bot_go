package api

import (
	"fmt"
	"math/rand"
	"sort"

	"github.com/bwmarrin/discordgo"
)

// UserRankEntry はユーザーランキングの1エントリを表す構造体
type UserRankEntry struct {
	Username string // ユーザー名
	Count    int    // 投稿数
}

// GetUserRanking は指定されたチャンネルでのユーザーアクティビティランキングを生成
func (c *Client) GetUserRanking(session *discordgo.Session, channelID string) (string, error) {
	maxCount := c.config.RankTotalCount // 設定で指定された件数（200件）
	maxPerPage := 100                   // Discord APIの制限：1回のリクエストで最大100件

	var allMessages []*discordgo.Message
	beforeID := ""
	remaining := maxCount

	// Ruby版と同様にページ分けしてメッセージを取得
	for remaining > 0 {
		// 今回のリクエストで取得する件数を決定
		limit := remaining
		if limit > maxPerPage {
			limit = maxPerPage
		}

		// Discord APIでメッセージを取得
		messages, err := session.ChannelMessages(channelID, limit, beforeID, "", "")
		if err != nil {
			return "", fmt.Errorf("チャンネルメッセージの取得に失敗: %w", err)
		}

		// メッセージがない場合は終了
		if len(messages) == 0 {
			break
		}

		// 取得したメッセージを追加
		allMessages = append(allMessages, messages...)

		// 次のページのための最古のメッセージIDを記録
		beforeID = messages[len(messages)-1].ID
		remaining -= len(messages)

		// 取得件数が想定より少ない場合（これ以上メッセージがない）は終了
		if len(messages) < limit {
			break
		}
	}

	// ユーザー別の投稿数をカウント
	userPostCount := make(map[string]int)
	for _, message := range allMessages {
		username := message.Author.Username
		userPostCount[username]++
	}

	// ランキング用にソート（投稿数の降順）
	var rankings []UserRankEntry
	for username, count := range userPostCount {
		rankings = append(rankings, UserRankEntry{
			Username: username,
			Count:    count,
		})
	}

	// 投稿数で降順ソート
	sort.Slice(rankings, func(i, j int) bool {
		return rankings[i].Count > rankings[j].Count
	})

	// 全投稿数を計算
	totalMessages := len(allMessages)

	// Ruby版と同じランダムなコメントを選択
	randomComments := []string{
		"人生は有意義にね！",
		"楽しそうだね！",
		"目指すならトップだよね！",
		"ねえねえ、仕事は？",
		fmt.Sprintf("最新%d件の結果だよ！", totalMessages),
		"この人たちに話しかけよう！",
		"他にやることないんだね〜",
		"かわいいね！",
		"これが最強戦士……！",
	}
	randomComment := randomComments[rand.Intn(len(randomComments))]

	// メッセージを構築
	message := fmt.Sprintf("ヒマな人ランキング in <#%s> だよ！ %s\n", channelID, randomComment)

	// 上位5位までを表示
	topFive := rankings
	if len(topFive) > 5 {
		topFive = rankings[:5]
	}

	for i, entry := range topFive {
		// 順位に応じた絵文字を追加
		var rankEmoji string
		switch i {
		case 0:
			rankEmoji = ":first_place: "
		case 1:
			rankEmoji = ":second_place: "
		case 2:
			rankEmoji = ":third_place: "
		default:
			rankEmoji = ""
		}

		// パーセンテージを計算（小数点以下2桁）
		percentage := float64(entry.Count) / float64(totalMessages) * 100

		// ランキング行を追加
		message += fmt.Sprintf("%s%d. %s (%d: %.2f%%)\n",
			rankEmoji, i+1, entry.Username, entry.Count, percentage)
	}

	return message, nil
}
