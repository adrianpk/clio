package ssg

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/adrianpk/clio/internal/am"

	"github.com/google/uuid"
)

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