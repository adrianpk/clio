package ssg

import (
	"time"

	"github.com/adrianpk/clio/internal/am"
	"github.com/google/uuid"
)

const (
	paramType = "param"
)

// Param represents a dynamic configuration entry.
type Param struct {
	// Common
	ID      uuid.UUID `json:"id" db:"id"`
	mType   string
	ShortID string `json:"-" db:"short_id"` // Note: short_id was removed from DB migration, but kept here for consistency with other models' Go struct definitions.

	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	Value       string `json:"value" db:"value"`
	RefKey      string `json:"ref_key" db:"ref_key"` // Should match xxx.yyy.zzz congfig property

	// Audit
	CreatedBy uuid.UUID `json:"-" db:"created_by"`
	UpdatedBy uuid.UUID `json:"-" db:"updated_by"`
	CreatedAt time.Time `json:"-" db:"created_at"`
	UpdatedAt time.Time `json:"-" db:"updated_at"`
}

// NewParam creates a new Param instance with default values.
func NewParam(name, value string) Param {
	p := Param{
		mType: paramType,
		Name:  name,
		Value: value,
	}
	return p
}

// Type returns the type of the entity.
func (p Param) Type() string {
	return am.DefaultType(p.mType)
}

// SetType sets the type of the entity.
func (p *Param) SetType(t string) {
	p.mType = t
}

// GetID returns the unique identifier of the entity.
func (p Param) GetID() uuid.UUID {
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

// GenCreateValues delegates to the functional helper.
func (p *Param) GenCreateValues(userID ...uuid.UUID) {
	am.SetCreateValues(p, userID...)
}

// GenUpdateValues delegates to the functional helper.
func (p *Param) GenUpdateValues(userID ...uuid.UUID) {
	am.SetUpdateValues(p, userID...)
}

// GetCreatedBy returns the UUID of the user who created the entity.
func (p *Param) GetCreatedBy() uuid.UUID {
	return p.CreatedBy
}

// GetUpdatedBy returns the UUID of the user who last updated the entity.
func (p *Param) GetUpdatedBy() uuid.UUID {
	return p.UpdatedBy
}

// GetCreatedAt returns the creation time of the entity.
func (p *Param) GetCreatedAt() time.Time {
	return p.CreatedAt
}

// GetUpdatedAt returns the last update time of the entity.
func (p *Param) GetUpdatedAt() time.Time {
	return p.UpdatedAt
}

// SetCreatedAt implements the Auditable interface.
func (p *Param) SetCreatedAt(t time.Time) {
	p.CreatedAt = t
}

// SetUpdatedAt implements the Auditable interface.
func (p *Param) SetUpdatedAt(t time.Time) {
	p.UpdatedAt = t
}

// SetCreatedBy implements the Auditable interface.
func (p *Param) SetCreatedBy(id uuid.UUID) {
	p.CreatedBy = id
}

// SetUpdatedBy implements the Auditable interface.
func (p *Param) SetUpdatedBy(id uuid.UUID) {
	p.UpdatedBy = id
}

// IsZero returns true if the Param is uninitialized.
func (p *Param) IsZero() bool {
	return p.ID == uuid.Nil
}

// Slug returns a slug for the param.
func (p *Param) Slug() string {
	return am.Normalize(p.Name) + "-" + p.GetShortID()
}
