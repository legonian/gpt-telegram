package main

type DefaultMode struct {
	prefix string
	prompt string
}

func NewDefaultMode() *DefaultMode {
	return &DefaultMode{
		prefix: "gpt ",
		prompt: "",
	}
}

func (m *DefaultMode) Prefix() string {
	return m.prefix
}

func (m *DefaultMode) HandleResponse(responce string) string {
	return responce
}

func (m *DefaultMode) HandlePrompt(prompt string) string {
	return m.prompt + prompt
}
