#!/bin/bash

export API_KEY="your-api-here"

export OPENAI_API_KEY="your-openai-key-here"
export GEMINI_API_KEY="your-gemini-key-here"

LDFLAGS="-X github.com/jorgevvs2/dockeryzer/src/config.APIKey=$API_KEY"

if [ -n "$OPENAI_API_KEY" ]; then
  LDFLAGS="$LDFLAGS -X github.com/jorgevvs2/dockeryzer/src/config.OpenAIKey=$OPENAI_API_KEY"
fi

if [ -n "$GEMINI_API_KEY" ]; then
  LDFLAGS="$LDFLAGS -X github.com/jorgevvs2/dockeryzer/src/config.GeminiKey=$GEMINI_API_KEY"
fi

go build -ldflags "$LDFLAGS" -o dockeryzer

echo "✅ Build finalizado"
