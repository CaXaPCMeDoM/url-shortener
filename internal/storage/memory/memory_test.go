package memory_test

import (
	"testing"
	"url-shortener/internal/storage/errdef"
	"url-shortener/internal/storage/memory"

	"github.com/stretchr/testify/assert"
)

func TestStorage_SaveUrl(t *testing.T) {
	storage := memory.New()

	url := "https://example.com"
	alias := "abc123"

	// Сохранение нового URL
	err := storage.SaveUrl(url, alias)
	assert.NoError(t, err, "unexpected error while saving URL")

	// Повторное сохранение того же alias должно вернуть ошибку
	err = storage.SaveUrl("https://another.com", alias)
	assert.ErrorIs(t, err, errdef.ErrURLExists, "expected ErrURLExists when saving duplicate alias")
}

func TestStorage_GetUrl(t *testing.T) {
	storage := memory.New()

	url := "https://example.com"
	alias := "abc123"

	// Запрос несуществующего alias
	_, err := storage.GetUrl(alias)
	assert.ErrorIs(t, err, errdef.ErrURLNotFound, "expected ErrURLNotFound for non-existing alias")

	// Сохранение и получение
	_ = storage.SaveUrl(url, alias)
	retrievedURL, err := storage.GetUrl(alias)
	assert.NoError(t, err, "unexpected error while getting URL")
	assert.Equal(t, url, retrievedURL, "retrieved URL does not match expected")
}

func TestStorage_GetAlias(t *testing.T) {
	storage := memory.New()

	url := "https://example.com"
	alias := "abc123"

	// Запрос несуществующего URL
	_, err := storage.GetAlias(url)
	assert.ErrorIs(t, err, errdef.ErrURLNotFound, "expected ErrURLNotFound for non-existing URL")

	// Сохранение и получение
	_ = storage.SaveUrl(url, alias)
	retrievedAlias, err := storage.GetAlias(url)
	assert.NoError(t, err, "unexpected error while getting alias")
	assert.Equal(t, alias, retrievedAlias, "retrieved alias does not match expected")
}

func TestStorage_CheckAliasURLExists(t *testing.T) {
	storage := memory.New()

	alias := "abc123"

	// Проверка несуществующего alias
	exists, err := storage.CheckAliasURLExists(alias)
	assert.NoError(t, err, "unexpected error while checking alias existence")
	assert.False(t, exists, "alias should not exist initially")

	// Сохранение и проверка существования
	_ = storage.SaveUrl("https://example.com", alias)
	exists, err = storage.CheckAliasURLExists(alias)
	assert.NoError(t, err, "unexpected error while checking alias existence")
	assert.True(t, exists, "alias should exist after saving")
}
