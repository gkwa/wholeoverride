package core

import (
	"fmt"
	"sort"
	"strings"

	"github.com/go-logr/logr"
)

type TableMarkdownGenerator struct{}

func NewTableMarkdownGenerator() *TableMarkdownGenerator {
	return &TableMarkdownGenerator{}
}

func (g *TableMarkdownGenerator) Generate(logger logr.Logger, recipes []*RecipeInfo, creators map[string]*CreatorInfo) (string, error) {
	var tableRows []string
	tableRows = append(tableRows, "| Recipe Image and Title | Creator's Image |")
	tableRows = append(tableRows, "|------------------------|-----------------|")

	sort.Slice(recipes, func(i, j int) bool {
		return strings.ToLower(recipes[i].Title) < strings.ToLower(recipes[j].Title)
	})

	for _, recipe := range recipes {
		creator, ok := creators[recipe.Creator]
		if !ok {
			logger.V(1).Info("Creator not found", "creator", recipe.Creator)
			continue
		}

		recipeImage := formatImage(recipe.Title, recipe.ImageURL, recipe.IsRemoteImage)
		creatorImage := formatImage(creator.Name, creator.ImageURL, creator.IsRemoteImage)

		tableRows = append(tableRows, fmt.Sprintf("| %s [[%s]] | %s [[%s]] |",
			recipeImage, recipe.Title,
			creatorImage, creator.Name))
	}

	return strings.Join(tableRows, "\n") + "\n\n[Back to top](#top)\n", nil
}
