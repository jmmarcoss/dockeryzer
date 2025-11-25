package utils

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/docker/docker/api/types"
)

// Helper to capture stdout
func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

// Tests for GetImageSizeInMBs
func TestGetImageSizeInMBs(t *testing.T) {
	tests := []struct {
		name     string
		size     int64
		expected float32
	}{
		{
			name:     "Small image - 10MB",
			size:     10000000,
			expected: 10.0,
		},
		{
			name:     "Medium image - 250MB",
			size:     250000000,
			expected: 250.0,
		},
		{
			name:     "Large image - 1GB",
			size:     1000000000,
			expected: 1000.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			imageInspect := types.ImageInspect{Size: tt.size}
			result := GetImageSizeInMBs(imageInspect)

			if result != tt.expected {
				t.Errorf("Expected %f MB, got %f MB", tt.expected, result)
			}
		})
	}
}

// Tests for GetImageSizeString
func TestGetImageSizeString(t *testing.T) {
	tests := []struct {
		name     string
		size     int64
		expected string
	}{
		{
			name:     "Small image in MB",
			size:     56900000,
			expected: "56.90 MB",
		},
		{
			name:     "Medium image in MB",
			size:     250000000,
			expected: "250.00 MB",
		},
		{
			name:     "Large image in GB",
			size:     1500000000,
			expected: "1.50 GB",
		},
		{
			name:     "Very small image",
			size:     13670000,
			expected: "13.67 MB",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			imageInspect := types.ImageInspect{Size: tt.size}
			result := GetImageSizeString(imageInspect)

			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

// Tests for GetImageNumberOfLayers
func TestGetImageNumberOfLayers(t *testing.T) {
	tests := []struct {
		name     string
		layers   int
		expected int
	}{
		{
			name:     "Few layers",
			layers:   5,
			expected: 5,
		},
		{
			name:     "Many layers",
			layers:   25,
			expected: 25,
		},
		{
			name:     "No layers",
			layers:   0,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			layers := make([]string, tt.layers)
			for i := 0; i < tt.layers; i++ {
				layers[i] = fmt.Sprintf("sha256:layer%d", i)
			}

			imageInspect := types.ImageInspect{
				RootFS: types.RootFS{
					Layers: layers,
				},
			}

			result := GetImageNumberOfLayers(imageInspect)

			if result != tt.expected {
				t.Errorf("Expected %d layers, got %d", tt.expected, result)
			}
		})
	}
}

// Tests for GetImageFormattedCreationDate
func TestGetImageFormattedCreationDate(t *testing.T) {
	tests := []struct {
		name     string
		created  string
		expected string
	}{
		{
			name:     "Valid date",
			created:  "2025-11-24T12:30:45.123456789Z",
			expected: "24 Nov 2025",
		},
		{
			name:     "Another valid date",
			created:  "2024-01-15T08:15:30.987654321Z",
			expected: "15 Jan 2024",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			imageInspect := types.ImageInspect{Created: tt.created}
			result := GetImageFormattedCreationDate(imageInspect)

			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

// Tests for GetImageAuthor
func TestGetImageAuthor(t *testing.T) {
	tests := []struct {
		name     string
		author   string
		expected string
	}{
		{
			name:     "With author",
			author:   "John Doe <john@example.com>",
			expected: "John Doe <john@example.com>",
		},
		{
			name:     "No author",
			author:   "",
			expected: "<none>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			imageInspect := types.ImageInspect{Author: tt.author}
			result := GetImageAuthor(imageInspect)

			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

// Tests for PrintImageCompareLanguageResults
func TestPrintImageCompareLanguageResults(t *testing.T) {
	tests := []struct {
		name           string
		image1Env      []string
		image2Env      []string
		expectedOutput string
	}{
		{
			name:           "Both images no language",
			image1Env:      []string{"PATH=/usr/bin"},
			image2Env:      []string{"PATH=/usr/bin"},
			expectedOutput: "No programming language runtime detected",
		},
		{
			name:           "Only first image has language",
			image1Env:      []string{"NODE_VERSION=18.0.0"},
			image2Env:      []string{"PATH=/usr/bin"},
			expectedOutput: "Only image image1 has detected language runtime",
		},
		{
			name:           "Only second image has language",
			image1Env:      []string{"PATH=/usr/bin"},
			image2Env:      []string{"PYTHON_VERSION=3.11.0"},
			expectedOutput: "Only image image2 has detected language runtime",
		},
		{
			name:           "Different languages",
			image1Env:      []string{"NODE_VERSION=18.0.0"},
			image2Env:      []string{"PYTHON_VERSION=3.11.0"},
			expectedOutput: "Images use different languages",
		},
		{
			name:           "Same language and version",
			image1Env:      []string{"NODE_VERSION=18.0.0"},
			image2Env:      []string{"NODE_VERSION=18.0.0"},
			expectedOutput: "Both images use the same Node.js version",
		},
		{
			name:           "Same language, first newer",
			image1Env:      []string{"NODE_VERSION=20.0.0"},
			image2Env:      []string{"NODE_VERSION=18.0.0"},
			expectedOutput: "uses newer Node.js",
		},
		{
			name:           "Same language, second newer",
			image1Env:      []string{"PYTHON_VERSION=3.9.0"},
			image2Env:      []string{"PYTHON_VERSION=3.11.0"},
			expectedOutput: "uses newer Python",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			image1Inspect := createMockImageInspect(tt.image1Env, []string{}, []string{}, "/app", 50000000)
			image2Inspect := createMockImageInspect(tt.image2Env, []string{}, []string{}, "/app", 50000000)

			output := captureOutput(func() {
				PrintImageCompareLanguageResults("image1", image1Inspect, "image2", image2Inspect)
			})

			if !strings.Contains(output, tt.expectedOutput) {
				t.Errorf("Expected output to contain '%s', got: %s", tt.expectedOutput, output)
			}
		})
	}
}

// Tests for PrintImageCompareSizeResults
func TestPrintImageCompareSizeResults(t *testing.T) {
	tests := []struct {
		name           string
		size1          int64
		size2          int64
		expectedOutput string
	}{
		{
			name:           "Same size",
			size1:          100000000,
			size2:          100000000,
			expectedOutput: "Images have the same size",
		},
		{
			name:           "First image smaller",
			size1:          50000000,
			size2:          100000000,
			expectedOutput: "smaller",
		},
		{
			name:           "Second image smaller",
			size1:          100000000,
			size2:          50000000,
			expectedOutput: "smaller",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			image1Inspect := types.ImageInspect{Size: tt.size1}
			image2Inspect := types.ImageInspect{Size: tt.size2}

			output := captureOutput(func() {
				PrintImageCompareSizeResults("image1", image1Inspect, "image2", image2Inspect)
			})

			if !strings.Contains(output, tt.expectedOutput) {
				t.Errorf("Expected output to contain '%s', got: %s", tt.expectedOutput, output)
			}
		})
	}
}

// // Tests for PrintImageCompareLayersResults
func TestPrintImageCompareLayersResults(t *testing.T) {
	tests := []struct {
		name           string
		layers1        int
		layers2        int
		expectedOutput string
	}{
		{
			name:           "Same number of layers",
			layers1:        10,
			layers2:        10,
			expectedOutput: "Images have the same number of layers",
		},
		{
			name:           "First image has fewer layers",
			layers1:        5,
			layers2:        10,
			expectedOutput: "less layers",
		},
		{
			name:           "Second image has fewer layers",
			layers1:        15,
			layers2:        8,
			expectedOutput: "less layers",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			layers1 := make([]string, tt.layers1)
			for i := 0; i < tt.layers1; i++ {
				layers1[i] = fmt.Sprintf("sha256:layer%d", i)
			}

			layers2 := make([]string, tt.layers2)
			for i := 0; i < tt.layers2; i++ {
				layers2[i] = fmt.Sprintf("sha256:layer%d", i)
			}

			image1Inspect := types.ImageInspect{
				RootFS: types.RootFS{Layers: layers1},
			}

			image2Inspect := types.ImageInspect{
				RootFS: types.RootFS{Layers: layers2},
			}

			output := captureOutput(func() {
				PrintImageCompareLayersResults("image1", image1Inspect, "image2", image2Inspect)
			})

			if !strings.Contains(output, tt.expectedOutput) {
				t.Errorf("Expected output to contain '%s', got: %s", tt.expectedOutput, output)
			}
		})
	}
}

// Integration test for PrintImageResults
// func TestPrintImageResults(t *testing.T) {
// 	imageInspect := types.ImageInspect{
// 		RepoTags: []string{"testimage:latest"},
// 		Size:     56900000,
// 		RootFS: types.RootFS{
// 			Layers: []string{"layer1", "layer2", "layer3", "layer4", "layer5", "layer6", "layer7"},
// 		},
// 		Config: &container.Config{
// 			Env:        []string{"NODE_VERSION=18.17.0"},
// 			Cmd:        []string{},
// 			Entrypoint: []string{"node", "index.js"},
// 			WorkingDir: "/app",
// 		},
// 		Author:  "Test Author",
// 		Created: "2025-11-24T12:30:45.123456789Z",
// 		Os:      "linux",
// 	}

// 	output := captureOutput(func() {
// 		PrintImageResults("testimage", imageInspect, false, true)
// 	})

// 	// Verify output contains expected elements
// 	expectedStrings := []string{
// 		"Details of image",
// 		"testimage",
// 		"Tags:",
// 		"Size:",
// 		"N. of Layers:",
// 		"Node.js version:",
// 		"Author:",
// 		"Creation date:",
// 		"OS:",
// 	}

// 	for _, expected := range expectedStrings {
// 		if !strings.Contains(output, expected) {
// 			t.Errorf("Expected output to contain '%s', got: %s", expected, output)
// 		}
// 	}
// }

// // Test for minimal output
// func TestPrintImageResultsMinimal(t *testing.T) {
// 	imageInspect := types.ImageInspect{
// 		RepoTags: []string{"testimage:latest"},
// 		Size:     56900000,
// 		RootFS: types.RootFS{
// 			Layers: []string{"layer1", "layer2"},
// 		},
// 		Config: &container.Config{
// 			Env:        []string{"GO_VERSION=1.21.0"},
// 			Cmd:        []string{},
// 			Entrypoint: []string{"/app/main"},
// 			WorkingDir: "/app",
// 		},
// 		Os: "linux",
// 	}

// 	output := captureOutput(func() {
// 		PrintImageResults("testimage", imageInspect, true, true)
// 	})

// 	// Should NOT contain author and creation date in minimal mode
// 	if strings.Contains(output, "Author:") {
// 		t.Error("Minimal output should not contain Author")
// 	}

// 	if strings.Contains(output, "Creation date:") {
// 		t.Error("Minimal output should not contain Creation date")
// 	}

// 	// Should contain basic info
// 	if !strings.Contains(output, "Size:") {
// 		t.Error("Minimal output should contain Size")
// 	}
// }

// Test suggestions are shown
// func TestPrintImageResultsWithSuggestions(t *testing.T) {
// 	imageInspect := types.ImageInspect{
// 		RepoTags: []string{"testimage:latest"},
// 		Size:     300000000, // Big image
// 		RootFS: types.RootFS{
// 			Layers: make([]string, 15), // Many layers
// 		},
// 		Config: &container.Config{
// 			Env:        []string{"NODE_VERSION=12.0.0"}, // Outdated
// 			Cmd:        []string{},
// 			Entrypoint: []string{},
// 			WorkingDir: "/app",
// 		},
// 		Os: "linux",
// 	}

// 	output := captureOutput(func() {
// 		PrintImageResults("testimage", imageInspect, false, false)
// 	})

// 	// Should show suggestions
// 	expectedSuggestions := []string{
// 		"Improvement suggestions",
// 		"reducing the size",
// 		"multiple layers",
// 		"outdated",
// 	}

// 	for _, expected := range expectedSuggestions {
// 		if !strings.Contains(output, expected) {
// 			t.Errorf("Expected output to contain suggestion '%s', got: %s", expected, output)
// 		}
// 	}
// }

// Test no suggestions when ignored
// func TestPrintImageResultsIgnoreSuggestions(t *testing.T) {
// 	imageInspect := types.ImageInspect{
// 		RepoTags: []string{"testimage:latest"},
// 		Size:     300000000, // Big image
// 		RootFS: types.RootFS{
// 			Layers: make([]string, 15), // Many layers
// 		},
// 		Config: &container.Config{
// 			Env:        []string{"NODE_VERSION=12.0.0"}, // Outdated
// 			Cmd:        []string{},
// 			Entrypoint: []string{},
// 			WorkingDir: "/app",
// 		},
// 		Os: "linux",
// 	}

// 	output := captureOutput(func() {
// 		PrintImageResults("testimage", imageInspect, false, true)
// 	})

// 	// Should NOT show suggestions
// 	if strings.Contains(output, "Improvement suggestions") {
// 		t.Error("Output should not contain suggestions when ignored")
// 	}
// }
