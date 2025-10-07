package auth

import (
	"context"

	"github.com/adrianpk/clio/internal/am"
	"github.com/google/uuid"
)

type Service interface {
	// User-related methods
	GetUsers(ctx context.Context) ([]User, error)
	GetUser(ctx context.Context, id uuid.UUID) (User, error)
	CreateUser(ctx context.Context, user *User) error
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	// NOTE: Check this, am.UserService and am.SessionStore implementations
	GetUserByID(ctx context.Context, userID uuid.UUID) (*am.UserCtxData, error)
}

type BaseService struct {
	*am.Service
	repo Repo
}

func NewService(repo Repo, opts ...am.Option) *BaseService {
	return &BaseService{
		Service: am.NewService("auth-svc", opts...),
		repo:    repo,
	}
}

// NewServiceWithParams creates an Auth Service with XParams.
func NewServiceWithParams(repo Repo, params am.XParams) *BaseService {
	return &BaseService{
		Service: am.NewServiceWithParams("auth-svc", params),
		repo:    repo,
	}
}

func (svc *BaseService) GetUserByID(ctx context.Context, userID uuid.UUID) (*am.UserCtxData, error) {
	user, err := svc.repo.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	userCtxData := &am.UserCtxData{
		ID: user.ID,
	}
	return userCtxData, nil
}

func (svc *BaseService) GetUsers(ctx context.Context) ([]User, error) {
	return svc.repo.GetUsers(ctx)
}

func (svc *BaseService) GetUser(ctx context.Context, id uuid.UUID) (User, error) {
	return svc.repo.GetUser(ctx, id)
}

func (svc *BaseService) CreateUser(ctx context.Context, user *User) error {
	user.GenCreateValues()
	return svc.repo.CreateUser(ctx, user)
}

func (svc *BaseService) UpdateUser(ctx context.Context, user *User) error {
	user.GenUpdateValues()
	return svc.repo.UpdateUser(ctx, user)
}

func (svc *BaseService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return svc.repo.DeleteUser(ctx, id)
}
