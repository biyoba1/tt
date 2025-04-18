package handler

import (
	"net/http"
	"valera/internal/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/auth/sign-up", h.signUp)
	mux.HandleFunc("/auth/refresh", h.refreshToken)

	return mux
}
