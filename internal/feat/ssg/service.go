package ssg

import (
	"bytes"
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/adrianpk/clio/internal/am"
	"github.com/google/uuid"
)

// Service defines the interface for the ssg service.
type Service interface {
	CreateContent(ctx context.Context, content *Content) error
	GetAllContentWithMeta(ctx context.Context) ([]Content, error)
	GetContent(ctx context.Context, id uuid.UUID) (Content, error)
	UpdateContent(ctx context.Context, content *Content) error
	DeleteContent(ctx context.Context, id uuid.UUID) error

	CreateSection(ctx context.Context, section Section) error
	GetSection(ctx context.Context, id uuid.UUID) (Section, error)
	GetSections(ctx context.Context) ([]Section, error)
	UpdateSection(ctx context.Context, section Section) error
	DeleteSection(ctx context.Context, id uuid.UUID) error

	CreateLayout(ctx context.Context, layout Layout) error
	GetLayout(ctx context.Context, id uuid.UUID) (Layout, error)
	GetAllLayouts(ctx context.Context) ([]Layout, error)
	UpdateLayout(ctx context.Context, layout Layout) error
	DeleteLayout(ctx context.Context, id uuid.UUID) error

	CreateTag(ctx context.Context, tag Tag) error
	GetTag(ctx context.Context, id uuid.UUID) (Tag, error)
	GetTagByName(ctx context.Context, name string) (Tag, error)
	GetAllTags(ctx context.Context) ([]Tag, error)
	UpdateTag(ctx context.Context, tag Tag) error
	DeleteTag(ctx context.Context, id uuid.UUID) error

	CreateParam(ctx context.Context, param *Param) error
	GetParam(ctx context.Context, id uuid.UUID) (Param, error)
	GetParamByName(ctx context.Context, name string) (Param, error)
	GetParamByRefKey(ctx context.Context, refKey string) (Param, error)
	ListParams(ctx context.Context) ([]Param, error)
	UpdateParam(ctx context.Context, param *Param) error
	DeleteParam(ctx context.Context, id uuid.UUID) error

	// Image related
	CreateImage(ctx context.Context, image *Image) error
	GetImage(ctx context.Context, id uuid.UUID) (Image, error)
	GetImageByShortID(ctx context.Context, shortID string) (Image, error)
	ListImages(ctx context.Context) ([]Image, error)
	UpdateImage(ctx context.Context, image *Image) error
	DeleteImage(ctx context.Context, id uuid.UUID) error

	// ImageVariant related
	CreateImageVariant(ctx context.Context, variant *ImageVariant) error
	GetImageVariant(ctx context.Context, id uuid.UUID) (ImageVariant, error)
	ListImageVariantsByImageID(ctx context.Context, imageID uuid.UUID) ([]ImageVariant, error)
	UpdateImageVariant(ctx context.Context, variant *ImageVariant) error
	DeleteImageVariant(ctx context.Context, id uuid.UUID) error

	// ContentTag related
	AddTagToContent(ctx context.Context, contentID uuid.UUID, tagName string) error
	RemoveTagFromContent(ctx context.Context, contentID, tagID uuid.UUID) error
	GetTagsForContent(ctx context.Context, contentID uuid.UUID) ([]Tag, error)
	GetContentForTag(ctx context.Context, tagID uuid.UUID) ([]Content, error)

	GenerateMarkdown(ctx context.Context) error
	GenerateHTMLFromContent(ctx context.Context) error
	Publish(ctx context.Context, commitMessage string) (string, error)
	Plan(ctx context.Context) (PlanReport, error)
}

// BaseService is the concrete implementation of the Service interface.
type BaseService struct {
	*am.Service
	assetsFS embed.FS
	repo     Repo
	gen      *Generator
	pub      Publisher
	pm       *ParamManager
}

// NewService creates a new BaseService.
func NewService(assetsFS embed.FS, repo Repo, gen *Generator, publisher Publisher, pm *ParamManager, opts ...am.Option) *BaseService {
	return &BaseService{
		Service:  am.NewService("ssg-svc", opts...),
		assetsFS: assetsFS,
		repo:     repo,
		gen:      gen,
		pub:      publisher,
		pm:       pm,
	}
}

// Publish delegates the publishing task to the underlying pub.
func (svc *BaseService) Publish(ctx context.Context, commitMessage string) (string, error) {
	svc.Log().Info("Service starting publish process")

	// For now, we build the config from the application's configuration.
	cfg := PublisherConfig{
		RepoURL: svc.pm.Get(ctx, am.Key.SSGPublishRepoURL, ""),
		Branch:  svc.pm.Get(ctx, am.Key.SSGPublishBranch, ""),
		Auth: am.GitAuth{
			// NOTE: This is oversimplified. We need to work out a bit more here.
			Method: am.AuthToken,
			Token:  svc.pm.Get(ctx, am.Key.SSGPublishAuthToken, ""),
		},
		CommitAuthor: am.GitCommit{
			UserName:  svc.pm.Get(ctx, am.Key.SSGPublishCommitUserName, ""),
			UserEmail: svc.pm.Get(ctx, am.Key.SSGPublishCommitUserEmail, ""),
			Message:   svc.pm.Get(ctx, am.Key.SSGPublishCommitMessage, ""),
		},
	}

	// Override commit message if provided in the request body
	if commitMessage != "" {
		cfg.CommitAuthor.Message = commitMessage
	}

	// Get the output directory for HTML files, which is the source for publishing
	sourceDir := svc.Cfg().StrValOrDef(am.Key.SSGHTMLPath, "_workspace/documents/html")

	commitURL, err := svc.pub.Publish(ctx, cfg, sourceDir)
	if err != nil {
		return "", fmt.Errorf("cannot publish site: %w", err)
	}

	svc.Log().Info("Service publish process finished successfully", "commit_url", commitURL)
	return commitURL, nil
}

// Plan delegates the plan task to the underlying pub.
func (svc *BaseService) Plan(ctx context.Context) (PlanReport, error) {
	svc.Log().Info("Service starting plan process")

	// For now, we build the config from the application's configuration.
	cfg := PublisherConfig{
		RepoURL: svc.pm.Get(ctx, am.Key.SSGPublishRepoURL, ""),
		Branch:  svc.pm.Get(ctx, am.Key.SSGPublishBranch, ""),
		Auth: am.GitAuth{
			// NOTE: This is oversimplified. We need to work out a bit more here.
			Method: am.AuthToken,
			Token:  svc.pm.Get(ctx, am.Key.SSGPublishAuthToken, ""),
		},
		CommitAuthor: am.GitCommit{
			UserName:  svc.pm.Get(ctx, am.Key.SSGPublishCommitUserName, ""),
			UserEmail: svc.pm.Get(ctx, am.Key.SSGPublishCommitUserEmail, ""),
			Message:   svc.pm.Get(ctx, am.Key.SSGPublishCommitMessage, ""),
		},
	}

	// Get the output directory for HTML files, which is the source for planning
	sourceDir := svc.Cfg().StrValOrDef(am.Key.SSGHTMLPath, "_workspace/documents/html")

	report, err := svc.pub.Plan(ctx, cfg, sourceDir)
	if err != nil {
		return PlanReport{}, fmt.Errorf("cannot plan site: %w", err)
	}

	svc.Log().Info("Service plan process finished successfully", "summary", report.Summary)
	return report, nil
}

// GenerateMarkdown generates markdown files from the content in the database.
func (svc *BaseService) GenerateMarkdown(ctx context.Context) error {
	svc.Log().Info("Service starting markdown generation")

	contents, err := svc.repo.GetAllContentWithMeta(ctx)
	if err != nil {
		return fmt.Errorf("cannot get all content with meta: %w", err)
	}

	if err := svc.gen.Generate(contents); err != nil {
		return fmt.Errorf("cannot generate markdown: %w", err)
	}

	svc.Log().Info("Service markdown generation finished")
	return nil
}

// GenerateHTMLFromContent generates HTML files from the content in the database.
func (svc *BaseService) GenerateHTMLFromContent(ctx context.Context) error {
	svc.Log().Info("Service starting HTML generation")

	contents, err := svc.repo.GetAllContentWithMeta(ctx)
	if err != nil {
		return fmt.Errorf("cannot get all content with meta: %w", err)
	}

	// Set placeholder for content without image
	for i := range contents {
		if contents[i].Image == "" {
			contents[i].Image = "/static/img/placeholder.png"
		}
	}

	sections, err := svc.repo.GetSections(ctx)
	if err != nil {
		return fmt.Errorf("cannot get sections: %w", err)
	}

	var menuSections []Section
	for _, s := range sections {
		if s.Name != "root" {
			menuSections = append(menuSections, s)
		}
	}

	layoutPath := svc.Cfg().StrValOrDef(am.Key.SSGLayoutPath, "assets/ssg/layout/layout.html")
	tmpl, err := template.ParseFS(svc.assetsFS,
		layoutPath,
		"assets/ssg/partial/list.tmpl",
		"assets/ssg/partial/blocks.tmpl",
		"assets/ssg/partial/article-blocks.tmpl",
		"assets/ssg/partial/blog-blocks.tmpl",
		"assets/ssg/partial/series-blocks.tmpl",
		"assets/ssg/partial/pagination.tmpl",
		"assets/ssg/partial/google-search.tmpl",
	)
	if err != nil {
		return fmt.Errorf("cannot parse template from embedded fs: %w", err)
	}

	processor := NewMarkdownProcessor()
	htmlPath := svc.Cfg().StrValOrDef(am.Key.SSGHTMLPath, "_workspace/documents/html")

	if err := CopyStaticAssets(svc.assetsFS, htmlPath); err != nil {
		return fmt.Errorf("cannot copy static assets: %w", err)
	}

	headerStyle := svc.Cfg().StrValOrDef(am.Key.SSGHeaderStyle, "boxed", true)
	imageExtensions := []string{".png", ".jpg", ".jpeg", ".webp"}

	// Prepare SearchData
	searchData := SearchData{
		Provider: "google", // O el proveedor que corresponda
		Enabled:  svc.Cfg().BoolVal(am.Key.SSGSearchGoogleEnabled, false),
		ID:       svc.Cfg().StrValOrDef(am.Key.SSGSearchGoogleID, ""),
	}
	svc.Log().Info("SearchData values", "enabled", searchData.Enabled, "id", searchData.ID) // Línea de log modificada

	for _, content := range contents {
		svc.Log().Debug("Processing content for HTML generation", "slug", content.Slug(), "section_path", content.SectionPath)
		if content.Draft {
			svc.Log().Debug("Skipping draft content", "slug", content.Slug())
			continue
		}

		// Paths and Asset Logic
		contentDir := filepath.Join(htmlPath, content.SectionPath, content.Slug())
		contentImgDir := filepath.Join(contentDir, "img")
		headerImagePath := ""

		foundSpecificHeader := false
		for _, ext := range imageExtensions {
			checkPath := filepath.Join("assets", "content", content.SectionPath, content.Slug(), "img", "header"+ext)
			if f, err := svc.assetsFS.Open(checkPath); err == nil {
				f.Close()
				if err := os.MkdirAll(contentImgDir, 0755); err != nil {
					return fmt.Errorf("cannot create img directory: %w", err)
				}
				dst := filepath.Join(contentImgDir, "header"+ext)
				if err := copyFile(svc.assetsFS, checkPath, dst); err != nil {
					return fmt.Errorf("cannot copy specific header: %w", err)
				}
				headerImagePath = "img/header" + ext
				foundSpecificHeader = true
				break
			}
		}

		if !foundSpecificHeader {
			headerImagePath = "/static/img/header.png"
		}

		assetPath := "/"

		htmlBody, err := processor.ToHTML([]byte(content.Body))
		if err != nil {
			svc.Log().Error("Error converting markdown to HTML", "slug", content.Slug(), "error", err)
			continue
		}

		if headerStyle == "boxed" || headerStyle == "overlay" {
			htmlBody = svc.removeFirstH1(htmlBody)
		}

		pageContent := PageContent{
			Heading:     content.Heading,
			HeaderImage: headerImagePath,
			Body:        template.HTML(htmlBody),
			Kind:        content.Kind,
		}

		blocks := BuildBlocks(content, contents, int(svc.Cfg().IntVal(am.Key.SSGBlocksMaxItems, 5)))

		data := PageData{
			HeaderStyle: headerStyle,
			AssetPath:   assetPath,
			Menu:        menuSections,
			Content:     pageContent,
			Blocks:      blocks,
			Search:      searchData,
		}

		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, data); err != nil {
			svc.Log().Error("Error executing template for content", "slug", content.Slug(), "error", err)
			continue
		}

		outputPath := filepath.Join(htmlPath, content.SectionPath, content.Slug(), "index.html")

		if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
			svc.Log().Error("Error creating directory for HTML file", "path", outputPath, "error", err)
			continue
		}

		if err := os.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
			svc.Log().Error("Error writing index HTML file", "path", outputPath, "error", err)
			continue
		}
	}

	// Generate index pages
	svc.Log().Info("Building site indexes...")
	indexes := BuildIndexes(contents, sections)

	// Create a lookup map for manual index pages
	manualIndexPages := make(map[string]bool)
	for _, c := range contents {
		if strings.ToLower(c.Kind) == "page" && c.Slug() == "index" {
			manualIndexPages[c.SectionPath] = true
		}
	}

	postsPerPage := int(svc.Cfg().IntVal(am.Key.SSGIndexMaxItems, 9))

	for _, index := range indexes {
		// Check if a manual index page exists for this path
		if manualIndexPages[index.Path] {
			svc.Log().Info(fmt.Sprintf("Skipping index generation for '%s': manual index page found.", index.Path))
			continue
		}

		// Paginate the content
		totalContent := len(index.Content)
		if totalContent == 0 {
			continue
		}
		totalPages := (totalContent + postsPerPage - 1) / postsPerPage

		for page := 1; page <= totalPages; page++ {
			start := (page - 1) * postsPerPage
			end := start + postsPerPage
			if end > totalContent {
				end = totalContent
			}
			pageContent := index.Content[start:end]

			// Determine output path for the index page
			var outputPath string
			if page == 1 {
				outputPath = filepath.Join(htmlPath, index.Path, "index.html")
			} else {
				outputPath = filepath.Join(htmlPath, index.Path, "page", fmt.Sprintf("%d", page), "index.html")
			}

			assetPath := "/"

			// Prepare pagination data
			pagination := &PaginationData{
				CurrentPage: page,
				TotalPages:  totalPages,
			}
			if page > 1 {
				if page == 2 {
					pagination.PrevPageURL = assetPath + strings.TrimSuffix(index.Path, "/")
				} else {
					pagination.PrevPageURL = fmt.Sprintf("%spage/%d", assetPath, page-1)
				}
			}
			if page < totalPages {
				pagination.NextPageURL = fmt.Sprintf("%spage/%d", assetPath, page+1)
			}

			data := PageData{
				HeaderStyle:     headerStyle,
				AssetPath:       assetPath,
				Menu:            menuSections,
				IsIndex:         true,
				ListPageContent: pageContent,
				Pagination:      pagination,
				Search:          searchData,
			}

			var buf bytes.Buffer
			if err := tmpl.Execute(&buf, data); err != nil {
				svc.Log().Error("Error executing template for index", "path", index.Path, "error", err)
				continue
			}

			if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
				svc.Log().Error("Error creating directory for index file", "path", outputPath, "error", err)
				continue
			}

			if err := os.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
				svc.Log().Error("Error writing index HTML file", "path", outputPath, "error", err)
				continue
			}
		}
	}

	svc.Log().Info("Service HTML generation finished")
	return nil
}

// Content related

func (svc *BaseService) CreateContent(ctx context.Context, content *Content) error {
	return svc.repo.CreateContent(ctx, content)
}

func (svc *BaseService) GetContent(ctx context.Context, id uuid.UUID) (Content, error) {
	return svc.repo.GetContent(ctx, id)
}

func (svc *BaseService) UpdateContent(ctx context.Context, content *Content) error {
	return svc.repo.UpdateContent(ctx, content)
}

func (svc *BaseService) DeleteContent(ctx context.Context, id uuid.UUID) error {
	return svc.repo.DeleteContent(ctx, id)
}

func (svc *BaseService) GetAllContentWithMeta(ctx context.Context) ([]Content, error) {
	return svc.repo.GetAllContentWithMeta(ctx)
}

// Section related
func (svc *BaseService) CreateSection(ctx context.Context, section Section) error {
	return svc.repo.CreateSection(ctx, section)
}

func (svc *BaseService) GetSection(ctx context.Context, id uuid.UUID) (Section, error) {
	return svc.repo.GetSection(ctx, id)
}

func (svc *BaseService) GetSections(ctx context.Context) ([]Section, error) {
	return svc.repo.GetSections(ctx)
}

func (svc *BaseService) UpdateSection(ctx context.Context, section Section) error {
	return svc.repo.UpdateSection(ctx, section)
}

func (svc *BaseService) DeleteSection(ctx context.Context, id uuid.UUID) error {
	return svc.repo.DeleteSection(ctx, id)
}

// Layout related
func (svc *BaseService) CreateLayout(ctx context.Context, layout Layout) error {
	return svc.repo.CreateLayout(ctx, layout)
}

func (svc *BaseService) GetLayout(ctx context.Context, id uuid.UUID) (Layout, error) {
	return svc.repo.GetLayout(ctx, id)
}

func (svc *BaseService) GetAllLayouts(ctx context.Context) ([]Layout, error) {
	return svc.repo.GetAllLayouts(ctx)
}

func (svc *BaseService) UpdateLayout(ctx context.Context, layout Layout) error {
	return svc.repo.UpdateLayout(ctx, layout)
}

func (svc *BaseService) DeleteLayout(ctx context.Context, id uuid.UUID) error {
	return svc.repo.DeleteLayout(ctx, id)
}

// Tag related
func (svc *BaseService) CreateTag(ctx context.Context, tag Tag) error {
	return svc.repo.CreateTag(ctx, tag)
}

func (svc *BaseService) GetTag(ctx context.Context, id uuid.UUID) (Tag, error) {
	return svc.repo.GetTag(ctx, id)
}

func (svc *BaseService) GetTagByName(ctx context.Context, name string) (Tag, error) {
	return svc.repo.GetTagByName(ctx, name)
}

func (svc *BaseService) GetAllTags(ctx context.Context) ([]Tag, error) {
	return svc.repo.GetAllTags(ctx)
}

func (svc *BaseService) UpdateTag(ctx context.Context, tag Tag) error {
	return svc.repo.UpdateTag(ctx, tag)
}

func (svc *BaseService) DeleteTag(ctx context.Context, id uuid.UUID) error {
	return svc.repo.DeleteTag(ctx, id)
}

// Param related
func (svc *BaseService) CreateParam(ctx context.Context, param *Param) error {
	return svc.repo.CreateParam(ctx, param)
}

func (svc *BaseService) GetParam(ctx context.Context, id uuid.UUID) (Param, error) {
	return svc.repo.GetParam(ctx, id)
}

func (svc *BaseService) GetParamByName(ctx context.Context, name string) (Param, error) {
	return svc.repo.GetParamByName(ctx, name)
}

func (svc *BaseService) GetParamByRefKey(ctx context.Context, refKey string) (Param, error) {
	return svc.repo.GetParamByRefKey(ctx, refKey)
}

func (svc *BaseService) ListParams(ctx context.Context) ([]Param, error) {
	return svc.repo.ListParams(ctx)
}

func (svc *BaseService) UpdateParam(ctx context.Context, param *Param) error {
	return svc.repo.UpdateParam(ctx, param)
}

func (svc *BaseService) DeleteParam(ctx context.Context, id uuid.UUID) error {
	return svc.repo.DeleteParam(ctx, id)
}

// Image related
func (svc *BaseService) CreateImage(ctx context.Context, image *Image) error {
	return svc.repo.CreateImage(ctx, image)
}

func (svc *BaseService) GetImage(ctx context.Context, id uuid.UUID) (Image, error) {
	return svc.repo.GetImage(ctx, id)
}

func (svc *BaseService) GetImageByShortID(ctx context.Context, shortID string) (Image, error) {
	return svc.repo.GetImageByShortID(ctx, shortID)
}

func (svc *BaseService) ListImages(ctx context.Context) ([]Image, error) {
	return svc.repo.ListImages(ctx)
}

func (svc *BaseService) UpdateImage(ctx context.Context, image *Image) error {
	return svc.repo.UpdateImage(ctx, image)
}

func (svc *BaseService) DeleteImage(ctx context.Context, id uuid.UUID) error {
	return svc.repo.DeleteImage(ctx, id)
}

// ImageVariant related
func (svc *BaseService) CreateImageVariant(ctx context.Context, variant *ImageVariant) error {
	return svc.repo.CreateImageVariant(ctx, variant)
}

func (svc *BaseService) GetImageVariant(ctx context.Context, id uuid.UUID) (ImageVariant, error) {
	return svc.repo.GetImageVariant(ctx, id)
}

func (svc *BaseService) ListImageVariantsByImageID(ctx context.Context, imageID uuid.UUID) ([]ImageVariant, error) {
	return svc.repo.ListImageVariantsByImageID(ctx, imageID)
}

func (svc *BaseService) UpdateImageVariant(ctx context.Context, variant *ImageVariant) error {
	return svc.repo.UpdateImageVariant(ctx, variant)
}

func (svc *BaseService) DeleteImageVariant(ctx context.Context, id uuid.UUID) error {
	return svc.repo.DeleteImageVariant(ctx, id)
}

// ContentTag related
func (svc *BaseService) AddTagToContent(ctx context.Context, contentID uuid.UUID, tagName string) error {
	tag, err := svc.repo.GetTagByName(ctx, tagName)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("error getting tag by name: %w", err)
	}

	if tag.IsZero() {
		newTag := NewTag(tagName)
		newTag.GenCreateValues()
		err = svc.repo.CreateTag(ctx, newTag)
		if err != nil {
			return fmt.Errorf("error creating tag: %w", err)
		}
		tag = newTag
	}

	return svc.repo.AddTagToContent(ctx, contentID, tag.ID)
}

func (svc *BaseService) RemoveTagFromContent(ctx context.Context, contentID, tagID uuid.UUID) error {
	return svc.repo.RemoveTagFromContent(ctx, contentID, tagID)
}

func (svc *BaseService) GetTagsForContent(ctx context.Context, contentID uuid.UUID) ([]Tag, error) {
	return svc.repo.GetTagsForContent(ctx, contentID)
}

func (svc *BaseService) GetContentForTag(ctx context.Context, tagID uuid.UUID) ([]Content, error) {
	return svc.repo.GetContentForTag(ctx, tagID)
}

var firstH1Regex = regexp.MustCompile(`(?i)<h1[^>]*>.*?</h1>`)

// removeFirstH1 removes the first <h1>...</h1> tag from an HTML string.
func (svc *BaseService) removeFirstH1(htmlContent string) string {
	return firstH1Regex.ReplaceAllStringFunc(htmlContent, func(match string) string {
		// Only replace the first occurrence
		if strings.HasPrefix(htmlContent, match) {
			return ""
		}
		return match
	})
}
