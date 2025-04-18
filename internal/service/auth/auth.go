package auth

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mail.v2"
	"log"
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

func hashWithSHA256(data string) []byte {
	hash := sha256.Sum256([]byte(data))
	return hash[:]
}

func GenerateTokens(ip, guid string) (string, string, error) {
	aToken := jwt.NewWithClaims(jwt.SigningMethodHS512, tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 15).Unix(),
		},
		Ip:   ip,
		Guid: guid,
	})

	rToken := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 15).Unix(),
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

func (s *AuthService) Login(guid, ip string) (*models.AuthResponse, error) {
	accessToken, refreshToken, err := GenerateTokens(ip, guid)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	hashedRefreshToken := hashWithSHA256(refreshToken)
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
		RefreshToken: base64.StdEncoding.EncodeToString([]byte(refreshToken)), // Передаем в Base64
	}, nil
}

func (s *AuthService) Refresh(ip, refreshTokenBase64 string) (*models.AuthResponse, error) {
	refreshTokenBytes, err := base64.StdEncoding.DecodeString(refreshTokenBase64)
	if err != nil {
		return nil, err
	}

	refreshToken := string(refreshTokenBytes)
	hashedRefreshToken := hashWithSHA256(refreshToken)
	parsedToken, err := jwt.ParseWithClaims(refreshToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(signingKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := parsedToken.Claims.(*tokenClaims); ok && parsedToken.Valid {
		if claims.Ip != ip {
			//Отправляем пользователю варнинг на почту
			//В предыдущих местах работы я это делал через горутины чтобы не ждать, но тут просто замокал
			m := mail.NewMessage()
			m.SetHeader("From", "pushkin@mail.ru")
			m.SetHeader("To", "kuda-to")
			m.SetHeader("Subject", "Your ip address changed")
			a := mail.NewDialer("smtp.mail.ru", 465, "pushkin@mail.ru", "12345")

			err := a.DialAndSend(m)
			if err != nil {
				log.Println("Ошибка отправки сообщения на почту:", err)
			}
		}

		storedHash, err := s.repo.CheckRefreshToken(claims.Guid)
		if err != nil {
			return nil, err
		}

		err = bcrypt.CompareHashAndPassword(storedHash, hashedRefreshToken)
		if err != nil {
			return nil, fmt.Errorf("invalid refresh token")
		}

		newAccessToken, newRefreshToken, err := GenerateTokens(ip, claims.Guid)
		if err != nil {
			return nil, err
		}

		newHashedRefreshToken := hashWithSHA256(newRefreshToken)
		newHashedToken, err := bcrypt.GenerateFromPassword(newHashedRefreshToken, bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}

		err = s.repo.SaveRefreshToken(claims.Guid, newHashedToken)
		if err != nil {
			return nil, err
		}

		return &models.AuthResponse{
			AccessToken:  newAccessToken,
			RefreshToken: base64.StdEncoding.EncodeToString([]byte(newRefreshToken)),
		}, nil
	}

	return nil, fmt.Errorf("invalid refresh token")
}
