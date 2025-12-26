package services

import (
	"starter-kit-fullstack-gonethttp-template/internal/models"
	"starter-kit-fullstack-gonethttp-template/pkg/utils"

	"github.com/google/uuid"
)

// DTOs
type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type CreateUserRequest struct {
	RegisterRequest
	Role string `json:"role" validate:"required,oneof=user admin"`
}

type UpdateUserRequest struct {
	Name     string `json:"name" validate:"omitempty"`
	Email    string `json:"email" validate:"omitempty,email"`
	Password string `json:"password" validate:"omitempty,min=8"`
}

type UserQueryOptions struct {
	Page         int
	Limit        int // -1 for all
	SortBy       string
	Search       string
	SearchScope  string // "all", "name", "email", etc.
	RoleFilter   string
}

// Interfaces

type AuthService interface {
	Login(email, password string) (*models.User, map[string]interface{}, error)
	Register(req RegisterRequest) (*models.User, map[string]interface{}, error)
	RefreshAuth(refreshToken string) (map[string]interface{}, error)
	Logout(refreshToken string) error
	
	ForgotPassword(email string) error
	ResetPassword(token, newPassword string) error
	
	VerifyEmail(token string) error
}

type UserService interface {
	CreateUser(req CreateUserRequest) (*models.User, error)
	GetUserByID(id uuid.UUID) (*models.User, error)
	GetUsers(options UserQueryOptions) (*utils.PaginationResult, error)
	UpdateUser(id uuid.UUID, req UpdateUserRequest) (*models.User, error)
	DeleteUser(id uuid.UUID) error
}

type EmailService interface {
	SendEmail(to, subject, body string) error
	SendResetPasswordEmail(to, token string) error
	SendVerificationEmail(to, token string) error
}