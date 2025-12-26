package api

import (
	"encoding/json"
	"net/http"

	"starter-kit-fullstack-gonethttp-template/internal/services"
	"starter-kit-fullstack-gonethttp-template/pkg/response"
	"starter-kit-fullstack-gonethttp-template/pkg/utils"
)

type AuthHandler struct {
	service services.AuthService
}

func NewAuthHandler(service services.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user account and return tokens
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body services.RegisterRequest true "Register Request"
// @Success 201 {object} response.APIResponse{data=map[string]interface{}}
// @Failure 400 {object} response.APIResponse
// @Router /v1/auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req services.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if errs := utils.ValidateStruct(req); errs != nil {
		response.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 400, "message": errs})
		return
	}

	user, tokens, err := h.service.Register(req)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, map[string]interface{}{
		"user":   user,
		"tokens": tokens,
	})
}

// Login godoc
// @Summary Login user
// @Description Authenticate user with email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body object{email=string,password=string} true "Login Request"
// @Success 200 {object} response.APIResponse{data=map[string]interface{}}
// @Failure 401 {object} response.APIResponse
// @Router /v1/auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if errs := utils.ValidateStruct(req); errs != nil {
		response.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 400, "message": errs})
		return
	}

	user, tokens, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	response.Success(w, http.StatusOK, map[string]interface{}{
		"user":   user,
		"tokens": tokens,
	})
}

// Logout godoc
// @Summary Logout user
// @Description Logout by invalidating the refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body object{refreshToken=string} true "Logout Request"
// @Success 204 "No Content"
// @Failure 400 {object} response.APIResponse
// @Router /v1/auth/logout [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refreshToken"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	if req.RefreshToken == "" {
		response.Error(w, http.StatusBadRequest, "refreshToken is required")
		return
	}

	h.service.Logout(req.RefreshToken)
	w.WriteHeader(http.StatusNoContent)
}

// RefreshTokens godoc
// @Summary Refresh auth tokens
// @Description Get new access and refresh tokens using a valid refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body object{refreshToken=string} true "Refresh Request"
// @Success 200 {object} response.APIResponse{data=map[string]interface{}}
// @Failure 401 {object} response.APIResponse
// @Router /v1/auth/refresh-tokens [post]
func (h *AuthHandler) RefreshTokens(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refreshToken"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	if req.RefreshToken == "" {
		response.Error(w, http.StatusBadRequest, "refreshToken is required")
		return
	}

	tokens, err := h.service.RefreshAuth(req.RefreshToken)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	response.Success(w, http.StatusOK, tokens)
}

// ForgotPassword godoc
// @Summary Request password reset
// @Description Send an email with a password reset link
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body object{email=string} true "Forgot Password Request"
// @Success 204 "No Content"
// @Failure 400 {object} response.APIResponse
// @Router /v1/auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email" validate:"required,email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request")
		return
	}

	if errs := utils.ValidateStruct(req); errs != nil {
		response.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 400, "message": errs})
		return
	}

	h.service.ForgotPassword(req.Email)
	w.WriteHeader(http.StatusNoContent)
}

// ResetPassword godoc
// @Summary Reset password
// @Description Reset password using a valid token
// @Tags Auth
// @Accept json
// @Produce json
// @Param token query string true "Reset Token"
// @Param request body object{password=string} true "Reset Password Request"
// @Success 204 "No Content"
// @Failure 400 {object} response.APIResponse
// @Router /v1/auth/reset-password [post]
func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	var req struct {
		Password string `json:"password" validate:"required,min=8"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	if token == "" || req.Password == "" {
		response.Error(w, http.StatusBadRequest, "Token and Password are required")
		return
	}

	if err := h.service.ResetPassword(token, req.Password); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}