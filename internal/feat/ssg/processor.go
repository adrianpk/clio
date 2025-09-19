package ssg

import (
	"bytes"

	"github.com/yuin/goldmark"
)

// MarkdownProcessor is responsible for converting Markdown text to HTML.
type MarkdownProcessor struct {
	parser goldmark.Markdown
}

// NewMarkdownProcessor creates and configures a new Markdown processor.
func NewMarkdownProcessor() *MarkdownProcessor {
	// For now, we use the default goldmark parser.
	// Extensions for syntax highlighting, etc., will be added later.
	md := goldmark.New(
		goldmark.WithExtensions(
			// Add extensions here, e.g., syntax.New()
		),
	)

	return &MarkdownProcessor{
		parser: md,
	}
}

// ToHTML converts a Markdown string to an HTML string.
func (p *MarkdownProcessor) ToHTML(markdown []byte) (string, error) {
	var buf bytes.Buffer
	if err := p.parser.Convert(markdown, &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}