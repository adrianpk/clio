package ssg

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/adrianpk/clio/internal/am"

	"github.com/google/uuid"
)

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