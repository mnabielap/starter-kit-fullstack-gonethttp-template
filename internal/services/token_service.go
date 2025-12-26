package services

import (
	"time"

	"starter-kit-fullstack-gonethttp-template/config"
	"starter-kit-fullstack-gonethttp-template/internal/models"
	"starter-kit-fullstack-gonethttp-template/internal/repository"
	"starter-kit-fullstack-gonethttp-template/pkg/utils"

	"github.com/google/uuid"
)

type TokenService struct {
	repo repository.TokenRepository
	cfg  *config.Config
}

func NewTokenService(repo repository.TokenRepository, cfg *config.Config) *TokenService {
	return &TokenService{repo: repo, cfg: cfg}
}

// GenerateAuthTokens creates access and refresh tokens
func (s *TokenService) GenerateAuthTokens(userID uuid.UUID) (map[string]interface{}, error) {
	// Access Token
	accessDur := time.Duration(s.cfg.JWT.AccessExpirationMinutes) * time.Minute
	accessToken, accessExp, err := utils.GenerateToken(userID, accessDur, "access", s.cfg.JWT.Secret)
	if err != nil {
		return nil, err
	}

	// Refresh Token
	refreshDur := time.Duration(s.cfg.JWT.RefreshExpirationDays) * 24 * time.Hour
	refreshToken, refreshExp, err := utils.GenerateToken(userID, refreshDur, models.TokenTypeRefresh, s.cfg.JWT.Secret)
	if err != nil {
		return nil, err
	}

	// Save Refresh Token to DB
	err = s.SaveToken(refreshToken, userID.String(), refreshExp, models.TokenTypeRefresh)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"access": map[string]interface{}{
			"token":   accessToken,
			"expires": accessExp,
		},
		"refresh": map[string]interface{}{
			"token":   refreshToken,
			"expires": refreshExp,
		},
	}, nil
}

func (s *TokenService) SaveToken(token, userID string, expires time.Time, tokenType string) error {
	tokenModel := &models.Token{
		Token:   token,
		UserID:  userID,
		Expires: expires,
		Type:    tokenType,
	}
	return s.repo.Create(tokenModel)
}

func (s *TokenService) VerifyToken(token string, tokenType string) (*models.Token, error) {
	// Verify signature
	if _, err := utils.ValidateToken(token, s.cfg.JWT.Secret); err != nil {
		return nil, err
	}
	// Verify existence in DB
	return s.repo.FindByToken(token, tokenType)
}