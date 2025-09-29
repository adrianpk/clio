package ssg

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/adrianpk/clio/internal/am"

	"github.com/google/uuid"
)

func (h *APIHandler) GetAllContent(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling GetAllContentWithMeta", h.Name())

	var contents []Content
	var err error
	contents, err = h.svc.GetAllContentWithMeta(r.Context())
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

	var err error
	var id uuid.UUID
	id, err = h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resContentNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	var content Content
	content, err = h.svc.GetContent(r.Context(), id)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotGetResource, resContentName)
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

	msg := fmt.Sprintf(am.MsgGetItem, resContentNameCap)
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

	var err error
	var id uuid.UUID
	id, err = h.ID(w, r)
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

	msg := fmt.Sprintf(am.MsgUpdateItem, resContentNameCap)
	h.OK(w, msg, content)
}

func (h *APIHandler) DeleteContent(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling DeleteContent", h.Name())

	var err error
	var id uuid.UUID
	id, err = h.ID(w, r)
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

func (h *APIHandler) AddTagToContent(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling AddTagToContent", h.Name())

	var err error
	var contentIDStr string
	contentIDStr, err = h.Param(w, r, "content_id")
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resContentNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}
	var contentID uuid.UUID
	contentID, err = uuid.Parse(contentIDStr)
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

	var err error
	var contentIDStr string
	contentIDStr, err = h.Param(w, r, "content_id")
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resContentNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}
	var contentID uuid.UUID
	contentID, err = uuid.Parse(contentIDStr)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resContentNameCap)
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