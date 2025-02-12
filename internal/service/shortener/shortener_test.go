package shortener_test

import (
	"context"
	"testing"

	"url-shortener/internal/config"
	"url-shortener/internal/service/shortener"
	"url-shortener/internal/storage/errdef"
	pb "url-shortener/protos/gen/go"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) GetAlias(originalURL string) (string, error) {
	args := m.Called(originalURL)
	return args.String(0), args.Error(1)
}

func (m *MockStorage) CheckAliasURLExists(alias string) (bool, error) {
	args := m.Called(alias)
	return args.Bool(0), args.Error(1)
}

func (m *MockStorage) SaveUrl(originalURL, alias string) error {
	args := m.Called(originalURL, alias)
	return args.Error(0)
}

func (m *MockStorage) GetUrl(alias string) (string, error) {
	args := m.Called(alias)
	return args.String(0), args.Error(1)
}

func TestCreateShortURL(t *testing.T) {
	mockStorage := new(MockStorage)
	cfg := config.Config{URLGenerator: config.URLGenerator{Length: 10, MaxAttempt: 5}}
	service := shortener.New(mockStorage, cfg)

	ctx := context.Background()
	originalURL := "https://example.com"
	alias := "abc123xyz_"

	// Тест: URL уже существует
	mockStorage.On("GetAlias", originalURL).Return(alias, nil).Once()
	resp, err := service.CreateShortURL(ctx, &pb.CreateShortURLRequest{OriginalUrl: originalURL})
	assert.NoError(t, err)
	assert.Equal(t, alias, resp.Alias)

	mockStorage.On("GetAlias", originalURL).Return("", errdef.ErrURLNotFound).Once()
	mockStorage.On("CheckAliasURLExists", mock.Anything).Return(false, nil).Once()
	mockStorage.On("SaveUrl", originalURL, mock.Anything).Return(nil).Once()

	resp, err = service.CreateShortURL(ctx, &pb.CreateShortURLRequest{OriginalUrl: originalURL})
	assert.NoError(t, err)
	assert.Len(t, resp.Alias, 10)

	mockStorage.AssertExpectations(t)
}

func TestGetOriginalURL(t *testing.T) {
	mockStorage := new(MockStorage)
	cfg := config.Config{URLGenerator: config.URLGenerator{Length: 10, MaxAttempt: 5}}
	service := shortener.New(mockStorage, cfg)

	ctx := context.Background()
	alias := "abc123xyz_"
	originalURL := "https://example.com"

	// Тест: Alias существует
	mockStorage.On("GetUrl", alias).Return(originalURL, nil).Once()
	resp, err := service.GetOriginalURL(ctx, &pb.GetOriginalURLRequest{Alias: alias})
	assert.NoError(t, err)
	assert.Equal(t, originalURL, resp.OriginalUrl)

	// Тест: Alias не найден
	mockStorage.On("GetUrl", alias).Return("", errdef.ErrURLNotFound).Once()
	resp, err = service.GetOriginalURL(ctx, &pb.GetOriginalURLRequest{Alias: alias})
	assert.Error(t, err)
	assert.Nil(t, resp)

	mockStorage.AssertExpectations(t)
}
