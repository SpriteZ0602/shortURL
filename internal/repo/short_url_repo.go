package repo

import (
	"context"
	"errors"
	"shortURL/internal/model"

	"gorm.io/gorm"
)

type ShortURLRepo struct{ db *gorm.DB }

// New 创建短链接仓库
func New(db *gorm.DB) *ShortURLRepo { return &ShortURLRepo{db} }

// Save 保存短链接
func (r *ShortURLRepo) Save(ctx context.Context, s *model.ShortURL) error {
	return r.db.WithContext(ctx).Create(s).Error
}

// FindByLong 通过长链接查询短链接
func (r *ShortURLRepo) FindByLong(ctx context.Context, url string) (*model.ShortURL, error) {
	var s model.ShortURL
	err := r.db.WithContext(ctx).Where("long_url = ?", url).First(&s).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &s, err
}

// FindByCode 通过短链接码查询短链接
func (r *ShortURLRepo) FindByCode(ctx context.Context, code string) (*model.ShortURL, error) {
	var su model.ShortURL
	err := r.db.WithContext(ctx).Where("short_code = ?", code).First(&su).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &su, err
}
