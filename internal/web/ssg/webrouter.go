package ssg

import (
	"github.com/adrianpk/clio/internal/am"
)

func NewWebRouter(handler *WebHandler, mw []am.Middleware, opts ...am.Option) *am.Router {
	core := am.NewWebRouter("app-ssg-web-router", opts...)
	core.SetMiddlewares(mw)

	// Content routes
	core.Get("/new-content", handler.NewContent)
	core.Post("/create-content", handler.CreateContent)
	core.Get("/edit-content", handler.EditContent)
	core.Post("/update-content", handler.UpdateContent)
	core.Get("/list-content", handler.ListContent)
	core.Get("/show-content", handler.ShowContent)
	core.Post("/delete-content", handler.DeleteContent)
	// Section routes
	core.Get("/new-section", handler.NewSection)
	core.Post("/create-section", handler.CreateSection)
	core.Get("/edit-section", handler.EditSection)
	core.Post("/update-section", handler.UpdateSection)
	core.Get("/list-sections", handler.ListSections)
	core.Get("/show-section", handler.ShowSection)
	core.Post("/delete-section", handler.DeleteSection)

	// Tag routes
	core.Get("/new-tag", handler.NewTag)
	core.Post("/create-tag", handler.CreateTag)
	core.Get("/edit-tag", handler.EditTag)
	core.Post("/update-tag", handler.UpdateTag)
	core.Get("/list-tags", handler.ListTags)
	core.Get("/show-tag", handler.ShowTag)
	core.Post("/delete-tag", handler.DeleteTag)

	// Layout routes
	core.Get("/new-layout", handler.NewLayout)
	core.Post("/create-layout", handler.CreateLayout)
	core.Get("/edit-layout", handler.EditLayout)
	core.Post("/update-layout", handler.UpdateLayout)
	core.Get("/list-layouts", handler.ListLayouts)
	core.Get("/show-layout", handler.ShowLayout)
	core.Post("/delete-layout", handler.DeleteLayout)

	// Param routes
	core.Get("/new-param", handler.NewParam)
	core.Post("/create-param", handler.CreateParam)
	core.Get("/edit-param", handler.EditParam)
	core.Post("/update-param", handler.UpdateParam)
	core.Get("/list-params", handler.ListParams)
	core.Get("/show-param", handler.ShowParam)
	core.Post("/delete-param", handler.DeleteParam)

	return core
}
