package core

import (
	"fmt"
	"sort"
	"strings"

	"github.com/go-logr/logr"
)

type SectionMarkdownGenerator struct{}

func NewSectionMarkdownGenerator() *SectionMarkdownGenerator {
	return &SectionMarkdownGenerator{}
}

func (g *SectionMarkdownGenerator) Generate(logger logr.Logger, recipes []*RecipeInfo, creators map[string]*CreatorInfo) (string, error) {
	var sections []string

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

		section := fmt.Sprintf(`## [[%s]]

[Back to top](#toc)


| Image | Creator |
|-|-|
| %s | %s [[%s]] |

`,
			recipe.Title,
			recipeImage,
			creatorImage,
			creator.Name,
		)

		sections = append(sections, section)
	}

	return strings.Join(sections, "\n"), nil
}
