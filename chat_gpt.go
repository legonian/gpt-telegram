package main

import (
	"context"
	"errors"
	"fmt"
	"log"
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

var ErrInvalidPrompt = errors.New("invalid prompt")
var ErrAPI = fmt.Errorf("API error")

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

	imamMode, err := NewImamMode()
	if err != nil {
		return nil, fmt.Errorf("NewImamMode: %w", err)
	}

	return &ChatGPT{
		client: openai.NewClient(apiKey),
		modes: []ChatMode{
			defaultMode,
			basedMode,
			antiMode,
			imamMode,
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
		return "", ErrInvalidPrompt
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
		log.Printf("cg.client.CreateChatCompletion: %v", err)
		return "", ErrAPI
	}

	responce := resp.Choices[0].Message.Content
	responce = responseHandler(responce)

	return responce, nil
}
