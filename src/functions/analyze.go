package functions

import (
	"fmt"
	"os"

	"github.com/jorgevvs2/dockeryzer/src/security"
	"github.com/jorgevvs2/dockeryzer/src/utils"
)

func AnalyzeImage(name string) {
	imageInspect := utils.GetDockerImageInspectByIdOrName(name)
	utils.PrintImageAnalyzeResults(name, imageInspect)
}

func AnalyzeDockerfile(path string) {
	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Failed to read Dockerfile:", err)
		return
	}

	analyzer := security.NewCISAnalyzer()
	results := analyzer.Analyze(string(content))

	security.PrintCISResults(results)
}
