package services

import (
	"errors"
	"time"

	"starter-kit-fullstack-gonethttp-template/config"
	"starter-kit-fullstack-gonethttp-template/internal/models"
	"starter-kit-fullstack-gonethttp-template/internal/repository"
	"starter-kit-fullstack-gonethttp-template/pkg/utils"

	"github.com/google/uuid"
)

type authService struct {
	userRepo     repository.UserRepository
	tokenRepo    repository.TokenRepository
	tokenService *TokenService
	emailService EmailService
	cfg          *config.Config
}

func NewAuthService(uRepo repository.UserRepository, tRepo repository.TokenRepository, tService *TokenService, eService EmailService, cfg *config.Config) AuthService {
	return &authService{
		userRepo:     uRepo,
		tokenRepo:    tRepo,
		tokenService: tService,
		emailService: eService,
		cfg:          cfg,
	}
}

func (s *authService) Login(email, password string) (*models.User, map[string]interface{}, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil || !user.ComparePassword(password) {
		return nil, nil, errors.New("incorrect email or password")
	}

	tokens, err := s.tokenService.GenerateAuthTokens(user.ID)
	if err != nil {
		return nil, nil, err
	}

	return user, tokens, nil
}

func (s *authService) Register(req RegisterRequest) (*models.User, map[string]interface{}, error) {
	if exists, _ := s.userRepo.ExistsByEmail(req.Email); exists {
		return nil, nil, errors.New("email already taken")
	}

	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role:     "user", // Default role
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, nil, err
	}

	tokens, err := s.tokenService.GenerateAuthTokens(user.ID)
	if err != nil {
		return nil, nil, err
	}

	return user, tokens, nil
}

func (s *authService) Logout(refreshToken string) error {
	tokenDoc, err := s.tokenService.VerifyToken(refreshToken, models.TokenTypeRefresh)
	if err != nil {
		return errors.New("not found")
	}
	return s.tokenRepo.Delete(tokenDoc)
}

func (s *authService) RefreshAuth(refreshToken string) (map[string]interface{}, error) {
	tokenDoc, err := s.tokenService.VerifyToken(refreshToken, models.TokenTypeRefresh)
	if err != nil {
		return nil, errors.New("please authenticate")
	}

	userUUID, err := uuid.Parse(tokenDoc.UserID)
	if err != nil {
		return nil, errors.New("invalid user data")
	}

	// Clean up old refresh token
	s.tokenRepo.Delete(tokenDoc)

	return s.tokenService.GenerateAuthTokens(userUUID)
}

func (s *authService) ForgotPassword(email string) error {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		// Return nil to avoid email enumeration
		return nil
	}

	expires := time.Duration(s.cfg.JWT.ResetPasswordExpiration) * time.Minute
	tokenStr, expTime, err := utils.GenerateToken(user.ID, expires, models.TokenTypeResetPassword, s.cfg.JWT.Secret)
	if err != nil {
		return err
	}

	if err := s.tokenService.SaveToken(tokenStr, user.ID.String(), expTime, models.TokenTypeResetPassword); err != nil {
		return err
	}

	return s.emailService.SendResetPasswordEmail(user.Email, tokenStr)
}

func (s *authService) ResetPassword(tokenStr, newPassword string) error {
	tokenDoc, err := s.tokenService.VerifyToken(tokenStr, models.TokenTypeResetPassword)
	if err != nil {
		return errors.New("password reset failed")
	}

	userUUID, err := uuid.Parse(tokenDoc.UserID)
	if err != nil {
		return errors.New("invalid user data")
	}

	user, err := s.userRepo.FindByID(userUUID)
	if err != nil {
		return errors.New("user not found")
	}

	user.Password = newPassword
	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	// Invalidate all reset tokens for this user
	return s.tokenRepo.DeleteByUserIDAndType(user.ID.String(), models.TokenTypeResetPassword)
}

func (s *authService) VerifyEmail(tokenStr string) error {
	tokenDoc, err := s.tokenService.VerifyToken(tokenStr, models.TokenTypeVerifyEmail)
	if err != nil {
		return errors.New("email verification failed")
	}

	userUUID, err := uuid.Parse(tokenDoc.UserID)
	if err != nil {
		return errors.New("invalid user data")
	}

	user, err := s.userRepo.FindByID(userUUID)
	if err != nil {
		return errors.New("user not found")
	}

	user.IsEmailVerified = true
	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	// Invalidate verify tokens
	return s.tokenRepo.DeleteByUserIDAndType(user.ID.String(), models.TokenTypeVerifyEmail)
}