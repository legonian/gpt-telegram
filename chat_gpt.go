package main

import (
	"context"
	"fmt"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

type ChatGPT struct {
	modes  []ChatMode
	client *openai.Client
}

type ChatMode interface {
	Prefix() string
	HandlePrompt(string) string
	HandleResponse(string) string
}

func NewChatGPT(apiKey string) (*ChatGPT, error) {
	defaultMode := NewDefaultMode()

	basedMode, err := NewBasedMode()
	if err != nil {
		return nil, fmt.Errorf("NewBasedMode: %w", err)
	}

	antiMode, err := NewAntiMode()
	if err != nil {
		return nil, fmt.Errorf("NewAntiMode: %w", err)
	}

	return &ChatGPT{
		client: openai.NewClient(apiKey),
		modes: []ChatMode{
			defaultMode,
			basedMode,
			antiMode,
		},
	}, nil
}

func (cg *ChatGPT) GenerateResponse(ctx context.Context, prompt string) (string, error) {
	var responseHandler func(string) string
	for _, m := range cg.modes {
		prefix := m.Prefix()
		if !strings.HasPrefix(prompt, prefix) {
			continue
		}
		prompt = strings.TrimPrefix(prompt, prefix)
		prompt = m.HandlePrompt(prompt)

		responseHandler = m.HandleResponse
	}
	if responseHandler == nil {
		return "", fmt.Errorf("no prefix")
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

	responce := resp.Choices[0].Message.Content
	responce = responseHandler(responce)

	return responce, nil
}
