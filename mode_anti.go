package main

import (
	"fmt"
	"os"
	"strings"
)

type AntiMode struct {
	prefix string
	prompt string
}

func NewAntiMode() (*AntiMode, error) {
	prompt, err := os.ReadFile("./prompt_anti.txt")
	if err != nil {
		return nil, fmt.Errorf("os.ReadFile: %w", err)
	}

	return &AntiMode{
		prefix: "antigpt ",
		prompt: string(prompt),
	}, nil
}

func (m *AntiMode) Prefix() string {
	return m.prefix
}

func (m *AntiMode) HandleResponse(responce string) string {
	divided := strings.Split(responce, "[AntiGPT]: ")

	if len(divided) == 2 {
		responce = divided[1]
	}

	return responce
}

func (m *AntiMode) HandlePrompt(prompt string) string {
	return m.prompt + prompt
}
