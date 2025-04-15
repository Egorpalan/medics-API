package model

import (
	"time"
)

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type RefreshToken struct {
	ID               int64     `db:"id"`
	UserID           string    `db:"user_id"`
	AccessJTI        string    `db:"access_jti"`
	RefreshTokenHash string    `db:"refresh_token_hash"`
	ClientIP         string    `db:"client_ip"`
	CreatedAt        time.Time `db:"created_at"`
	ExpiresAt        time.Time `db:"expires_at"`
	Used             bool      `db:"used"`
}
