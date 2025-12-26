package services

import (
	"errors"

	"starter-kit-fullstack-gonethttp-template/internal/models"
	"starter-kit-fullstack-gonethttp-template/internal/repository"
	"starter-kit-fullstack-gonethttp-template/pkg/utils"

	"github.com/google/uuid"
)

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(req CreateUserRequest) (*models.User, error) {
	if exists, _ := s.repo.ExistsByEmail(req.Email); exists {
		return nil, errors.New("email already taken")
	}

	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) GetUserByID(id uuid.UUID) (*models.User, error) {
	return s.repo.FindByID(id)
}

func (s *userService) GetUsers(opts UserQueryOptions) (*utils.PaginationResult, error) {
	// 1. Prepare Pagination Scope
	if opts.Page < 1 {
		opts.Page = 1
	}
	if opts.Limit == 0 {
		opts.Limit = 10
	}
	
	paginationScope := &utils.PaginationScope{
		Page:  opts.Page,
		Limit: opts.Limit,
		Sort:  opts.SortBy,
	}

	// 2. Prepare Filters (Strict)
	filters := make(map[string]interface{})
	if opts.RoleFilter == "user" || opts.RoleFilter == "admin" {
		filters["role"] = opts.RoleFilter
	}

	// 3. Prepare Search Logic (Fuzzy)
	var searchFields []string
	if opts.Search != "" {
		if opts.SearchScope == "all" || opts.SearchScope == "" {
			searchFields = []string{"name", "email", "role", "id"}
		} else {
			// Validate specific scope
			validColumns := map[string]bool{"name": true, "email": true, "role": true, "id": true}
			if validColumns[opts.SearchScope] {
				searchFields = []string{opts.SearchScope}
			} else {
				// Default fallback
				searchFields = []string{"name"}
			}
		}
	}

	// 4. Query Repository
	users, totalRows, err := s.repo.FindAll(filters, opts.Search, searchFields, paginationScope)
	if err != nil {
		return nil, err
	}

	result := utils.GetPaginationResult(totalRows, opts.Page, opts.Limit, users)
	return &result, nil
}

func (s *userService) UpdateUser(id uuid.UUID, req UpdateUserRequest) (*models.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if req.Email != "" && req.Email != user.Email {
		if exists, _ := s.repo.ExistsByEmail(req.Email); exists {
			return nil, errors.New("email already taken")
		}
		user.Email = req.Email
	}
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Password != "" {
		user.Password = req.Password
	}

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) DeleteUser(id uuid.UUID) error {
	if _, err := s.repo.FindByID(id); err != nil {
		return errors.New("user not found")
	}
	return s.repo.Delete(id)
}