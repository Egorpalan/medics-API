package model

import (
	"time"
)

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type RefreshToken struct {
	ID               int64
	UserID           string
	AccessJTI        string
	RefreshTokenHash string
	ClientIP         string
	CreatedAt        time.Time
	ExpiresAt        time.Time
	Used             bool
}
