package ai

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAIProvider_APIKeyObrigatoria(t *testing.T) {
	config := ProviderConfig{
		Type:   ProviderOpenAI,
		APIKey: "",
		Model:  "gpt-4",
	}

	provider, err := NewAIProvider(config)

	assert.Nil(t, provider)
	assert.Error(t, err)
	assert.Equal(t, "API key is required", err.Error())
}

func TestNewAIProvider_OpenAIComSucesso(t *testing.T) {
	config := ProviderConfig{
		Type:   ProviderOpenAI,
		APIKey: "fake-key",
		Model:  "gpt-4",
	}

	provider, err := NewAIProvider(config)

	assert.NoError(t, err)
	assert.NotNil(t, provider)
}

func TestNewAIProvider_GeminiComSucesso(t *testing.T) {
	config := ProviderConfig{
		Type:   ProviderGemini,
		APIKey: "fake-key",
		Model:  "gemini-pro",
	}

	provider, err := NewAIProvider(config)

	assert.NoError(t, err)
	assert.NotNil(t, provider)
}

func TestNewAIProvider_TipoNaoSuportado(t *testing.T) {
	config := ProviderConfig{
		Type:   ProviderType("unknown"),
		APIKey: "key",
		Model:  "model-x",
	}

	provider, err := NewAIProvider(config)

	assert.Nil(t, provider)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported provider type")
}
