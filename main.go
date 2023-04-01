package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"

	"gpt-telegram/quran"
)

func main() {
	log.Printf("Starting...")

	err := godotenv.Load()
	if err != nil {
		log.Printf("godotenv.Load: %v", err)
	}

	err = app()
	if err != nil {
		log.Fatalf("app: %v", err)
	}
}

func app() error {
	q, err := quran.NewQuran()
	if err != nil {
		return fmt.Errorf("quran.NewQuran: %w", err)
	}
	chatGPT, err := NewChatGPT(os.Getenv("OPEN_AI_KEY"))
	if err != nil {
		return fmt.Errorf("NewChatGPT: %w", err)
	}

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
		if update.Message == nil || update.Message.Text == "" {
			continue
		}
		messageText := update.Message.Text

		log.Printf("[%s] (%v) '%v'",
			update.Message.From.UserName,
			update.UpdateID,
			messageText,
		)
		var reply string
		switch messageText {
		case "/random":
			reply = q.RandomAyah(2)
		default:
			reply, err = chatGPT.GenerateResponse(ctx, messageText)
			if err != nil {
				if !errors.Is(err, ErrAPI) {
					continue
				}
				reply = "API error"
			}
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		msg.ReplyToMessageID = update.Message.MessageID

		if _, err := bot.Send(msg); err != nil {
			log.Printf("(%v) %v", update.UpdateID, err)
		}
	}

	return nil
}
