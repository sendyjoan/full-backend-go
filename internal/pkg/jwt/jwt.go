package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Secrets struct {
	Access  []byte
	Refresh []byte
}

type Claims struct {
	UserID string `json:"uid"`
	jwt.RegisteredClaims
}

func GenerateAccess(userID string, secret []byte, ttl time.Duration) (string, error) {
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
}

func GenerateRefresh(userID string, secret []byte, ttl time.Duration) (string, error) {
	// refresh bisa pakai claims minimal
	claims := jwt.MapClaims{"uid": userID, "exp": time.Now().Add(ttl).Unix()}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
}

func ParseAccess(tokenStr string, secret []byte) (*Claims, error) {
	t, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	return t.Claims.(*Claims), nil
}
