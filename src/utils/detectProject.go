// File: src/utils/detector.go
package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ProjectTechnology representa uma tecnologia detectada
type ProjectTechnology struct {
	Language        string            `json:"language"`                 // Ex: "go", "python", "javascript"
	Framework       string            `json:"framework,omitempty"`      // Ex: "nextjs", "django", "gin"
	BuildTool       string            `json:"buildTool,omitempty"`      // Ex: "vite", "webpack", "maven"
	PackageManager  string            `json:"packageManager,omitempty"` // Ex: "npm", "yarn", "pip"
	Version         string            `json:"version,omitempty"`        // Versão se detectada
	ConfigFiles     []string          `json:"configFiles"`              // Arquivos de config encontrados
	Dependencies    map[string]string `json:"dependencies,omitempty"`
	DevDependencies map[string]string `json:"devDependencies,omitempty"`
	RootFiles       []string          `json:"rootFiles"`
	FileExtensions  map[string]int    `json:"fileExtensions"` // Contagem de arquivos por extensão
}

// FileExtensionStats coleta estatísticas de extensões de arquivo
type FileExtensionStats struct {
	Extensions map[string]int // extensão -> contagem
	TotalFiles int
}

// DetectProject analisa o projeto e retorna informações completas
func DetectProject() *ProjectTechnology {
	tech := &ProjectTechnology{
		ConfigFiles:    []string{},
		RootFiles:      GetRootFiles(),
		FileExtensions: make(map[string]int),
	}

	// 1. Coletar estatísticas de extensões de arquivo
	stats := analyzeFileExtensions()
	tech.FileExtensions = stats.Extensions

	// 2. Detectar linguagem principal baseado em extensões
	tech.Language = detectLanguageFromExtensions(stats)

	// 3. Detectar arquivos de configuração
	tech.ConfigFiles = detectConfigFiles()

	// 4. Detectar tecnologias específicas baseado na linguagem
	switch tech.Language {
	case "javascript", "typescript":
		detectNodeJSProject(tech)
	case "python":
		detectPythonProject(tech)
	case "go":
		detectGoProject(tech)
	case "java":
		detectJavaProject(tech)
	case "rust":
		detectRustProject(tech)
	case "php":
		detectPHPProject(tech)
	case "ruby":
		detectRubyProject(tech)
	case "csharp":
		detectCSharpProject(tech)
	default:
		// Tenta detectar por arquivos de configuração conhecidos
		detectByConfigFiles(tech)
	}

	return tech
}

// ShowProjectInfo exibe informações detectadas do projeto (útil para debug)
func ShowProjectInfo(useAI bool) {
	var tech *ProjectTechnology

	// if useAI {
	// 	apiKey := config.APIKey
	// 	if apiKey == "" {
	// 		log.Fatal("API key not set in binary. Please rebuild with -ldflags")
	// 	}
	// 	tech = DetectProjectSmart(apiKey)
	// } else {
	tech = DetectProject()
	// }

	fmt.Println("\n📊 Project Detection Results")
	fmt.Println("========================================")
	fmt.Printf("Language:        %s\n", tech.Language)
	fmt.Printf("Framework:       %s\n", tech.Framework)
	fmt.Printf("Package Manager: %s\n", tech.PackageManager)
	fmt.Printf("Build Tool:      %s\n", tech.BuildTool)
	fmt.Printf("Version:         %s\n", tech.Version)

	fmt.Println("\n📁 Config Files Found:")
	if len(tech.ConfigFiles) > 0 {
		for _, file := range tech.ConfigFiles {
			fmt.Printf("  • %s\n", file)
		}
	} else {
		fmt.Println("  (none)")
	}

	fmt.Println("\n📄 File Extensions Distribution:")
	if len(tech.FileExtensions) > 0 {
		// Ordenar por quantidade
		type extCount struct {
			ext   string
			count int
		}
		var sorted []extCount
		for ext, count := range tech.FileExtensions {
			sorted = append(sorted, extCount{ext, count})
		}
		// Simple bubble sort
		for i := 0; i < len(sorted); i++ {
			for j := i + 1; j < len(sorted); j++ {
				if sorted[i].count < sorted[j].count {
					sorted[i], sorted[j] = sorted[j], sorted[i]
				}
			}
		}
		// Mostrar top 10
		max := 10
		if len(sorted) < max {
			max = len(sorted)
		}
		for i := 0; i < max; i++ {
			fmt.Printf("  %s: %d files\n", sorted[i].ext, sorted[i].count)
		}
	} else {
		fmt.Println("  (none)")
	}

	if len(tech.Dependencies) > 0 {
		fmt.Println("\n📦 Dependencies (sample):")
		count := 0
		for dep, version := range tech.Dependencies {
			if count >= 5 {
				fmt.Printf("  ... and %d more\n", len(tech.Dependencies)-5)
				break
			}
			fmt.Printf("  • %s: %s\n", dep, version)
			count++
		}
	}

	fmt.Println("========================================")
}

// analyzeFileExtensions percorre recursivamente o projeto
func analyzeFileExtensions() FileExtensionStats {
	stats := FileExtensionStats{
		Extensions: make(map[string]int),
		TotalFiles: 0,
	}

	// Diretórios a ignorar
	ignoreDirs := map[string]bool{
		"node_modules": true,
		".git":         true,
		"vendor":       true,
		"venv":         true,
		".venv":        true,
		"__pycache__":  true,
		"dist":         true,
		"build":        true,
		"target":       true,
		".next":        true,
		".nuxt":        true,
	}

	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		// Ignorar diretórios específicos
		if info.IsDir() {
			if ignoreDirs[info.Name()] {
				return filepath.SkipDir
			}
			return nil
		}

		// Ignorar arquivos ocultos
		if strings.HasPrefix(info.Name(), ".") && path != "." {
			return nil
		}

		ext := filepath.Ext(path)
		if ext != "" {
			ext = strings.ToLower(ext)
			stats.Extensions[ext]++
			stats.TotalFiles++
		}

		return nil
	})

	return stats
}

// detectLanguageFromExtensions determina a linguagem principal
func detectLanguageFromExtensions(stats FileExtensionStats) string {
	// Mapeamento de extensões para linguagens
	langMap := map[string]string{
		".js":    "javascript",
		".jsx":   "javascript",
		".ts":    "typescript",
		".tsx":   "typescript",
		".py":    "python",
		".go":    "go",
		".java":  "java",
		".kt":    "kotlin",
		".rs":    "rust",
		".php":   "php",
		".rb":    "ruby",
		".cs":    "csharp",
		".cpp":   "cpp",
		".c":     "c",
		".swift": "swift",
		".dart":  "dart",
	}

	// Contar arquivos por linguagem
	langCount := make(map[string]int)
	for ext, count := range stats.Extensions {
		if lang, ok := langMap[ext]; ok {
			langCount[lang] += count
		}
	}

	// Retornar linguagem com mais arquivos
	maxCount := 0
	primaryLang := "unknown"
	for lang, count := range langCount {
		if count > maxCount {
			maxCount = count
			primaryLang = lang
		}
	}

	return primaryLang
}

// detectConfigFiles encontra arquivos de configuração conhecidos
func detectConfigFiles() []string {
	knownConfigs := []string{
		// Node.js
		"package.json", "package-lock.json", "yarn.lock", "pnpm-lock.yaml",
		"tsconfig.json", "webpack.config.js", "vite.config.js", "vite.config.ts",
		"next.config.js", "nuxt.config.js", "svelte.config.js",

		// Python
		"requirements.txt", "Pipfile", "Pipfile.lock", "pyproject.toml", "setup.py",
		"poetry.lock", "conda.yml", "environment.yml",

		// Go
		"go.mod", "go.sum",

		// Java
		"pom.xml", "build.gradle", "build.gradle.kts", "settings.gradle",

		// Rust
		"Cargo.toml", "Cargo.lock",

		// PHP
		"composer.json", "composer.lock",

		// Ruby
		"Gemfile", "Gemfile.lock",

		// .NET
		"*.csproj", "*.sln", "packages.config",

		// Docker
		"Dockerfile", "docker-compose.yml", "docker-compose.yaml",

		// Others
		"Makefile", "CMakeLists.txt",
	}

	var found []string
	for _, config := range knownConfigs {
		if fileExists(config) {
			found = append(found, config)
		}
	}

	return found
}

// detectNodeJSProject detecta projetos Node.js/JavaScript/TypeScript
func detectNodeJSProject(tech *ProjectTechnology) {
	if !fileExists("package.json") {
		return
	}

	// Ler package.json
	data, err := os.ReadFile("package.json")
	if err != nil {
		return
	}

	var pkgJson map[string]interface{}
	if err := json.Unmarshal(data, &pkgJson); err != nil {
		return
	}

	// Detectar package manager
	if fileExists("yarn.lock") {
		tech.PackageManager = "yarn"
	} else if fileExists("pnpm-lock.yaml") {
		tech.PackageManager = "pnpm"
	} else {
		tech.PackageManager = "npm"
	}

	// Extrair dependências
	if deps, ok := pkgJson["dependencies"].(map[string]interface{}); ok {
		tech.Dependencies = make(map[string]string)
		for key, value := range deps {
			if strValue, ok := value.(string); ok {
				tech.Dependencies[key] = strValue
			}
		}
	}

	if devDeps, ok := pkgJson["devDependencies"].(map[string]interface{}); ok {
		tech.DevDependencies = make(map[string]string)
		for key, value := range devDeps {
			if strValue, ok := value.(string); ok {
				tech.DevDependencies[key] = strValue
			}
		}
	}

	// Detectar framework
	allDeps := mergeMaps(tech.Dependencies, tech.DevDependencies)

	if _, ok := allDeps["next"]; ok {
		tech.Framework = "nextjs"
	} else if _, ok := allDeps["nuxt"]; ok {
		tech.Framework = "nuxt"
	} else if _, ok := allDeps["react"]; ok {
		tech.Framework = "react"
	} else if _, ok := allDeps["vue"]; ok {
		tech.Framework = "vue"
	} else if _, ok := allDeps["svelte"]; ok {
		tech.Framework = "svelte"
	} else if _, ok := allDeps["express"]; ok {
		tech.Framework = "express"
	} else if _, ok := allDeps["nestjs"]; ok {
		tech.Framework = "nestjs"
	}

	// Detectar build tool
	if fileExists("vite.config.js") || fileExists("vite.config.ts") {
		tech.BuildTool = "vite"
	} else if fileExists("webpack.config.js") {
		tech.BuildTool = "webpack"
	} else if _, ok := allDeps["vite"]; ok {
		tech.BuildTool = "vite"
	} else if _, ok := allDeps["webpack"]; ok {
		tech.BuildTool = "webpack"
	}
}

// detectPythonProject detecta projetos Python
func detectPythonProject(tech *ProjectTechnology) {
	tech.Language = "python"

	// Detectar package manager
	if fileExists("Pipfile") {
		tech.PackageManager = "pipenv"
	} else if fileExists("poetry.lock") {
		tech.PackageManager = "poetry"
	} else if fileExists("requirements.txt") {
		tech.PackageManager = "pip"
	} else if fileExists("conda.yml") || fileExists("environment.yml") {
		tech.PackageManager = "conda"
	}

	// Detectar framework (básico - pode ser expandido)
	if fileExists("manage.py") {
		tech.Framework = "django"
	} else if fileExists("app.py") || fileExists("main.py") {
		// Tentar detectar Flask/FastAPI lendo imports (simplificado)
		if data, err := os.ReadFile("app.py"); err == nil {
			content := string(data)
			if strings.Contains(content, "from flask") || strings.Contains(content, "import flask") {
				tech.Framework = "flask"
			} else if strings.Contains(content, "from fastapi") || strings.Contains(content, "import fastapi") {
				tech.Framework = "fastapi"
			}
		}
	}
}

// detectGoProject detecta projetos Go
func detectGoProject(tech *ProjectTechnology) {
	tech.Language = "go"
	tech.PackageManager = "go modules"

	if !fileExists("go.mod") {
		return
	}

	// Ler go.mod para versão e dependências
	data, err := os.ReadFile("go.mod")
	if err != nil {
		return
	}

	content := string(data)
	lines := strings.Split(content, "\n")

	// Detectar framework comum
	if strings.Contains(content, "github.com/gin-gonic/gin") {
		tech.Framework = "gin"
	} else if strings.Contains(content, "github.com/gofiber/fiber") {
		tech.Framework = "fiber"
	} else if strings.Contains(content, "github.com/labstack/echo") {
		tech.Framework = "echo"
	}

	// Extrair versão do Go
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "go ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				tech.Version = parts[1]
			}
			break
		}
	}
}

// detectJavaProject detecta projetos Java
func detectJavaProject(tech *ProjectTechnology) {
	tech.Language = "java"

	if fileExists("pom.xml") {
		tech.PackageManager = "maven"
		tech.BuildTool = "maven"
	} else if fileExists("build.gradle") || fileExists("build.gradle.kts") {
		tech.PackageManager = "gradle"
		tech.BuildTool = "gradle"
	}

	// Detectar framework Spring Boot
	if tech.PackageManager == "maven" {
		if data, err := os.ReadFile("pom.xml"); err == nil {
			if strings.Contains(string(data), "spring-boot") {
				tech.Framework = "spring-boot"
			}
		}
	}
}

// detectRustProject detecta projetos Rust
func detectRustProject(tech *ProjectTechnology) {
	tech.Language = "rust"
	tech.PackageManager = "cargo"
	tech.BuildTool = "cargo"

	if fileExists("Cargo.toml") {
		data, err := os.ReadFile("Cargo.toml")
		if err == nil {
			content := string(data)
			// Detectar frameworks comuns
			if strings.Contains(content, "actix-web") {
				tech.Framework = "actix-web"
			} else if strings.Contains(content, "rocket") {
				tech.Framework = "rocket"
			} else if strings.Contains(content, "axum") {
				tech.Framework = "axum"
			}
		}
	}
}

// detectPHPProject detecta projetos PHP
func detectPHPProject(tech *ProjectTechnology) {
	tech.Language = "php"

	if fileExists("composer.json") {
		tech.PackageManager = "composer"

		data, err := os.ReadFile("composer.json")
		if err == nil {
			content := string(data)
			if strings.Contains(content, "laravel/framework") {
				tech.Framework = "laravel"
			} else if strings.Contains(content, "symfony/symfony") {
				tech.Framework = "symfony"
			}
		}
	}
}

// detectRubyProject detecta projetos Ruby
func detectRubyProject(tech *ProjectTechnology) {
	tech.Language = "ruby"

	if fileExists("Gemfile") {
		tech.PackageManager = "bundler"

		data, err := os.ReadFile("Gemfile")
		if err == nil {
			content := string(data)
			if strings.Contains(content, "rails") {
				tech.Framework = "rails"
			} else if strings.Contains(content, "sinatra") {
				tech.Framework = "sinatra"
			}
		}
	}
}

// detectCSharpProject detecta projetos C#/.NET
func detectCSharpProject(tech *ProjectTechnology) {
	tech.Language = "csharp"
	tech.PackageManager = "nuget"

	// Procurar por arquivos .csproj
	files, _ := filepath.Glob("*.csproj")
	if len(files) > 0 {
		data, err := os.ReadFile(files[0])
		if err == nil {
			content := string(data)
			if strings.Contains(content, "Microsoft.AspNetCore") {
				tech.Framework = "aspnet-core"
			}
		}
	}
}

// detectByConfigFiles tenta detectar quando linguagem não foi identificada
func detectByConfigFiles(tech *ProjectTechnology) {
	for _, configFile := range tech.ConfigFiles {
		switch configFile {
		case "package.json":
			tech.Language = "javascript"
			detectNodeJSProject(tech)
		case "go.mod":
			tech.Language = "go"
			detectGoProject(tech)
		case "requirements.txt", "Pipfile", "pyproject.toml":
			tech.Language = "python"
			detectPythonProject(tech)
		case "Cargo.toml":
			tech.Language = "rust"
			detectRustProject(tech)
		case "composer.json":
			tech.Language = "php"
			detectPHPProject(tech)
		}
	}
}

// Funções auxiliares

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func mergeMaps(maps ...map[string]string) map[string]string {
	result := make(map[string]string)
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

func GetRootFiles() []string {
	files, err := os.ReadDir(".")
	if err != nil {
		return []string{}
	}

	var fileNames []string
	for _, file := range files {
		if !file.IsDir() {
			fileNames = append(fileNames, file.Name())
		}
	}
	return fileNames
}

// Funções mantidas para compatibilidade (deprecated)

func hasPackageJson() bool {
	return fileExists("package.json")
}

func GetPackageJsonScripts() map[string]string {
	tech := &ProjectTechnology{}
	detectNodeJSProject(tech)

	if tech.Dependencies == nil {
		return nil
	}

	// Ler scripts do package.json
	data, err := os.ReadFile("package.json")
	if err != nil {
		return nil
	}

	var pkgJson map[string]interface{}
	if err := json.Unmarshal(data, &pkgJson); err != nil {
		return nil
	}

	scripts := make(map[string]string)
	if scriptsSection, ok := pkgJson["scripts"].(map[string]interface{}); ok {
		for key, value := range scriptsSection {
			if strValue, ok := value.(string); ok {
				scripts[key] = strValue
			}
		}
	}

	return scripts
}

func IsViteProject() bool {
	return fileExists("vite.config.js") || fileExists("vite.config.ts")
}

func HasBuildCommand() bool {
	scripts := GetPackageJsonScripts()
	if scripts == nil {
		return false
	}
	_, exists := scripts["build"]
	return exists
}

func GetPackageJsonDependencies() (map[string]string, map[string]string) {
	tech := &ProjectTechnology{}
	detectNodeJSProject(tech)
	return tech.Dependencies, tech.DevDependencies
}
