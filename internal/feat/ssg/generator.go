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

	// TODO: We should use config, now is only working as dev-mode
	basePath := filepath.Join(wd, "_workspace", "documents", "markdown")

	for _, content := range contents {
		fileName := content.Slug() + ".md"
		filePath := filepath.Join(basePath, content.SectionPath, fileName)

		frontMatter := make(map[string]interface{})
		frontMatter["title"] = content.Heading
		frontMatter["date"] = content.CreatedAt
		frontMatter["status"] = content.Status

		var tags []string
		for _, t := range content.Tags {
			tags = append(tags, t.Name)
		}
		if len(tags) > 0 {
			frontMatter["tags"] = tags
		}

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
