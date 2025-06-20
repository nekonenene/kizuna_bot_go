package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"kizuna_bot_go/internal/bot"
)

func main() {
	// .envファイルから環境変数を読み込み（APIキーなどの設定情報）
	if err := godotenv.Load(); err != nil {
		log.Println(".envファイルが見つからないため、システム環境変数を使用します")
	}

	// Kizuna Botのインスタンスを作成
	kizunaBot, err := bot.NewKizunaBot()
	if err != nil {
		log.Fatalf("ボットの作成に失敗しました: %v", err)
	}

	// Discordサーバーへの接続を開始
	if err := kizunaBot.Start(); err != nil {
		log.Fatalf("ボットの起動に失敗しました: %v", err)
	}

	log.Println("ボットが正常に起動しました。終了するにはCTRL-Cを押してください。")

	// システムからの終了シグナル（CTRL-Cなど）を待機
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Discordとの接続を安全に切断
	kizunaBot.Close()
	log.Println("ボットを正常に終了しました。")
}
