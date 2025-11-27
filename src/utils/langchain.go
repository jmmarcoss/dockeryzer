package utils

import (
	"context"
	"log"
	"os"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

func CreateDockerfileWithLangChain(ignoreComments bool) {

	projectTree, err := GetProjectStructure()
	if err != nil {
		log.Fatal(err)
	}

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
