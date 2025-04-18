package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"valera/models"
)

func (h *Handler) signUp(c *gin.Context) {
	var input models.LoginReq
	err := c.BindJSON(&input)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if input.Guid == "" {
		errorResponse(c, http.StatusBadRequest, "guid are required")
		return
	}
	clientIP := c.ClientIP()
	tokens, err := h.services.Login(input.Guid, clientIP)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
	})
}

func (h *Handler) refreshToken(c *gin.Context) {
	var input models.RefreshReq
	c.BindJSON(&input)
	if input.RefreshToken == "" {
		errorResponse(c, http.StatusBadRequest, "refresh token are required")
		return
	}

	tokens, err := h.services.Refresh(c.ClientIP(), input.RefreshToken)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
	})
}
