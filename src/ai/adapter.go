package ai

import "context"

type AIProvider interface {
	GenerateContent(ctx context.Context, systemPrompt, userPrompt string, temperature float32) (string, error)
	Close() error
}
