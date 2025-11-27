package utils

import "fmt"

func BuildDockerfilePrompt(projectTree string, ignoreComments bool) string {

	commentRule := "Include explanatory comments."
	if ignoreComments {
		commentRule = "Do not include comments."
	}

	return fmt.Sprintf(`
You are a Docker expert.

Generate a production-ready optimized Dockerfile for a project with the following project struture:
%s

Technical requirements:
- Detect the primary language and framework from the provided information
- Use appropriate base image for the detected language/framework:
  * Node.js projects: node:alpine or node:lts-alpine
  * Python projects: python:3.11-slim or python:alpine
  * Go projects: golang:alpine for build, alpine for runtime
  * Java projects: eclipse-temurin:8u472-b08-jre-alpine-3.22 or openjdk:21-ea-slim-bookworm
  * Rust projects: rust:alpine for build, alpine for runtime
  * PHP projects: php:8.2-fpm-alpine or php:apache
  * Ruby projects: ruby:3.2-alpine
  * .NET projects: mcr.microsoft.com/dotnet/sdk for build, runtime for production
- The Dockerfile must be optimized for production use
- Use multi-stage builds to optimize the final image size whenever possible
- Try to keep the number of layers as low as possible
- Follow security best practices (non-root user, minimal base image)
- Include only necessary files (use .dockerignore patterns in comments if helpful)
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
Respond with only the raw Dockerfile content, starting with FROM (or the comment block) and no other text or formatting.
`, projectTree, commentRule)
}
