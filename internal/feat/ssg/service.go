package ssg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"

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

	AddTagToContent(ctx context.Context, contentID uuid.UUID, tagName string) error
	RemoveTagFromContent(ctx context.Context, contentID, tagID uuid.UUID) error
	GetTagsForContent(ctx context.Context, contentID uuid.UUID) ([]Tag, error)
	GetContentForTag(ctx context.Context, tagID uuid.UUID) ([]Content, error)

	GenerateMarkdown(ctx context.Context) error
	GenerateHTMLFromContent(ctx context.Context) error
}

// BaseService is the concrete implementation of the Service interface.
type BaseService struct {
	*am.Service
	repo Repo
	gen  *Generator
}

// NewService creates a new BaseService.
func NewService(repo Repo, gen *Generator, opts ...am.Option) *BaseService {
	return &BaseService{
		Service: am.NewService("ssg-svc", opts...),
		repo:    repo,
		gen:     gen,
	}
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

	processor := NewMarkdownProcessor()
	htmlPath := svc.Cfg().StrValOrDef(am.Key.SSGHTMLPath, "_workspace/documents/html")

	for _, content := range contents {
		if content.Draft {
			svc.Log().Debug("Skipping draft content", "slug", content.Slug())
			continue
		}

		htmlBody, err := processor.ToHTML([]byte(content.Body))
		if err != nil {
			svc.Log().Error("Error converting markdown to HTML", "slug", content.Slug(), "error", err)
			continue // or return err, depending on desired behavior
		}

		// This is a simplified pathing logic. We will need to replicate the logic from the drafts for content types.
		outputPath := filepath.Join(htmlPath, content.SectionPath, content.Slug()+".html")

		if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
			svc.Log().Error("Error creating directory for HTML file", "path", outputPath, "error", err)
			continue
		}

		if err := os.WriteFile(outputPath, []byte(htmlBody), 0644); err != nil {
			svc.Log().Error("Error writing HTML file", "path", outputPath, "error", err)
			continue
		}
	}

	svc.Log().Info("Service HTML generation finished")
	return nil
}


// Content related

func (svc *BaseService) CreateContent(ctx context.Context, content *Content) error {
	return svc.repo.CreateContent(ctx, content)
}

func (svc *BaseService) GetAllContentWithMeta(ctx context.Context) ([]Content, error) {
	return svc.repo.GetAllContentWithMeta(ctx)
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
