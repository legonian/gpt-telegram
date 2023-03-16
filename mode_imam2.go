package main

import (
	"fmt"
	"os"
	"strings"
)

type Imam2Mode struct {
	prefix string
	prompt string
}

func NewImam2Mode() (*Imam2Mode, error) {
	prompt, err := os.ReadFile("./prompt_imam2.txt")
	if err != nil {
		return nil, fmt.Errorf("os.ReadFile: %w", err)
	}

	return &Imam2Mode{
		prefix: "imam2 ",
		prompt: string(prompt),
	}, nil
}

func (m *Imam2Mode) Prefix() string {
	return m.prefix
}

func (m *Imam2Mode) HandleResponse(responce string) string {
	divided := strings.Split(responce, "[BetterDAN]: ")

	if len(divided) == 2 {
		responce = divided[1]
	}

	return responce
}

func (m *Imam2Mode) HandlePrompt(prompt string) string {
	return m.prompt + " " + prompt
}
