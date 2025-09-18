package ssg

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/adrianpk/clio/internal/am"
)

type Generator struct {
	am.Core
}

func NewGenerator(opts ...am.Option) *Generator {
	core := am.NewCore("ssg-generator", opts...)
	g := &Generator{
		Core: core,
	}
	return g
}

func (g *Generator) Generate(contents []Content) error {
	g.Log().Info("Starting markdown generation")

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	basePath := filepath.Join(wd, "_workspace", "documents", "markdown")

	for _, content := range contents {
		fileName := content.Slug() + ".md"
		filePath := filepath.Join(basePath, content.SectionPath, fileName)

		// --- Frontmatter Generation ---
		frontMatter := make(map[string]interface{})

		// Core Identity & Content
		frontMatter["slug"] = content.Slug()
		frontMatter["permalink"] = "" // TODO: Construct full permalink
		frontMatter["excerpt"] = content.Meta.Description // Using description as a stand-in
		frontMatter["summary"] = "" // TODO: Add a summary field if needed

		// Presentation & Layout
		frontMatter["image"] = "" // TODO: Add image field to a model
		frontMatter["social-image"] = "" // TODO: Add social-image field
		frontMatter["layout"] = content.SectionName // Assuming layout is related to section
		frontMatter["table-of-contents"] = content.Meta.TableOfContents

		// SEO & Discovery
		frontMatter["title"] = content.Heading
		frontMatter["description"] = content.Meta.Description
		frontMatter["keywords"] = content.Meta.Keywords
		frontMatter["robots"] = content.Meta.Robots
		frontMatter["canonical-url"] = content.Meta.CanonicalURL
		frontMatter["sitemap"] = content.Meta.Sitemap

		// Features & Engagement
		frontMatter["featured"] = content.Featured
		frontMatter["share"] = content.Meta.Share
		frontMatter["comments"] = content.Meta.Comments

		// Localization
		frontMatter["locale"] = "" // TODO: Add locale field

		// Timestamps & Status
		frontMatter["draft"] = content.Draft
		frontMatter["published-at"] = content.PublishedAt
		frontMatter["created-at"] = content.CreatedAt
		frontMatter["updated-at"] = content.UpdatedAt

		// Taxonomy
		var tags []string
		for _, t := range content.Tags {
			tags = append(tags, t.Name)
		}
		if len(tags) > 0 {
			frontMatter["tags"] = tags
		}
		// --- End of Frontmatter ---

		yamlBytes, err := yaml.Marshal(frontMatter)
		if err != nil {
			g.Log().Error("Cannot marshal front matter", "error", err, "content_id", content.GetShortID())
			continue
		}

		fileContent := fmt.Sprintf("---\n%s---\n%s", string(yamlBytes), content.Body)

		dir := filepath.Dir(filePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			g.Log().Error("Cannot create directory", "error", err, "path", dir)
			continue
		}

		if err := os.WriteFile(filePath, []byte(fileContent), 0644); err != nil {
			g.Log().Error("Cannot write file", "error", err, "path", filePath)
			continue
		}

		g.Log().Debug("Generated file", "path", filePath)
	}

	g.Log().Info("Markdown generation finished")

	return nil
}
