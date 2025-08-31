package repo

import (
	"errors"
	"shortURL/internal/model"
)
import "gorm.io/gorm"

type UserRepo struct{ db *gorm.DB }

// NewUser 创建用户仓库
func NewUser(db *gorm.DB) *UserRepo { return &UserRepo{db} }

// Save 保存用户
func (r *UserRepo) Save(u *model.User) error { return r.db.Create(u).Error }

// FindByEmail 根据邮箱查询用户
func (r *UserRepo) FindByEmail(email string) (*model.User, error) {
	var u model.User
	err := r.db.Where("email = ?", email).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &u, err
}
