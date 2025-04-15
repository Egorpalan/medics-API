package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Egorpalan/medods-api/internal/usecase"
	"github.com/Egorpalan/medods-api/pkg/logger"
	"go.uber.org/zap"
)

type AuthHandler struct {
	authUsecase *usecase.AuthUsecase
}

func NewAuthHandler(authUsecase *usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUsecase: authUsecase}
}

func (h *AuthHandler) GenerateTokenPair(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "user_id required", http.StatusBadRequest)
		return
	}
	clientIP := getClientIP(r)
	ctx := r.Context()

	tokens, err := h.authUsecase.GenerateTokenPair(ctx, userID, clientIP)
	if err != nil {
		logger.Log.Error("Failed to generate token pair", zap.Error(err))
		http.Error(w, "failed to generate token pair", http.StatusInternalServerError)
		return
	}
	writeJSON(w, tokens)
}

type RefreshRequest struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (h *AuthHandler) RefreshTokenPair(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	clientIP := getClientIP(r)
	ctx := r.Context()

	tokens, err := h.authUsecase.RefreshTokenPair(ctx, req.AccessToken, req.RefreshToken, clientIP)
	if err != nil {
		logger.Log.Warn("Failed to refresh token pair", zap.Error(err))
		http.Error(w, "failed to refresh token pair: "+err.Error(), http.StatusUnauthorized)
		return
	}
	writeJSON(w, tokens)
}

func getClientIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.RemoteAddr
	}
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		ip = ip[:idx]
	}
	return ip
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}
