package model

type ShortURL struct {
	ID        uint   `gorm:"primaryKey"`
	ShortCode string `gorm:"type:varchar(8);uniqueIndex"`
	LongURL   string `gorm:"type:varchar(2048)"`
}
