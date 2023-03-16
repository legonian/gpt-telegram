package main

import (
	"fmt"
	"os"
	"strings"
)

type BasedMode struct {
	prefix string
	prompt string
}

func NewBasedMode() (*BasedMode, error) {
	prompt, err := os.ReadFile("./prompt_based.txt")
	if err != nil {
		return nil, fmt.Errorf("os.ReadFile: %w", err)
	}

	return &BasedMode{
		prefix: "basedgpt ",
		prompt: string(prompt),
	}, nil
}

func (m *BasedMode) Prefix() string {
	return m.prefix
}

func (m *BasedMode) HandleResponse(responce string) string {
	divided := strings.Split(responce, "BasedGPT: ")

	if len(divided) == 2 {
		responce = divided[1]
	}

	return responce
}

func (m *BasedMode) HandlePrompt(prompt string) string {
	return m.prompt + prompt
}
