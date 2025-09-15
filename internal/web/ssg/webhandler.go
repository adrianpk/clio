package ssg

import (
	"github.com/adrianpk/clio/internal/am"
)

const (
	// WIP: This will be obtained from configuration.
	defaultAPIBaseURL = "http://localhost:8081/api/v1"
)

const (
	ssgFeat = "ssg"
	ssgPath = "/ssg"
)

type WebHandler struct {
	*am.WebHandler
	apiClient *am.APIClient
}

func NewWebHandler(tm *am.TemplateManager, flash *am.FlashManager, opts ...am.Option) *WebHandler {
	handler := am.NewWebHandler(tm, flash, opts...)
	apiClient := am.NewAPIClient("web-api-client", func() string { return "" }, defaultAPIBaseURL, opts...)
	return &WebHandler{
		WebHandler: handler,
		apiClient:  apiClient,
	}
}
