package utils

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types/image"
)

type LanguageInfo struct {
	Name    string
	Version string
	Color   string // "success", "warning", "error"
}

// Detecta a linguagem principal da imagem
func DetectPrimaryLanguage(imageInspect image.InspectResponse) *LanguageInfo {
	envVars := imageInspect.Config.Env
	cmd := imageInspect.Config.Cmd
	entrypoint := imageInspect.Config.Entrypoint
	workingDir := imageInspect.Config.WorkingDir

	// Ordem de prioridade baseada em especificidade das variáveis
	// Quanto mais específica a variável, maior a confiança

	// 1. Node.js - verifica NODE_VERSION (muito específico)
	if nodeVersion := detectNodeJSVersion(envVars); nodeVersion != "" {
		return &LanguageInfo{
			Name:    "Node.js",
			Version: nodeVersion,
			Color:   getNodeJSColor(nodeVersion),
		}
	}

	// 2. Python - verifica PYTHON_VERSION
	if pythonVersion := detectPythonVersion(envVars); pythonVersion != "" {
		return &LanguageInfo{
			Name:    "Python",
			Version: pythonVersion,
			Color:   getPythonColor(pythonVersion),
		}
	}

	// 3. Java - verifica JAVA_VERSION ou JAVA_HOME
	if javaVersion := detectJavaVersion(envVars); javaVersion != "" {
		return &LanguageInfo{
			Name:    "Java",
			Version: javaVersion,
			Color:   getJavaColor(javaVersion),
		}
	}

	// 4. Go - verifica GOLANG_VERSION, GO_VERSION ou GOPATH
	if goVersion := detectGoVersion(envVars); goVersion != "" {
		return &LanguageInfo{
			Name:    "Go",
			Version: goVersion,
			Color:   getGoColor(goVersion),
		}
	}

	// 5. PHP - verifica PHP_VERSION
	if phpVersion := detectPHPVersion(envVars); phpVersion != "" {
		return &LanguageInfo{
			Name:    "PHP",
			Version: phpVersion,
			Color:   getPHPColor(phpVersion),
		}
	}

	// 6. Ruby - verifica RUBY_VERSION
	if rubyVersion := detectRubyVersion(envVars); rubyVersion != "" {
		return &LanguageInfo{
			Name:    "Ruby",
			Version: rubyVersion,
			Color:   getRubyColor(rubyVersion),
		}
	}

	// 7. .NET - verifica DOTNET_VERSION ou ASPNETCORE_VERSION
	if dotnetVersion := detectDotNetVersion(envVars); dotnetVersion != "" {
		return &LanguageInfo{
			Name:    ".NET",
			Version: dotnetVersion,
			Color:   getDotNetColor(dotnetVersion),
		}
	}

	// 8. Rust - verifica RUST_VERSION ou CARGO_HOME
	if rustVersion := detectRustVersion(envVars); rustVersion != "" {
		return &LanguageInfo{
			Name:    "Rust",
			Version: rustVersion,
			Color:   "success",
		}
	}

	// 9. Detecção por CMD/Entrypoint (para linguagens interpretadas)
	if lang := detectByCommand(cmd, entrypoint); lang != nil {
		return lang
	}

	// 10. Detecção por padrões de binários compilados (Go, Rust, C/C++)
	if lang := detectCompiledBinary(entrypoint, cmd, workingDir, imageInspect.Size); lang != nil {
		return lang
	}

	return nil
}

// Detectores de versão específicos (apenas por variáveis de ambiente)
func detectNodeJSVersion(envVars []string) string {
	for _, envVar := range envVars {
		if strings.HasPrefix(envVar, "NODE_VERSION=") {
			return strings.TrimPrefix(envVar, "NODE_VERSION=")
		}
	}
	return ""
}

func detectPythonVersion(envVars []string) string {
	for _, envVar := range envVars {
		if strings.HasPrefix(envVar, "PYTHON_VERSION=") {
			return strings.TrimPrefix(envVar, "PYTHON_VERSION=")
		}
	}
	return ""
}

func detectJavaVersion(envVars []string) string {
	for _, envVar := range envVars {
		if strings.HasPrefix(envVar, "JAVA_VERSION=") {
			return strings.TrimPrefix(envVar, "JAVA_VERSION=")
		}
		if strings.HasPrefix(envVar, "JAVA_HOME=") {
			path := strings.TrimPrefix(envVar, "JAVA_HOME=")
			if strings.Contains(path, "java-") {
				parts := strings.Split(path, "java-")
				if len(parts) > 1 {
					version := strings.Split(parts[1], "/")[0]
					return version
				}
			}
			return "detected"
		}
	}
	return ""
}

func detectGoVersion(envVars []string) string {
	for _, envVar := range envVars {
		if strings.HasPrefix(envVar, "GOLANG_VERSION=") {
			return strings.TrimPrefix(envVar, "GOLANG_VERSION=")
		}
		if strings.HasPrefix(envVar, "GO_VERSION=") {
			return strings.TrimPrefix(envVar, "GO_VERSION=")
		}
	}

	// Verifica GOPATH apenas se não houver outras linguagens mais óbvias
	for _, envVar := range envVars {
		if strings.HasPrefix(envVar, "GOPATH=") {
			return "detected"
		}
	}

	return ""
}

func detectPHPVersion(envVars []string) string {
	for _, envVar := range envVars {
		if strings.HasPrefix(envVar, "PHP_VERSION=") {
			return strings.TrimPrefix(envVar, "PHP_VERSION=")
		}
	}
	return ""
}

func detectRubyVersion(envVars []string) string {
	for _, envVar := range envVars {
		if strings.HasPrefix(envVar, "RUBY_VERSION=") {
			return strings.TrimPrefix(envVar, "RUBY_VERSION=")
		}
	}
	return ""
}

func detectDotNetVersion(envVars []string) string {
	for _, envVar := range envVars {
		if strings.HasPrefix(envVar, "DOTNET_VERSION=") {
			return strings.TrimPrefix(envVar, "DOTNET_VERSION=")
		}
		if strings.HasPrefix(envVar, "ASPNETCORE_VERSION=") {
			return strings.TrimPrefix(envVar, "ASPNETCORE_VERSION=")
		}
	}
	return ""
}

func detectRustVersion(envVars []string) string {
	for _, envVar := range envVars {
		if strings.HasPrefix(envVar, "RUST_VERSION=") {
			return strings.TrimPrefix(envVar, "RUST_VERSION=")
		}
		if strings.HasPrefix(envVar, "CARGO_HOME=") {
			return "detected"
		}
	}
	return ""
}

// Detecção por comando (para linguagens interpretadas)
func detectByCommand(cmd []string, entrypoint []string) *LanguageInfo {
	allCommands := append(entrypoint, cmd...)
	commandStr := strings.Join(allCommands, " ")

	// Ordem de prioridade
	if strings.Contains(commandStr, "node") || strings.Contains(commandStr, "npm") {
		return &LanguageInfo{Name: "Node.js", Version: "unknown", Color: "warning"}
	}

	if strings.Contains(commandStr, "python") {
		return &LanguageInfo{Name: "Python", Version: "unknown", Color: "warning"}
	}

	if strings.Contains(commandStr, "java -jar") || strings.Contains(commandStr, "java ") {
		return &LanguageInfo{Name: "Java", Version: "unknown", Color: "warning"}
	}

	if strings.Contains(commandStr, "php") {
		return &LanguageInfo{Name: "PHP", Version: "unknown", Color: "warning"}
	}

	if strings.Contains(commandStr, "ruby") {
		return &LanguageInfo{Name: "Ruby", Version: "unknown", Color: "warning"}
	}

	if strings.Contains(commandStr, "dotnet") {
		return &LanguageInfo{Name: ".NET", Version: "unknown", Color: "warning"}
	}

	return nil
}

// Detecção de binários compilados (Go, Rust, C/C++)
func detectCompiledBinary(entrypoint []string, cmd []string, workingDir string, imageSize int64) *LanguageInfo {
	// Imagens muito pequenas com binário no entrypoint geralmente são Go
	// Imagens Go compiladas costumam ter entre 5MB e 50MB
	sizeInMB := float64(imageSize) / (1024 * 1024)

	// Verifica se há um binário no entrypoint
	if len(entrypoint) > 0 {
		binary := entrypoint[0]

		// Padrões comuns de binários Go
		// Geralmente são caminhos como /app/main, /app/server, /bin/app, etc.
		isLikelyGoBinary := (strings.HasPrefix(binary, "/app/") ||
			strings.HasPrefix(binary, "/usr/local/bin/") ||
			strings.HasPrefix(binary, "/bin/")) &&
			!strings.HasSuffix(binary, ".sh") &&
			!strings.Contains(binary, "python") &&
			!strings.Contains(binary, "node") &&
			!strings.Contains(binary, "java") &&
			!strings.Contains(binary, "ruby") &&
			!strings.Contains(binary, "php")

		// Working dir comum em projetos Go
		hasGoWorkingDir := workingDir == "/app" || workingDir == "/go/src/app"

		// Imagem pequena é um forte indicador de Go (ou Rust)
		isSmallImage := sizeInMB > 5 && sizeInMB < 100

		// Se tem todas as características de um binário Go
		if isLikelyGoBinary && (hasGoWorkingDir || isSmallImage) {
			// Se a imagem é extremamente pequena (< 20MB), é muito provável que seja Go
			if sizeInMB < 20 {
				return &LanguageInfo{
					Name:    "Go",
					Version: "compiled",
					Color:   "success",
				}
			}
			// Se tem working dir /app e é pequena, também é provável Go
			if hasGoWorkingDir && isSmallImage {
				return &LanguageInfo{
					Name:    "Go",
					Version: "compiled",
					Color:   "success",
				}
			}
		}
	}

	return nil
}

// Funções de colorização por linguagem
func getNodeJSColor(version string) string {
	if version == "detected" || version == "unknown" {
		return "warning"
	}
	major := getMajorVersion(version)
	if major < 14 {
		return "error"
	}
	if major >= 14 && major <= 16 {
		return "warning"
	}
	return "success"
}

func getPythonColor(version string) string {
	if version == "detected" || version == "unknown" {
		return "warning"
	}
	major := getMajorVersion(version)
	if major < 3 {
		return "error"
	}
	if major == 3 {
		minor := getMinorVersion(version)
		if minor < 8 {
			return "warning"
		}
	}
	return "success"
}

func getJavaColor(version string) string {
	if version == "detected" || version == "unknown" {
		return "warning"
	}
	major := getMajorVersion(version)
	if major < 11 {
		return "error"
	}
	if major >= 11 && major < 17 {
		return "warning"
	}
	return "success"
}

func getGoColor(version string) string {
	if version == "detected" || version == "unknown" || version == "compiled" {
		return "success"
	}
	major := getMajorVersion(version)
	minor := getMinorVersion(version)

	if major < 1 {
		return "error"
	}
	if major == 1 && minor < 19 {
		return "warning"
	}
	return "success"
}

func getPHPColor(version string) string {
	if version == "detected" || version == "unknown" {
		return "warning"
	}
	major := getMajorVersion(version)
	if major < 7 {
		return "error"
	}
	if major == 7 {
		return "warning"
	}
	return "success"
}

func getRubyColor(version string) string {
	if version == "detected" || version == "unknown" {
		return "warning"
	}
	major := getMajorVersion(version)
	if major < 2 {
		return "error"
	}
	if major == 2 {
		return "warning"
	}
	return "success"
}

func getDotNetColor(version string) string {
	if version == "detected" || version == "unknown" {
		return "warning"
	}
	major := getMajorVersion(version)
	if major < 6 {
		return "warning"
	}
	return "success"
}

// Utilitários para extrair versões
func getMajorVersion(version string) int {
	parts := strings.Split(version, ".")
	if len(parts) == 0 {
		return 0
	}
	num, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0
	}
	return num
}

func getMinorVersion(version string) int {
	parts := strings.Split(version, ".")
	if len(parts) < 2 {
		return 0
	}
	num, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0
	}
	return num
}

// Função para imprimir linguagem detectada com cor
func PrintLanguageWithColor(imageInspect image.InspectResponse) {
	lang := DetectPrimaryLanguage(imageInspect)

	if lang == nil {
		fmt.Printf("  - Language: ")
		fmt.Println(WarningSprintf("<none detected>"))
		return
	}

	fmt.Printf("  - %s version: ", lang.Name)
	switch lang.Color {
	case "success":
		fmt.Println(SuccessSprintf("%s", lang.Version))
	case "warning":
		fmt.Println(WarningSprintf("%s", lang.Version))
	case "error":
		fmt.Println(ErrorSprintf("%s", lang.Version))
	default:
		fmt.Println(lang.Version)
	}
}

// Função para verificar se a linguagem está desatualizada
func HasOutdatedLanguage(imageInspect image.InspectResponse) bool {
	lang := DetectPrimaryLanguage(imageInspect)

	if lang == nil {
		return false
	}

	return lang.Color == "error" || lang.Color == "warning"
}

// Função para obter sugestões de melhoria de linguagem
func GetLanguageImprovementSuggestions(imageInspect image.InspectResponse) []string {
	lang := DetectPrimaryLanguage(imageInspect)
	suggestions := []string{}

	if lang == nil {
		return suggestions
	}

	if lang.Color == "error" {
		suggestions = append(suggestions,
			fmt.Sprintf("  - %s version %s is outdated and may have security vulnerabilities. Consider upgrading to a newer version.",
				lang.Name, lang.Version))
	} else if lang.Color == "warning" {
		if lang.Version == "unknown" {
			suggestions = append(suggestions,
				fmt.Sprintf("  - %s runtime detected but version could not be determined. Consider using official base images with explicit version tags.",
					lang.Name))
		} else {
			suggestions = append(suggestions,
				fmt.Sprintf("  - %s version %s is approaching end-of-life. Consider upgrading to ensure continued support.",
					lang.Name, lang.Version))
		}
	}

	return suggestions
}

// Compatibilidade com código antigo
func GetImageNodeJsMajorVersionNumber(imageInspect image.InspectResponse) int {
	lang := DetectPrimaryLanguage(imageInspect)
	if lang == nil || lang.Name != "Node.js" {
		return 0
	}
	return getMajorVersion(lang.Version)
}
