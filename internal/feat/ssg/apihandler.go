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
	resContentName    = "content"
	resContentNameCap = "Content"
	resSectionName    = "section"
	resSectionNameCap = "Section"
	resLayoutName     = "layout"
	resLayoutNameCap  = "Layout"
	resTagName        = "tag"
	resTagNameCap     = "Tag"
	resParamName      = "param"
	resParamNameCap   = "Param"
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
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil && err != io.EOF { // io.EOF means empty body, which is fine for optional commit_message
		h.Err(w, http.StatusBadRequest, am.ErrInvalidBody, err)
		return
	}

	commitURL, err := h.svc.Publish(r.Context(), req.CommitMessage)
	if err != nil {
		h.Err(w, http.StatusInternalServerError, "error publishing site", err)
		return
	}

	h.OK(w, fmt.Sprintf("Site published successfully: %s", commitURL), nil)
}

// Layout related API handlers

func (h *APIHandler) GetLayout(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling GetLayout", h.Name())

	id, err := h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resLayoutNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	layout, err := h.svc.GetLayout(r.Context(), id)
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
	err := json.NewDecoder(r.Body).Decode(&layout)
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

	id, err := h.ID(w, r)
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

	id, err := h.ID(w, r)
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

	layouts, err := h.svc.GetAllLayouts(r.Context())
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

	sections, err := h.svc.GetSections(r.Context())
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

	id, err := h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resSectionNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	section, err := h.svc.GetSection(r.Context(), id)
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
	err := json.NewDecoder(r.Body).Decode(&section)
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

	id, err := h.ID(w, r)
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

	id, err := h.ID(w, r)
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

	contents, err := h.svc.GetAllContentWithMeta(r.Context())
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotGetResources, resContentName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgGetAllItems, resContentNameCap)
	h.OK(w, msg, contents)
}

func (h *APIHandler) GetContent(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling GetContent", h.Name())

	id, err := h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resContentNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	content, err := h.svc.GetContent(r.Context(), id)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotGetResource, resContentName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	// The new GetContent service method should already include tags.
	// If not, this logic is a fallback.
	if len(content.Tags) == 0 {
		tags, err := h.svc.GetTagsForContent(r.Context(), id)
		if err != nil {
			msg := fmt.Sprintf("Cannot get tags for content %s", id)
			h.Err(w, http.StatusInternalServerError, msg, err)
			return
		}
		content.Tags = tags
	}

	msg := fmt.Sprintf(am.MsgGetItem, resContentNameCap)
	h.OK(w, msg, content)
}

func (h *APIHandler) CreateContent(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling CreateContent", h.Name())

	var content Content
	err := json.NewDecoder(r.Body).Decode(&content)
	if err != nil {
		h.Err(w, http.StatusBadRequest, am.ErrInvalidBody, err)
		return
	}

	content.GenCreateValues()

	err = h.svc.CreateContent(r.Context(), &content)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotCreateResource, resContentName)
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

	msg := fmt.Sprintf(am.MsgCreateItem, resContentNameCap)
	h.Created(w, msg, content)
}

func (h *APIHandler) UpdateContent(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling UpdateContent", h.Name())

	id, err := h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resContentNameCap)
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
		msg := fmt.Sprintf(am.ErrCannotUpdateResource, resContentName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	// Remove existing tags
	existingTags, err := h.svc.GetTagsForContent(r.Context(), id)
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

	msg := fmt.Sprintf(am.MsgUpdateItem, resContentNameCap)
	h.OK(w, msg, content)
}

func (h *APIHandler) DeleteContent(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling DeleteContent", h.Name())

	id, err := h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resContentNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	err = h.svc.DeleteContent(r.Context(), id)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotDeleteResource, resContentName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgDeleteItem, resContentNameCap)
	h.OK(w, msg, json.RawMessage("null"))
}

// Tag related API handlers

func (h *APIHandler) CreateTag(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling CreateTag", h.Name())

	var tag Tag
	err := json.NewDecoder(r.Body).Decode(&tag)
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

	id, err := h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resTagNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	tag, err := h.svc.GetTag(r.Context(), id)
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

	name, err := h.Param(w, r, "name")
	if err != nil {
		msg := fmt.Sprintf("%s: %s", am.ErrInvalidParam, "name")
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	tag, err := h.svc.GetTagByName(r.Context(), name)
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

	tags, err := h.svc.GetAllTags(r.Context())
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

	id, err := h.ID(w, r)
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

	id, err := h.ID(w, r)
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
	err := json.NewDecoder(r.Body).Decode(&param)
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

	id, err := h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resParamNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	param, err := h.svc.GetParam(r.Context(), id)
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

	name, err := h.Param(w, r, "name")
	if err != nil {
		msg := fmt.Sprintf("%s: %s", am.ErrInvalidParam, "name")
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	param, err := h.svc.GetParamByName(r.Context(), name)
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

	refKey, err := h.Param(w, r, "ref_key")
	if err != nil {
		msg := fmt.Sprintf("%s: %s", am.ErrInvalidParam, "ref_key")
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	param, err := h.svc.GetParamByRefKey(r.Context(), refKey)
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

	params, err := h.svc.ListParams(r.Context())
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

	id, err := h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resParamNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	var param Param
	err = json.NewDecoder(r.Body).Decode(&param)
	if err != nil {
		h.Err(w, http.StatusBadRequest, am.ErrInvalidBody, err)
		return
	}

	updatedParam := NewParam(param.Name, param.Value)
	updatedParam.Description = param.Description
	updatedParam.RefKey = param.RefKey
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

	id, err := h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resParamNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
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

// AddTagToContentForm represents the data for adding a tag to content.
type AddTagToContentForm struct {
	Name string `json:"name"`
}

// Content-Tag related API handlers

func (h *APIHandler) AddTagToContent(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling AddTagToContent", h.Name())

	contentIDStr, err := h.Param(w, r, "content_id")
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resContentNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}
	contentID, err := uuid.Parse(contentIDStr)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resContentNameCap)
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

	contentIDStr, err := h.Param(w, r, "content_id")
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resContentNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}
	contentID, err := uuid.Parse(contentIDStr)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resContentNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	tagIDStr, err := h.Param(w, r, "tag_id")
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resTagNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}
	tagID, err := uuid.Parse(tagIDStr)
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

	err := h.svc.GenerateMarkdown(r.Context())
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

	err := h.svc.GenerateHTMLFromContent(r.Context())
	if err != nil {
		msg := fmt.Sprintf("Cannot generate HTML: %v", err)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := "HTML generation process started successfully"
	h.OK(w, msg, nil)
}
