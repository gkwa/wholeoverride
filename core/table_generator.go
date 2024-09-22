package core

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/go-logr/logr"
)

type TableRow struct {
	RecipeImage  string
	RecipeTitle  string
	CreatorImage string
	CreatorName  string
}

func GenerateTable(logger logr.Logger, baseDir string) error {
	logger.V(1).Info("Starting table generation", "baseDir", baseDir)

	files, err := FindMarkdownFiles(logger, baseDir)
	if err != nil {
		return fmt.Errorf("error finding markdown files: %w", err)
	}

	logger.Info("Found markdown files", "count", len(files))

	var tableRows []TableRow
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

		creator, err := ParseCreatorFile(logger, baseDir, recipe.Creator)
		if err != nil {
			logger.Error(err, "Failed to parse creator file, skipping", "creator", recipe.Creator)
			skippedCount++
			continue
		}

		logger.V(1).Info("Parsed creator file", "name", creator.Name)

		recipeImage := formatImage(recipe.Title, recipe.ImageURL, recipe.IsRemoteImage)
		creatorImage := formatImage(creator.Name, creator.ImageURL, creator.IsRemoteImage)

		tableRows = append(tableRows, TableRow{
			RecipeImage:  recipeImage,
			RecipeTitle:  recipe.Title,
			CreatorImage: creatorImage,
			CreatorName:  creator.Name,
		})

		logger.V(1).Info("Added row to table", "recipe", recipe.Title, "creator", creator.Name)
		processedCount++
	}

	// Sort tableRows by creator name (case-insensitive)
	sort.Slice(tableRows, func(i, j int) bool {
		return strings.ToLower(tableRows[i].CreatorName) < strings.ToLower(tableRows[j].CreatorName)
	})

	// Generate the final sorted table
	var finalTableRows []string
	finalTableRows = append(finalTableRows, "| Recipe Image and Title | Creator's Image |")
	finalTableRows = append(finalTableRows, "|------------------------|-----------------|")

	for _, row := range tableRows {
		finalTableRows = append(finalTableRows, fmt.Sprintf("| %s [[%s]] | %s [[%s]] |",
			row.RecipeImage, row.RecipeTitle,
			row.CreatorImage, row.CreatorName))
	}

	logger.Info("Table generation summary",
		"totalFiles", len(files),
		"processedFiles", processedCount,
		"skippedFiles", skippedCount,
		"tableRows", len(tableRows))

	table := strings.Join(finalTableRows, "\n")
	outputPath := filepath.Join(baseDir, "recipeindex.md")
	err = WriteFile(logger, outputPath, []byte(table))
	if err != nil {
		return fmt.Errorf("error writing output file: %w", err)
	}

	logger.V(1).Info("Table generation completed", "outputFile", outputPath)
	return nil
}

func formatImage(name, url string, isRemote bool) string {
	if isRemote {
		return fmt.Sprintf("![%s](%s)", name, url)
	}
	return fmt.Sprintf("![[%s]]", url)
}
