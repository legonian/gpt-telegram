package main

import (
	"context"
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := app()
	if err != nil {
		log.Fatalf("app: %v", err)
	}
}

func app() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("godotenv.Load: %w", err)
	}
	chatGPT := NewChatGPT(os.Getenv("OPEN_AI_KEY"))

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		return fmt.Errorf("tgbotapi.NewBotAPI: %w", err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := bot.GetUpdatesChan(updateConfig)

	ctx := context.Background()
	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] (%v) '%v'",
			update.Message.From.UserName,
			update.UpdateID,
			update.Message.Text,
		)

		reply, err := chatGPT.GenerateResponse(ctx, update.Message.Text)
		if err != nil {
			log.Printf("(%v) %v", update.UpdateID, err)
			reply = "GPT error"
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		msg.ReplyToMessageID = update.Message.MessageID

		if _, err := bot.Send(msg); err != nil {
			log.Printf("(%v) %v", update.UpdateID, err)
		}
	}

	return nil
}
