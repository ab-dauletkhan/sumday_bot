package bot

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ab-dauletkhan/sumday_bot/internal/models"
	"github.com/ab-dauletkhan/sumday_bot/internal/repository"
	"github.com/ab-dauletkhan/sumday_bot/internal/service"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

type Bot struct {
	bot        *telego.Bot
	repo       *repository.MapRepo
	llmService *service.LLMService
}

func NewBot(botToken string, repo *repository.MapRepo, llm *service.LLMService) (*Bot, error) {
	telebot, err := telego.NewBot(botToken)
	if err != nil {
		return nil, err
	}

	return &Bot{bot: telebot, repo: repo, llmService: llm}, nil
}

func (b *Bot) Start() error {
	updates, err := b.bot.UpdatesViaLongPolling(nil)
	if err != nil {
		log.Fatal(err)
	}

	bh, err := th.NewBotHandler(b.bot, updates)
	if err != nil {
		log.Fatal(err)
	}
	defer bh.Stop()
	defer b.bot.StopLongPolling()

	// Register `/start` command handler
	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		_, err = bot.SendMessage(tu.Message(
			tu.ID(update.Message.Chat.ID),
			fmt.Sprintf("Hello %s!", update.Message.From.FirstName),
		))
		if err != nil {
			log.Println("error sending message from /start")
		}
	}, th.CommandEqual("start"))

	// Register `/summary` command handler
	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		ctx := context.Background()
		summary, err := b.llmService.GenerateSummary(ctx, update.Message.From.ID)
		if err != nil {
			_, err = b.bot.SendMessage(tu.Message(update.Message.Chat.ChatID(), "Error generating summary."))
			if err != nil {
				log.Println("error sending message from /start")
			}
			return
		}
		_, err = b.bot.SendMessage(tu.Message(update.Message.Chat.ChatID(), summary))
		if err != nil {
			log.Println("error sending message from /start")
		}
	}, th.CommandEqual("summary"))

	// Handle all messages
	bh.HandleMessage(func(bot *telego.Bot, msg telego.Message) {
		text := "Your message was captured"

		b.repo.Save(msg.From.ID, models.Message{
			Timestamp: time.Unix(msg.Date, 0),
			Text:      msg.Text,
		})

		_, err = bot.SendMessage(tu.Message(msg.Chat.ChatID(), text))
		if err != nil {
			log.Println("error sending message from /start")
		}
	})

	bh.Start()
	return nil
}
