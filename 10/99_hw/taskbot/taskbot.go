package main

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"net/http"
	handlers "taskbot/handlers/telegram"
)

var (
	BotToken   string
	WebhookURL string
)

func startTaskBot(ctx context.Context, httpListenAddr string) error {
	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		return err
	}

	_, err = bot.SetWebhook(tgbotapi.NewWebhook(WebhookURL))
	if err != nil {
		return err
	}

	telegramHandler := handlers.TelegramBotHandler{
		Bot: bot,
	}
	telegramHandler.SetDefault()

	updates := bot.ListenForWebhook("/")

	go func() {
		err := http.ListenAndServe(httpListenAddr, nil)
		if err != nil {

		}
	}()

	for {
		select {
		case <-ctx.Done():
			break
		case update := <-updates:
			telegramHandler.Route(update)
		}
	}
}

func main() {
	err := startTaskBot(context.Background(), ":8081")
	if err != nil {
		log.Fatalln(err)
	}
}

// это заглушка чтобы импорт сохранился
func __dummy() {
	tgbotapi.APIEndpoint = "_dummy"
}
