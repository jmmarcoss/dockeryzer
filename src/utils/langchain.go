package utils

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jorgevvs2/dockeryzer/src/config"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

func CreateDockerfileWithLangChain(ignoreComments bool) {

	apiKey := config.APIKey
	if apiKey == "" {
		log.Fatal("API key not set in binary. Please rebuild with -ldflags")
	}

	projectTree, err := GetProjectStructure()
	if err != nil {
		log.Fatal(err)
	}

	tech := DetectProjectSmart(apiKey)

	fmt.Printf("🔍 Detected: %s", tech.Language)
	if tech.Framework != "" {
		fmt.Printf(" (%s)", tech.Framework)
	}
	if tech.PackageManager != "" {
		fmt.Printf(" [%s]", tech.PackageManager)
	}
	fmt.Println()

	prompt := BuildDockerfilePrompt(projectTree, ignoreComments)

	llm, err := openai.New(
		openai.WithToken(os.Getenv("OPENAI_API_KEY")),
		openai.WithModel("gpt-4.1-mini"),
	)
	if err != nil {
		log.Fatal(err)
	}

	response, err := llms.GenerateFromSinglePrompt(
		context.Background(),
		llm,
		prompt,
	)

	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("Dockeryzer.Dockerfile", []byte(response), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
