package ssg

import "html/template"

// PageData holds all the data needed to render a complete HTML page.
type PageData struct {
	Menu    []Section
	Content PageContent
}

// PageContent holds the specific content to be rendered in the template.
type PageContent struct {
	Heading string
	Body    template.HTML
}
