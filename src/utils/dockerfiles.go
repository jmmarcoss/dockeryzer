package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jorgevvs2/dockeryzer/src/ai"
	"github.com/jorgevvs2/dockeryzer/src/config"
)

func generateAIPrompt(tech *ProjectTechnology, ignoreComments bool) string {
	// Convert project info to a concise JSON string
	techJson, _ := json.MarshalIndent(tech, "", "  ")

	basePrompt := `Generate a production-ready optimized Dockerfile for a project with the following characteristics:
%s

Technical requirements:
- Detect the primary language and framework from the provided information
- Use appropriate base image for the detected language/framework:
  * Node.js projects: node:alpine or node:lts-alpine
  * Python projects: python:3.12-slim or python:alpine
  * Go projects: golang:1.25.1 for build, alpine for runtime
  * Java or Spring Boot projects: openjdk:25-ea-slim-bookworm
  * Rust projects: rust:alpine for build, alpine for runtime
  * PHP projects: php:8.2-fpm-alpine or php:apache
  * Ruby projects: ruby:3.2-alpine
  * .NET projects: mcr.microsoft.com/dotnet/sdk for build, runtime for production
- The Dockerfile must be optimized for production use
- Use multi-stage builds to optimize the final image size whenever possible
- Try to keep the number of layers as low as possible
- Follow security best practices (non-root user, minimal base image)
- Include only necessary files (use .dockerignore patterns in comments if helpful)
- Include Health Check instruction
- Make sure the application starts correctly
- Copy all necessary configuration and dependency files
- Install the correct package manager if needed (npm, yarn, pnpm, pip, poetry, cargo, composer, etc.)
- Expose appropriate ports based on the framework
- At the end of the Dockerfile, add a comment with the "docker run" example command to start the application

Formatting requirements:
- Return ONLY the raw Dockerfile content without any markdown formatting, code blocks, or explanations
- Start directly with the FROM instruction or the comment block
- Do not include any markdown backticks or formatting
%s

Remember:
Respond with only the raw Dockerfile content, starting with FROM (or the comment block) and no other text or formatting.`

	commentInstruction := ""
	if ignoreComments {
		commentInstruction = "- Do not include any comments in the Dockerfile"
	} else {
		commentInstruction = "- Each instruction must be preceded by a comment explaining its purpose\n- Comments must be on their own lines, above their related instructions"
	}

	return fmt.Sprintf(basePrompt, string(techJson), commentInstruction)
}

func getFallbackDockerfile(tech *ProjectTechnology, ignoreComments bool) string {
	// Fallback baseado na linguagem detectada
	switch tech.Language {
	case "javascript", "typescript":
		if tech.BuildTool == "vite" || tech.Framework == "react" || tech.Framework == "vue" {
			return getViteDockerfileContent(ignoreComments)
		}
		if HasBuildCommand() {
			return getGenericDockerfileContentWithBuildStep(ignoreComments)
		}
		return getGenericDockerfileContent(ignoreComments)

	case "python":
		return getPythonDockerfileContent(tech, ignoreComments)

	case "go":
		return getGoDockerfileContent(tech, ignoreComments)

	case "java":
		return getJavaDockerfileContent(tech, ignoreComments)

	case "rust":
		return getRustDockerfileContent(ignoreComments)

	case "php":
		return getPHPDockerfileContent(tech, ignoreComments)

	case "ruby":
		return getRubyDockerfileContent(tech, ignoreComments)

	default:
		// Fallback genérico para Node.js (compatibilidade)
		return getGenericDockerfileContent(ignoreComments)
	}
}

func getDockerfileContent(ignoreComments bool) string {
	// Use the embedded API key
	apiKey := config.APIKey
	if apiKey == "" {
		log.Fatal("API key not set in binary. Please rebuild with -ldflags")
	}

	// Detectar tecnologias do projeto (heurística + AI se necessário)
	tech := DetectProjectSmart(apiKey)

	fmt.Printf("🔍 Detected: %s", tech.Language)
	if tech.Framework != "" {
		fmt.Printf(" (%s)", tech.Framework)
	}
	if tech.PackageManager != "" {
		fmt.Printf(" [%s]", tech.PackageManager)
	}
	fmt.Println()

	// Generate AI prompt
	systemPrompt := "You are a Docker expert. Respond only with Dockerfile content, no explanations."
	userPrompt := generateAIPrompt(tech, ignoreComments)

	fmt.Println("🤖 AI is analyzing your project and generating a Dockerfile...")

	// Create AI provider using factory
	providerConfig := ai.ProviderConfig{
		Type:   ai.ProviderGemini, // Change to ai.ProviderOpenAI or ai.ProviderClaude
		APIKey: apiKey,
		Model:  "", // Empty string uses default model
	}

	provider, err := ai.NewAIProvider(providerConfig)
	if err != nil {
		fmt.Printf("❌ Error creating AI provider: %v\n", err)
		fmt.Println("❌ Falling back to default logic...")
		return getFallbackDockerfile(tech, ignoreComments)
	}
	defer provider.Close()

	// Generate content
	ctx := context.Background()
	dockerfile, err := provider.GenerateContent(ctx, systemPrompt, userPrompt, 0.2)
	if err != nil {
		fmt.Printf("❌ Error generating content: %v\n", err)
		fmt.Println("❌ Falling back to default logic...")
		return getFallbackDockerfile(tech, ignoreComments)
	}

	fmt.Println("✅ Dockerfile generated successfully!")

	// Clean up the response
	dockerfile = strings.TrimSpace(dockerfile)
	dockerfile = strings.TrimPrefix(dockerfile, "```dockerfile")
	dockerfile = strings.TrimPrefix(dockerfile, "```")
	dockerfile = strings.TrimSuffix(dockerfile, "```")
	dockerfile = strings.TrimSpace(dockerfile)

	return dockerfile
}

func CreateDockerfileContent(ignoreComments bool) {
	f, err := os.Create("Dockeryzer.Dockerfile")
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()
	content := getDockerfileContent(ignoreComments)

	_, err2 := f.WriteString(content)

	if err2 != nil {
		log.Fatal(err2)
	}
}

// Fallback templates para diferentes linguagens

func getPythonDockerfileContent(tech *ProjectTechnology, ignoreComments bool) string {
	if ignoreComments {
		return `FROM python:3.11-slim

WORKDIR /app

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY . .

CMD ["python", "app.py"]
`
	}

	return `# Use Python slim image
FROM python:3.11-slim

# Set working directory
WORKDIR /app

# Copy and install dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy application code
COPY . .

# Run the application
CMD ["python", "app.py"]

# Example: docker run -p 8000:8000 image-name
`
}

func getGoDockerfileContent(tech *ProjectTechnology, ignoreComments bool) string {
	if ignoreComments {
		return `FROM golang:alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
CMD ["./main"]
`
	}

	return `# Build stage
FROM golang:alpine AS builder

WORKDIR /app

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Build the application
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Production stage
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/main .

# Run the application
CMD ["./main"]

# Example: docker run -p 8080:8080 image-name
`
}

func getJavaDockerfileContent(tech *ProjectTechnology, ignoreComments bool) string {
	if tech.PackageManager == "gradle" {
		return `# Build stage
FROM gradle:jdk17-alpine AS builder
WORKDIR /app
COPY . .
RUN gradle build --no-daemon

# Production stage
FROM eclipse-temurin:17-jre-alpine
WORKDIR /app
COPY --from=builder /app/build/libs/*.jar app.jar
CMD ["java", "-jar", "app.jar"]
`
	}

	// Maven
	return `# Build stage
FROM maven:3.9-eclipse-temurin-17-alpine AS builder
WORKDIR /app
COPY pom.xml .
RUN mvn dependency:go-offline
COPY src ./src
RUN mvn package -DskipTests

# Production stage
FROM eclipse-temurin:17-jre-alpine
WORKDIR /app
COPY --from=builder /app/target/*.jar app.jar
CMD ["java", "-jar", "app.jar"]
`
}

func getRustDockerfileContent(ignoreComments bool) string {
	return `# Build stage
FROM rust:alpine AS builder
WORKDIR /app
COPY Cargo.toml Cargo.lock ./
RUN mkdir src && echo "fn main() {}" > src/main.rs && cargo build --release && rm -rf src
COPY . .
RUN cargo build --release

# Production stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/target/release/app .
CMD ["./app"]
`
}

func getPHPDockerfileContent(tech *ProjectTechnology, ignoreComments bool) string {
	if tech.Framework == "laravel" {
		return `FROM php:8.2-fpm-alpine

WORKDIR /app

RUN docker-php-ext-install pdo pdo_mysql

COPY --from=composer:latest /usr/bin/composer /usr/bin/composer

COPY composer.json composer.lock ./
RUN composer install --no-dev --optimize-autoloader

COPY . .

CMD ["php-fpm"]
`
	}

	return `FROM php:8.2-apache

WORKDIR /var/www/html

RUN docker-php-ext-install pdo pdo_mysql

COPY . .

RUN chown -R www-data:www-data /var/www/html

CMD ["apache2-foreground"]
`
}

func getRubyDockerfileContent(tech *ProjectTechnology, ignoreComments bool) string {
	if tech.Framework == "rails" {
		return `FROM ruby:3.2-alpine

WORKDIR /app

COPY Gemfile Gemfile.lock ./
RUN bundle install --without development test

COPY . .

RUN bundle exec rake assets:precompile

CMD ["rails", "server", "-b", "0.0.0.0"]
`
	}

	return `FROM ruby:3.2-alpine

WORKDIR /app

COPY Gemfile Gemfile.lock ./
RUN bundle install

COPY . .

CMD ["ruby", "app.rb"]
`
}

// Templates Node.js originais (mantidos para compatibilidade)

func getViteDockerfileContent(ignoreComments bool) string {
	if ignoreComments {
		return `FROM node:alpine AS builder
WORKDIR /workspace/app
COPY --chown=node:node . /workspace/app
RUN npm ci --only=production && npm run build && npm cache clean --force

FROM node:alpine
COPY --from=builder --chown=node:node /workspace/app/dist /app
WORKDIR /app
CMD ["npx", "serve", "-p", "3000", "-s", "/app"]
`
	}

	return `# Build stage
FROM node:alpine AS builder

WORKDIR /workspace/app

COPY --chown=node:node . /workspace/app

RUN npm ci --only=production && npm run build && npm cache clean --force

# Production stage
FROM node:alpine

COPY --from=builder --chown=node:node /workspace/app/dist /app

WORKDIR /app

CMD ["npx", "serve", "-p", "3000", "-s", "/app"]

# Example: docker run -p 3000:3000 image-name
`
}

// DetectProjectWithAI usa LLM quando heurística falha
func DetectProjectWithAI(tech *ProjectTechnology, apiKey string) error {
	// Se já detectamos tudo, não precisa de AI
	if tech.Language != "unknown" && tech.Language != "" {
		return nil
	}

	fmt.Println("🤖 Using AI to detect project technology...")

	// Preparar contexto para a AI
	dataTemplate := map[string]interface{}{
		"rootFiles":      tech.RootFiles,
		"configFiles":    tech.ConfigFiles,
		"fileExtensions": tech.FileExtensions,
	}

	contextJson, _ := json.MarshalIndent(dataTemplate, "", "  ")

	prompt := fmt.Sprintf(`Analyze the following project structure and identify:
1. Primary programming language
2. Framework (if any)
3. Package manager
4. Build tool (if any)

Project information:
%s

Respond ONLY with a JSON object in this exact format:
{
  "language": "language-name",
  "framework": "framework-name",
  "packageManager": "package-manager-name",
  "buildTool": "build-tool-name"
}

Use lowercase for all values. If something is not detected, use empty string.`, string(contextJson))

	// Criar provedor AI
	providerConfig := ai.ProviderConfig{
		Type:   ai.ProviderGemini,
		APIKey: apiKey,
		Model:  "",
	}

	provider, err := ai.NewAIProvider(providerConfig)
	if err != nil {
		return fmt.Errorf("failed to create AI provider: %w", err)
	}
	defer provider.Close()

	ctx := context.Background()
	response, err := provider.GenerateContent(ctx, "You are a project analysis expert. Always respond with valid JSON only.", prompt, 0.1)
	if err != nil {
		return fmt.Errorf("failed to generate content: %w", err)
	}

	// Parse resposta JSON
	response = strings.TrimSpace(response)
	response = strings.TrimPrefix(response, "```json")
	response = strings.TrimPrefix(response, "```")
	response = strings.TrimSuffix(response, "```")
	response = strings.TrimSpace(response)

	var aiResult struct {
		Language       string `json:"language"`
		Framework      string `json:"framework"`
		PackageManager string `json:"packageManager"`
		BuildTool      string `json:"buildTool"`
	}

	if err := json.Unmarshal([]byte(response), &aiResult); err != nil {
		return fmt.Errorf("failed to parse AI response: %w", err)
	}

	// Atualizar tecnologia detectada
	if aiResult.Language != "" {
		tech.Language = aiResult.Language
	}
	if aiResult.Framework != "" {
		tech.Framework = aiResult.Framework
	}
	if aiResult.PackageManager != "" {
		tech.PackageManager = aiResult.PackageManager
	}
	if aiResult.BuildTool != "" {
		tech.BuildTool = aiResult.BuildTool
	}

	fmt.Printf("✅ AI detected: %s", tech.Language)
	if tech.Framework != "" {
		fmt.Printf(" (%s)", tech.Framework)
	}
	fmt.Println()

	return nil
}

// DetectProjectSmart usa heurística primeiro, AI como fallback
func DetectProjectSmart(apiKey string) *ProjectTechnology {
	tech := DetectProject()

	// Se linguagem não detectada ou desconhecida, usar AI
	if tech.Language == "unknown" || tech.Language == "" {
		if apiKey != "" {
			fmt.Println("⚠️  Could not detect project type automatically")
			err := DetectProjectWithAI(tech, apiKey)
			if err != nil {
				fmt.Printf("⚠️  AI detection failed: %v\n", err)
			}
		}
	}

	return tech
}

func getGenericDockerfileContent(ignoreComments bool) string {
	if ignoreComments {
		return `FROM node:alpine AS builder
WORKDIR /workspace/app
COPY --chown=node:node package*.json ./
RUN npm ci --only=production && npm cache clean --force
COPY --chown=node:node . .

FROM node:alpine
WORKDIR /workspace/app
COPY --from=builder --chown=node:node /workspace/app .
ENTRYPOINT ["npm", "run", "start"]
`
	}

	return `# Build stage
FROM node:alpine AS builder

WORKDIR /workspace/app

COPY --chown=node:node package*.json ./
RUN npm ci --only=production && npm cache clean --force

COPY --chown=node:node . .

# Production stage
FROM node:alpine

WORKDIR /workspace/app

COPY --from=builder --chown=node:node /workspace/app .

ENTRYPOINT ["npm", "run", "start"]

# Example: docker run -p 3000:3000 image-name
`
}

func getGenericDockerfileContentWithBuildStep(ignoreComments bool) string {
	if ignoreComments {
		return `FROM node:alpine AS builder
WORKDIR /workspace/app
COPY --chown=node:node . .
RUN npm ci --only=production && npm run build && npm cache clean --force

FROM node:alpine
WORKDIR /workspace/app
COPY --from=builder --chown=node:node /workspace/app/dist .
ENTRYPOINT ["npm", "run", "start"]
`
	}

	return `# Build stage
FROM node:alpine AS builder

WORKDIR /workspace/app

COPY --chown=node:node . .

RUN npm ci --only=production && npm run build && npm cache clean --force

# Production stage
FROM node:alpine

WORKDIR /workspace/app

COPY --from=builder --chown=node:node /workspace/app/dist .

ENTRYPOINT ["npm", "run", "start"]

# Example: docker run -p 3000:3000 image-name
`
}
