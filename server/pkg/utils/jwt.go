package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// BỔ SUNG ĐOẠN NÀY: Khai báo struct CustomClaims
type CustomClaims struct {
	OrgID string `json:"org_id,omitempty"`
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

var jwtSecret []byte
var accessTokenExpire time.Duration
var refreshTokenExpire time.Duration

func InitJWT(secret string, accessExpire, refreshExpire time.Duration) {
	jwtSecret = []byte(secret)
	accessTokenExpire = accessExpire
	refreshTokenExpire = refreshExpire
}

func GenerateToken(userID, orgID string) (string, error) {


	regClaims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenExpire)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	claims := CustomClaims{
		OrgID:			  orgID,
		UserID:           userID,
		RegisteredClaims: regClaims,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func GenerateRefreshToken(userID string, orgID string) (string, error) {


	regClaims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTokenExpire)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	claims := CustomClaims{
		OrgID:			  orgID,
		UserID:           userID,
		RegisteredClaims: regClaims,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ValidateToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("thuật toán ký không hợp lệ")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, errors.New("token không hợp lệ")
	}

	return claims, nil
}