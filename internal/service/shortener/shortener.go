package shortener

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"url-shortener/internal/config"
	"url-shortener/internal/storage"
	"url-shortener/internal/storage/errdef"
	pb "url-shortener/protos/gen/go"
)

const (
	chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
)

var (
	aliasLength int
	maxAttempts int
)

type Service struct {
	pb.UnimplementedUrlShortenerServer
	storage storage.Storage
}

func New(storage storage.Storage, config config.Config) *Service {
	aliasLength = config.URLGenerator.Length
	maxAttempts = config.URLGenerator.MaxAttempt
	return &Service{storage: storage}
}

func (s *Service) CreateShortURL(ctx context.Context, req *pb.CreateShortURLRequest) (*pb.CreateShortURLResponse, error) {
	originalURL := req.GetOriginalUrl()

	if alias, err := s.storage.GetAlias(originalURL); err == nil {
		return &pb.CreateShortURLResponse{Alias: alias}, nil
	}

	for i := 0; i < maxAttempts; i++ {
		alias, err := generateAlias()
		if err != nil {
			return nil, fmt.Errorf("failed to generate alias: %w", err)
		}

		if exists, _ := s.storage.CheckAliasURLExists(alias); !exists {
			if err := s.storage.SaveUrl(originalURL, alias); err != nil {
				if errors.Is(err, errdef.ErrURLExists) {
					continue
				}
				return nil, err
			}
			return &pb.CreateShortURLResponse{Alias: alias}, nil
		}
	}
	return nil, errors.New("failed to generate unique alias")
}

func (s *Service) GetOriginalURL(ctx context.Context, req *pb.GetOriginalURLRequest) (*pb.GetOriginalURLResponse, error) {
	alias := req.GetAlias()

	originalURL, err := s.storage.GetUrl(alias)
	if err != nil {
		return nil, err
	}

	return &pb.GetOriginalURLResponse{OriginalUrl: originalURL}, nil
}

func generateAlias() (string, error) {
	result := make([]byte, aliasLength)
	for i := 0; i < aliasLength; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", err
		}
		result[i] = chars[num.Int64()]
	}
	return string(result), nil
}
