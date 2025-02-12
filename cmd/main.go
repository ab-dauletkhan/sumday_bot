package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/ab-dauletkhan/sumday_bot/internal/bot"
	"github.com/ab-dauletkhan/sumday_bot/internal/repository"
	"github.com/ab-dauletkhan/sumday_bot/internal/service"
)

func main() {
	// Initialize logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Load environment variables
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Fatal("BOT_TOKEN environment variable is required")
	}

	geminiAPIKey := os.Getenv("GEMINI_API_KEY")
	if geminiAPIKey == "" {
		log.Fatal("GEMINI_API_KEY environment variable is required")
	}

	repo := repository.NewMapRepo()

	llmService := service.NewLLMService(repo)

	telegramBot, err := bot.NewBot(botToken, repo, llmService)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	logger.Info("Bot is starting...")

	if err := telegramBot.Start(); err != nil {
		log.Fatalf("Bot stopped unexpectedly: %v", err)
	}
}
