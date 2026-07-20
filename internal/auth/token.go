package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalidToken = errors.New("invalid token")

type TokenService struct {
	secret []byte
	issuer string
	ttl    time.Duration
}

func NewTokenService(secret, issuer string, ttl time.Duration) *TokenService {
	return &TokenService{secret: []byte(secret), issuer: issuer, ttl: ttl}
}

func (s *TokenService) Generate(subject string) (string, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Issuer:    s.issuer,
		Subject:   subject,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(s.ttl)),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.secret)
}

func (s *TokenService) Parse(tokenString string) (*jwt.RegisteredClaims, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return s.secret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}), jwt.WithIssuer(s.issuer), jwt.WithExpirationRequired())
	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}
	return claims, nil
}
