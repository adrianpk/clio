package ssg

import "html/template"

// PageData holds all the data needed to render a complete HTML page.
type PageData struct {
	HeaderStyle string
	AssetPath   string
	Menu        []Section
	Blocks      *GeneratedBlocks

	// Flags and data for index pages
	IsIndex         bool
	ListPageContent []Content
	Pagination      *PaginationData

	// Content for a single page
	Content PageContent
}

// PageContent holds the specific content to be rendered in the template for a single page.
type PageContent struct {
	Heading     string
	HeaderImage string
	Body        template.HTML
	Kind        string
}

// PaginationData holds data for rendering pagination controls.
type PaginationData struct {
	CurrentPage int
	TotalPages  int
	NextPageURL string
	PrevPageURL string
}