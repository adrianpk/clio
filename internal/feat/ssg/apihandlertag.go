package ssg

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/adrianpk/clio/internal/am"

	"github.com/google/uuid"
)

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

	msg := fmt.Sprintf(am.MsgCreateItem, am.Cap(resTagName))
	h.Created(w, msg, newTag)
}

func (h *APIHandler) GetTag(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling GetTag", h.Name())

	var err error
	var id uuid.UUID
	id, err = h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, am.Cap(resTagName))
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

	msg := fmt.Sprintf(am.MsgGetItem, am.Cap(resTagName))
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

	msg := fmt.Sprintf(am.MsgGetItem, am.Cap(resTagName))
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

	msg := fmt.Sprintf(am.MsgGetAllItems, am.Cap(resTagName))
	h.OK(w, msg, tags)
}

func (h *APIHandler) UpdateTag(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling UpdateTag", h.Name())

	var err error
	var id uuid.UUID
	id, err = h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, am.Cap(resTagName))
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

	msg := fmt.Sprintf(am.MsgUpdateItem, am.Cap(resTagName))
	h.OK(w, msg, updatedTag)
}

func (h *APIHandler) DeleteTag(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling DeleteTag", h.Name())

	var err error
	var id uuid.UUID
	id, err = h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, am.Cap(resTagName))
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	err = h.svc.DeleteTag(r.Context(), id)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotDeleteResource, resTagName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgDeleteItem, am.Cap(resTagName))
	h.OK(w, msg, json.RawMessage("null"))
}