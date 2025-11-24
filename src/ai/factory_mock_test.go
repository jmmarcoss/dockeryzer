//go:build test

package ai

import "github.com/stretchr/testify/mock"

type MockAIProvider struct {
	mock.Mock
}

// Essas funções substituem as reais SOMENTE em testes
func NewGeminiProvider(apiKey, model string) (AIProvider, error) {
	return &MockAIProvider{}, nil
}

func NewOpenAIProvider(apiKey, model string) (AIProvider, error) {
	return &MockAIProvider{}, nil
}
