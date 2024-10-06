package core

import (
	"strings"
	"testing"

	"github.com/go-logr/logr/testr"
)

func TestSectionMarkdownGenerator_Generate(t *testing.T) {
	logger := testr.New(t)
	generator := NewSectionMarkdownGenerator()

	recipes := []*RecipeInfo{
		{
			Title:         "Test Recipe",
			ImageURL:      "test.jpg",
			Creator:       "Test Creator",
			IsRemoteImage: true,
			UUID:          "test-uuid-1234",
		},
	}

	creators := map[string]*CreatorInfo{
		"Test Creator": {
			Name:          "Test Creator",
			ImageURL:      "creator.jpg",
			IsRemoteImage: false,
		},
	}

	result, err := generator.Generate(logger, recipes, creators)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedTOCLink := "[[#^test-uuid-1234|toc]]"
	if !strings.Contains(result, expectedTOCLink) {
		t.Errorf("Expected TOC link %s not found in result", expectedTOCLink)
	}

	expectedSection := `## [[Test Recipe]]
[[#^test-uuid-1234|toc]]

| Image | Creator |
|-|-|
| ![Test Recipe](test.jpg) | ![[creator.jpg]] [[Test Creator]] |
`
	if !strings.Contains(result, expectedSection) {
		t.Errorf("Expected section not found in result. Got:\n%s", result)
	}
}
