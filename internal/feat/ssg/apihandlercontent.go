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

	msg := fmt.Sprintf(am.MsgGetAllItems, am.Cap(resContentName))
	h.OK(w, msg, contents)
}

func (h *APIHandler) GetContent(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling GetContent", h.Name())

	var err error
	var id uuid.UUID
	id, err = h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, am.Cap(resContentName))
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

	msg := fmt.Sprintf(am.MsgGetItem, am.Cap(resContentName))
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

	msg := fmt.Sprintf(am.MsgCreateItem, am.Cap(resContentName))
	h.Created(w, msg, content)
}

func (h *APIHandler) UpdateContent(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling UpdateContent", h.Name())

	var err error
	var id uuid.UUID
	id, err = h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, am.Cap(resContentName))
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

	msg := fmt.Sprintf(am.MsgUpdateItem, am.Cap(resContentName))
	h.OK(w, msg, content)
}

func (h *APIHandler) DeleteContent(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling DeleteContent", h.Name())

	var err error
	var id uuid.UUID
	id, err = h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, am.Cap(resContentName))
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	err = h.svc.DeleteContent(r.Context(), id)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotDeleteResource, resContentName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgDeleteItem, am.Cap(resContentName))
	h.OK(w, msg, json.RawMessage("null"))
}

func (h *APIHandler) AddTagToContent(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling AddTagToContent", h.Name())

	var err error
	var contentIDStr string
	contentIDStr, err = h.Param(w, r, "content_id")
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, am.Cap(resContentName))
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}
	var contentID uuid.UUID
	contentID, err = uuid.Parse(contentIDStr)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, am.Cap(resContentName))
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
		msg := fmt.Sprintf(am.ErrInvalidID, am.Cap(resContentName))
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}
	var contentID uuid.UUID
	contentID, err = uuid.Parse(contentIDStr)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, am.Cap(resContentName))
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	var tagIDStr string
	tagIDStr, err = h.Param(w, r, "tag_id")
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, am.Cap(resTagName))
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}
	var tagID uuid.UUID
	tagID, err = uuid.Parse(tagIDStr)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, am.Cap(resTagName))
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

// UploadContentImage handles image upload for content (header or content images)
func (h *APIHandler) UploadContentImage(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling UploadContentImage", h.Name())

	// Parse content ID from URL
	contentIDStr, err := h.Param(w, r, "content_id")
	if err != nil {
		h.Err(w, http.StatusBadRequest, "Invalid content ID", err)
		return
	}

	contentID, err := uuid.Parse(contentIDStr)
	if err != nil {
		h.Err(w, http.StatusBadRequest, "Invalid content ID format", err)
		return
	}

	// Parse image type from form
	imageTypeStr := r.FormValue("image_type")
	if imageTypeStr == "" {
		h.Err(w, http.StatusBadRequest, "Missing image_type parameter", nil)
		return
	}

	imageType := ImageType(imageTypeStr)
	if imageType != ImageTypeContent && imageType != ImageTypeHeader {
		h.Err(w, http.StatusBadRequest, "Invalid image_type for content", nil)
		return
	}

	// Parse uploaded file
	file, header, err := r.FormFile("image")
	if err != nil {
		h.Err(w, http.StatusBadRequest, "Failed to parse uploaded file", err)
		return
	}
	defer file.Close()

	// Process upload through service
	result, err := h.svc.UploadContentImage(r.Context(), contentID, file, header, imageType)
	if err != nil {
		h.Err(w, http.StatusInternalServerError, "Failed to upload image", err)
		return
	}

	// Return success response
	msg := fmt.Sprintf("Image uploaded successfully: %s", result.Filename)
	h.OK(w, msg, map[string]interface{}{
		"filename":      result.Filename,
		"relative_path": result.RelativePath,
		"metadata":      result.Metadata,
	})
}

// GetContentImages returns all images for a specific content
func (h *APIHandler) GetContentImages(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling GetContentImages", h.Name())

	// Parse content ID from URL
	contentIDStr, err := h.Param(w, r, "content_id")
	if err != nil {
		h.Err(w, http.StatusBadRequest, "Invalid content ID", err)
		return
	}

	contentID, err := uuid.Parse(contentIDStr)
	if err != nil {
		h.Err(w, http.StatusBadRequest, "Invalid content ID format", err)
		return
	}

	// Get images through service
	images, err := h.svc.GetContentImages(r.Context(), contentID)
	if err != nil {
		h.Err(w, http.StatusInternalServerError, "Failed to get content images", err)
		return
	}

	// Return images list
	msg := fmt.Sprintf("Retrieved %d images for content", len(images))
	h.OK(w, msg, map[string]interface{}{
		"images": images,
	})
}

// DeleteContentImageRequest represents the request body for deleting content images
type DeleteContentImageRequest struct {
	ImagePath string `json:"image_path"`
}

// DeleteContentImage handles deletion of content images by path
func (h *APIHandler) DeleteContentImage(w http.ResponseWriter, r *http.Request) {
	h.Log().Infof("%s: Handling DeleteContentImage", h.Name())

	// Parse content ID from URL
	contentIDStr, err := h.Param(w, r, "content_id")
	if err != nil {
		h.Log().Infof("Failed to parse content_id: %v", err)
		h.Err(w, http.StatusBadRequest, "Invalid content ID", err)
		return
	}
	h.Log().Infof("Content ID: %s", contentIDStr)

	contentID, err := uuid.Parse(contentIDStr)
	if err != nil {
		h.Log().Infof("Failed to parse UUID: %v", err)
		h.Err(w, http.StatusBadRequest, "Invalid content ID format", err)
		return
	}

	// Parse image path from request body
	var req DeleteContentImageRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.Log().Infof("Failed to parse request body: %v", err)
		h.Err(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	h.Log().Infof("Image path: %s", req.ImagePath)

	// Delete image through service
	err = h.svc.DeleteContentImage(r.Context(), contentID, req.ImagePath)
	if err != nil {
		h.Log().Infof("Service delete failed: %v", err)
		h.Err(w, http.StatusInternalServerError, "Failed to delete image", err)
		return
	}

	h.Log().Infof("Image deleted successfully: %s", req.ImagePath)
	// Return success response
	msg := fmt.Sprintf("Content image deleted successfully")
	h.OK(w, msg, nil)
}