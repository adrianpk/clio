package ssg

import (
	"sort"
	"strings"
)

// Index represents a single generated index page, containing the list of content
// that belongs to it.
type Index struct {
	Path    string    // The output path for the index, e.g., "/news/" or "/blog/".
	Type    string    // Type of index (section, blog, series) to determine sorting.
	Content []Content // The list of content items for this index.
}

// BuildIndexes analyzes all site content and sections to generate the data for all
// required index pages (global, section, blog, and series).
func BuildIndexes(allContent []Content, allSections []Section) []*Index {
	// Use a map for efficient lookup and to avoid duplicate index paths.
	indexes := make(map[string]*Index)

	// Ensure the root index always exists.
	indexes["/"] = &Index{Path: "/", Type: "section", Content: []Content{}}

	// Ensure an index exists for every section defined in the database.
	for _, section := range allSections {
		if _, exists := indexes[section.Path]; !exists {
			indexes[section.Path] = &Index{Path: section.Path, Type: "section", Content: []Content{}}
		}
	}

	// Distribute content into the appropriate indexes.
	for _, content := range allContent {
		kind := strings.ToLower(content.Kind)
		// NOTE: Only these kinds are included in any index.
		if kind != "article" && kind != "blog" && kind != "series" {
			continue
		}

		// Add to its local section index.
		if sectionIndex, ok := indexes[content.SectionPath]; ok {
			sectionIndex.Content = append(sectionIndex.Content, content)
		}

		// Add to the global root index.
		indexes["/"].Content = append(indexes["/"].Content, content)

		// Add to a dedicated blog index if it's a blog post.
		if kind == "blog" {
			basePath := strings.TrimSuffix(content.SectionPath, "/")
			blogPath := basePath + "/blog/"
			if content.SectionPath == "/" {
				blogPath = "/blog/"
			}

			if _, ok := indexes[blogPath]; !ok {
				indexes[blogPath] = &Index{Path: blogPath, Type: "blog", Content: []Content{}}
			}
			indexes[blogPath].Content = append(indexes[blogPath].Content, content)
		}

		// Add to a dedicated series index if it's a series post.
		if kind == "series" && content.Series != "" {
			basePath := strings.TrimSuffix(content.SectionPath, "/")
			seriesPath := basePath + "/" + content.Series + "/"
			if content.SectionPath == "/" {
				seriesPath = "/" + content.Series + "/"
			}

			if _, ok := indexes[seriesPath]; !ok {
				indexes[seriesPath] = &Index{Path: seriesPath, Type: "series", Content: []Content{}}
			}
			indexes[seriesPath].Content = append(indexes[seriesPath].Content, content)
		}
	}

	// Sort each index based on its type.
	for _, index := range indexes {
		switch index.Type {
		case "series":
			// Series are ordered by their predefined sequence number.
			sort.Slice(index.Content, func(i, j int) bool {
				return index.Content[i].SeriesOrder < index.Content[j].SeriesOrder
			})
		default:
			// Section and Blog indexes are ordered chronologically, newest first.
			sort.Slice(index.Content, func(i, j int) bool {
				if index.Content[i].PublishedAt == nil || index.Content[j].PublishedAt == nil {
					return false // Keep original order if dates are missing
				}
				return index.Content[i].PublishedAt.After(*index.Content[j].PublishedAt)
			})
		}
	}

	var result []*Index
	for _, index := range indexes {
		if len(index.Content) > 0 {
			result = append(result, index)
		}
	}

	return result
}
