package auth

import (
	"crypto/sha512"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
	"valera/internal/repository"
	"valera/models"
)

var signingKey = os.Getenv("signing_key")

type AuthService struct {
	repo repository.AuthRepository
}

type tokenClaims struct {
	jwt.StandardClaims
	Ip   string `json:"ip"`
	Guid string `json:"guid"`
}

func NewAuthService(repo repository.AuthRepository) *AuthService {
	return &AuthService{repo: repo}
}

func hashWithSHA512(data string) []byte {
	hash := sha512.Sum512([]byte(data))
	return hash[:]
}

func (s *AuthService) Login(guid, ip string) (*models.AuthResponse, error) {
	accessToken, refreshToken, err := GenerateTokens(ip, guid)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	hashedRefreshToken := hashWithSHA512(refreshToken)
	hashedToken, err := bcrypt.GenerateFromPassword(hashedRefreshToken, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	err = s.repo.SaveRefreshToken(guid, hashedToken)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func GenerateTokens(ip, guid string) (string, string, error) {
	aToken := jwt.NewWithClaims(jwt.SigningMethodHS512, tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(3600 * 24 * 15).Unix(),
		},
		Ip:   ip,
		Guid: guid,
	})

	rToken := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(3600 * 24 * 15).Unix(),
		},
		Ip:   ip,
		Guid: guid,
	})

	accessToken, err := aToken.SignedString([]byte(signingKey))
	if err != nil {
		return "", "", err
	}

	refreshToken, err := rToken.SignedString([]byte(signingKey))
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) Refresh(ip, refreshToken string) (*models.AuthResponse, error) {
	hashedRefreshToken := hashWithSHA512(refreshToken)

	storedHash, guid, err := s.repo.CheckRefreshToken(hashedRefreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve refresh token: %w", err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(storedHash), hashedRefreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	newAccessToken, newRefreshToken, err := GenerateTokens(ip, guid)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new tokens: %w", err)
	}

	newHashedRefreshToken := hashWithSHA512(newRefreshToken)
	newHashedToken, err := bcrypt.GenerateFromPassword(newHashedRefreshToken, bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash new refresh token with bcrypt: %w", err)
	}

	err = s.repo.UpdateRefreshToken(guid, newHashedToken)
	if err != nil {
		return nil, fmt.Errorf("failed to update refresh token: %w", err)
	}

	return &models.AuthResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
