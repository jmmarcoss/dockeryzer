package ai

import (
	"fmt"
	"strings"
)

type ProviderType string

const (
	ProviderGemini ProviderType = "gemini"
	ProviderOpenAI ProviderType = "openai"
	// ProviderClaude ProviderType = "claude"
)

// ProviderConfig holds configuration for AI provider
type ProviderConfig struct {
	Type   ProviderType
	APIKey string
	Model  string
}

// NewAIProvider creates a new AI provider based on the config
func NewAIProvider(config ProviderConfig) (AIProvider, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("API key is required")
	}

	switch strings.ToLower(string(config.Type)) {
	case string(ProviderGemini):
		return NewGeminiProvider(config.APIKey, config.Model)
	case string(ProviderOpenAI):
		return NewOpenAIProvider(config.APIKey, config.Model)
	// case string(ProviderClaude):
	// 	return NewClaudeProvider(config.APIKey, config.Model)
	default:
		return nil, fmt.Errorf("unsupported provider type: %s", config.Type)
	}
}
