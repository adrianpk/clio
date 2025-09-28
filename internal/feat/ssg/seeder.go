package ssg

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/adrianpk/clio/internal/am"
	"github.com/adrianpk/clio/internal/feat/auth"
	"github.com/google/uuid"
)

const (
	ssgFeat = "ssg"
)

type Seeder struct {
	*am.JSONSeeder
	repo Repo
}

type SeedFile struct {
	Layouts     []map[string]any `json:"layouts"`
	Sections    []map[string]any `json:"sections"`
	Contents    []map[string]any `json:"contents"`
	Metas       []map[string]any `json:"metas"`
	Tags        []map[string]any `json:"tags"`
	ContentTags []map[string]any `json:"content_tags"`
	Params      []map[string]any `json:"params"`
}

func NewSeeder(assetsFS embed.FS, engine string, repo Repo) *Seeder {
	return &Seeder{
		JSONSeeder: am.NewJSONSeeder(ssgFeat, assetsFS, engine),
		repo:       repo,
	}
}

func (s *Seeder) Setup(ctx context.Context) error {
	if err := s.JSONSeeder.Setup(ctx); err != nil {
		return err
	}
	return s.SeedAll(ctx)
}

func (s *Seeder) SeedAll(ctx context.Context) error {
	s.Log().Info("Seeding SSG data...")
	byFeature, err := s.JSONSeeder.LoadJSONSeeds()
	if err != nil {
		return fmt.Errorf("failed to load JSON seeds: %w", err)
	}
	for feature, seeds := range byFeature {
		if feature != ssgFeat {
			continue
		}

		for _, seed := range seeds {
			applied, err := s.JSONSeeder.SeedApplied(seed.Datetime, seed.Name, feature)
			if err != nil {
				return fmt.Errorf("failed to check if seed was applied: %w", err)
			}
			if applied {
				s.Log().Debugf("Seed already applied: %s-%s [%s]", seed.Datetime, seed.Name, feature)
				continue
			}

			var data SeedFile
			err = json.Unmarshal([]byte(seed.Content), &data)
			if err != nil {
				return fmt.Errorf("failed to unmarshal %s seed: %w", feature, err)
			}

			err = s.seedData(ctx, &data)
			if err != nil {
				return err
			}

			err = s.JSONSeeder.ApplyJSONSeed(seed.Datetime, seed.Name, feature, seed.Content)
			if err != nil {
				s.Log().Errorf("error recording JSON seed: %v", err)
			}
		}
	}
	return nil
}

func (s *Seeder) seedData(ctx context.Context, data *SeedFile) error {
	userCache := make(map[string]auth.User)

	layoutRefToID := make(map[string]uuid.UUID)
	for _, lMap := range data.Layouts {
		l := Layout{
			Name:        lMap["name"].(string),
			Description: lMap["description"].(string),
			Code:        lMap["code"].(string),
		}
		l.GenCreateValues()
		if err := s.repo.CreateLayout(ctx, l); err != nil {
			return fmt.Errorf("error inserting layout: %w", err)
		}
		if ref, ok := lMap["ref"].(string); ok {
			layoutRefToID[ref] = l.GetID()
		}
	}

	sectionRefToID := make(map[string]uuid.UUID)
	for _, sMap := range data.Sections {
		sec := Section{
			Name:        sMap["name"].(string),
			Description: sMap["description"].(string),
			Path:        sMap["path"].(string),
		}
		if image, ok := sMap["image"].(string); ok {
			sec.Image = image
		}
		if header, ok := sMap["header"].(string); ok {
			sec.Header = header
		}
		if layoutRef, ok := sMap["layout_ref"].(string); ok {
			if id, found := layoutRefToID[layoutRef]; found {
				sec.LayoutID = id
			}
		}
		sec.GenCreateValues()
		if err := s.repo.CreateSection(ctx, sec); err != nil {
			return fmt.Errorf("error inserting section: %w", err)
		}
		if ref, ok := sMap["ref"].(string); ok {
			sectionRefToID[ref] = sec.GetID()
		}
	}

	metaByContentRef := make(map[string]map[string]any)
	for _, mMap := range data.Metas {
		if contentRef, ok := mMap["content_ref"].(string); ok {
			metaByContentRef[contentRef] = mMap
		}
	}

	contentRefToID := make(map[string]uuid.UUID)
	for _, cMap := range data.Contents {
		con := Content{
			Heading:  cMap["heading"].(string),
			Body:     cMap["body"].(string),
			Draft:    cMap["draft"].(bool),
			Featured: cMap["featured"].(bool),
		}
		if pubAt, ok := cMap["published_at"].(string); ok {
			t, err := time.Parse(time.RFC3339, pubAt)
			if err == nil {
				con.PublishedAt = &t
			}
		} else {
			con.PublishedAt = nil
		}

		if sectionRef, ok := cMap["section_ref"].(string); ok {
			if id, found := sectionRefToID[sectionRef]; found {
				con.SectionID = id
			}
		}

		// Handle user reference
		var author auth.User
		var err error
		userRef, ok := cMap["user_ref"].(string)
		if !ok || userRef == "" {
			userRef = "user-superadmin" // Default to superadmin
		}

		username := strings.TrimPrefix(userRef, "user-")
		author, ok = userCache[username]
		if !ok {
			author, err = s.repo.GetUserByUsername(ctx, username)
			if err != nil {
				return fmt.Errorf("error getting user '%s' for content seeding: %w", username, err)
			}
			userCache[username] = author
		}

		contentRef := cMap["ref"].(string)
		if mMap, ok := metaByContentRef[contentRef]; ok {
			meta := Meta{
				Description:     mMap["description"].(string),
				Keywords:        mMap["keywords"].(string),
				Robots:          mMap["robots"].(string),
				CanonicalURL:    mMap["canonical_url"].(string),
				Sitemap:         mMap["sitemap"].(string),
				TableOfContents: mMap["table_of_contents"].(bool),
				Share:           mMap["share"].(bool),
				Comments:        mMap["comments"].(bool),
			}
			con.Meta = meta
		}

		con.GenID()
		con.GenShortID()
		now := time.Now()
		con.SetCreatedAt(now)
		con.SetUpdatedAt(now)
		con.UserID = author.ID
		con.SetCreatedBy(author.ID)
		con.SetUpdatedBy(author.ID)

		if err := s.repo.CreateContent(ctx, &con); err != nil {
			return fmt.Errorf("error inserting content '%s': %w", contentRef, err)
		}
		contentRefToID[contentRef] = con.GetID()
	}

	tagRefToID := make(map[string]uuid.UUID)
	for _, tMap := range data.Tags {
		tag := Tag{
			Name: tMap["name"].(string),
		}
		tag.GenCreateValues()
		existingTag, err := s.repo.GetTagByName(ctx, tag.Name)
		if err == nil {
			tag = existingTag
		} else {
			if err := s.repo.CreateTag(ctx, tag); err != nil {
				return fmt.Errorf("error inserting tag: %w", err)
			}
		}
		if ref, ok := tMap["ref"].(string); ok {
			tagRefToID[ref] = tag.GetID()
		}
	}

	for _, ctMap := range data.ContentTags {
		contentRef := ctMap["content_ref"].(string)
		tagRef := ctMap["tag_ref"].(string)
		contentID, cOK := contentRefToID[contentRef]
		tagID, tOK := tagRefToID[tagRef]

		if cOK && tOK {
			if err := s.repo.AddTagToContent(ctx, contentID, tagID); err != nil {
				return fmt.Errorf("error adding tag '%s' to content '%s': %w", tagRef, contentRef, err)
			}
		}
	}

	for _, pMap := range data.Params {
		var systemVal int
		if system, ok := pMap["system"].(float64); ok {
			systemVal = int(system)
		}

		p := Param{
			Name:        pMap["name"].(string),
			Description: pMap["description"].(string),
			Value:       pMap["value"].(string),
			RefKey:      pMap["ref_key"].(string),
			System:      systemVal, // Assign System here
		}
		p.GenCreateValues()
		if err := s.repo.CreateParam(ctx, &p); err != nil {
			return fmt.Errorf("error inserting param: %w", err)
		}
	}

	return nil
}
