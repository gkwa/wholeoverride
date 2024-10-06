package core

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/go-logr/logr"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

type RecipeInfo struct {
	Title         string
	ImageURL      string
	Creator       string
	IsRemoteImage bool
	UUID          string
}

type CreatorInfo struct {
	Name          string
	ImageURL      string
	IsRemoteImage bool
}

func ParseRecipeFile(logger logr.Logger, path string) (*RecipeInfo, error) {
	content, err := ReadFile(logger, path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	markdown := goldmark.New(goldmark.WithExtensions(meta.Meta))
	context := parser.NewContext()
	markdown.Parser().Parse(text.NewReader(content), parser.WithContext(context))

	metaData := meta.Get(context)
	logger.V(2).Info("Parsed frontmatter", "file", path, "metadata", metaData)

	fileType, ok := metaData["filetype"].(string)
	if !ok || fileType != "recipe" {
		logger.V(2).Info("File is not a recipe", "file", path, "filetype", fileType)
		return nil, nil
	}

	pic, _ := metaData["pic"].(string)
	creator, _ := metaData["creator"].(string)

	if pic == "" {
		logger.V(2).Info("Recipe file has no 'pic' field", "file", path)
	}

	if creator == "" {
		logger.V(2).Info("Recipe file has no 'creator' field", "file", path)
	}

	isRemoteImage := isRemoteURL(pic)

	return &RecipeInfo{
		Title:         strings.TrimSuffix(filepath.Base(path), ".md"),
		ImageURL:      pic,
		Creator:       strings.Trim(creator, "[]"),
		IsRemoteImage: isRemoteImage,
	}, nil
}

func ParseCreatorFile(logger logr.Logger, baseDir, creatorName string) (*CreatorInfo, error) {
	path := filepath.Join(baseDir, creatorName+".md")
	content, err := ReadFile(logger, path)
	if err != nil {
		return nil, err
	}

	markdown := goldmark.New(goldmark.WithExtensions(meta.Meta))
	context := parser.NewContext()
	markdown.Parser().Parse(text.NewReader(content), parser.WithContext(context))

	metaData := meta.Get(context)
	logger.V(2).Info("Parsed creator frontmatter", "file", path, "metadata", metaData)

	pic, _ := metaData["pic"].(string)
	isRemoteImage := isRemoteURL(pic)

	return &CreatorInfo{
		Name:          creatorName,
		ImageURL:      pic,
		IsRemoteImage: isRemoteImage,
	}, nil
}

func isRemoteURL(urlString string) bool {
	u, err := url.Parse(urlString)
	if err != nil {
		return false
	}
	return u.Scheme == "http" || u.Scheme == "https"
}
