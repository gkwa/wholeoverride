package core

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/go-logr/logr"
)

func FindMarkdownFiles(logger logr.Logger, baseDir string) ([]string, error) {
	var files []string
	var skippedCount int

	err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logger.Error(err, "Error accessing path", "path", path)
			return nil // Continue to next file instead of stopping
		}

		logger.V(2).Info("Encountered file", "path", path, "isDir", info.IsDir())

		if info.IsDir() {
			if strings.Contains(strings.ToLower(path), ".git") ||
				strings.Contains(strings.ToLower(path), ".trash") {
				logger.V(2).
					Info("Skipping directory", "path", path, "reason", "contains .git or .trash")
				return filepath.SkipDir
			}
			return nil
		}

		if strings.HasSuffix(strings.ToLower(info.Name()), ".md") {
			if strings.HasPrefix(info.Name(), ".#") {
				logger.V(2).Info("Skipping temporary file", "path", path)
				skippedCount++
			} else {
				logger.V(2).Info("Including markdown file", "path", path)
				files = append(files, path)
			}
		} else {
			logger.V(2).Info("Skipping non-markdown file", "path", path)
			skippedCount++
		}

		return nil
	})

	logger.V(1).Info("Finished searching for markdown files",
		"fileCount", len(files),
		"skippedCount", skippedCount,
		"totalProcessed", len(files)+skippedCount)

	return files, err
}

func ReadFile(logger logr.Logger, path string) ([]byte, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		logger.Error(err, "Failed to read file", "path", path)
		return nil, err
	}
	logger.V(2).Info("File read successfully", "path", path, "size", len(content))
	return content, nil
}

func WriteFile(logger logr.Logger, path string, content []byte) error {
	err := os.WriteFile(path, content, 0o644)
	if err != nil {
		logger.Error(err, "Failed to write file", "path", path)
	} else {
		logger.V(2).Info("File written successfully", "path", path, "size", len(content))
	}
	return err
}
