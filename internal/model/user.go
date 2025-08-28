package model

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Email    string `gorm:"uniqueIndex;size:255"`
	Password string `json:"-"` // 不序列化
}
