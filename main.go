package main

import (
	"github.com/jorgevvs2/dockeryzer/src/ai"
	"github.com/jorgevvs2/dockeryzer/src/cmd"
)

func main() {
	ai.InitAIEnv()
	cmd.Execute()
}
