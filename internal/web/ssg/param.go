package ssg

import (
	"github.com/adrianpk/clio/internal/am"
	feat "github.com/adrianpk/clio/internal/feat/ssg"
	"github.com/google/uuid"
)

const (
	paramType = "param"
)

// Param model for web layer.
type Param struct {
	ID          uuid.UUID `json:"id"`
	ShortID     string    `json:"-"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Value       string    `json:"value"`
	RefKey      string    `json:"ref_key"`
}

// NewParam creates a new Param for the web layer.
func NewParam(name, value string) Param {
	return Param{
		Name:  name,
		Value: value,
	}
}

// Type returns the type of the entity.
func (p *Param) Type() string {
	return am.DefaultType(paramType)
}

// GetID returns the unique identifier of the entity.
func (p *Param) GetID() uuid.UUID {
	return p.ID
}

// GenID delegates to the functional helper.
func (p *Param) GenID() {
	am.GenID(p)
}

// SetID sets the unique identifier of the entity.
func (p *Param) SetID(id uuid.UUID, force ...bool) {
	shouldForce := len(force) > 0 && force[0]
	if p.ID == uuid.Nil || (shouldForce && id != uuid.Nil) {
		p.ID = id
	}
}

// GetShortID returns the short ID portion of the slug.
func (p *Param) GetShortID() string {
	return p.ShortID
}

// GenShortID delegates to the functional helper.
func (p *Param) GenShortID() {
	am.GenShortID(p)
}

// SetShortID sets the short ID of the entity.
func (p *Param) SetShortID(shortID string, force ...bool) {
	shouldForce := len(force) > 0 && force[0]
	if p.ShortID == "" || shouldForce {
		p.ShortID = shortID
	}
}

// TypeID returns a universal identifier for a specific model instance.
func (p *Param) TypeID() string {
	return am.Normalize(p.Type()) + "-" + p.GetShortID()
}

// IsZero returns true if the Param is uninitialized.
func (p *Param) IsZero() bool {
	return p.ID == uuid.Nil
}

// Slug returns a slug for the param.
func (p *Param) Slug() string {
	return am.Normalize(p.Name) + "-" + p.GetShortID()
}

// ToWebParam converts a feat.Param model to a web.Param model.
func ToWebParam(featParam feat.Param) Param {
	return Param{
		ID:          featParam.ID,
		ShortID:     featParam.ShortID,
		Name:        featParam.Name,
		Description: featParam.Description,
		Value:       featParam.Value,
		RefKey:      featParam.RefKey,
	}
}

// ToWebParams converts a slice of feat.Param models to a slice of web.Param models.
func ToWebParams(featParams []feat.Param) []Param {
	webParams := make([]Param, len(featParams))
	for i, featParam := range featParams {
		webParams[i] = ToWebParam(featParam)
	}
	return webParams
}
