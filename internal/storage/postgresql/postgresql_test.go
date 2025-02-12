package postgresql

import (
	"database/sql"
	"github.com/lib/pq"
	"testing"
	"url-shortener/internal/storage/errdef"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestStorage_SaveUrl(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := &Storage{db: db}

	url := "https://example.com"
	alias := "abc123"

	mock.ExpectPrepare("INSERT INTO url").ExpectExec().
		WithArgs(alias, url).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = storage.SaveUrl(url, alias)
	assert.NoError(t, err)

	mock.ExpectPrepare("INSERT INTO url").ExpectExec().
		WithArgs(alias, "https://another.com").
		WillReturnError(&pq.Error{Code: UniqueViolation})

	err = storage.SaveUrl("https://another.com", alias)
	assert.ErrorIs(t, err, errdef.ErrURLExists)
}

func TestStorage_GetUrl(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := &Storage{db: db}

	alias := "abc123"
	expectedURL := "https://example.com"

	mock.ExpectPrepare("SELECT url FROM url WHERE alias").
		ExpectQuery().
		WithArgs(alias).
		WillReturnRows(sqlmock.NewRows([]string{"url"}).AddRow(expectedURL))

	url, err := storage.GetUrl(alias)
	assert.NoError(t, err)
	assert.Equal(t, expectedURL, url)

	mock.ExpectPrepare("SELECT url FROM url WHERE alias").
		ExpectQuery().
		WithArgs("unknown").
		WillReturnError(sql.ErrNoRows)

	_, err = storage.GetUrl("unknown")
	assert.ErrorIs(t, err, errdef.ErrURLNotFound)
}

func TestStorage_GetAlias(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := &Storage{db: db}

	url := "https://example.com"
	expectedAlias := "abc123"

	mock.ExpectPrepare("SELECT alias FROM url WHERE url").
		ExpectQuery().
		WithArgs(url).
		WillReturnRows(sqlmock.NewRows([]string{"alias"}).AddRow(expectedAlias))

	alias, err := storage.GetAlias(url)
	assert.NoError(t, err)
	assert.Equal(t, expectedAlias, alias)

	mock.ExpectPrepare("SELECT alias FROM url WHERE url").
		ExpectQuery().
		WithArgs("https://unknown.com").
		WillReturnError(sql.ErrNoRows)

	_, err = storage.GetAlias("https://unknown.com")
	assert.ErrorIs(t, err, errdef.ErrURLNotFound)
}

func TestStorage_CheckAliasURLExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := &Storage{db: db}

	alias := "abc123"

	mock.ExpectPrepare("SELECT EXISTS").
		ExpectQuery().
		WithArgs(alias).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	exists, err := storage.CheckAliasURLExists(alias)
	assert.NoError(t, err)
	assert.True(t, exists)

	mock.ExpectPrepare("SELECT EXISTS").
		ExpectQuery().
		WithArgs("unknown").
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	exists, err = storage.CheckAliasURLExists("unknown")
	assert.NoError(t, err)
	assert.False(t, exists)
}
