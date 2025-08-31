package service

import (
	"context"
	"fmt"
	"shortURL/internal/model"
	"shortURL/internal/repo"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type ShortService struct {
	repo  *repo.ShortURLRepo
	genID func() string
}

// New 创建一个 ShortService 实例
func New(r *repo.ShortURLRepo, gen func() string) *ShortService {
	return &ShortService{repo: r, genID: gen}
}

// Shorten 创建一个短链接
func (s *ShortService) Shorten(ctx context.Context, longURL string) (string, error) {
	tracer := otel.Tracer("shorturl")
	ctx, span := tracer.Start(ctx, "service.Shorten")
	span.SetAttributes(attribute.String("long_url", longURL))
	defer span.End()

	if su, _ := s.repo.FindByLong(ctx, longURL); su != nil {
		fmt.Println("found existing short URL:", su.ShortCode)
		return su.ShortCode, nil
	}

	code := s.genID()

	if err := s.repo.Save(ctx, &model.ShortURL{ShortCode: code, LongURL: longURL}); err != nil {
		fmt.Println("save short URL error:", err)
		return "", err
	}
	return code, nil
}
