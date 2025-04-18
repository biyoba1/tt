package service

import (
	"valera/internal/repository"
	"valera/internal/service/auth"
	"valera/models"
)

type AuthService interface {
	Login(guid, ip string) (*models.AuthResponse, error)
	Refresh(ip, refreshToken string) (*models.AuthResponse, error)
}

type Service struct {
	AuthService
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		AuthService: auth.NewAuthService(repos),
	}
}
