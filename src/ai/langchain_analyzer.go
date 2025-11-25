package ai

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

// AIAnalyzer encapsula a lógica de análise com IA
type AIAnalyzer struct {
	llm llms.Model
}

// NewAIAnalyzer cria uma nova instância do analisador AI
func NewAIAnalyzer() (*AIAnalyzer, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable not set")
	}

	llm, err := openai.New(
		openai.WithModel("gpt-4"),
		openai.WithToken(apiKey),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OpenAI client: %w", err)
	}

	return &AIAnalyzer{llm: llm}, nil
}

// ImageAnalysisResult representa o resultado da análise de IA
type ImageAnalysisResult struct {
	SecurityScore      int // 0-100
	OptimizationScore  int // 0-100
	BestPracticesScore int // 0-100
	Recommendations    []string
	SecurityIssues     []string
	OptimizationTips   []string
	Summary            string
}

// AnalyzeImage analisa uma imagem Docker usando IA
func (a *AIAnalyzer) AnalyzeImage(imageInspect types.ImageInspect, imageName string) (*ImageAnalysisResult, error) {
	ctx := context.Background()

	// Prepara o contexto da imagem
	imageContext := a.buildImageContext(imageInspect, imageName)

	// Cria o prompt para análise
	prompt := a.buildAnalysisPrompt(imageContext)

	// Chama o LLM
	response, err := llms.GenerateFromSinglePrompt(ctx, a.llm, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate analysis: %w", err)
	}

	// Parse da resposta
	result := a.parseAnalysisResponse(response)

	return result, nil
}

// buildImageContext constrói o contexto da imagem para análise
func (a *AIAnalyzer) buildImageContext(imageInspect types.ImageInspect, imageName string) string {
	var context strings.Builder

	context.WriteString(fmt.Sprintf("Image Name: %s\n", imageName))
	context.WriteString(fmt.Sprintf("Tags: %v\n", imageInspect.RepoTags))

	// Tamanho
	sizeInMB := float64(imageInspect.Size) / 1000000
	if sizeInMB > 1000 {
		context.WriteString(fmt.Sprintf("Size: %.2f GB\n", sizeInMB/1000))
	} else {
		context.WriteString(fmt.Sprintf("Size: %.2f MB\n", sizeInMB))
	}

	context.WriteString(fmt.Sprintf("Layers: %d\n", len(imageInspect.RootFS.Layers)))
	context.WriteString(fmt.Sprintf("OS: %s\n", imageInspect.Os))
	context.WriteString(fmt.Sprintf("Architecture: %s\n", imageInspect.Architecture))

	// Linguagem detectada - detectar aqui mesmo sem importar utils
	lang := detectLanguageFromEnv(imageInspect.Config.Env)
	if lang != "" {
		context.WriteString(fmt.Sprintf("Language: %s\n", lang))
	} else {
		context.WriteString("Language: None detected\n")
	}

	// Variáveis de ambiente (filtradas)
	context.WriteString("\nEnvironment Variables:\n")
	for _, env := range imageInspect.Config.Env {
		// Filtra variáveis sensíveis
		if !strings.Contains(strings.ToUpper(env), "PASSWORD") &&
			!strings.Contains(strings.ToUpper(env), "SECRET") &&
			!strings.Contains(strings.ToUpper(env), "TOKEN") {
			context.WriteString(fmt.Sprintf("  - %s\n", env))
		}
	}

	// CMD e Entrypoint
	if len(imageInspect.Config.Cmd) > 0 {
		context.WriteString(fmt.Sprintf("\nCmd: %v\n", imageInspect.Config.Cmd))
	}
	if len(imageInspect.Config.Entrypoint) > 0 {
		context.WriteString(fmt.Sprintf("Entrypoint: %v\n", imageInspect.Config.Entrypoint))
	}

	// Working directory
	if imageInspect.Config.WorkingDir != "" {
		context.WriteString(fmt.Sprintf("Working Directory: %s\n", imageInspect.Config.WorkingDir))
	}

	// Exposed ports
	if len(imageInspect.Config.ExposedPorts) > 0 {
		context.WriteString("\nExposed Ports:\n")
		for port := range imageInspect.Config.ExposedPorts {
			context.WriteString(fmt.Sprintf("  - %s\n", port))
		}
	}

	// User
	if imageInspect.Config.User != "" {
		context.WriteString(fmt.Sprintf("\nUser: %s\n", imageInspect.Config.User))
	} else {
		context.WriteString("\nUser: root (⚠️  running as root)\n")
	}

	return context.String()
}

// detectLanguageFromEnv detecta linguagem básica das variáveis de ambiente
func detectLanguageFromEnv(envVars []string) string {
	for _, env := range envVars {
		if strings.HasPrefix(env, "NODE_VERSION=") {
			version := strings.TrimPrefix(env, "NODE_VERSION=")
			return fmt.Sprintf("Node.js %s", version)
		}
		if strings.HasPrefix(env, "PYTHON_VERSION=") {
			version := strings.TrimPrefix(env, "PYTHON_VERSION=")
			return fmt.Sprintf("Python %s", version)
		}
		if strings.HasPrefix(env, "JAVA_VERSION=") {
			version := strings.TrimPrefix(env, "JAVA_VERSION=")
			return fmt.Sprintf("Java %s", version)
		}
		if strings.HasPrefix(env, "GO_VERSION=") || strings.HasPrefix(env, "GOLANG_VERSION=") {
			version := strings.TrimPrefix(strings.TrimPrefix(env, "GO_VERSION="), "GOLANG_VERSION=")
			return fmt.Sprintf("Go %s", version)
		}
		if strings.HasPrefix(env, "PHP_VERSION=") {
			version := strings.TrimPrefix(env, "PHP_VERSION=")
			return fmt.Sprintf("PHP %s", version)
		}
		if strings.HasPrefix(env, "RUBY_VERSION=") {
			version := strings.TrimPrefix(env, "RUBY_VERSION=")
			return fmt.Sprintf("Ruby %s", version)
		}
		if strings.HasPrefix(env, "DOTNET_VERSION=") {
			version := strings.TrimPrefix(env, "DOTNET_VERSION=")
			return fmt.Sprintf(".NET %s", version)
		}
		if strings.HasPrefix(env, "RUST_VERSION=") {
			version := strings.TrimPrefix(env, "RUST_VERSION=")
			return fmt.Sprintf("Rust %s", version)
		}
	}
	return ""
}

// buildAnalysisPrompt cria o prompt para análise
func (a *AIAnalyzer) buildAnalysisPrompt(imageContext string) string {
	return fmt.Sprintf(`You are a Docker expert analyzing container images for security, optimization, and best practices.

Analyze the following Docker image and provide:

1. Security Score (0-100): Evaluate security aspects like running as root, exposed ports, base image, etc.
2. Optimization Score (0-100): Evaluate size, layers, and efficiency.
3. Best Practices Score (0-100): Evaluate adherence to Docker best practices.
4. Top 3-5 Security Issues (if any)
5. Top 3-5 Optimization Tips
6. Top 3-5 General Recommendations
7. Brief Summary (2-3 sentences)

Image Details:
%s

Respond in the following format:
SECURITY_SCORE: <number>
OPTIMIZATION_SCORE: <number>
BEST_PRACTICES_SCORE: <number>

SECURITY_ISSUES:
- <issue 1>
- <issue 2>
...

OPTIMIZATION_TIPS:
- <tip 1>
- <tip 2>
...

RECOMMENDATIONS:
- <recommendation 1>
- <recommendation 2>
...

SUMMARY:
<summary text>
`, imageContext)
}

// parseAnalysisResponse faz o parse da resposta do LLM
func (a *AIAnalyzer) parseAnalysisResponse(response string) *ImageAnalysisResult {
	result := &ImageAnalysisResult{
		SecurityScore:      0,
		OptimizationScore:  0,
		BestPracticesScore: 0,
		Recommendations:    []string{},
		SecurityIssues:     []string{},
		OptimizationTips:   []string{},
		Summary:            "",
	}

	lines := strings.Split(response, "\n")
	currentSection := ""

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		// Parse scores
		if strings.HasPrefix(line, "SECURITY_SCORE:") {
			fmt.Sscanf(line, "SECURITY_SCORE: %d", &result.SecurityScore)
		} else if strings.HasPrefix(line, "OPTIMIZATION_SCORE:") {
			fmt.Sscanf(line, "OPTIMIZATION_SCORE: %d", &result.OptimizationScore)
		} else if strings.HasPrefix(line, "BEST_PRACTICES_SCORE:") {
			fmt.Sscanf(line, "BEST_PRACTICES_SCORE: %d", &result.BestPracticesScore)
		} else if strings.HasPrefix(line, "SECURITY_ISSUES:") {
			currentSection = "security"
		} else if strings.HasPrefix(line, "OPTIMIZATION_TIPS:") {
			currentSection = "optimization"
		} else if strings.HasPrefix(line, "RECOMMENDATIONS:") {
			currentSection = "recommendations"
		} else if strings.HasPrefix(line, "SUMMARY:") {
			currentSection = "summary"
		} else if strings.HasPrefix(line, "- ") {
			// Parse bullet points
			item := strings.TrimPrefix(line, "- ")
			switch currentSection {
			case "security":
				result.SecurityIssues = append(result.SecurityIssues, item)
			case "optimization":
				result.OptimizationTips = append(result.OptimizationTips, item)
			case "recommendations":
				result.Recommendations = append(result.Recommendations, item)
			}
		} else if currentSection == "summary" {
			result.Summary += line + " "
		}
	}

	result.Summary = strings.TrimSpace(result.Summary)

	return result
}

// CompareImages compara duas imagens usando IA
func (a *AIAnalyzer) CompareImages(
	image1Name string,
	image1Inspect types.ImageInspect,
	image2Name string,
	image2Inspect types.ImageInspect,
) (string, error) {
	ctx := context.Background()

	// Constrói contexto de ambas as imagens
	context1 := a.buildImageContext(image1Inspect, image1Name)
	context2 := a.buildImageContext(image2Inspect, image2Name)

	prompt := fmt.Sprintf(`You are a Docker expert comparing two container images.

Compare these images and provide:
1. Which image is better overall and why
2. Key differences between them
3. Specific recommendations for each image

Image 1:
%s

Image 2:
%s

Provide a clear, concise comparison focusing on security, optimization, and best practices.
`, context1, context2)

	response, err := llms.GenerateFromSinglePrompt(ctx, a.llm, prompt)
	if err != nil {
		return "", fmt.Errorf("failed to generate comparison: %w", err)
	}

	return response, nil
}

// SuggestOptimizations sugere otimizações específicas para a imagem
func (a *AIAnalyzer) SuggestOptimizations(imageInspect types.ImageInspect, imageName string) ([]string, error) {
	ctx := context.Background()

	imageContext := a.buildImageContext(imageInspect, imageName)

	prompt := fmt.Sprintf(`You are a Docker optimization expert.

Analyze this image and provide 5-10 specific, actionable Dockerfile improvements.
Focus on: multi-stage builds, layer optimization, security, and size reduction.

Image Details:
%s

Provide each optimization as a bullet point with:
- The specific change to make
- Why it helps
- Example Dockerfile snippet if applicable

Format:
- Optimization 1: <description>
- Optimization 2: <description>
...
`, imageContext)

	response, err := llms.GenerateFromSinglePrompt(ctx, a.llm, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate optimizations: %w", err)
	}

	// Parse bullet points
	optimizations := []string{}
	lines := strings.Split(response, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "- ") {
			optimizations = append(optimizations, strings.TrimPrefix(line, "- "))
		}
	}

	return optimizations, nil
}

// GenerateDockerfile gera um Dockerfile otimizado baseado na análise
func (a *AIAnalyzer) GenerateDockerfile(imageInspect types.ImageInspect, imageName string) (string, error) {
	ctx := context.Background()

	lang := detectLanguageFromEnv(imageInspect.Config.Env)
	if lang == "" {
		lang = "unknown"
	}

	sizeInMB := float64(imageInspect.Size) / 1000000
	sizeStr := fmt.Sprintf("%.2f MB", sizeInMB)
	if sizeInMB > 1000 {
		sizeStr = fmt.Sprintf("%.2f GB", sizeInMB/1000)
	}

	prompt := fmt.Sprintf(`You are a Docker expert. Generate an optimized Dockerfile for an application.

Current Image: %s
Detected Language/Runtime: %s
Current Size: %s
Current Layers: %d

Generate a production-ready, multi-stage Dockerfile that:
1. Uses the detected language/runtime
2. Follows best practices
3. Optimizes for size and security
4. Uses non-root user
5. Implements proper layer caching

Return ONLY the Dockerfile content, nothing else.
`, imageName, lang, sizeStr, len(imageInspect.RootFS.Layers))

	response, err := llms.GenerateFromSinglePrompt(ctx, a.llm, prompt)
	if err != nil {
		return "", fmt.Errorf("failed to generate Dockerfile: %w", err)
	}

	return response, nil
}
