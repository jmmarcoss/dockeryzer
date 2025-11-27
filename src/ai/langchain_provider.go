package ai

import (
	"context"
	"os"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

type LangChainProvider struct {
	llm llms.Model
}

func NewLangChainProvider() (*LangChainProvider, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")

	llm, err := openai.New(
		openai.WithToken(apiKey),
		openai.WithModel("gpt-4.1-mini"),
	)
	if err != nil {
		return nil, err
	}

	return &LangChainProvider{llm: llm}, nil
}

func (p *LangChainProvider) GenerateDockerfile(prompt string) (string, error) {
	ctx := context.Background()

	response, err := llms.GenerateFromSinglePrompt(ctx, p.llm, prompt)
	if err != nil {
		return "", err
	}

	return response, nil
}
