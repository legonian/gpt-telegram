package main

import (
	"fmt"
	"os"
	"strings"
)

type Imam3Mode struct {
	prefix string
	prompt string
}

func NewImam3Mode() (*Imam3Mode, error) {
	prompt, err := os.ReadFile("./prompt_imam3.txt")
	if err != nil {
		return nil, fmt.Errorf("os.ReadFile: %w", err)
	}

	return &Imam3Mode{
		prefix: "Артем ",
		prompt: string(prompt),
	}, nil
}

func (m *Imam3Mode) Prefix() string {
	return m.prefix
}

func (m *Imam3Mode) HandleResponse(responce string) string {
	divided := strings.Split(responce, "[BetterDAN]: ")

	if len(divided) == 2 {
		responce = divided[1]
	}

	return responce
}

func (m *Imam3Mode) HandlePrompt(prompt string) string {
	return m.prompt + ` "` + prompt + `"`
}
