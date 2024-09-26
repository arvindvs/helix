package core

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/sashabaranov/go-openai"
)

type Assistant struct {
	client *openai.Client
}

func NewAssistant() (*Assistant, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY is not set")
	}

	client := openai.NewClient(apiKey)

	return &Assistant{
		client: client,
	}, nil
}

func (a *Assistant) ProcessText(ctx context.Context, input string) (string) {
	resp, err := a.client.CreateChatCompletionStream(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT4o,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: input,
				},
			},
			Stream: true,
		},
	)

	if err != nil {
		fmt.Printf("error calling OpenAI API: %v\n", err)
		return ""
	}
	defer resp.Close()

	var totalResponse string

	fmt.Println("GPT response:")
	for {
		select {
		case <-ctx.Done():
			return ""
		default:
			response, err := resp.Recv()

			if errors.Is(err, io.EOF) {
				fmt.Println("\nYour input:")
				return totalResponse
			}

			if err != nil {
				slog.Error("Stream error:", "error", err.Error())
				return ""
			}
			
			text := response.Choices[0].Delta.Content
			totalResponse += text

			fmt.Print(text)
		}	
	}
}
