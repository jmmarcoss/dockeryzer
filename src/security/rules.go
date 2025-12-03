package security

import (
	"os"
	"strings"
)

// CIS-1.1 Uso de imagem base oficial
type OfficialBaseImageRule struct{}

func (r OfficialBaseImageRule) Check(df string) CISResult {
	lines := strings.Split(df, "\n")
	for _, line := range lines {
		if strings.HasPrefix(strings.ToUpper(strings.TrimSpace(line)), "FROM") {
			image := strings.Fields(line)[1]
			if strings.Contains(image, "/") && !strings.HasPrefix(image, "library/") {
				return CISResult{
					RuleID:      "CIS-1.1",
					Description: "Use official base images",
					Passed:      false,
					Severity:    "MEDIUM",
					Message:     "Base image does not appear to be official",
				}
			}
			return CISResult{RuleID: "CIS-1.1", Passed: true}
		}
	}

	return CISResult{
		RuleID:      "CIS-1.1",
		Description: "Use official base images",
		Passed:      false,
		Severity:    "HIGH",
		Message:     "No FROM instruction found",
	}
}

// CIS-1.2 Uso de tag explícita (não latest)
type ExplicitTagRule struct{}

func (r ExplicitTagRule) Check(df string) CISResult {
	lines := strings.Split(df, "\n")
	for _, line := range lines {
		if strings.HasPrefix(strings.ToUpper(strings.TrimSpace(line)), "FROM") {
			image := strings.Fields(line)[1]
			if !strings.Contains(image, ":") || strings.HasSuffix(image, ":latest") {
				return CISResult{
					RuleID:      "CIS-1.2",
					Description: "Use explicit image tag (not latest)",
					Passed:      false,
					Severity:    "HIGH",
					Message:     "Image uses latest or no tag",
				}
			}
			return CISResult{RuleID: "CIS-1.2", Passed: true}
		}
	}
	return CISResult{RuleID: "CIS-1.2", Passed: false}
}

// CIS-4.1: Container n;ao deve rodar como root
type NoRootUserRule struct{}

func (r NoRootUserRule) Check(df string) CISResult {
	if !strings.Contains(strings.ToLower(df), "user") {
		return CISResult{
			RuleID:      "CIS-4.1",
			Description: "Container should not run as root",
			Passed:      false,
			Severity:    "HIGH",
			Message:     "Missing USER instruction",
		}
	}
	return CISResult{RuleID: "CIS-4.1", Passed: true}
}

type CleanCacheRule struct{}

// CIS-5.1 Limpeza de cache e arquivos temporários
func (r CleanCacheRule) Check(df string) CISResult {
	lower := strings.ToLower(df)
	if strings.Contains(lower, "apt-get clean") ||
		strings.Contains(lower, "rm -rf /var/lib/apt/lists") ||
		strings.Contains(lower, "apk --no-cache") {
		return CISResult{RuleID: "CIS-5.1", Passed: true}
	}

	return CISResult{
		RuleID:      "CIS-5.1",
		Description: "Remove cache and temporary files",
		Passed:      false,
		Severity:    "MEDIUM",
		Message:     "No cache cleanup detected",
	}
}

// CIS-4.6: HEALTHCHECK
type HealthcheckRule struct{}

func (r HealthcheckRule) Check(df string) CISResult {
	if !strings.Contains(strings.ToLower(df), "healthcheck") {
		return CISResult{
			RuleID:      "CIS-4.6",
			Description: "Container must define HEALTHCHECK",
			Passed:      false,
			Severity:    "LOW",
			Message:     "Missing HEALTHCHECK instruction",
		}
	}
	return CISResult{RuleID: "CIS-4.6", Passed: true}
}

type DockerIgnoreRule struct{}

// CIS-5.2 Uso de .dockerignore
func (r DockerIgnoreRule) Check(df string) CISResult {
	if _, err := os.Stat(".dockerignore"); err == nil {
		return CISResult{RuleID: "CIS-5.2", Passed: true}
	}

	return CISResult{
		RuleID:      "CIS-5.2",
		Description: "Use .dockerignore",
		Passed:      false,
		Severity:    "LOW",
		Message:     ".dockerignore file not found",
	}
}

type MinimalPortExposureRule struct{}

// CIS-6.1 Exposição mínima de portas
func (r MinimalPortExposureRule) Check(df string) CISResult {
	count := 0
	for _, line := range strings.Split(df, "\n") {
		if strings.HasPrefix(strings.ToUpper(strings.TrimSpace(line)), "EXPOSE") {
			count++
		}
	}

	if count > 1 {
		return CISResult{
			RuleID:      "CIS-6.1",
			Description: "Expose minimum ports",
			Passed:      false,
			Severity:    "LOW",
			Message:     "Multiple exposed ports detected",
		}
	}

	return CISResult{RuleID: "CIS-6.1", Passed: true}
}

type MultiStageBuildRule struct{}

// CIS-7.1 Uso de Multi-stage build
func (r MultiStageBuildRule) Check(df string) CISResult {
	count := 0
	for _, line := range strings.Split(df, "\n") {
		if strings.HasPrefix(strings.ToUpper(strings.TrimSpace(line)), "FROM") {
			count++
		}
	}

	if count <= 1 {
		return CISResult{
			RuleID:      "CIS-7.1",
			Description: "Use multi-stage builds when appropriate",
			Passed:      false,
			Severity:    "MEDIUM",
			Message:     "Single stage build detected",
		}
	}

	return CISResult{RuleID: "CIS-7.1", Passed: true}
}

type CombinedRunCommandRule struct{}

// CIS-8.1 Combinação de comandos RUN
func (r CombinedRunCommandRule) Check(df string) CISResult {
	for _, line := range strings.Split(df, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "RUN") && !strings.Contains(line, "&&") {
			return CISResult{
				RuleID:      "CIS-8.1",
				Description: "Combine RUN instructions",
				Passed:      false,
				Severity:    "LOW",
				Message:     "RUN instruction not combined",
			}
		}
	}
	return CISResult{RuleID: "CIS-8.1", Passed: true}
}

type OptimizedOrderRule struct{}

// CIS-9.1 Ordem otimizada para cache
func (r OptimizedOrderRule) Check(df string) CISResult {
	installIndex := -1
	copyIndex := -1

	lines := strings.Split(df, "\n")

	for i, line := range lines {
		upper := strings.ToUpper(line)
		if strings.Contains(upper, "INSTALL") && strings.HasPrefix(strings.TrimSpace(upper), "RUN") {
			installIndex = i
		}
		if strings.HasPrefix(strings.TrimSpace(upper), "COPY") {
			copyIndex = i
		}
	}

	if installIndex > copyIndex {
		return CISResult{
			RuleID:      "CIS-9.1",
			Description: "Optimize instruction order",
			Passed:      false,
			Severity:    "LOW",
			Message:     "COPY before dependency install",
		}
	}

	return CISResult{RuleID: "CIS-9.1", Passed: true}
}
