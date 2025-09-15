package ssg

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"github.com/adrianpk/clio/internal/am"
)

const (
	tagType = "tag"
)

// Tag model.
type Tag struct {
	// Common
	ID      uuid.UUID `json:"id" db:"id"`
	mType   string
	ShortID string `json:"-" db:"short_id"`

	// Tag specific fields
	Name string `json:"name" db:"name"`
	SlugField string `json:"slug" db:"slug"`

	// Audit
	CreatedBy uuid.UUID `json:"-" db:"created_by"`
	UpdatedBy uuid.UUID `json:"-" db:"updated_by"`
	CreatedAt time.Time `json:"-" db:"created_at"`
	UpdatedAt time.Time `json:"-" db:"updated_at"`
}

// NewTag creates a new Tag.
func NewTag(name string) Tag {
	t := Tag{
		mType: tagType,
		Name:  name,
	}

	return t
}

// Type returns the type of the entity.
func (t *Tag) Type() string {
	return am.DefaultType(t.mType)
}

// SetType sets the type of the entity.
func (t *Tag) SetType(typ string) {
	t.mType = typ
}

// GetID returns the unique identifier of the entity.
func (t *Tag) GetID() uuid.UUID {
	return t.ID
}

// GenID delegates to the functional helper.
func (t *Tag) GenID() {
	am.GenID(t)
}

// SetID sets the unique identifier of the entity.
func (t *Tag) SetID(id uuid.UUID, force ...bool) {
	shouldForce := len(force) > 0 && force[0]
	if t.ID == uuid.Nil || (shouldForce && id != uuid.Nil) {
		t.ID = id
	}
}

// ShortID returns the short ID portion of the slug.
func (t *Tag) GetShortID() string {
	return t.ShortID
}

// GenShortID delegates to the functional helper.
func (t *Tag) GenShortID() {
	am.GenShortID(t)
}

// SetShortID sets the short ID of the entity.
func (t *Tag) SetShortID(shortID string, force ...bool) {
	shouldForce := len(force) > 0 && force[0]
	if t.ShortID == "" || shouldForce {
		t.ShortID = shortID
	}
}

// TypeID returns a universal identifier for a specific model instance.
func (t *Tag) TypeID() string {
	return am.Normalize(t.Type()) + "-" + t.GetShortID()
}

// GenCreateValues delegates to the functional helper.
func (t *Tag) GenCreateValues(userID ...uuid.UUID) {
	am.SetCreateValues(t, userID...)
}

// GenUpdateValues delegates to the functional helper.
func (t *Tag) GenUpdateValues(userID ...uuid.UUID) {
	am.SetUpdateValues(t, userID...)
}

// CreatedBy returns the UUID of the user who created the entity.
func (t *Tag) GetCreatedBy() uuid.UUID {
	return t.CreatedBy
}

// UpdatedBy returns the UUID of the user who last updated the entity.
func (t *Tag) GetUpdatedBy() uuid.UUID {
	return t.UpdatedBy
}

// CreatedAt returns the creation time of the entity.
func (t *Tag) GetCreatedAt() time.Time {
	return t.CreatedAt
}

// UpdatedAt returns the last update time of the entity.
func (t *Tag) GetUpdatedAt() time.Time {
	return t.UpdatedAt
}

// SetCreatedAt implements the Auditable interface.
func (t *Tag) SetCreatedAt(createdAt time.Time) {
	t.CreatedAt = createdAt
}

// SetUpdatedAt implements the Auditable interface.
func (t *Tag) SetUpdatedAt(updatedAt time.Time) {
	t.UpdatedAt = updatedAt
}

// SetCreatedBy implements the Auditable interface.
func (t *Tag) SetCreatedBy(createdBy uuid.UUID) {
	t.CreatedBy = createdBy
}

// SetUpdatedBy implements the Auditable interface.
func (t *Tag) SetUpdatedBy(updatedBy uuid.UUID) {
	t.UpdatedBy = updatedBy
}

// IsZero returns true if the Tag is uninitialized.
func (t *Tag) IsZero() bool {
	return t.ID == uuid.Nil
}

// Slug returns a human-readable, URL-friendly string identifier for the entity.
func (t *Tag) Slug() string {
	if t.SlugField != "" {
		return t.SlugField
	}
	return am.Normalize(t.Name) + "-" + t.GetShortID()
}

func (t *Tag) OptValue() string {
	return t.GetID().String()
}

func (t *Tag) OptLabel() string {
	return t.Name
}

// UnmarshalJSON ensures model fields are initialized after unmarshal.
func (t *Tag) UnmarshalJSON(data []byte) error {
	type Alias Tag
	temp := &struct {
		*Alias
	}{
		Alias: (*Alias)(t),
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	if t.mType == "" {
		t.mType = tagType
	}

	return nil
}
