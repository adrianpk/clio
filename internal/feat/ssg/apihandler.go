package ssg

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/adrianpk/clio/internal/am"
)

const (
	resContentName         = "content"
	resContentNameCap      = "Content"
	resSectionName         = "section"
	resSectionNameCap      = "Section"
	resLayoutName          = "layout"
	resLayoutNameCap       = "Layout"
	resTagName             = "tag"
	resTagNameCap          = "Tag"
	resParamName           = "param"
	resParamNameCap        = "Param"
	resImageName           = "image"
	resImageNameCap        = "Image"
	resImageVariantName    = "image variant"
	resImageVariantNameCap = "Image variant"
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
	if data == nil {
		return map[string]interface{}{
			"status": "success",
			"data":   nil,
		}
	}
	return map[string]interface{}{
		"status": "success",
		"data":   data,
	}
}

func (h *APIHandler) Publish(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling Publish", h.Name())

	var err error

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.Err(w, http.StatusBadRequest, am.ErrInvalidBody, err)
		return
	}

	var data PublishRequest
	err = json.Unmarshal(body, &data)
	if err != nil {
		h.Err(w, http.StatusBadRequest, am.ErrInvalidBody, err)
		return
	}

	// Run the publish process
	commitURL, err := h.svc.Publish(r.Context(), data.Message)
	if err != nil {
		msg := fmt.Sprintf("Cannot publish: %v", err)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := "Publish process started successfully"
	result := map[string]string{"commitURL": commitURL}
	h.OK(w, msg, result)
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

// PublishRequest represents the data for a publish request.
type PublishRequest struct {
	Message string `json:"message"`
}

// AddTagToContentForm represents the data for adding a tag to content.
type AddTagToContentForm struct {
	Name string `json:"name"`
}