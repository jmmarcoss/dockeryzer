package ai

import (
	"os"

	"github.com/jorgevvs2/dockeryzer/src/config"
)

func InitAIEnv() {
	if config.OpenAIKey != "" {
		os.Setenv("OPENAI_API_KEY", config.OpenAIKey)
	}

	if config.GeminiKey != "" {
		os.Setenv("GEMINI_API_KEY", config.GeminiKey)
	}
}
