package core

import (
	"bytes"

	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

type FrontmatterParser interface {
	Extract(content []byte) (map[string]interface{}, []byte, error)
}

type GoldmarkFrontmatterParser struct{}

func NewGoldmarkFrontmatterParser() *GoldmarkFrontmatterParser {
	return &GoldmarkFrontmatterParser{}
}

func (gfp *GoldmarkFrontmatterParser) Extract(
	content []byte,
) (map[string]interface{}, []byte, error) {
	markdown := goldmark.New(goldmark.WithExtensions(meta.Meta))
	context := parser.NewContext()
	reader := text.NewReader(content)
	doc := markdown.Parser().Parse(reader, parser.WithContext(context))

	frontmatter := meta.Get(context)
	if frontmatter == nil {
		frontmatter = make(map[string]interface{})
	}

	var buf bytes.Buffer
	err := markdown.Renderer().Render(&buf, content, doc)
	if err != nil {
		return nil, nil, err
	}

	return frontmatter, buf.Bytes(), nil
}
