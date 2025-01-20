package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"google.golang.org/api/option"
)

func main() {
	botToken := os.Getenv("TG_BOT_TOKEN")

	bot, err := telego.NewBot(botToken, telego.WithDefaultDebugLogger())
	if err != nil {
		log.Fatal(err)
	}

	updates, err := bot.UpdatesViaLongPolling(nil)
	if err != nil {
		log.Fatal(err)
	}

	bh, err := th.NewBotHandler(bot, updates)
	if err != nil {
		log.Fatal(err)
	}
	defer bh.Stop()
	defer bot.StopLongPolling()

	userMessages := make(map[int64][]string)

	// Register new handler with match on command `/start`
	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		// Send message
		_, _ = bot.SendMessage(tu.Message(
			tu.ID(update.Message.Chat.ID),
			fmt.Sprintf("Hello %s!", update.Message.From.FirstName),
		))
	}, th.CommandEqual("start"))

	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		ctx := context.Background()

		apiKey, ok := os.LookupEnv("GEMINI_API_KEY")
		if !ok {
			log.Fatalln("Environment variable GEMINI_API_KEY not set")
		}

		client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
		if err != nil {
			log.Fatalf("Error creating client: %v", err)
		}
		defer client.Close()

		model := client.GenerativeModel("gemini-2.0-flash-exp")

		model.SetTemperature(1)
		model.SetTopK(40)
		model.SetTopP(0.95)
		model.SetMaxOutputTokens(8192)
		model.ResponseMIMEType = "text/plain"

		session := model.StartChat()
		session.History = []*genai.Content{}

		resp, err := session.SendMessage(ctx, genai.Text("INSERT_INPUT_HERE"))
		if err != nil {
			log.Fatalf("Error sending message: %v", err)
		}

		for _, part := range resp.Candidates[0].Content.Parts {
			fmt.Printf("%v\n", part)
		}
	}, th.CommandEqual("summary"))

	bh.HandleMessage(func(bot *telego.Bot, msg telego.Message) {
		text := "your message was captured"

		userMessages[msg.From.ID] = append(userMessages[msg.From.ID], msg.Text)
		_, _ = bot.SendMessage(tu.Message(msg.Chat.ChatID(), text))
	})

	bh.Start()
}
