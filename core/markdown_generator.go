package core

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
)

type MarkdownGenerator interface {
	Generate(
		logger logr.Logger,
		recipes []*RecipeInfo,
		creators map[string]*CreatorInfo,
	) (string, error)
}

func GenerateMarkdownWithFormat(logger logr.Logger, baseDir, format string) error {
	var generator MarkdownGenerator
	switch format {
	case "sections":
		generator = NewSectionMarkdownGenerator()
	case "table":
		generator = NewTableMarkdownGenerator()
	default:
		return fmt.Errorf("invalid format specified: %s", format)
	}

	return GenerateMarkdown(logger, baseDir, generator)
}

func GenerateMarkdown(logger logr.Logger, baseDir string, generator MarkdownGenerator) error {
	logger.V(1).Info("Starting markdown generation", "baseDir", baseDir)

	files, err := FindMarkdownFiles(logger, baseDir)
	if err != nil {
		return fmt.Errorf("error finding markdown files: %w", err)
	}

	logger.Info("Found markdown files", "count", len(files))

	var recipes []*RecipeInfo
	creators := make(map[string]*CreatorInfo)
	processedCount := 0
	skippedCount := 0

	for _, file := range files {
		logger.V(1).Info("Processing file", "file", file)

		recipe, err := ParseRecipeFile(logger, file)
		if err != nil {
			logger.Error(err, "Failed to parse recipe file, skipping", "file", file)
			skippedCount++
			continue
		}
		if recipe == nil {
			logger.V(2).Info("Skipping non-recipe file", "file", file)
			skippedCount++
			continue
		}

		logger.V(1).Info("Parsed recipe file", "title", recipe.Title, "creator", recipe.Creator)

		if recipe.Creator == "" {
			logger.V(2).Info("Skipping recipe with no creator", "file", file)
			skippedCount++
			continue
		}

		if _, ok := creators[recipe.Creator]; !ok {
			creator, err := ParseCreatorFile(logger, baseDir, recipe.Creator)
			if err != nil {
				logger.Error(
					err,
					"Failed to parse creator file, skipping",
					"creator",
					recipe.Creator,
				)
				skippedCount++
				continue
			}
			creators[recipe.Creator] = creator
		}

		if recipe.UUID == "" {
			recipe.UUID = uuid.New().String()
		}

		recipes = append(recipes, recipe)
		processedCount++
	}

	logger.Info("Markdown generation summary",
		"totalFiles", len(files),
		"processedFiles", processedCount,
		"skippedFiles", skippedCount,
		"recipeCount", len(recipes))

	content, err := generator.Generate(logger, recipes, creators)
	if err != nil {
		return fmt.Errorf("error generating markdown: %w", err)
	}

	toc := generateTOC(recipes)
	content = "\n\n\n\n\n\n" + "# TOC\n" + toc + "\n" + content

	outputPath := filepath.Join(baseDir, "recipeindex.md")
	err = WriteFile(logger, outputPath, []byte(content))
	if err != nil {
		return fmt.Errorf("error writing output file: %w", err)
	}

	logger.V(1).Info("Markdown generation completed", "outputFile", outputPath)
	return nil
}

func generateTOC(recipes []*RecipeInfo) string {
	var toc []string
	for _, recipe := range recipes {
		toc = append(toc, fmt.Sprintf("- [[#%s|%s]] ^%s", recipe.Title, recipe.Title, recipe.UUID))
	}

	sort.Slice(toc, func(i, j int) bool {
		return strings.ToLower(toc[i]) < strings.ToLower(toc[j])
	})

	return strings.Join(toc, "\n")
}
