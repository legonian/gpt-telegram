package main

import (
	"fmt"
	"os"
)

type ImamMode struct {
	prefix string
	prompt string
}

func NewImamMode() (*ImamMode, error) {
	prompt, err := os.ReadFile("./prompt_imam.txt")
	if err != nil {
		return nil, fmt.Errorf("os.ReadFile: %w", err)
	}

	return &ImamMode{
		prefix: "imam ",
		prompt: string(prompt),
	}, nil
}

func (m *ImamMode) Prefix() string {
	return m.prefix
}

func (m *ImamMode) HandleResponse(responce string) string {
	return responce
}

func (m *ImamMode) HandlePrompt(prompt string) string {
	return m.prompt + `"` + prompt + `"`
}
