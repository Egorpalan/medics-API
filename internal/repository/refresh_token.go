package repository

import (
	"context"
	"database/sql"
	"github.com/Egorpalan/medods-api/internal/domain/model"
	"github.com/Egorpalan/medods-api/pkg/logger"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type RefreshTokenRepository interface {
	Save(ctx context.Context, token *model.RefreshToken) error
	FindByAccessJTI(ctx context.Context, accessJTI string) (*model.RefreshToken, error)
	MarkUsed(ctx context.Context, id int64) error
}

type RefreshTokenRepo struct {
	db *sqlx.DB
}

func NewRefreshTokenRepo(db *sqlx.DB) *RefreshTokenRepo {
	return &RefreshTokenRepo{db: db}
}

func (r *RefreshTokenRepo) Save(ctx context.Context, token *model.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (user_id, access_jti, refresh_token_hash, client_ip, expires_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.ExecContext(ctx, query, token.UserID, token.AccessJTI, token.RefreshTokenHash, token.ClientIP, token.ExpiresAt)
	if err != nil {
		logger.Log.Error("Failed to save refresh token", zap.Error(err))
	}
	return err
}

func (r *RefreshTokenRepo) FindByAccessJTI(ctx context.Context, accessJTI string) (*model.RefreshToken, error) {
	query := `
		SELECT id, user_id, access_jti, refresh_token_hash, client_ip, created_at, expires_at, used
		FROM refresh_tokens
		WHERE access_jti = $1
	`
	var token model.RefreshToken
	err := r.db.GetContext(ctx, &token, query, accessJTI)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.Log.Error("Failed to find refresh token", zap.Error(err))
		}
		return nil, err
	}
	return &token, nil
}

func (r *RefreshTokenRepo) MarkUsed(ctx context.Context, id int64) error {
	query := `UPDATE refresh_tokens SET used = TRUE WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		logger.Log.Error("Failed to mark refresh token as used", zap.Error(err))
	}
	return err
}
