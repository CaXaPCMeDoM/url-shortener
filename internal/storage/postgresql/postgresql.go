package postgresql

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"log/slog"
	"url-shortener/internal/storage/errdef"
)

const (
	UniqueViolation pq.ErrorCode = "23505" // Ошибка нарушения уникальности
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	pool, err := sql.Open("postgres", storagePath)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	return &Storage{db: pool}, nil
}

func (s *Storage) SaveUrl(urlAbsolute string, alias string) error {
	stmt, err := s.db.Prepare("INSERT INTO url(alias, url) VALUES ($1, $2)")
	if err != nil {
		slog.Error("failed to prepare statement", slog.String("error", err.Error()))
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			slog.Error("failed to close statement", slog.String("error", err.Error()))
		}
	}()

	_, err = stmt.Exec(alias, urlAbsolute)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == UniqueViolation {
			slog.Warn("duplicate alias", slog.String("alias", alias))
			return errdef.ErrURLExists
		}

		slog.Error("failed to execute statement", slog.String("error", err.Error()))
		return fmt.Errorf("failed to execute statement: %w", err)
	}

	return nil
}

func (s *Storage) GetUrl(alias string) (url string, err error) {
	stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias = $1")
	if err != nil {
		slog.Error("failed to prepare statement", slog.String("error", err.Error()))
		return "", fmt.Errorf("failed to prepare statement: %w", err)
	}

	defer func() {
		if err := stmt.Close(); err != nil {
			slog.Error("failed to close statement", slog.String("error", err.Error()))
		}
	}()

	var urlAbsolute string

	err = stmt.QueryRow(alias).Scan(&urlAbsolute)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Warn("alias not found", slog.String("alias", alias))
			return "", errdef.ErrURLNotFound
		}

		slog.Error("failed to execute query", slog.String("error", err.Error()))
		return "", fmt.Errorf("failed to execute query: %w", err)
	}

	return urlAbsolute, nil
}

func (s *Storage) GetAlias(urlAbsolute string) (string, error) {
	stmt, err := s.db.Prepare("SELECT alias FROM url WHERE url = $1")
	if err != nil {
		slog.Error("failed to prepare statement", slog.String("error", err.Error()))
		return "", fmt.Errorf("failed to prepare statement: %w", err)
	}

	defer func() {
		if err := stmt.Close(); err != nil {
			slog.Error("failed to close statement", slog.String("error", err.Error()))
		}
	}()

	var alias string

	err = stmt.QueryRow(urlAbsolute).Scan(&alias)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Warn("url not found", slog.String("url", urlAbsolute))
			return "", errdef.ErrURLNotFound
		}

		slog.Error("failed to execute query", slog.String("error", err.Error()))
		return "", fmt.Errorf("failed to execute query: %w", err)
	}

	return alias, nil
}

func (s *Storage) CheckAliasURLExists(alias string) (bool, error) {
	var exists bool
	stmt, err := s.db.Prepare("SELECT EXISTS (SELECT 1 FROM url WHERE alias = $1)")
	if err != nil {
		slog.Error("failed to prepare statement", slog.String("error", err.Error()))
		return false, fmt.Errorf("failed to prepare statement: %w", err)
	}
	err = stmt.QueryRow(alias).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Warn("alias not found", slog.String("alias", alias))
			return false, nil
		}
		slog.Error("failed to execute query", slog.String("error", err.Error()))
		return false, fmt.Errorf("failed to execute query: %w", err)
	}
	return exists, nil
}
