package ssg

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/adrianpk/clio/internal/am"
	"github.com/google/uuid"
	feat "github.com/adrianpk/clio/internal/feat/ssg"
)

// ContentForm represents the form data for a content.
type ContentForm struct {
	*am.BaseForm
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	SectionID string `json:"section_id"`
	Heading   string `json:"heading"`
	Body      string `json:"body"`
	Status    string `json:"status"`
	Tags      string `json:"tags"`
	Errors    map[string]string
}

// NewContentForm creates a new ContentForm from a request.
func NewContentForm(r *http.Request) ContentForm {
	return ContentForm{
		BaseForm: am.NewBaseForm(r),
		Errors:   make(map[string]string),
	}
}

// ContentFormFromRequest creates a ContentForm from an HTTP request.
func ContentFormFromRequest(r *http.Request) (ContentForm, error) {
	if err := r.ParseForm(); err != nil {
		return ContentForm{}, fmt.Errorf("error parsing form: %w", err)
	}

	form := NewContentForm(r) // Initialize with BaseForm
	form.ID = r.Form.Get("id")
	form.UserID = r.Form.Get("user_id")
	form.SectionID = r.Form.Get("section_id")
	form.Heading = r.Form.Get("heading")
	form.Body = r.Form.Get("body")
	form.Status = r.Form.Get("status")
	form.Tags = r.Form.Get("tags")

	return form, nil
}

// ToFeatContent converts a ContentForm to a feat.Content model.
func ToFeatContent(form ContentForm) feat.Content {
	content := feat.NewContent(form.Heading, form.Body)
	content.Status = form.Status

	if form.ID != "" {
		id, err := uuid.Parse(form.ID)
		if err == nil {
			content.ID = id
		}
	}

	if form.UserID != "" {
		userID, err := uuid.Parse(form.UserID)
		if err == nil {
			content.UserID = userID
		}
	}

	if form.SectionID != "" {
		sectionID, err := uuid.Parse(form.SectionID)
		if err == nil {
			content.SectionID = sectionID
		}
	}

	// New part for tags
	if form.Tags != "" {
		var tags []struct {
			Value string `json:"value"`
		}
		if err := json.Unmarshal([]byte(form.Tags), &tags); err == nil {
			for _, t := range tags {
				content.Tags = append(content.Tags, feat.Tag{Name: t.Value})
			}
		}
	}

	return content
}

// ToContentForm converts a feat.Content model to a ContentForm.
func ToContentForm(r *http.Request, content feat.Content) ContentForm {
	form := NewContentForm(r) // Initialize with BaseForm
	form.ID = content.GetID().String()
	form.UserID = content.UserID.String()
	form.SectionID = content.SectionID.String()
	form.Heading = content.Heading
	form.Body = content.Body
	form.Status = content.Status

	// Create a comma-separated string of tag names
	tagNames := make([]string, len(content.Tags))
	for i, tag := range content.Tags {
		tagNames[i] = tag.Name
	}
	form.Tags = strings.Join(tagNames, ",")

	return form
}

// Validate validates the ContentForm.
func (f *ContentForm) Validate() error {
	if f.Heading == "" {
		f.Errors["heading"] = "Heading cannot be empty"
	}
	if f.Body == "" {
		f.Errors["body"] = "Body cannot be empty"
	}

	return nil
}

// HasErrors returns true if the form has validation errors.
func (f *ContentForm) HasErrors() bool { return len(f.Errors) > 0 }

// LayoutForm represents the form for creating or updating a layout.
type LayoutForm struct {
	*am.BaseForm
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Code        string `json:"code"`
	Errors      map[string]string
}

// NewLayoutForm creates a new LayoutForm.
func NewLayoutForm(r *http.Request) LayoutForm {
	return LayoutForm{
		BaseForm: am.NewBaseForm(r),
		Errors:   make(map[string]string),
	}
}

// LayoutFormFromRequest creates a LayoutForm from an HTTP request.
func LayoutFormFromRequest(r *http.Request) (LayoutForm, error) {
	if err := r.ParseForm(); err != nil {
		return LayoutForm{}, fmt.Errorf("error parsing form: %w", err)
	}

	form := NewLayoutForm(r)
	form.ID = r.Form.Get("id")
	form.Name = r.Form.Get("name")
	form.Description = r.Form.Get("description")
	form.Code = r.Form.Get("code")

	return form, nil
}

// ToFeatLayout converts a LayoutForm to a feat.Layout model.
func ToFeatLayout(form LayoutForm) feat.Layout {
	layout := feat.Newlayout(form.Name, form.Description, form.Code)
	if form.ID != "" {
		id, err := uuid.Parse(form.ID)
		if err == nil {
			layout.ID = id
		}
	}
	return layout
}

// ToLayoutForm converts a feat.Layout model to a LayoutForm.
func ToLayoutForm(r *http.Request, layout feat.Layout) LayoutForm {
	form := NewLayoutForm(r)
	form.ID = layout.GetID().String()
	form.Name = layout.Name
	form.Description = layout.Description
	form.Code = layout.Code

	return form
}

// Validate validates the LayoutForm.
func (f *LayoutForm) Validate() error {
	if f.Name == "" {
		f.Errors["name"] = "Name is required"
	}
	if f.Code == "" {
		f.Errors["code"] = "Code is required"
	}
	return nil
}

// HasErrors returns true if the form has validation errors.
func (f *LayoutForm) HasErrors() bool {
	return len(f.Errors) > 0
}

// SectionForm represents the form for creating or updating a section.
type SectionForm struct {
	*am.BaseForm
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Path        string `json:"path"`
	LayoutID    string `json:"layout_id"`
	Image       string `json:"image"`
	Header      string `json:"header"`
	Errors      map[string]string
}

// NewSectionForm creates a new SectionForm.
func NewSectionForm(r *http.Request) SectionForm {
	return SectionForm{
		BaseForm: am.NewBaseForm(r),
		Errors:   make(map[string]string),
	}
}

// SectionFormFromRequest creates a SectionForm from an HTTP request.
func SectionFormFromRequest(r *http.Request) (SectionForm, error) {
	if err := r.ParseForm(); err != nil {
		return SectionForm{}, fmt.Errorf("error parsing form: %w", err)
	}

	form := NewSectionForm(r)
	form.ID = r.Form.Get("id")
	form.Name = r.Form.Get("name")
	form.Description = r.Form.Get("description")
	form.Path = r.Form.Get("path")
	form.LayoutID = r.Form.Get("layout_id")
	form.Image = r.Form.Get("image")
	form.Header = r.Form.Get("header")

	return form, nil
}

// ToFeatSection converts a SectionForm to a feat.Section model.
func ToFeatSection(form SectionForm) feat.Section {
	layoutID, _ := uuid.Parse(form.LayoutID)
	section := feat.NewSection(form.Name, form.Description, form.Path, layoutID)
	section.Image = form.Image
	section.Header = form.Header
	if form.ID != "" {
		id, err := uuid.Parse(form.ID)
		if err == nil {
			section.ID = id
		}
	}

	return section
}

// ToSectionForm converts a feat.Section model to a SectionForm.
func ToSectionForm(r *http.Request, section feat.Section) SectionForm {
	form := NewSectionForm(r)
	form.ID = section.GetID().String()
	form.Name = section.Name
	form.Description = section.Description
	form.Path = section.Path
	form.LayoutID = section.LayoutID.String()
	form.Image = section.Image
	form.Header = section.Header
	return form
}

// Validate validates the SectionForm.
func (f *SectionForm) Validate() error {
	if f.Name == "" {
		f.Errors["name"] = "Name is required"
	}

	if f.Path == "" {
		f.Errors["path"] = "Path is required"
	}

	if f.LayoutID == "" {
		f.Errors["layout_id"] = "Layout is required"
	}
	return nil
}

// HasErrors returns true if the form has validation errors.
func (f *SectionForm) HasErrors() bool {
	return len(f.Errors) > 0
}

// TagForm represents the form data for a tag.
type TagForm struct {
	*am.BaseForm
	ID     string `json:"id"`
	Name   string `json:"name"`
	Errors map[string]string
}

// NewTagForm creates a new TagForm from a request.
func NewTagForm(r *http.Request) TagForm {
	return TagForm{
		BaseForm: am.NewBaseForm(r),
		Errors:   make(map[string]string),
	}
}

// TagFormFromRequest creates a TagForm from an HTTP request.
func TagFormFromRequest(r *http.Request) (TagForm, error) {
	if err := r.ParseForm(); err != nil {
		return TagForm{}, fmt.Errorf("error parsing form: %w", err)
	}

	form := NewTagForm(r)
	form.ID = r.Form.Get("id")
	form.Name = r.Form.Get("name")

	return form, nil
}

// ToFeatTag converts a TagForm to a feat.Tag model.
func ToFeatTag(form TagForm) feat.Tag {
	tag := feat.NewTag(form.Name)
	if form.ID != "" {
		id, err := uuid.Parse(form.ID)
		if err == nil {
			tag.ID = id
		}
	}
	return tag
}

// ToTagForm converts a feat.Tag model to a TagForm.
func ToTagForm(r *http.Request, featTag feat.Tag) TagForm {
	form := NewTagForm(r)
	form.ID = featTag.GetID().String()
	form.Name = featTag.Name
	return form
}

// Validate validates the TagForm.
func (f *TagForm) Validate() error {
	if f.Name == "" {
		f.Errors["name"] = "Name cannot be empty"
	}
	return nil
}

// HasErrors returns true if the form has validation errors.
func (f *TagForm) HasErrors() bool { return len(f.Errors) > 0 }
