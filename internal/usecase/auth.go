package usecase

import (
	"context"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"time"

	"github.com/Egorpalan/medods-api/internal/domain/model"
	"github.com/Egorpalan/medods-api/internal/repository"
	"github.com/Egorpalan/medods-api/pkg/logger"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthUseСase interface {
	GenerateTokenPair(ctx context.Context, userID, clientIP string) (*model.TokenPair, error)
	RefreshTokenPair(ctx context.Context, accessToken, refreshToken, clientIP string) (*model.TokenPair, error)
}

type AuthUsecase struct {
	repo      repository.RefreshTokenRepository
	jwtSecret []byte
}

func NewAuthUsecase(repo repository.RefreshTokenRepository, jwtSecret []byte) *AuthUsecase {
	return &AuthUsecase{repo: repo, jwtSecret: jwtSecret}
}

type CustomClaims struct {
	UserID   string `json:"user_id"`
	ClientIP string `json:"client_ip"`
	jwt.RegisteredClaims
}

func (u *AuthUsecase) GenerateTokenPair(ctx context.Context, userID, clientIP string) (*model.TokenPair, error) {
	accessJTI, err := generateJTI()
	if err != nil {
		logger.Log.Error("Failed to generate JTI", zap.Error(err))
		return nil, err
	}
	accessToken, err := u.generateJWT(userID, clientIP, accessJTI, 15*time.Minute)
	if err != nil {
		logger.Log.Error("Failed to generate access token", zap.Error(err))
		return nil, err
	}

	rawRefresh, err := generateRandomString(32)
	if err != nil {
		logger.Log.Error("Failed to generate refresh token", zap.Error(err))
		return nil, err
	}
	refreshToken := base64.StdEncoding.EncodeToString([]byte(rawRefresh))

	refreshHash, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		logger.Log.Error("Failed to hash refresh token", zap.Error(err))
		return nil, err
	}

	refreshModel := &model.RefreshToken{
		UserID:           userID,
		AccessJTI:        accessJTI,
		RefreshTokenHash: string(refreshHash),
		ClientIP:         clientIP,
		ExpiresAt:        time.Now().Add(7 * 24 * time.Hour),
	}
	if err := u.repo.Save(ctx, refreshModel); err != nil {
		return nil, err
	}

	return &model.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (u *AuthUsecase) RefreshTokenPair(ctx context.Context, accessToken, refreshToken, clientIP string) (*model.TokenPair, error) {
	token, claims, err := u.parseJWT(accessToken)
	if err != nil || !token.Valid {
		logger.Log.Warn("Invalid access token", zap.Error(err))
		return nil, errors.New("invalid access token")
	}

	refreshModel, err := u.repo.FindByAccessJTI(ctx, claims.ID)
	if err != nil {
		return nil, errors.New("refresh token not found")
	}
	if refreshModel.Used {
		return nil, errors.New("refresh token already used")
	}
	if time.Now().After(refreshModel.ExpiresAt) {
		return nil, errors.New("refresh token expired")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(refreshModel.RefreshTokenHash), []byte(refreshToken)); err != nil {
		return nil, errors.New("invalid refresh token")
	}

	if refreshModel.ClientIP != clientIP {
		logger.Log.Warn("Client IP changed, sending warning email", zap.String("user_id", refreshModel.UserID))
	}

	if err := u.repo.MarkUsed(ctx, refreshModel.ID); err != nil {
		return nil, err
	}

	return u.GenerateTokenPair(ctx, refreshModel.UserID, clientIP)
}

func (u *AuthUsecase) generateJWT(userID, clientIP, jti string, ttl time.Duration) (string, error) {
	claims := CustomClaims{
		UserID:   userID,
		ClientIP: clientIP,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        jti,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString(u.jwtSecret)
}

func (u *AuthUsecase) parseJWT(tokenStr string) (*jwt.Token, *jwt.RegisteredClaims, error) {
	claims := &CustomClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return u.jwtSecret, nil
	})
	return token, &claims.RegisteredClaims, err
}

func generateJTI() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	hash := sha512.Sum512(b)
	return base64.URLEncoding.EncodeToString(hash[:]), nil
}

func generateRandomString(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
