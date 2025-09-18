package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/adrianpk/clio/internal/feat/ssg"
	"github.com/google/uuid"
)

var (
	featSSG    = "ssg"
	resLayout  = "layout"
	resContent = "content"
	resMeta    = "meta"
	resSection = "section"
	resTag     = "tag"
)

// Content related

func (repo *ClioRepo) CreateContent(ctx context.Context, c *ssg.Content) (err error) {
	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("cannot begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				err = fmt.Errorf("cannot rollback transaction: %v (original error: %w)", rbErr, err)
			}
			return
		}
		err = tx.Commit()
	}()

	// Create Content
	contentQuery, err := repo.Query().Get(featSSG, resContent, "Create")
	if err != nil {
		return fmt.Errorf("cannot get create content query: %w", err)
	}
	if _, err = tx.NamedExecContext(ctx, contentQuery, c); err != nil {
		return fmt.Errorf("cannot create content: %w", err)
	}

	// Create Meta
	c.Meta.ContentID = c.ID
	c.Meta.GenID()
	c.Meta.GenCreateValues(c.CreatedBy)
	metaQuery, err := repo.Query().Get(featSSG, resMeta, "Create")
	if err != nil {
		return fmt.Errorf("cannot get create meta query: %w", err)
	}
	if _, err = tx.NamedExecContext(ctx, metaQuery, c.Meta); err != nil {
		return fmt.Errorf("cannot create meta: %w", err)
	}

	return nil
}

func (repo *ClioRepo) GetContent(ctx context.Context, id uuid.UUID) (ssg.Content, error) {
	// This is a placeholder. A specific query "GetWithMeta" is needed for optimal performance.
	// For now, we will filter from the large GetAll query.
	contents, err := repo.GetAllContentWithMeta(ctx)
	if err != nil {
		return ssg.Content{}, err
	}
	for _, content := range contents {
		if content.ID == id {
			return content, nil
		}
	}
	return ssg.Content{}, errors.New("content not found")
}

func (repo *ClioRepo) UpdateContent(ctx context.Context, c *ssg.Content) (err error) {
	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("cannot begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				err = fmt.Errorf("cannot rollback transaction: %v (original error: %w)", rbErr, err)
			}
			return
		}
		err = tx.Commit()
	}()

	// Update Content
	contentQuery, err := repo.Query().Get(featSSG, resContent, "Update")
	if err != nil {
		return fmt.Errorf("cannot get update content query: %w", err)
	}
	if _, err = tx.NamedExecContext(ctx, contentQuery, c); err != nil {
		return fmt.Errorf("cannot update content: %w", err)
	}

	// Update Meta
	metaQuery, err := repo.Query().Get(featSSG, resMeta, "Update")
	if err != nil {
		return fmt.Errorf("cannot get update meta query: %w", err)
	}
	if _, err = tx.NamedExecContext(ctx, metaQuery, c.Meta); err != nil {
		return fmt.Errorf("cannot update meta: %w", err)
	}

	return nil
}

func (repo *ClioRepo) DeleteContent(ctx context.Context, id uuid.UUID) error {
	query, err := repo.Query().Get(featSSG, resContent, "Delete")
	if err != nil {
		return err
	}

	_, err = repo.db.ExecContext(ctx, query, id)
	return err
}

func (repo *ClioRepo) GetAllContentWithMeta(ctx context.Context) ([]ssg.Content, error) {
	query, err := repo.Query().Get(featSSG, resContent, "GetAllContentWithMeta")
	if err != nil {
		return nil, err
	}

	rows, err := repo.db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	contentMap := make(map[uuid.UUID]*ssg.Content)
	var contentOrder []uuid.UUID

	for rows.Next() {
		var c ssg.Content
		var m ssg.Meta
		var t ssg.Tag
		var sectionPath, sectionName sql.NullString
		var publishedAt sql.NullTime

		var metaID sql.NullString
		var description, keywords, robots, canonicalURL, sitemap sql.NullString
		var tableOfContents, share, comments sql.NullBool

		var tagID, tagShortID, tagName, tagSlug sql.NullString

		err := rows.Scan(
			&c.ID, &c.UserID, &c.SectionID, &c.Heading, &c.Body, &c.Draft, &c.Featured, &publishedAt, &c.ShortID,
			&c.CreatedBy, &c.UpdatedBy, &c.CreatedAt, &c.UpdatedAt,
			&sectionPath, &sectionName,
			&metaID, &description, &keywords, &robots, &canonicalURL, &sitemap, &tableOfContents, &share, &comments,
			&tagID, &tagShortID, &tagName, &tagSlug,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}

		if _, ok := contentMap[c.ID]; !ok {
			c.SetType(resContent)
			c.SectionPath = sectionPath.String
			c.SectionName = sectionName.String
			if publishedAt.Valid {
				c.PublishedAt = &publishedAt.Time
			}

			if metaID.Valid {
				m.ID, _ = uuid.Parse(metaID.String)
				m.ContentID = c.ID
				m.Description = description.String
				m.Keywords = keywords.String
				m.Robots = robots.String
				m.CanonicalURL = canonicalURL.String
				m.Sitemap = sitemap.String
				m.TableOfContents = tableOfContents.Bool
				m.Share = share.Bool
				m.Comments = comments.Bool
				c.Meta = m
			}

			contentMap[c.ID] = &c
			contentOrder = append(contentOrder, c.ID)
		}

		if tagID.Valid {
			t.ID, _ = uuid.Parse(tagID.String)
			t.SetShortID(tagShortID.String)
			t.Name = tagName.String
			t.SlugField = tagSlug.String
			t.SetType("tag")
			contentMap[c.ID].Tags = append(contentMap[c.ID].Tags, t)
		}
	}

	contents := make([]ssg.Content, len(contentOrder))
	for i, id := range contentOrder {
		contents[i] = *contentMap[id]
	}

	return contents, nil
}

// Section related

func (repo *ClioRepo) CreateSection(ctx context.Context, section ssg.Section) error {
	query, err := repo.Query().Get(featSSG, resSection, "Create")
	if err != nil {
		return err
	}

	_, err = repo.db.ExecContext(ctx, query,
		section.GetID(),
		section.GetShortID(),
		section.Name,
		section.Description,
		section.Path,
		section.LayoutID,
		section.Image,
		section.Header,
		section.GetCreatedBy(),
		section.GetUpdatedBy(),
		section.GetCreatedAt(),
		section.GetUpdatedAt(),
	)
	return err
}

func (repo *ClioRepo) GetSections(ctx context.Context) ([]ssg.Section, error) {
	query, err := repo.Query().Get(featSSG, resSection, "GetAll")
	if err != nil {
		return nil, err
	}
	rows, err := repo.db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sections []ssg.Section
	for rows.Next() {
		var s ssg.Section
		var layoutName sql.NullString
		err := rows.Scan(
			&s.ID, &s.ShortID, &s.Name, &s.Description, &s.Path, &s.LayoutID, &s.Image, &s.Header,
			&s.CreatedBy, &s.UpdatedBy, &s.CreatedAt, &s.UpdatedAt, &layoutName,
		)
		if err != nil {
			return nil, err
		}
		s.LayoutName = layoutName.String
		sections = append(sections, s)
	}
	return sections, nil
}

func (repo *ClioRepo) GetSection(ctx context.Context, id uuid.UUID) (ssg.Section, error) {
	query, err := repo.Query().Get(featSSG, resSection, "Get")
	if err != nil {
		return ssg.Section{}, err
	}

	row := repo.db.QueryRowxContext(ctx, query, id)

	var (
		sectionID   uuid.UUID
		name        string
		description string
		path        string
		layoutID    uuid.UUID
		image       string
		header      string
		shortID     string
		createdBy   uuid.UUID
		updatedBy   uuid.UUID
		createdAt   time.Time
		updatedAt   time.Time
		layoutName  sql.NullString
	)

	err = row.Scan(
		&sectionID, &shortID, &name, &description, &path, &layoutID, &image, &header,
		&createdBy, &updatedBy, &createdAt, &updatedAt, &layoutName,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ssg.Section{}, errors.New("section not found")
		}
		return ssg.Section{}, err
	}

	section := ssg.NewSection(name, description, path, layoutID)
	section.SetID(sectionID)
	section.Image = image
	section.Header = header
	section.LayoutName = layoutName.String
	section.SetShortID(shortID)
	section.SetCreatedBy(createdBy)
	section.SetUpdatedBy(updatedBy)
	section.SetCreatedAt(createdAt)
	section.SetUpdatedAt(updatedAt)

	return section, nil
}

func (repo *ClioRepo) UpdateSection(ctx context.Context, section ssg.Section) error {
	query, err := repo.Query().Get(featSSG, resSection, "Update")
	if err != nil {
		return err
	}

	_, err = repo.db.ExecContext(ctx, query,
		section.Name,
		section.Description,
		section.Path,
		section.LayoutID,
		section.Image,
		section.Header,
		section.GetShortID(),
		section.GetUpdatedBy(),
		section.GetUpdatedAt(),
		section.GetID(),
	)
	return err
}

func (repo *ClioRepo) DeleteSection(ctx context.Context, id uuid.UUID) error {
	query, err := repo.Query().Get(featSSG, resSection, "Delete")
	if err != nil {
		return err
	}

	_, err = repo.db.ExecContext(ctx, query, id)
	return err
}

// Layout related

func (repo *ClioRepo) CreateLayout(ctx context.Context, layout ssg.Layout) error {
	query, err := repo.Query().Get(featSSG, "layout", "Create")
	if err != nil {
		return err
	}
	_, err = repo.db.ExecContext(ctx, query,
		layout.GetID(),
		layout.GetShortID(),
		layout.Name,
		layout.Description,
		layout.Code,
		layout.GetCreatedBy(),
		layout.GetUpdatedBy(),
		layout.GetCreatedAt(),
		layout.GetUpdatedAt(),
	)
	return err
}

func (repo *ClioRepo) GetAllLayouts(ctx context.Context) ([]ssg.Layout, error) {
	query, err := repo.Query().Get(featSSG, resLayout, "GetAll")
	if err != nil {
		return nil, err
	}

	rows, err := repo.db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var layouts []ssg.Layout
	for rows.Next() {
		var (
			id          uuid.UUID
			shortID     string
			name        string
			description string
			code        string
			createdBy   uuid.UUID
			updatedBy   uuid.UUID
			createdAt   time.Time
			updatedAt   time.Time
		)

		err := rows.Scan(
			&id, &shortID, &name, &description, &code,
			&createdBy, &updatedBy, &createdAt, &updatedAt,
		)
		if err != nil {
			return nil, err
		}

		layout := ssg.Newlayout(name, description, code)
		layout.SetID(id)
		layout.SetShortID(shortID)
		layout.SetCreatedBy(createdBy)
		layout.SetUpdatedBy(updatedBy)
		layout.SetCreatedAt(createdAt)
		layout.SetUpdatedAt(updatedAt)

		layouts = append(layouts, layout)
	}

	return layouts, nil
}

func (repo *ClioRepo) GetLayout(ctx context.Context, id uuid.UUID) (ssg.Layout, error) {
	query, err := repo.Query().Get(featSSG, "layout", "Get")
	if err != nil {
		return ssg.Layout{}, err
	}

	var layout ssg.Layout
	err = repo.db.GetContext(ctx, &layout, query, id)
	if err != nil {
		return ssg.Layout{}, err
	}

	return layout, nil
}

func (repo *ClioRepo) UpdateLayout(ctx context.Context, layout ssg.Layout) error {
	query, err := repo.Query().Get(featSSG, resLayout, "Update")
	if err != nil {
		return err
	}

	_, err = repo.db.ExecContext(ctx, query,
		layout.Name,
		layout.Description,
		layout.Code,
		layout.GetUpdatedBy(),
		layout.GetUpdatedAt(),
		layout.GetID(),
	)
	return err
}

	func (repo *ClioRepo) DeleteLayout(ctx context.Context, id uuid.UUID) error {
	query, err := repo.Query().Get(featSSG, resLayout, "Delete")
	if err != nil {
		return err
	}

	_, err = repo.db.ExecContext(ctx, query, id)
	return err
}


// Tag related

func (repo *ClioRepo) CreateTag(ctx context.Context, tag ssg.Tag) error {
	query, err := repo.Query().Get(featSSG, resTag, "Create")
	if err != nil {
		return err
	}

	_, err = repo.db.NamedExecContext(ctx, query, tag)
	return err
}

func (repo *ClioRepo) GetTag(ctx context.Context, id uuid.UUID) (ssg.Tag, error) {
	query, err := repo.Query().Get(featSSG, resTag, "Get")
	if err != nil {
		return ssg.Tag{}, err
	}

	var tag ssg.Tag
	err = repo.db.GetContext(ctx, &tag, query, id)
	if err != nil {
		return ssg.Tag{}, err
	}

	return tag, nil
}

func (repo *ClioRepo) GetTagByName(ctx context.Context, name string) (ssg.Tag, error) {
	query, err := repo.Query().Get(featSSG, resTag, "GetByName")
	if err != nil {
		return ssg.Tag{}, err
	}

	var tag ssg.Tag
	err = repo.db.GetContext(ctx, &tag, query, name)
	if err != nil {
		return ssg.Tag{}, err
	}

	return tag, nil
}

func (repo *ClioRepo) GetAllTags(ctx context.Context) ([]ssg.Tag, error) {
	query, err := repo.Query().Get(featSSG, resTag, "GetAll")
	if err != nil {
		return nil, err
	}

	var tags []ssg.Tag
	err = repo.db.SelectContext(ctx, &tags, query)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func (repo *ClioRepo) UpdateTag(ctx context.Context, tag ssg.Tag) error {
	query, err := repo.Query().Get(featSSG, resTag, "Update")
	if err != nil {
		return err
	}

	_, err = repo.db.NamedExecContext(ctx, query, tag)
	return err
}

func (repo *ClioRepo) DeleteTag(ctx context.Context, id uuid.UUID) error {
	query, err := repo.Query().Get(featSSG, resTag, "Delete")
	if err != nil {
		return err
	}

	_, err = repo.db.ExecContext(ctx, query, id)
	return err
}

// ContentTag related

func (repo *ClioRepo) AddTagToContent(ctx context.Context, contentID, tagID uuid.UUID) error {
	query, err := repo.Query().Get(featSSG, resTag, "AddTagToContent")
	if err != nil {
		return err
	}

	_, err = repo.db.ExecContext(ctx, query, contentID, tagID)
	return err
}

func (repo *ClioRepo) RemoveTagFromContent(ctx context.Context, contentID, tagID uuid.UUID) error {
	query, err := repo.Query().Get(featSSG, resTag, "RemoveTagFromContent")
	if err != nil {
		return err
	}

	_, err = repo.db.ExecContext(ctx, query, contentID, tagID)
	return err
}

func (repo *ClioRepo) GetTagsForContent(ctx context.Context, contentID uuid.UUID) ([]ssg.Tag, error) {
	query, err := repo.Query().Get(featSSG, resTag, "GetTagsForContent")
	if err != nil {
		return nil, err
	}

	var tags []ssg.Tag
	err = repo.db.SelectContext(ctx, &tags, query, contentID)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func (repo *ClioRepo) GetContentForTag(ctx context.Context, tagID uuid.UUID) ([]ssg.Content, error) {
	query, err := repo.Query().Get(featSSG, resTag, "GetContentForTag")
	if err != nil {
		return nil, err
	}

	var contents []ssg.Content
	err = repo.db.SelectContext(ctx, &contents, query, tagID)
	if err != nil {
		return nil, err
	}

	return contents, nil
}