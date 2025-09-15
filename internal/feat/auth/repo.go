package auth

import (
	"context"

	"github.com/adrianpk/clio/internal/am"

	"github.com/google/uuid"
)

type Repo interface {
	am.Repo

	// SECTION: User-related methods

	GetUserByUsername(ctx context.Context, username string) (User, error)
	GetUsers(ctx context.Context) ([]User, error)
	GetUser(ctx context.Context, id uuid.UUID) (User, error)
	CreateUser(ctx context.Context, user *User) error
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
}
