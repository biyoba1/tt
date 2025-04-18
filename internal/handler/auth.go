package handler

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"valera/models"
)

func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {
	var input models.LoginReq
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		errorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if input.Guid == "" {
		errorResponse(w, http.StatusBadRequest, "guid is required")
		return
	}

	clientIP := r.RemoteAddr
	tokens, err := h.services.Login(input.Guid, clientIP)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
	})
}

func (h *Handler) refreshToken(w http.ResponseWriter, r *http.Request) {
	var input models.RefreshReq
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		errorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if input.RefreshToken == "" {
		errorResponse(w, http.StatusBadRequest, "refresh token is required")
		return
	}
	refreshTokenBytes, err := base64.StdEncoding.DecodeString(input.RefreshToken)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "invalid refresh token format")
		return
	}

	clientIP := r.RemoteAddr
	tokens, err := h.services.Refresh(clientIP, string(refreshTokenBytes))
	if err != nil {
		errorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
	})
}
