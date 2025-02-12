package service

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/ab-dauletkhan/sumday_bot/internal/repository"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type LLMService struct {
	repo *repository.MapRepo
}

func NewLLMService(mr *repository.MapRepo) *LLMService {
	return &LLMService{repo: mr}
}

func (s *LLMService) GenerateSummary(ctx context.Context, userID int64) (string, error) {
	apikey, ok := os.LookupEnv("GEMINI_API_KEY")
	if !ok {
		return "", fmt.Errorf("GEMINI_API_KEY not set")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apikey))
	if err != nil {
		return "", fmt.Errorf("error creating client: %w", err)
	}
	defer client.Close()

	messages := s.repo.GetMessages(userID)
	if len(messages) == 0 {
		return "No messages to summarize", nil
	}

	var input strings.Builder
	input.WriteString("\n")

	for _, msg := range messages {
		input.WriteString(fmt.Sprintf("[%s] %s\n", msg.Timestamp.Format("15:04"), msg.Text))
	}

	model := client.GenerativeModel("gemini-2.0-pro-exp-02-05")
	resp, err := model.GenerateContent(ctx, genai.Text(input.String()))
	if err != nil {
		return "", fmt.Errorf("error sending message: %w", err)
	}

	return fmt.Sprintln(resp.Candidates[0].Content.Parts[0]), nil
}
