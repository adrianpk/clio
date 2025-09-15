package auth

import (
	"github.com/adrianpk/clio/internal/am"
)

func NewAPIRouter(handler *APIHandler, mw []am.Middleware, opts ...am.Option) *am.Router {
	core := am.NewAPIRouter("api-router", opts...)
	core.SetMiddlewares(mw)

	// User API routes
	core.Get("/users", handler.GetAllUsers)
	core.Get("/users/{id}", handler.GetUser)
	core.Post("/users", handler.CreateUser)
	core.Put("/users/{id}", handler.UpdateUser)
	core.Delete("/users/{id}", handler.DeleteUser)

	return core
}
