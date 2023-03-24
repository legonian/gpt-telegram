package main

import (
	"fmt"
	"os"
	"strings"
)

type Imam4Mode struct {
	prefix string
	prompt string
}

func NewImam4Mode() (*Imam4Mode, error) {
	prompt, err := os.ReadFile("./prompt_imam4.txt")
	if err != nil {
		return nil, fmt.Errorf("os.ReadFile: %w", err)
	}

	return &Imam4Mode{
		prefix: "Artem ",
		prompt: string(prompt),
	}, nil
}

func (m *Imam4Mode) Prefix() string {
	return m.prefix
}

func (m *Imam4Mode) HandleResponse(responce string) string {
	divided := strings.Split(responce, "[Artem]: ")

	if len(divided) == 2 {
		responce = divided[1]
	}

	return responce
}

func (m *Imam4Mode) HandlePrompt(prompt string) string {
	return m.prompt + " " + prompt
}
