package utils

import (
	"os"
	"path/filepath"
	"strings"
)

var ignoredDirs = map[string]bool{
	".git":         true,
	"node_modules": true,
	"vendor":       true,
	"dist":         true,
	"build":        true,
	".idea":        true,
	".vscode":      true,
	".DS_Store":    true,
}

func GetProjectStructure() (string, error) {
	root, err := os.Getwd()
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	err = walkDirectory(root, "", &builder)
	if err != nil {
		return "", err
	}

	return builder.String(), nil
}

func walkDirectory(path string, prefix string, builder *strings.Builder) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for i, entry := range entries {
		name := entry.Name()

		if ignoredDirs[name] {
			continue
		}

		isLast := i == len(entries)-1
		var connector string

		if isLast {
			connector = "└── "
		} else {
			connector = "├── "
		}

		builder.WriteString(prefix + connector + name + "\n")

		if entry.IsDir() {
			newPrefix := prefix
			if isLast {
				newPrefix += "    "
			} else {
				newPrefix += "│   "
			}
			err := walkDirectory(filepath.Join(path, name), newPrefix, builder)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
