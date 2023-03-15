package main

import (
	"context"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

const defaultmessageLimit = 1000

type ChatGPT struct {
	messageLimit int
	client       *openai.Client
}

func NewChatGPT(apiKey string) *ChatGPT {
	return &ChatGPT{
		messageLimit: defaultmessageLimit,
		client:       openai.NewClient(apiKey),
	}
}

func (cg *ChatGPT) GenerateResponse(ctx context.Context, prompt string) (string, error) {
	if cg.messageLimit < len(prompt) {
		return "", fmt.Errorf("prompt is to big: %v", len(prompt))
	}

	resp, err := cg.client.CreateChatCompletion(ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		return "", fmt.Errorf("cg.client.CreateChatCompletion: %w", err)
	}

	return resp.Choices[0].Message.Content, nil
}
