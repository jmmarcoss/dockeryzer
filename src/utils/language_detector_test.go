package utils

import (
	"testing"

	"github.com/docker/docker/api/types/image"
	specs "github.com/moby/docker-image-spec/specs-go/v1"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

// Helper function to create mock ImageInspect
func createMockImageInspect(envVars []string, cmd []string, entrypoint []string, workingDir string, size int64) image.InspectResponse {
	return image.InspectResponse{
		Config: &specs.DockerOCIImageConfig{
			ImageConfig: ocispec.ImageConfig{
				Env:        envVars,
				Cmd:        cmd,
				Entrypoint: entrypoint,
				WorkingDir: workingDir,
			},
		},
	}
}

// Tests for Node.js Detection
func TestDetectNodeJS(t *testing.T) {
	tests := []struct {
		name          string
		envVars       []string
		cmd           []string
		entrypoint    []string
		workingDir    string
		size          int64
		expectedLang  string
		expectedVer   string
		expectedColor string
	}{
		{
			name:          "Node.js with explicit version",
			envVars:       []string{"NODE_VERSION=18.17.0", "PATH=/usr/local/bin:/usr/bin"},
			cmd:           []string{},
			entrypoint:    []string{"node", "index.js"},
			workingDir:    "/app",
			size:          100000000,
			expectedLang:  "Node.js",
			expectedVer:   "18.17.0",
			expectedColor: "success",
		},
		{
			name:          "Node.js old version",
			envVars:       []string{"NODE_VERSION=12.0.0"},
			cmd:           []string{},
			entrypoint:    []string{},
			workingDir:    "/app",
			size:          100000000,
			expectedLang:  "Node.js",
			expectedVer:   "12.0.0",
			expectedColor: "error",
		},
		{
			name:          "Node.js warning version",
			envVars:       []string{"NODE_VERSION=14.20.0"},
			cmd:           []string{},
			entrypoint:    []string{},
			workingDir:    "/app",
			size:          100000000,
			expectedLang:  "Node.js",
			expectedVer:   "14.20.0",
			expectedColor: "warning",
		},
		{
			name:          "Node.js detected by command",
			envVars:       []string{"PATH=/usr/local/bin:/usr/bin"},
			cmd:           []string{},
			entrypoint:    []string{"node", "server.js"},
			workingDir:    "/app",
			size:          100000000,
			expectedLang:  "Node.js",
			expectedVer:   "unknown",
			expectedColor: "warning",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			imageInspect := createMockImageInspect(tt.envVars, tt.cmd, tt.entrypoint, tt.workingDir, tt.size)
			result := DetectPrimaryLanguage(imageInspect)

			if result == nil {
				t.Fatalf("Expected language to be detected, got nil")
			}

			if result.Name != tt.expectedLang {
				t.Errorf("Expected language %s, got %s", tt.expectedLang, result.Name)
			}

			if result.Version != tt.expectedVer {
				t.Errorf("Expected version %s, got %s", tt.expectedVer, result.Version)
			}

			if result.Color != tt.expectedColor {
				t.Errorf("Expected color %s, got %s", tt.expectedColor, result.Color)
			}
		})
	}
}

// Tests for Go Detection
func TestDetectGo(t *testing.T) {
	tests := []struct {
		name          string
		envVars       []string
		cmd           []string
		entrypoint    []string
		workingDir    string
		size          int64
		expectedLang  string
		expectedVer   string
		expectedColor string
	}{
		{
			name:          "Go with explicit version",
			envVars:       []string{"GOLANG_VERSION=1.21.0", "PATH=/usr/local/go/bin:/usr/bin"},
			cmd:           []string{},
			entrypoint:    []string{"/app/main"},
			workingDir:    "/app",
			size:          15000000,
			expectedLang:  "Go",
			expectedVer:   "1.21.0",
			expectedColor: "success",
		},
		{
			name:          "Go compiled binary - small image",
			envVars:       []string{"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"},
			cmd:           []string{},
			entrypoint:    []string{"/app/main"},
			workingDir:    "/app",
			size:          13670000, // 13.67 MB
			expectedLang:  "Go",
			expectedVer:   "compiled",
			expectedColor: "success",
		},
		{
			name:          "Go with GOPATH",
			envVars:       []string{"GOPATH=/go", "PATH=/go/bin:/usr/local/go/bin:/usr/bin"},
			cmd:           []string{},
			entrypoint:    []string{"/app/server"},
			workingDir:    "/go/src/app",
			size:          20000000,
			expectedLang:  "Go",
			expectedVer:   "detected",
			expectedColor: "success",
		},
		{
			name:          "Go old version",
			envVars:       []string{"GOLANG_VERSION=1.16.0"},
			cmd:           []string{},
			entrypoint:    []string{},
			workingDir:    "/app",
			size:          20000000,
			expectedLang:  "Go",
			expectedVer:   "1.16.0",
			expectedColor: "warning",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			imageInspect := createMockImageInspect(tt.envVars, tt.cmd, tt.entrypoint, tt.workingDir, tt.size)
			result := DetectPrimaryLanguage(imageInspect)

			if result == nil {
				t.Fatalf("Expected language to be detected, got nil")
			}

			if result.Name != tt.expectedLang {
				t.Errorf("Expected language %s, got %s", tt.expectedLang, result.Name)
			}

			if result.Version != tt.expectedVer {
				t.Errorf("Expected version %s, got %s", tt.expectedVer, result.Version)
			}

			if result.Color != tt.expectedColor {
				t.Errorf("Expected color %s, got %s", tt.expectedColor, result.Color)
			}
		})
	}
}

// Tests for Python Detection
func TestDetectPython(t *testing.T) {
	tests := []struct {
		name          string
		envVars       []string
		cmd           []string
		entrypoint    []string
		workingDir    string
		size          int64
		expectedLang  string
		expectedVer   string
		expectedColor string
	}{
		{
			name:          "Python 3.11",
			envVars:       []string{"PYTHON_VERSION=3.11.2"},
			cmd:           []string{},
			entrypoint:    []string{"python", "app.py"},
			workingDir:    "/app",
			size:          50000000,
			expectedLang:  "Python",
			expectedVer:   "3.11.2",
			expectedColor: "success",
		},
		{
			name:          "Python 2.7 - old",
			envVars:       []string{"PYTHON_VERSION=2.7.18"},
			cmd:           []string{},
			entrypoint:    []string{},
			workingDir:    "/app",
			size:          50000000,
			expectedLang:  "Python",
			expectedVer:   "2.7.18",
			expectedColor: "error",
		},
		{
			name:          "Python 3.7 - warning",
			envVars:       []string{"PYTHON_VERSION=3.7.0"},
			cmd:           []string{},
			entrypoint:    []string{},
			workingDir:    "/app",
			size:          50000000,
			expectedLang:  "Python",
			expectedVer:   "3.7.0",
			expectedColor: "warning",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			imageInspect := createMockImageInspect(tt.envVars, tt.cmd, tt.entrypoint, tt.workingDir, tt.size)
			result := DetectPrimaryLanguage(imageInspect)

			if result == nil {
				t.Fatalf("Expected language to be detected, got nil")
			}

			if result.Name != tt.expectedLang {
				t.Errorf("Expected language %s, got %s", tt.expectedLang, result.Name)
			}

			if result.Version != tt.expectedVer {
				t.Errorf("Expected version %s, got %s", tt.expectedVer, result.Version)
			}

			if result.Color != tt.expectedColor {
				t.Errorf("Expected color %s, got %s", tt.expectedColor, result.Color)
			}
		})
	}
}

// Tests for Java Detection
func TestDetectJava(t *testing.T) {
	tests := []struct {
		name          string
		envVars       []string
		expectedLang  string
		expectedVer   string
		expectedColor string
	}{
		{
			name:          "Java 17",
			envVars:       []string{"JAVA_VERSION=17.0.1"},
			expectedLang:  "Java",
			expectedVer:   "17.0.1",
			expectedColor: "success",
		},
		{
			name:          "Java 11 - warning",
			envVars:       []string{"JAVA_VERSION=11.0.1"},
			expectedLang:  "Java",
			expectedVer:   "11.0.1",
			expectedColor: "warning",
		},
		{
			name:          "Java 8 - old",
			envVars:       []string{"JAVA_VERSION=8"},
			expectedLang:  "Java",
			expectedVer:   "8",
			expectedColor: "error",
		},
		{
			name:          "Java with JAVA_HOME",
			envVars:       []string{"JAVA_HOME=/usr/lib/jvm/java-17-openjdk"},
			expectedLang:  "Java",
			expectedVer:   "17-openjdk",
			expectedColor: "error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			imageInspect := createMockImageInspect(tt.envVars, []string{}, []string{}, "/app", 50000000)
			result := DetectPrimaryLanguage(imageInspect)

			if result == nil {
				t.Fatalf("Expected language to be detected, got nil")
			}

			if result.Name != tt.expectedLang {
				t.Errorf("Expected language %s, got %s", tt.expectedLang, result.Name)
			}

			if result.Version != tt.expectedVer {
				t.Errorf("Expected version %s, got %s", tt.expectedVer, result.Version)
			}

			if result.Color != tt.expectedColor {
				t.Errorf("Expected color %s, got %s", tt.expectedColor, result.Color)
			}
		})
	}
}

// Tests for PHP Detection
func TestDetectPHP(t *testing.T) {
	tests := []struct {
		name          string
		envVars       []string
		expectedVer   string
		expectedColor string
	}{
		{
			name:          "PHP 8.2",
			envVars:       []string{"PHP_VERSION=8.2.0"},
			expectedVer:   "8.2.0",
			expectedColor: "success",
		},
		{
			name:          "PHP 7.4",
			envVars:       []string{"PHP_VERSION=7.4.0"},
			expectedVer:   "7.4.0",
			expectedColor: "warning",
		},
		{
			name:          "PHP 5.6",
			envVars:       []string{"PHP_VERSION=5.6.0"},
			expectedVer:   "5.6.0",
			expectedColor: "error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			imageInspect := createMockImageInspect(tt.envVars, []string{}, []string{}, "/var/www", 50000000)
			result := DetectPrimaryLanguage(imageInspect)

			if result == nil {
				t.Fatalf("Expected language to be detected, got nil")
			}

			if result.Name != "PHP" {
				t.Errorf("Expected language PHP, got %s", result.Name)
			}

			if result.Version != tt.expectedVer {
				t.Errorf("Expected version %s, got %s", tt.expectedVer, result.Version)
			}

			if result.Color != tt.expectedColor {
				t.Errorf("Expected color %s, got %s", tt.expectedColor, result.Color)
			}
		})
	}
}

// Tests for Ruby Detection
func TestDetectRuby(t *testing.T) {
	tests := []struct {
		name          string
		envVars       []string
		expectedVer   string
		expectedColor string
	}{
		{
			name:          "Ruby 3.2",
			envVars:       []string{"RUBY_VERSION=3.2.0"},
			expectedVer:   "3.2.0",
			expectedColor: "success",
		},
		{
			name:          "Ruby 2.7",
			envVars:       []string{"RUBY_VERSION=2.7.0"},
			expectedVer:   "2.7.0",
			expectedColor: "warning",
		},
		{
			name:          "Ruby 1.9",
			envVars:       []string{"RUBY_VERSION=1.9.3"},
			expectedVer:   "1.9.3",
			expectedColor: "error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			imageInspect := createMockImageInspect(tt.envVars, []string{}, []string{}, "/app", 50000000)
			result := DetectPrimaryLanguage(imageInspect)

			if result == nil {
				t.Fatalf("Expected language to be detected, got nil")
			}

			if result.Name != "Ruby" {
				t.Errorf("Expected language Ruby, got %s", result.Name)
			}

			if result.Version != tt.expectedVer {
				t.Errorf("Expected version %s, got %s", tt.expectedVer, result.Version)
			}

			if result.Color != tt.expectedColor {
				t.Errorf("Expected color %s, got %s", tt.expectedColor, result.Color)
			}
		})
	}
}

// Tests for .NET Detection
func TestDetectDotNet(t *testing.T) {
	tests := []struct {
		name          string
		envVars       []string
		expectedVer   string
		expectedColor string
	}{
		{
			name:          ".NET 8.0",
			envVars:       []string{"DOTNET_VERSION=8.0"},
			expectedVer:   "8.0",
			expectedColor: "success",
		},
		{
			name:          ".NET 6.0",
			envVars:       []string{"DOTNET_VERSION=6.0"},
			expectedVer:   "6.0",
			expectedColor: "success",
		},
		{
			name:          "ASP.NET Core 5.0",
			envVars:       []string{"ASPNETCORE_VERSION=5.0"},
			expectedVer:   "5.0",
			expectedColor: "warning",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			imageInspect := createMockImageInspect(tt.envVars, []string{}, []string{}, "/app", 50000000)
			result := DetectPrimaryLanguage(imageInspect)

			if result == nil {
				t.Fatalf("Expected language to be detected, got nil")
			}

			if result.Name != ".NET" {
				t.Errorf("Expected language .NET, got %s", result.Name)
			}

			if result.Version != tt.expectedVer {
				t.Errorf("Expected version %s, got %s", tt.expectedVer, result.Version)
			}

			if result.Color != tt.expectedColor {
				t.Errorf("Expected color %s, got %s", tt.expectedColor, result.Color)
			}
		})
	}
}

// Tests for Rust Detection
func TestDetectRust(t *testing.T) {
	tests := []struct {
		name        string
		envVars     []string
		expectedVer string
	}{
		{
			name:        "Rust with version",
			envVars:     []string{"RUST_VERSION=1.70.0"},
			expectedVer: "1.70.0",
		},
		{
			name:        "Rust with CARGO_HOME",
			envVars:     []string{"CARGO_HOME=/usr/local/cargo"},
			expectedVer: "detected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			imageInspect := createMockImageInspect(tt.envVars, []string{}, []string{}, "/app", 50000000)
			result := DetectPrimaryLanguage(imageInspect)

			if result == nil {
				t.Fatalf("Expected language to be detected, got nil")
			}

			if result.Name != "Rust" {
				t.Errorf("Expected language Rust, got %s", result.Name)
			}

			if result.Version != tt.expectedVer {
				t.Errorf("Expected version %s, got %s", tt.expectedVer, result.Version)
			}

			if result.Color != "success" {
				t.Errorf("Expected color success, got %s", result.Color)
			}
		})
	}
}

// Test for No Language Detection
func TestNoLanguageDetected(t *testing.T) {
	imageInspect := createMockImageInspect(
		[]string{"PATH=/usr/local/bin:/usr/bin"},
		[]string{},
		[]string{"/bin/sh"},
		"/",
		200000000,
	)

	result := DetectPrimaryLanguage(imageInspect)

	if result != nil {
		t.Errorf("Expected no language to be detected, but got %s", result.Name)
	}
}

// Test Language Priority
func TestLanguagePriority(t *testing.T) {
	// Node.js should take priority over Go when NODE_VERSION is present
	imageInspect := createMockImageInspect(
		[]string{"NODE_VERSION=18.0.0", "GOPATH=/go"},
		[]string{},
		[]string{"/app/main"},
		"/app",
		15000000,
	)

	result := DetectPrimaryLanguage(imageInspect)

	if result == nil {
		t.Fatalf("Expected language to be detected, got nil")
	}

	if result.Name != "Node.js" {
		t.Errorf("Expected Node.js to take priority, got %s", result.Name)
	}
}

// Test HasOutdatedLanguage
func TestHasOutdatedLanguage(t *testing.T) {
	tests := []struct {
		name     string
		envVars  []string
		expected bool
	}{
		{
			name:     "Outdated Node.js",
			envVars:  []string{"NODE_VERSION=12.0.0"},
			expected: true,
		},
		{
			name:     "Current Node.js",
			envVars:  []string{"NODE_VERSION=20.0.0"},
			expected: false,
		},
		{
			name:     "Warning Node.js",
			envVars:  []string{"NODE_VERSION=16.0.0"},
			expected: true,
		},
		{
			name:     "No language",
			envVars:  []string{"PATH=/usr/bin"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			imageInspect := createMockImageInspect(tt.envVars, []string{}, []string{}, "/app", 50000000)
			result := HasOutdatedLanguage(imageInspect)

			if result != tt.expected {
				t.Errorf("Expected HasOutdatedLanguage to be %v, got %v", tt.expected, result)
			}
		})
	}
}

// Test GetLanguageImprovementSuggestions
func TestGetLanguageImprovementSuggestions(t *testing.T) {
	tests := []struct {
		name                string
		envVars             []string
		expectedSuggestions int
	}{
		{
			name:                "Outdated language",
			envVars:             []string{"NODE_VERSION=12.0.0"},
			expectedSuggestions: 1,
		},
		{
			name:                "Warning language",
			envVars:             []string{"NODE_VERSION=14.0.0"},
			expectedSuggestions: 1,
		},
		{
			name:                "Current language",
			envVars:             []string{"NODE_VERSION=20.0.0"},
			expectedSuggestions: 0,
		},
		{
			name:                "No language",
			envVars:             []string{"PATH=/usr/bin"},
			expectedSuggestions: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			imageInspect := createMockImageInspect(tt.envVars, []string{}, []string{}, "/app", 50000000)
			suggestions := GetLanguageImprovementSuggestions(imageInspect)

			if len(suggestions) != tt.expectedSuggestions {
				t.Errorf("Expected %d suggestions, got %d", tt.expectedSuggestions, len(suggestions))
			}
		})
	}
}

// Test Version Extraction Helpers
func TestGetMajorVersion(t *testing.T) {
	tests := []struct {
		version  string
		expected int
	}{
		{"18.17.0", 18},
		{"3.11.2", 3},
		{"1.21.0", 1},
		{"8", 8},
		{"invalid", 0},
	}

	for _, tt := range tests {
		t.Run(tt.version, func(t *testing.T) {
			result := getMajorVersion(tt.version)
			if result != tt.expected {
				t.Errorf("Expected major version %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestGetMinorVersion(t *testing.T) {
	tests := []struct {
		version  string
		expected int
	}{
		{"18.17.0", 17},
		{"3.11.2", 11},
		{"1.21.0", 21},
		{"8", 0},
		{"invalid", 0},
	}

	for _, tt := range tests {
		t.Run(tt.version, func(t *testing.T) {
			result := getMinorVersion(tt.version)
			if result != tt.expected {
				t.Errorf("Expected minor version %d, got %d", tt.expected, result)
			}
		})
	}
}
