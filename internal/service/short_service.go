package service

import (
	"fmt"
	"shortURL/internal/model"
	"shortURL/internal/repo"
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
func (s *ShortService) Shorten(longURL string) (string, error) {
	if su, _ := s.repo.FindByLong(longURL); su != nil {
		fmt.Println("found existing short URL:", su.ShortCode)
		return su.ShortCode, nil
	}

	code := s.genID()

	if err := s.repo.Save(&model.ShortURL{ShortCode: code, LongURL: longURL}); err != nil {
		fmt.Println("save short URL error:", err)
		return "", err
	}
	return code, nil
}
