package ssg

import "html/template"

// PageData holds all the data needed to render a complete HTML page.
type PageData struct {
	HeaderStyle string
	AssetPath   string
	Menu        []Section
	Content     PageContent
	Blocks      *GeneratedBlocks
}

// PageContent holds the specific content to be rendered in the template.
type PageContent struct {
	Heading     string
	HeaderImage string
	Body        template.HTML
	Kind        string
}
