package utils

import (
	"fmt"
	"math"
	"time"

	"github.com/docker/docker/api/types"
)

func GetImageSizeInMBs(imageInspect types.ImageInspect) float32 {
	sizeInMbs := float32(imageInspect.Size) / float32(math.Pow(10.0, 6))
	return sizeInMbs
}

func GetImageSizeString(imageInspect types.ImageInspect) string {
	sizeUnit := "MB"
	sizeInMbs := float32(imageInspect.Size) / float32(math.Pow(10.0, 6))
	sizeInGbs := float32(0.0)

	finalSize := sizeInMbs
	isMoreThanOneGb := sizeInMbs > 1000
	if isMoreThanOneGb {
		sizeUnit = "GB"
		sizeInGbs = sizeInMbs / float32(math.Pow(10.0, 3))
		finalSize = sizeInGbs
	}

	return fmt.Sprintf("%.2f %s", finalSize, sizeUnit)
}

func GetImageNumberOfLayers(imageInspect types.ImageInspect) int {
	return len(imageInspect.RootFS.Layers)
}

func GetImageFormattedCreationDate(imageInspect types.ImageInspect) string {
	parsedTime, err := time.Parse(time.RFC3339Nano, imageInspect.Created)
	if err != nil {
		fmt.Println("Failed to parsing date:", err)
		return ""
	}

	return parsedTime.Format("02 Jan 2006")
}

func GetImageAuthor(imageInspect types.ImageInspect) string {
	if imageInspect.Author == "" {
		return "<none>"
	}
	return imageInspect.Author
}

func GetImageSizeWithColor(imageInspect types.ImageInspect) string {
	sizeInMBs := GetImageSizeInMBs(imageInspect)

	fmt.Printf("  - Size: ")
	if sizeInMBs < 250 {
		return SuccessSprintf("%s", GetImageSizeString(imageInspect))
	}

	if sizeInMBs >= 250 && sizeInMBs <= 500 {
		return WarningSprintf("%s", GetImageSizeString(imageInspect))
	}

	return ErrorSprintf("%s", GetImageSizeString(imageInspect))
}

func GetImageLayersWithColor(imageInspect types.ImageInspect) string {
	numberOfLayers := GetImageNumberOfLayers(imageInspect)

	fmt.Printf("  - N. of Layers: ")
	if numberOfLayers < 10 {
		return SuccessSprintf("%d", numberOfLayers)
	}

	if numberOfLayers >= 10 && numberOfLayers <= 20 {
		return WarningSprintf("%d", numberOfLayers)
	}

	return ErrorSprintf("%d", numberOfLayers)
}

func PrintImageResults(name string, imageInspect types.ImageInspect, minimal bool, ignoreSuggestions bool) {
	fmt.Printf("Details of image ")
	BoldPrintf("%s:\n", name)
	fmt.Printf("  - Tags: %s\n", imageInspect.RepoTags)
	fmt.Println(GetImageSizeWithColor(imageInspect))
	fmt.Println(GetImageLayersWithColor(imageInspect))

	// Nova função para detectar a linguagem principal
	PrintLanguageWithColor(imageInspect)

	if !minimal {
		fmt.Printf("  - Author: %s\n", GetImageAuthor(imageInspect))
		fmt.Printf("  - Creation date: %s\n", GetImageFormattedCreationDate(imageInspect))
		fmt.Printf("  - OS: %s\n", imageInspect.Os)
	}

	sizeInMBs := GetImageSizeInMBs(imageInspect)
	numberOfLayers := GetImageNumberOfLayers(imageInspect)
	hasOutdatedLanguage := HasOutdatedLanguage(imageInspect)

	isBigImage := sizeInMBs > 250
	hasManyLayers := numberOfLayers > 10

	shouldShowSuggestions := isBigImage || hasManyLayers || hasOutdatedLanguage

	if ignoreSuggestions {
		return
	}

	if shouldShowSuggestions {
		fmt.Println("\n Improvement suggestions:")
	}

	if isBigImage {
		fmt.Println("  - Consider reducing the size of your image. Try using smaller base images and ensure that no unnecessary files are included.")
	}

	if hasManyLayers {
		fmt.Println("  - Your image has multiple layers. Consider applying a multi-build stage strategy or combining commands to reduce the number of layers.")
	}

	// Sugestões específicas por linguagem
	languageSuggestions := GetLanguageImprovementSuggestions(imageInspect)
	for _, suggestion := range languageSuggestions {
		fmt.Println(suggestion)
	}

	// Se nenhuma linguagem foi detectada
	lang := DetectPrimaryLanguage(imageInspect)
	if lang == nil && !ignoreSuggestions && shouldShowSuggestions {
		fmt.Println("  - No programming language runtime detected. Ensure your image is configured correctly if it requires a runtime environment.")
	}
}

func PrintImageAnalyzeResults(name string, imageInspect types.ImageInspect) {
	PrintImageResults(name, imageInspect, false, false)
}

func DebugImageInfo(imageInspect types.ImageInspect) {
	fmt.Println("\n=== DEBUG IMAGE INFO ===")

	fmt.Println("\nEnvironment Variables:")
	if len(imageInspect.Config.Env) == 0 {
		fmt.Println("  (empty)")
	} else {
		for i, env := range imageInspect.Config.Env {
			fmt.Printf("  [%d] %s\n", i, env)
		}
	}

	fmt.Println("\nCmd:")
	if len(imageInspect.Config.Cmd) == 0 {
		fmt.Println("  (empty)")
	} else {
		for i, cmd := range imageInspect.Config.Cmd {
			fmt.Printf("  [%d] %s\n", i, cmd)
		}
	}

	fmt.Println("\nEntrypoint:")
	if len(imageInspect.Config.Entrypoint) == 0 {
		fmt.Println("  (empty)")
	} else {
		for i, ep := range imageInspect.Config.Entrypoint {
			fmt.Printf("  [%d] %s\n", i, ep)
		}
	}

	fmt.Printf("\nWorking Directory: %s\n", imageInspect.Config.WorkingDir)

	fmt.Println("\nLabels:")
	if len(imageInspect.Config.Labels) == 0 {
		fmt.Println("  (empty)")
	} else {
		for key, value := range imageInspect.Config.Labels {
			fmt.Printf("  %s = %s\n", key, value)
		}
	}

	fmt.Println("\nImage Architecture:", imageInspect.Architecture)
	fmt.Println("Image OS:", imageInspect.Os)

	fmt.Println("\n=== END DEBUG ===\n")
}

func PrintImageCompareResults(name string, imageInspect types.ImageInspect) {
	PrintImageResults(name, imageInspect, true, true)
}

func PrintImageCompareLayersResults(image1 string, image1Inspect types.ImageInspect, image2 string, image2Inspect types.ImageInspect) {
	numberOfLayers1 := len(image1Inspect.RootFS.Layers)
	numberOfLayers2 := len(image2Inspect.RootFS.Layers)

	minorImage := image1
	minorLayers := numberOfLayers1
	if numberOfLayers2 < numberOfLayers1 {
		minorImage = image2
		minorLayers = numberOfLayers2
	}

	biggerImage := image1
	mostLayers := numberOfLayers1
	if numberOfLayers2 > numberOfLayers1 {
		biggerImage = image2
		mostLayers = numberOfLayers2
	}

	layersDiff := numberOfLayers1 - numberOfLayers2
	if numberOfLayers2 > numberOfLayers1 {
		layersDiff = numberOfLayers2 - numberOfLayers1
	}

	if layersDiff == 0 {
		fmt.Printf("  - Images have the same number of layers: %d\n", numberOfLayers2)
		return
	}
	fmt.Printf("  - Image ")
	SuccessPrintf("%s", minorImage)
	fmt.Printf(" has ")
	SuccessPrintf("%d less layers", layersDiff)
	fmt.Printf(" than image ")
	ErrorPrintf("%s", biggerImage)
	fmt.Printf(" (")
	SuccessPrintf("%d", minorLayers)
	fmt.Printf(" < ")
	ErrorPrintf("%d", mostLayers)
	fmt.Println(").")
}

func PrintImageCompareSizeResults(image1 string, image1Inspect types.ImageInspect, image2 string, image2Inspect types.ImageInspect) {
	size1 := image1Inspect.Size
	size2 := image2Inspect.Size

	size1String := GetImageSizeString(image1Inspect)
	size2String := GetImageSizeString(image2Inspect)

	minorImage := image1
	minorImageString := size1String
	minorSize := size1
	if size2 < size1 {
		minorImage = image2
		minorImageString = size2String
		minorSize = size2
	}

	biggerImage := image1
	biggerImageString := size1String
	biggerSize := size1
	if size2 > size1 {
		biggerImage = image2
		biggerImageString = size2String
		biggerSize = size2
	}

	sizeDiff := size1 - size2
	if size2 > size1 {
		sizeDiff = size2 - size1
	}

	if sizeDiff == 0 {
		fmt.Printf("  - Images have the same size: %s\n", GetImageSizeString(image1Inspect))
		return
	}

	percent := 100 - (float32(minorSize)/float32(biggerSize))*100

	fmt.Printf("  - Image ")
	SuccessPrintf("%s", minorImage)
	fmt.Printf(" is ")
	SuccessPrintf("%.2f%% smaller", percent)
	fmt.Printf(" than image ")
	ErrorPrintf("%s", biggerImage)
	fmt.Printf(" (")
	SuccessPrintf(minorImageString)
	fmt.Printf(" < ")
	ErrorPrintf(biggerImageString)
	fmt.Println(").")
}

func PrintImageCompareLanguageResults(image1 string, image1Inspect types.ImageInspect, image2 string, image2Inspect types.ImageInspect) {
	lang1 := DetectPrimaryLanguage(image1Inspect)
	lang2 := DetectPrimaryLanguage(image2Inspect)

	if lang1 == nil && lang2 == nil {
		fmt.Println("  - No programming language runtime detected in either image.")
		return
	}

	if lang1 == nil {
		fmt.Printf("  - Only image ")
		SuccessPrintf("%s", image2)
		fmt.Printf(" has detected language runtime: %s %s\n", lang2.Name, lang2.Version)
		return
	}

	if lang2 == nil {
		fmt.Printf("  - Only image ")
		SuccessPrintf("%s", image1)
		fmt.Printf(" has detected language runtime: %s %s\n", lang1.Name, lang1.Version)
		return
	}

	// Comparar as linguagens
	if lang1.Name != lang2.Name {
		fmt.Printf("  - Images use different languages: ")
		fmt.Printf("%s (%s) vs %s (%s)\n",
			lang1.Name, lang1.Version,
			lang2.Name, lang2.Version)
		return
	}

	// Mesma linguagem, comparar versões
	major1 := getMajorVersion(lang1.Version)
	major2 := getMajorVersion(lang2.Version)

	if lang1.Version == lang2.Version {
		fmt.Printf("  - Both images use the same %s version: %s\n", lang1.Name, lang1.Version)
		return
	}

	if major1 == major2 {
		fmt.Printf("  - Both images use %s version %s (minor version may differ)\n",
			lang1.Name, lang1.Version)
		return
	}

	if major1 > major2 {
		fmt.Printf("  - Image ")
		SuccessPrintf("%s", image1)
		fmt.Printf(" uses newer %s (", lang1.Name)
		SuccessPrintf("%s", lang1.Version)
		fmt.Printf(" > ")
		ErrorPrintf("%s", lang2.Version)
		fmt.Println(")")
	} else {
		fmt.Printf("  - Image ")
		SuccessPrintf("%s", image2)
		fmt.Printf(" uses newer %s (", lang2.Name)
		SuccessPrintf("%s", lang2.Version)
		fmt.Printf(" > ")
		ErrorPrintf("%s", lang1.Version)
		fmt.Println(")")
	}
}
