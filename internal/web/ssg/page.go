package ssg

import (
	"net/http"

	"github.com/adrianpk/clio/internal/am"
)

// ParamPage extends am.Page to include a specific Param for templates.
type ParamPage struct {
	am.Page
	Param Param
}

// NewParamPage creates a new ParamPage.
func NewParamPage(r *http.Request, param Param) *ParamPage {
	page := am.NewPage(r, nil) // Pass nil for Data, as we'll use Param field
	return &ParamPage{
		Page:  *page,
		Param: param,
	}
}
