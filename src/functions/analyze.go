package functions

import (
	"github.com/jorgevvs2/dockeryzer/src/utils"
)

func Analyze(name string) {
	imageInspect := utils.GetDockerImageInspectByIdOrName(name)
	utils.DebugImageInfo(imageInspect) // Adicione esta linha
	utils.PrintImageAnalyzeResults(name, imageInspect)
}
