package ssg

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"

	"github.com/adrianpk/clio/internal/am"
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
	Layouts  []map[string]any `json:"layouts"`
	Sections []map[string]any `json:"sections"`
	Contents []map[string]any `json:"contents"`
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
	ctx, tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("error at beginning tx for seedLayouts: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	layoutRefToID := make(map[string]uuid.UUID)
	for _, lMap := range data.Layouts {
		l := Layout{
			Name:        lMap["name"].(string),
			Description: lMap["description"].(string),
			Code:        lMap["code"].(string),
		}
		l.GenCreateValues()
		err := s.repo.CreateLayout(ctx, l)
		if err != nil {
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
			Image:       sMap["image"].(string),
			Header:      sMap["header"].(string),
		}
		if layoutRef, ok := sMap["layout_ref"].(string); ok {
			if id, found := layoutRefToID[layoutRef]; found {
				sec.LayoutID = id
			}
		}

		sec.GenCreateValues()
		err := s.repo.CreateSection(ctx, sec)
		if err != nil {
			return fmt.Errorf("error inserting section: %w", err)
		}
		if ref, ok := sMap["ref"].(string); ok {
			sectionRefToID[ref] = sec.GetID()
		}
	}

	for _, cMap := range data.Contents {
		con := Content{
			Heading: cMap["heading"].(string),
			Body:    cMap["body"].(string),
			Status:  "draft",
		}

		if sectionRef, ok := cMap["section_ref"].(string); ok {
			if id, found := sectionRefToID[sectionRef]; found {
				con.SectionID = id
			}
		}

		con.GenCreateValues()
		err := s.repo.CreateContent(ctx, con)
		if err != nil {
			return fmt.Errorf("error inserting content: %w", err)
		}
	}

	return tx.Commit()
}
