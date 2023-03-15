package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

const promptPrefixFile = "./prompt_prefix.txt"

type ChatGPT struct {
	promptPrefix string
	client       *openai.Client
}

func NewChatGPT(apiKey string) (*ChatGPT, error) {
	promptPrefix, err := os.ReadFile(promptPrefixFile)
	if err != nil {
		return nil, fmt.Errorf("os.ReadFile: %w", err)
	}

	return &ChatGPT{
		promptPrefix: string(promptPrefix),
		client:       openai.NewClient(apiKey),
	}, nil
}

func (cg *ChatGPT) GenerateResponse(ctx context.Context, prompt string) (string, error) {
	isBased := false
	switch {
	case strings.HasPrefix(prompt, "gpt "):
		prompt = strings.TrimPrefix(prompt, "gpt ")
	case strings.HasPrefix(prompt, "basedgpt "):
		isBased = true
		prompt = strings.TrimPrefix(prompt, "basedgpt ")
		prompt = cg.promptPrefix + prompt
	default:
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

	if isBased {
		divided := strings.Split(responce, "BasedGPT: ")

		if len(divided) == 2 {
			responce = divided[1]
		}
	}

	return responce, nil
}
