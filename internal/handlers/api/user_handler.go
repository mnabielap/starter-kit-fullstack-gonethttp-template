package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"starter-kit-fullstack-gonethttp-template/internal/services"
	"starter-kit-fullstack-gonethttp-template/pkg/response"
	"starter-kit-fullstack-gonethttp-template/pkg/utils"

	"github.com/google/uuid"
)

type UserHandler struct {
	service services.UserService
}

func NewUserHandler(service services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// CreateUser godoc
// @Summary Create a new user (Admin)
// @Description Create a user with specific role
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body services.CreateUserRequest true "Create User Request"
// @Success 201 {object} models.User
// @Failure 400 {object} response.APIResponse
// @Router /v1/users [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req services.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if errs := utils.ValidateStruct(req); errs != nil {
		response.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 400, "message": errs})
		return
	}

	user, err := h.service.CreateUser(req)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, user)
}

// GetUsers godoc
// @Summary Get all users
// @Description Get users with pagination, sorting, and filtering
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param limit query int false "Items per page (-1 for all)"
// @Param sortBy query string false "Sort by field (e.g. name:asc)"
// @Param search query string false "Search term"
// @Param scope query string false "Search scope (name, email, role, id)"
// @Param role query string false "Filter by role"
// @Success 200 {object} utils.PaginationResult
// @Router /v1/users [get]
func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	page, _ := strconv.Atoi(query.Get("page"))
	limit, _ := strconv.Atoi(query.Get("limit"))
	
	// Default limit if not provided
	if query.Get("limit") == "" {
		limit = 10
	}

	scope := query.Get("scope")
	if scope == "" {
		scope = "all"
	}

	opts := services.UserQueryOptions{
		Page:        page,
		Limit:       limit,
		SortBy:      query.Get("sortBy"),
		Search:      query.Get("name"),
		SearchScope: scope,
		RoleFilter:  query.Get("role"),
	}
	
	if s := query.Get("search"); s != "" {
		opts.Search = s
	}

	result, err := h.service.GetUsers(opts)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, http.StatusOK, result)
}

// GetUser godoc
// @Summary Get user by ID
// @Description Get a specific user details
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} models.User
// @Failure 404 {object} response.APIResponse
// @Router /v1/users/{id} [get]
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid User ID")
		return
	}

	user, err := h.service.GetUserByID(id)
	if err != nil {
		response.Error(w, http.StatusNotFound, "User not found")
		return
	}

	response.Success(w, http.StatusOK, user)
}

// UpdateUser godoc
// @Summary Update user
// @Description Update user details
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param request body services.UpdateUserRequest true "Update Request"
// @Success 200 {object} models.User
// @Failure 400 {object} response.APIResponse
// @Router /v1/users/{id} [patch]
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid User ID")
		return
	}

	var req services.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := h.service.UpdateUser(id, req)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(w, http.StatusOK, user)
}

// DeleteUser godoc
// @Summary Delete user
// @Description Permanently delete a user
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 204 "No Content"
// @Failure 404 {object} response.APIResponse
// @Router /v1/users/{id} [delete]
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid User ID")
		return
	}

	if err := h.service.DeleteUser(id); err != nil {
		response.Error(w, http.StatusNotFound, "User not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}