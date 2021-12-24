package cred

import (
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type UserClaims struct {
	jwt.RegisteredClaims
	UserID string `json:"user_id"`
}

const (
	ExpireDayAccessToken  = 7
	ExpireDayRefreshToken = 30
)

var privateKey *rsa.PrivateKey
var publicKey *rsa.PublicKey

func InitJwt() error {
	sKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}
	privateKey = sKey
	publicKey = &sKey.PublicKey

	return nil
}

func IssueJWT(userID string, days time.Duration) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodRS512, UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{
				Time: time.Now().Add(24 * days * time.Hour),
			},
		},
		UserID: userID,
	})
	return t.SignedString(privateKey)
}

func ParseJWT(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("invalid sign method")
		}

		return publicKey, nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["user_id"].(string), nil
	}

	return "", errors.New("invalid token")
}
