package repo

import (
	"errors"
	"gorm.io/gorm"
	"shortURL/internal/model"
)

type ShortURLRepo struct{ db *gorm.DB }

func New(db *gorm.DB) *ShortURLRepo { return &ShortURLRepo{db} }

func (r *ShortURLRepo) Save(s *model.ShortURL) error { return r.db.Create(s).Error }
func (r *ShortURLRepo) FindByLong(url string) (*model.ShortURL, error) {
	var s model.ShortURL
	err := r.db.Where("long_url = ?", url).First(&s).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &s, err
}
