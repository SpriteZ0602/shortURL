package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var Secret = []byte("change_me")

type Claims struct {
	UserID uint
	jwt.RegisteredClaims
}

// Generate 生成 token
func Generate(uid uint) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}).SignedString(Secret)
}

// Parse 解析 token
func Parse(token string) (*Claims, error) {
	var claims Claims
	_, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (any, error) { return Secret, nil })
	return &claims, err
}
