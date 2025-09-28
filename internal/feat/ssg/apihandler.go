package ssg

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/adrianpk/clio/internal/am"

	"github.com/google/uuid"
)

const (
	resContent         = "content"
	resContentCap      = "Content"
	resSection         = "section"
	resSectionCap      = "Section"
	resLayout          = "layout"
	resLayoutCap       = "Layout"
	resTag             = "tag"
	resTagCap          = "Tag"
	resParam           = "param"
	resParamCap        = "Param"
	resImage           = "image"
	resImageCap        = "Image"
	resImageVariant    = "image variant"
	resImageVariantCap = "Image variant"
)

type APIHandler struct {
	*am.APIHandler
	svc Service
}

func NewAPIHandler(name string, service Service, options ...am.Option) *APIHandler {
	return &APIHandler{
		APIHandler: am.NewAPIHandler(name, options...),
		svc:        service,
	}
}

func (h *APIHandler) OK(w http.ResponseWriter, message string, data interface{}) {
	wrappedData := h.wrapData(data)
	h.APIHandler.OK(w, message, wrappedData)
}

func (h *APIHandler) Created(w http.ResponseWriter, message string, data interface{}) {
	wrappedData := h.wrapData(data)
	h.APIHandler.Created(w, message, wrappedData)
}

func (h *APIHandler) wrapData(data interface{}) interface{} {
	switch v := data.(type) {
	// Single entities
	case Layout:
		return map[string]interface{}{"layout": v}
	case Section:
		return map[string]interface{}{"section": v}
	case Content:
		return map[string]interface{}{"content": v}
	case Tag:
		return map[string]interface{}{"tag": v}
	case Param:
		return map[string]interface{}{"param": v}
	case Image:
		return map[string]interface{}{"image": v}
	case ImageVariant:
		return map[string]interface{}{"image_variant": v}

	// Slices of entities
	case []Layout:
		return map[string]interface{}{"layouts": v}
	case []Section:
		return map[string]interface{}{"sections": v}
	case []Content:
		return map[string]interface{}{"contents": v}
	case []Tag:
		return map[string]interface{}{"tags": v}
	case []Param:
		return map[string]interface{}{"params": v}
	case []Image:
		return map[string]interface{}{"images": v}
	case []ImageVariant:
		return map[string]interface{}{"image_variants": v}

	// Default case for nil, maps, or other types
	default:
		return data
	}
}

// Publish handles the API request to publish the site.
// PublishRequest represents the request body for the Publish endpoint.
type PublishRequest struct {
	CommitMessage string `json:"commit_message"`
}

// Publish handles the API request to publish the site.
func (h *APIHandler) Publish(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Handling publish request")

	var req PublishRequest
	var err error
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil && err != io.EOF { // io.EOF means empty body, which is fine for optional commit_message
		h.Err(w, http.StatusBadRequest, am.ErrInvalidBody, err)
		return
	}

	var commitURL string
	commitURL, err = h.svc.Publish(r.Context(), req.CommitMessage)
	if err != nil {
		h.Err(w, http.StatusInternalServerError, "error publishing site", err)
		return
	}

	h.OK(w, fmt.Sprintf("Site published successfully: %s", commitURL), nil)
}

// Layout related API handlers

func (h *APIHandler) GetLayout(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling GetLayout", h.Name())

	var err error
	var id uuid.UUID
	id, err = h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resLayoutNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	var layout Layout
	layout, err = h.svc.GetLayout(r.Context(), id)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotGetResource, resLayoutName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgGetItem, resLayoutNameCap)
	h.OK(w, msg, layout)
}

func (h *APIHandler) CreateLayout(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling CreateLayout", h.Name())

	var layout Layout
	var err error
	err = json.NewDecoder(r.Body).Decode(&layout)
	if err != nil {
		h.Err(w, http.StatusBadRequest, am.ErrInvalidBody, err)
		return
	}

	newLayout := Newlayout(layout.Name, layout.Description, layout.Code)
	newLayout.GenCreateValues()

	err = h.svc.CreateLayout(r.Context(), newLayout)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotCreateResource, resLayoutName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgCreateItem, resLayoutNameCap)
	h.Created(w, msg, newLayout)
}

func (h *APIHandler) UpdateLayout(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling UpdateLayout", h.Name())

	var err error
	var id uuid.UUID
	id, err = h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resLayoutNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	var layout Layout
	err = json.NewDecoder(r.Body).Decode(&layout)
	if err != nil {
		h.Err(w, http.StatusBadRequest, am.ErrInvalidBody, err)
		return
	}

	updatedLayout := Newlayout(layout.Name, layout.Description, layout.Code)
	updatedLayout.SetID(id, true)
	updatedLayout.GenUpdateValues()

	err = h.svc.UpdateLayout(r.Context(), updatedLayout)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotUpdateResource, resLayoutName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgUpdateItem, resLayoutNameCap)
	h.OK(w, msg, updatedLayout)
}

func (h *APIHandler) DeleteLayout(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling DeleteLayout", h.Name())

	var err error
	var id uuid.UUID
	id, err = h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resLayoutNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	err = h.svc.DeleteLayout(r.Context(), id)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotDeleteResource, resLayoutName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgDeleteItem, resLayoutNameCap)
	h.OK(w, msg, json.RawMessage("null"))
}

func (h *APIHandler) GetAllLayouts(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling GetAllLayouts", h.Name())

	var layouts []Layout
	var err error
	layouts, err = h.svc.GetAllLayouts(r.Context())
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotGetResources, resLayoutName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgGetAllItems, resLayoutNameCap)
	h.OK(w, msg, layouts)
}

// Section related API handlers

func (h *APIHandler) GetAllSections(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling GetAllSections", h.Name())

	var sections []Section
	var err error
	sections, err = h.svc.GetSections(r.Context())
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotGetResources, resSectionName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgGetAllItems, resSectionNameCap)
	h.OK(w, msg, sections)
}

func (h *APIHandler) GetSection(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling GetSection", h.Name())

	var err error
	var id uuid.UUID
	id, err = h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resSectionNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	var section Section
	section, err = h.svc.GetSection(r.Context(), id)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotGetResource, resSectionName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgGetItem, resSectionNameCap)
	h.OK(w, msg, section)
}

func (h *APIHandler) CreateSection(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling CreateSection", h.Name())

	var section Section
	var err error
	err = json.NewDecoder(r.Body).Decode(&section)
	if err != nil {
		h.Err(w, http.StatusBadRequest, am.ErrInvalidBody, err)
		return
	}

	newSection := NewSection(section.Name, section.Description, section.Path, section.LayoutID)
	newSection.GenCreateValues()

	err = h.svc.CreateSection(r.Context(), newSection)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotCreateResource, resSectionName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgCreateItem, resSectionNameCap)
	h.Created(w, msg, newSection)
}

func (h *APIHandler) UpdateSection(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling UpdateSection", h.Name())

	var err error
	var id uuid.UUID
	id, err = h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resSectionNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	var section Section
	err = json.NewDecoder(r.Body).Decode(&section)
	if err != nil {
		h.Err(w, http.StatusBadRequest, am.ErrInvalidBody, err)
		return
	}

	updatedSection := NewSection(section.Name, section.Description, section.Path, section.LayoutID)
	updatedSection.SetID(id, true)
	updatedSection.GenUpdateValues()

	err = h.svc.UpdateSection(r.Context(), updatedSection)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotUpdateResource, resSectionName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgUpdateItem, resSectionNameCap)
	h.OK(w, msg, updatedSection)
}

func (h *APIHandler) DeleteSection(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling DeleteSection", h.Name())

	var err error
	var id uuid.UUID
	id, err = h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resSectionNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	err = h.svc.DeleteSection(r.Context(), id)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotDeleteResource, resSectionName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgDeleteItem, resSectionNameCap)
	h.OK(w, msg, json.RawMessage("null"))
}

// Content related API handlers

func (h *APIHandler) GetAllContent(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling GetAllContentWithMeta", h.Name())

	var contents []Content
	var err error
	contents, err = h.svc.GetAllContentWithMeta(r.Context())
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotGetResources, resContent)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgGetAllItems, resContentCap)
	h.OK(w, msg, contents)
}

func (h *APIHandler) GetContent(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling GetContent", h.Name())

	var err error
	var id uuid.UUID
	id, err = h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resContentCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	var content Content
	content, err = h.svc.GetContent(r.Context(), id)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotGetResource, resContent)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	// The new GetContent service method should already include tags.
	// If not, this logic is a fallback.
	if len(content.Tags) == 0 {
		var tags []Tag
		tags, err = h.svc.GetTagsForContent(r.Context(), id)
		if err != nil {
			msg := fmt.Sprintf("Cannot get tags for content %s", id)
			h.Err(w, http.StatusInternalServerError, msg, err)
			return
		}
		content.Tags = tags
	}

	msg := fmt.Sprintf(am.MsgGetItem, resContentCap)
	h.OK(w, msg, content)
}

func (h *APIHandler) CreateContent(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling CreateContent", h.Name())

	var content Content
	var err error
	err = json.NewDecoder(r.Body).Decode(&content)
	if err != nil {
		h.Err(w, http.StatusBadRequest, am.ErrInvalidBody, err)
		return
	}

	content.GenCreateValues()

	err = h.svc.CreateContent(r.Context(), &content)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotCreateResource, resContent)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	for _, tag := range content.Tags {
		err = h.svc.AddTagToContent(r.Context(), content.ID, tag.Name)
		if err != nil {
			msg := fmt.Sprintf("Cannot add tag %s to content %s", tag.Name, content.ID)
			h.Err(w, http.StatusInternalServerError, msg, err)
			return
		}
	}

	msg := fmt.Sprintf(am.MsgCreateItem, resContentCap)
	h.Created(w, msg, content)
}

func (h *APIHandler) UpdateContent(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling UpdateContent", h.Name())

	var err error
	var id uuid.UUID
	id, err = h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resContentCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	var content Content
	err = json.NewDecoder(r.Body).Decode(&content)
	if err != nil {
		h.Err(w, http.StatusBadRequest, am.ErrInvalidBody, err)
		return
	}

	content.SetID(id, true)   // Set the ID from the URL on the decoded content
	content.GenUpdateValues() // Generate audit values

	err = h.svc.UpdateContent(r.Context(), &content)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotUpdateResource, resContent)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	// Remove existing tags
	var existingTags []Tag
	existingTags, err = h.svc.GetTagsForContent(r.Context(), id)
	if err != nil {
		msg := fmt.Sprintf("Cannot get existing tags for content %s", id)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}
	for _, tag := range existingTags {
		err = h.svc.RemoveTagFromContent(r.Context(), id, tag.ID)
		if err != nil {
			msg := fmt.Sprintf("Cannot remove tag %s from content %s", tag.ID, id)
			h.Err(w, http.StatusInternalServerError, msg, err)
			return
		}
	}

	// Add new tags
	for _, tag := range content.Tags {
		err = h.svc.AddTagToContent(r.Context(), id, tag.Name)
		if err != nil {
			msg := fmt.Sprintf("Cannot add tag %s to content %s", tag.Name, id)
			h.Err(w, http.StatusInternalServerError, msg, err)
			return
		}
	}

	msg := fmt.Sprintf(am.MsgUpdateItem, resContentCap)
	h.OK(w, msg, content)
}

func (h *APIHandler) DeleteContent(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling DeleteContent", h.Name())

	var err error
	var id uuid.UUID
	id, err = h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resContentCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	err = h.svc.DeleteContent(r.Context(), id)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotDeleteResource, resContent)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgDeleteItem, resContentCap)
	h.OK(w, msg, json.RawMessage("null"))
}

// Tag related API handlers

func (h *APIHandler) CreateTag(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling CreateTag", h.Name())

	var tag Tag
	var err error
	err = json.NewDecoder(r.Body).Decode(&tag)
	if err != nil {
		h.Err(w, http.StatusBadRequest, am.ErrInvalidBody, err)
		return
	}

	newTag := NewTag(tag.Name)
	newTag.GenCreateValues()

	err = h.svc.CreateTag(r.Context(), newTag)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotCreateResource, resTagName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgCreateItem, resTagNameCap)
	h.Created(w, msg, newTag)
}

func (h *APIHandler) GetTag(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling GetTag", h.Name())

	var err error
	var id uuid.UUID
	id, err = h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resTagNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	var tag Tag
	tag, err = h.svc.GetTag(r.Context(), id)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotGetResource, resTagName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgGetItem, resTagNameCap)
	h.OK(w, msg, tag)
}

func (h *APIHandler) GetTagByName(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling GetTagByName", h.Name())

	var err error
	var name string
	name, err = h.Param(w, r, "name")
	if err != nil {
		msg := fmt.Sprintf("%s: %s", am.ErrInvalidParam, "name")
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	var tag Tag
	tag, err = h.svc.GetTagByName(r.Context(), name)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotGetResource, resTagName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgGetItem, resTagNameCap)
	h.OK(w, msg, tag)
}

func (h *APIHandler) GetAllTags(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling GetAllTags", h.Name())

	var tags []Tag
	var err error
	tags, err = h.svc.GetAllTags(r.Context())
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotGetResources, resTagName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgGetAllItems, resTagNameCap)
	h.OK(w, msg, tags)
}

func (h *APIHandler) UpdateTag(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling UpdateTag", h.Name())

	var err error
	var id uuid.UUID
	id, err = h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resTagNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	var tag Tag
	err = json.NewDecoder(r.Body).Decode(&tag)
	if err != nil {
		h.Err(w, http.StatusBadRequest, am.ErrInvalidBody, err)
		return
	}

	updatedTag := NewTag(tag.Name)
	updatedTag.SetID(id, true)
	updatedTag.GenUpdateValues()

	err = h.svc.UpdateTag(r.Context(), updatedTag)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotUpdateResource, resTagName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgUpdateItem, resTagNameCap)
	h.OK(w, msg, updatedTag)
}

func (h *APIHandler) DeleteTag(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling DeleteTag", h.Name())

	var err error
	var id uuid.UUID
	id, err = h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resTagNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	err = h.svc.DeleteTag(r.Context(), id)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotDeleteResource, resTagName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgDeleteItem, resTagNameCap)
	h.OK(w, msg, json.RawMessage("null"))
}

// Param related API handlers

func (h *APIHandler) CreateParam(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling CreateParam", h.Name())

	var param Param
	var err error
	err = json.NewDecoder(r.Body).Decode(&param)
	if err != nil {
		h.Err(w, http.StatusBadRequest, am.ErrInvalidBody, err)
		return
	}

	newParam := NewParam(param.Name, param.Value)
	newParam.Description = param.Description
	newParam.RefKey = param.RefKey
	newParam.GenCreateValues()

	err = h.svc.CreateParam(r.Context(), &newParam)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotCreateResource, resParamName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgCreateItem, resParamNameCap)
	h.Created(w, msg, newParam)
}

func (h *APIHandler) GetParam(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling GetParam", h.Name())

	var err error
	var id uuid.UUID
	id, err = h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resParamNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	var param Param
	param, err = h.svc.GetParam(r.Context(), id)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotGetResource, resParamName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgGetItem, resParamNameCap)
	h.OK(w, msg, param)
}

func (h *APIHandler) GetParamByName(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling GetParamByName", h.Name())

	var err error
	var name string
	name, err = h.Param(w, r, "name")
	if err != nil {
		msg := fmt.Sprintf("%s: %s", am.ErrInvalidParam, "name")
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	var param Param
	param, err = h.svc.GetParamByName(r.Context(), name)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotGetResource, resParamName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgGetItem, resParamNameCap)
	h.OK(w, msg, param)
}

func (h *APIHandler) GetParamByRefKey(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling GetParamByRefKey", h.Name())

	var err error
	var refKey string
	refKey, err = h.Param(w, r, "ref_key")
	if err != nil {
		msg := fmt.Sprintf("%s: %s", am.ErrInvalidParam, "ref_key")
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	var param Param
	param, err = h.svc.GetParamByRefKey(r.Context(), refKey)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotGetResource, resParamName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgGetItem, resParamNameCap)
	h.OK(w, msg, param)
}

func (h *APIHandler) ListParams(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling ListParams", h.Name())

	var params []Param
	var err error
	params, err = h.svc.ListParams(r.Context())
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotGetResources, resParamName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgGetAllItems, resParamNameCap)
	h.OK(w, msg, params)
}

func (h *APIHandler) UpdateParam(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling UpdateParam", h.Name())

	var err error
	var id uuid.UUID
	id, err = h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resParamNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	// Get existing param to check if it's a system param
	var existingParam Param
	existingParam, err = h.svc.GetParam(r.Context(), id)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotGetResource, resParamName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	var param Param
	err = json.NewDecoder(r.Body).Decode(&param)
	if err != nil {
		h.Err(w, http.StatusBadRequest, am.ErrInvalidBody, err)
		return
	}

	if existingParam.System == 1 {
		// For system params, only Value can be updated
		if param.Name != existingParam.Name {
			h.Err(w, http.StatusBadRequest, "cannot change name of system parameter", nil)
			return
		}
		if param.RefKey != existingParam.RefKey {
			h.Err(w, http.StatusBadRequest, "cannot change ref key of system parameter", nil)
			return
		}
		if param.Description != existingParam.Description {
			h.Err(w, http.StatusBadRequest, "cannot change description of system parameter", nil)
			return
		}
		// Use existing param's Name, RefKey, and Description, only update Value
		param.Name = existingParam.Name
		param.RefKey = existingParam.RefKey
		param.Description = existingParam.Description
	}

	updatedParam := NewParam(param.Name, param.Value)
	updatedParam.Description = param.Description
	updatedParam.RefKey = param.RefKey
	updatedParam.System = existingParam.System // Preserve system flag
	updatedParam.SetID(id, true)
	updatedParam.GenUpdateValues()

	err = h.svc.UpdateParam(r.Context(), &updatedParam)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotUpdateResource, resParamName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgUpdateItem, resParamNameCap)
	h.OK(w, msg, updatedParam)
}

func (h *APIHandler) DeleteParam(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling DeleteParam", h.Name())

	var err error
	var id uuid.UUID
	id, err = h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resParamNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	// Get existing param to check if it's a system param
	var existingParam Param
	existingParam, err = h.svc.GetParam(r.Context(), id)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotGetResource, resParamName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	if existingParam.System == 1 {
		h.Err(w, http.StatusForbidden, "cannot delete system parameter", nil)
		return
	}

	err = h.svc.DeleteParam(r.Context(), id)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotDeleteResource, resParamName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgDeleteItem, resParamNameCap)
	h.OK(w, msg, json.RawMessage("null"))
}

// Image related API handlers

func (h *APIHandler) CreateImage(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling CreateImage", h.Name())

	var image Image
	var err error
	err = json.NewDecoder(r.Body).Decode(&image)
	if err != nil {
		h.Err(w, http.StatusBadRequest, am.ErrInvalidBody, err)
		return
	}

	newImage := NewImage() // Call constructor without arguments
	// Assign fields from the decoded JSON to the newImage instance
	newImage.ContentHash = image.ContentHash
	newImage.Mime = image.Mime
	newImage.Width = image.Width
	newImage.Height = image.Height
	newImage.FilesizeByte = image.FilesizeByte
	newImage.Etag = image.Etag
	newImage.Title = image.Title
	newImage.AltText = image.AltText
	newImage.AltLang = image.AltLang
	newImage.LongDescription = image.LongDescription
	newImage.Caption = image.Caption
	newImage.Decorative = image.Decorative
	newImage.DescribedByID = image.DescribedByID

	newImage.GenCreateValues()

	err = h.svc.CreateImage(r.Context(), &newImage)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotCreateResource, resImageName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgCreateItem, resImageNameCap)
	h.Created(w, msg, newImage)
}

func (h *APIHandler) GetImage(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling GetImage", h.Name())

	var err error
	var id uuid.UUID
	id, err = h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resImageNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	var image Image
	image, err = h.svc.GetImage(r.Context(), id)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotGetResource, resImageName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgGetItem, resImageNameCap)
	h.OK(w, msg, image)
}

func (h *APIHandler) GetImageByShortID(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling GetImageByShortID", h.Name())

	var err error
	var shortID string
	shortID, err = h.Param(w, r, "short_id")
	if err != nil {
		msg := fmt.Sprintf("%s: %s", am.ErrInvalidParam, "short_id")
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	var image Image
	image, err = h.svc.GetImageByShortID(r.Context(), shortID)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotGetResource, resImageName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgGetItem, resImageNameCap)
	h.OK(w, msg, image)
}

func (h *APIHandler) ListImages(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling ListImages", h.Name())

	var images []Image
	var err error
	images, err = h.svc.ListImages(r.Context())
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotGetResources, resImageName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgGetAllItems, resImageNameCap)
	h.OK(w, msg, images)
}

func (h *APIHandler) UpdateImage(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling UpdateImage", h.Name())

	var err error
	var id uuid.UUID
	id, err = h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resImageNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	var image Image
	err = json.NewDecoder(r.Body).Decode(&image)
	if err != nil {
		h.Err(w, http.StatusBadRequest, am.ErrInvalidBody, err)
		return
	}

	updatedImage := NewImage()   // Call constructor without arguments
	updatedImage.SetID(id, true) // Set the ID from the URL on the decoded content

	// Assign fields from the decoded JSON to the updatedImage instance
	updatedImage.ContentHash = image.ContentHash
	updatedImage.Mime = image.Mime
	updatedImage.Width = image.Width
	updatedImage.Height = image.Height
	updatedImage.FilesizeByte = image.FilesizeByte
	updatedImage.Etag = image.Etag
	updatedImage.Title = image.Title
	updatedImage.AltText = image.AltText
	updatedImage.AltLang = image.AltLang
	updatedImage.LongDescription = image.LongDescription
	updatedImage.Caption = image.Caption
	updatedImage.Decorative = image.Decorative
	updatedImage.DescribedByID = image.DescribedByID

	updatedImage.GenUpdateValues()

	err = h.svc.UpdateImage(r.Context(), &updatedImage)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotUpdateResource, resImageName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgUpdateItem, resImageNameCap)
	h.OK(w, msg, updatedImage)
}

func (h *APIHandler) DeleteImage(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling DeleteImage", h.Name())

	var err error
	var id uuid.UUID
	id, err = h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resImageNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	err = h.svc.DeleteImage(r.Context(), id)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotDeleteResource, resImageName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgDeleteItem, resImageNameCap)
	h.OK(w, msg, json.RawMessage("null"))
}

// ImageVariant related API handlers

func (h *APIHandler) CreateImageVariant(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling CreateImageVariant", h.Name())

	var variant ImageVariant
	var err error
	err = json.NewDecoder(r.Body).Decode(&variant)
	if err != nil {
		h.Err(w, http.StatusBadRequest, am.ErrInvalidBody, err)
		return
	}

	newVariant := NewImageVariant() // Call constructor without arguments
	// Assign fields from the decoded JSON to the newVariant instance
	newVariant.ImageID = variant.ImageID
	newVariant.Kind = variant.Kind
	newVariant.Width = variant.Width
	newVariant.Height = variant.Height
	newVariant.FilesizeByte = variant.FilesizeByte
	newVariant.Mime = variant.Mime
	newVariant.BlobRef = variant.BlobRef

	newVariant.GenCreateValues()

	err = h.svc.CreateImageVariant(r.Context(), &newVariant)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotCreateResource, resImageVariantName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgCreateItem, resImageVariantNameCap)
	h.Created(w, msg, newVariant)
}

func (h *APIHandler) GetImageVariant(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling GetImageVariant", h.Name())

	var err error
	var id uuid.UUID
	id, err = h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resImageVariantNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	var variant ImageVariant
	variant, err = h.svc.GetImageVariant(r.Context(), id)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotGetResource, resImageVariantName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgGetItem, resImageVariantNameCap)
	h.OK(w, msg, variant)
}

func (h *APIHandler) ListImageVariantsByImageID(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling ListImageVariantsByImageID", h.Name())

	var err error
	var imageIDStr string
	imageIDStr, err = h.Param(w, r, "image_id")
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resImageNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}
	var imageID uuid.UUID
	imageID, err = uuid.Parse(imageIDStr)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resImageNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	var variants []ImageVariant
	variants, err = h.svc.ListImageVariantsByImageID(r.Context(), imageID)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotGetResources, resImageVariantName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgGetAllItems, resImageVariantNameCap)
	h.OK(w, msg, variants)
}

func (h *APIHandler) UpdateImageVariant(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling UpdateImageVariant", h.Name())

	var err error
	var id uuid.UUID
	id, err = h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resImageVariantNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	var variant ImageVariant
	err = json.NewDecoder(r.Body).Decode(&variant)
	if err != nil {
		h.Err(w, http.StatusBadRequest, am.ErrInvalidBody, err)
		return
	}

	updatedVariant := NewImageVariant() // Call constructor without arguments
	updatedVariant.SetID(id, true)      // Set the ID from the URL on the decoded content

	// Assign fields from the decoded JSON to the updatedVariant instance
	updatedVariant.ImageID = variant.ImageID
	updatedVariant.Kind = variant.Kind
	updatedVariant.Width = variant.Width
	updatedVariant.Height = variant.Height
	updatedVariant.FilesizeByte = variant.FilesizeByte
	updatedVariant.Mime = variant.Mime
	updatedVariant.BlobRef = variant.BlobRef

	updatedVariant.GenUpdateValues()

	err = h.svc.UpdateImageVariant(r.Context(), &updatedVariant)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotUpdateResource, resImageVariantName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgUpdateItem, resImageVariantNameCap)
	h.OK(w, msg, updatedVariant)
}

func (h *APIHandler) DeleteImageVariant(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling DeleteImageVariant", h.Name())

	var err error
	var id uuid.UUID
	id, err = h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resImageVariantNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	err = h.svc.DeleteImageVariant(r.Context(), id)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotDeleteResource, resImageVariantName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgDeleteItem, resImageVariantNameCap)
	h.OK(w, msg, json.RawMessage("null"))
}

// AddTagToContentForm represents the data for adding a tag to content.
type AddTagToContentForm struct {
	Name string `json:"name"`
}

// Content-Tag related API handlers

func (h *APIHandler) AddTagToContent(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling AddTagToContent", h.Name())

	var err error
	var contentIDStr string
	contentIDStr, err = h.Param(w, r, "content_id")
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resContentCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}
	var contentID uuid.UUID
	contentID, err = uuid.Parse(contentIDStr)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resContentCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	var form AddTagToContentForm
	err = json.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		h.Err(w, http.StatusBadRequest, am.ErrInvalidBody, err)
		return
	}

	err = h.svc.AddTagToContent(r.Context(), contentID, form.Name)
	if err != nil {
		msg := fmt.Sprintf("Cannot add tag %s to content %s", form.Name, contentID)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf("Tag %s added to content %s", form.Name, contentID)
	h.Created(w, msg, nil)
}

func (h *APIHandler) RemoveTagFromContent(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling RemoveTagFromContent", h.Name())

	var err error
	var contentIDStr string
	contentIDStr, err = h.Param(w, r, "content_id")
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resContentCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}
	var contentID uuid.UUID
	contentID, err = uuid.Parse(contentIDStr)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resContentCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	var tagIDStr string
	tagIDStr, err = h.Param(w, r, "tag_id")
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resTagNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}
	var tagID uuid.UUID
	tagID, err = uuid.Parse(tagIDStr)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resTagNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	err = h.svc.RemoveTagFromContent(r.Context(), contentID, tagID)
	if err != nil {
		msg := fmt.Sprintf("Cannot remove tag %s from content %s", tagID, contentID)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf("Tag %s removed from content %s", tagID, contentID)
	h.OK(w, msg, json.RawMessage("null"))
}

func (h *APIHandler) GenerateMarkdown(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling GenerateMarkdown", h.Name())

	var err error
	err = h.svc.GenerateMarkdown(r.Context())
	if err != nil {
		msg := fmt.Sprintf("Cannot generate markdown: %v", err)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := "Markdown generation process started successfully"
	h.OK(w, msg, nil)
}

func (h *APIHandler) GenerateHTML(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling GenerateHTML", h.Name())

	var err error
	err = h.svc.GenerateHTMLFromContent(r.Context())
	if err != nil {
		msg := fmt.Sprintf("Cannot generate HTML: %v", err)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := "HTML generation process started successfully"
	h.OK(w, msg, nil)
}
