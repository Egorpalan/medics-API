package usecase

import (
	"context"
	"github.com/Egorpalan/medods-api/pkg/logger"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"time"

	"github.com/Egorpalan/medods-api/internal/domain/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRefreshTokenRepo struct {
	mock.Mock
}

func (m *MockRefreshTokenRepo) Save(ctx context.Context, token *model.RefreshToken) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}
func (m *MockRefreshTokenRepo) FindByAccessJTI(ctx context.Context, accessJTI string) (*model.RefreshToken, error) {
	args := m.Called(ctx, accessJTI)
	return args.Get(0).(*model.RefreshToken), args.Error(1)
}
func (m *MockRefreshTokenRepo) MarkUsed(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestGenerateTokenPair(t *testing.T) {
	logger.InitLogger()
	mockRepo := new(MockRefreshTokenRepo)
	uc := NewAuthUsecase(mockRepo, []byte("testsecret"))

	mockRepo.On("Save", mock.Anything, mock.AnythingOfType("*model.RefreshToken")).Return(nil)

	tokens, err := uc.GenerateTokenPair(context.Background(), "user-guid", "127.0.0.1")
	assert.NoError(t, err)
	assert.NotEmpty(t, tokens.AccessToken)
	assert.NotEmpty(t, tokens.RefreshToken)
	mockRepo.AssertExpectations(t)
}

func TestRefreshTokenPair_IPChange(t *testing.T) {
	mockRepo := new(MockRefreshTokenRepo)
	uc := NewAuthUsecase(mockRepo, []byte("testsecret"))

	mockRepo.On("Save", mock.Anything, mock.AnythingOfType("*model.RefreshToken")).Return(nil)
	tokens, err := uc.GenerateTokenPair(context.Background(), "user-guid", "127.0.0.1")
	assert.NoError(t, err)

	refreshModel := &model.RefreshToken{
		ID:               1,
		UserID:           "user-guid",
		AccessJTI:        "",
		RefreshTokenHash: "",
		ClientIP:         "127.0.0.1",
		ExpiresAt:        time.Now().Add(1 * time.Hour),
		Used:             false,
	}
	_, claims, _ := uc.parseJWT(tokens.AccessToken)
	refreshModel.AccessJTI = claims.ID
	refreshHash, _ := bcrypt.GenerateFromPassword([]byte(tokens.RefreshToken), bcrypt.DefaultCost)
	refreshModel.RefreshTokenHash = string(refreshHash)

	mockRepo.On("FindByAccessJTI", mock.Anything, claims.ID).Return(refreshModel, nil)
	mockRepo.On("MarkUsed", mock.Anything, int64(1)).Return(nil)
	mockRepo.On("Save", mock.Anything, mock.AnythingOfType("*model.RefreshToken")).Return(nil)

	newTokens, err := uc.RefreshTokenPair(context.Background(), tokens.AccessToken, tokens.RefreshToken, "10.0.0.1")
	assert.NoError(t, err)
	assert.NotEmpty(t, newTokens.AccessToken)
	assert.NotEmpty(t, newTokens.RefreshToken)
	mockRepo.AssertExpectations(t)
}

func TestRefreshTokenPair_UsedToken(t *testing.T) {
	logger.InitLogger()
	mockRepo := new(MockRefreshTokenRepo)
	uc := NewAuthUsecase(mockRepo, []byte("testsecret"))

	mockRepo.On("Save", mock.Anything, mock.AnythingOfType("*model.RefreshToken")).Return(nil)
	tokens, err := uc.GenerateTokenPair(context.Background(), "user-guid", "127.0.0.1")
	assert.NoError(t, err)

	_, claims, _ := uc.parseJWT(tokens.AccessToken)
	refreshModel := &model.RefreshToken{
		ID:               1,
		UserID:           "user-guid",
		AccessJTI:        claims.ID,
		RefreshTokenHash: "",
		ClientIP:         "127.0.0.1",
		ExpiresAt:        time.Now().Add(1 * time.Hour),
		Used:             true,
	}
	refreshHash, _ := bcrypt.GenerateFromPassword([]byte(tokens.RefreshToken), bcrypt.DefaultCost)
	refreshModel.RefreshTokenHash = string(refreshHash)

	mockRepo.On("FindByAccessJTI", mock.Anything, claims.ID).Return(refreshModel, nil)

	_, err = uc.RefreshTokenPair(context.Background(), tokens.AccessToken, tokens.RefreshToken, "127.0.0.1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already used")
}
