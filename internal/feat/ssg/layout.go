package ssg

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"github.com/adrianpk/clio/internal/am"
)

const (
	layoutType = "layout"
)

// Layout model.
type Layout struct {
	// Common
	ID       uuid.UUID `json:"id" db:"id"`
	mType    string
	ShortID  string `json:"-" db:"short_id"`
	RefValue string `json:"ref"`

	// Layout specific fields
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	Code        string `json:"code" db:"code"`

	// Audit
	CreatedBy uuid.UUID `json:"-" db:"created_by"`
	UpdatedBy uuid.UUID `json:"-" db:"updated_by"`
	CreatedAt time.Time `json:"-" db:"created_at"`
	UpdatedAt time.Time `json:"-" db:"updated_at"`
}

// Newlayout creates a new Layout.
func Newlayout(name, description, code string) Layout {
	l := Layout{
		mType:       layoutType,
		Name:        name,
		Description: description,
		Code:        code,
	}

	return l
}

// Type returns the type of the entity.
func (l *Layout) Type() string {
	return am.DefaultType(l.mType)
}

// SetType sets the type of the entity.
func (l *Layout) SetType(t string) {
	l.mType = t
}

// GetID returns the unique identifier of the entity.
func (l *Layout) GetID() uuid.UUID {
	return l.ID
}

// GenID delegates to the functional helper.
func (l *Layout) GenID() {
	am.GenID(l)
}

// SetID sets the unique identifier of the entity.
func (l *Layout) SetID(id uuid.UUID, force ...bool) {
	shouldForce := len(force) > 0 && force[0]
	if l.ID == uuid.Nil || (shouldForce && id != uuid.Nil) {
		l.ID = id
	}
}

// ShortID returns the short ID portion of the slug.
func (l *Layout) GetShortID() string {
	return l.ShortID
}

// GenShortID delegates to the functional helper.
func (l *Layout) GenShortID() {
	am.GenShortID(l)
}

// SetShortID sets the short ID of the entity.
func (l *Layout) SetShortID(shortID string, force ...bool) {
	shouldForce := len(force) > 0 && force[0]
	if l.ShortID == "" || shouldForce {
		l.ShortID = shortID
	}
}

// TypeID returns a universal identifier for a specific model instance.
func (l *Layout) TypeID() string {
	return am.Normalize(l.Type()) + "-" + l.GetShortID()
}

// GenCreateValues delegates to the functional helper.
func (l *Layout) GenCreateValues(userID ...uuid.UUID) {
	am.SetCreateValues(l, userID...)
}

// GenUpdateValues delegates to the functional helper.
func (l *Layout) GenUpdateValues(userID ...uuid.UUID) {
	am.SetUpdateValues(l, userID...)
}

// CreatedBy returns the UUID of the user who created the entity.
func (l *Layout) GetCreatedBy() uuid.UUID {
	return l.CreatedBy
}

// UpdatedBy returns the UUID of the user who last updated the entity.
func (l *Layout) GetUpdatedBy() uuid.UUID {
	return l.UpdatedBy
}

// CreatedAt returns the creation time of the entity.
func (l *Layout) GetCreatedAt() time.Time {
	return l.CreatedAt
}

// UpdatedAt returns the last update time of the entity.
func (l *Layout) GetUpdatedAt() time.Time {
	return l.UpdatedAt
}

// SetCreatedAt implements the Auditable interface.
func (l *Layout) SetCreatedAt(t time.Time) {
	l.CreatedAt = t
}

// SetUpdatedAt implements the Auditable interface.
func (l *Layout) SetUpdatedAt(t time.Time) {
	l.UpdatedAt = t
}

// SetCreatedBy implements the Auditable interface.
func (l *Layout) SetCreatedBy(u uuid.UUID) {
	l.CreatedBy = u
}

// SetUpdatedBy implements the Auditable interface.
func (l *Layout) SetUpdatedBy(u uuid.UUID) {
	l.UpdatedBy = u
}

// IsZero returns true if the Layout is uninitialized.
func (l *Layout) IsZero() bool {
	return l.ID == uuid.Nil
}

// Slug returns a human-readable, URL-friendly string identifier for the entity.
func (l *Layout) Slug() string {
	return am.Normalize(l.Name) + "-" + l.GetShortID()
}

func (l *Layout) OptValue() string {
	return l.GetID().String()
}

func (l *Layout) OptLabel() string {
	return l.Name
}

// Ref returns the reference string for this entity.
func (l *Layout) Ref() string {
	return l.RefValue
}

// SetRef sets the reference string for this entity.
func (l *Layout) SetRef(ref string) {
	l.RefValue = ref
}

// UnmarshalJSON ensures model fields are initialized after unmarshal.
func (l *Layout) UnmarshalJSON(data []byte) error {
	type Alias Layout
	temp := &struct {
		*Alias
	}{
		Alias: (*Alias)(l),
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	if l.mType == "" {
		l.mType = layoutType
	}

	return nil
}

// StringID returns the unique identifier of the entity as a string.
func (l *Layout) StringID() string {
	return l.GetID().String()
}
