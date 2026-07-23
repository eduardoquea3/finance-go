package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalidToken = errors.New("invalid token")

type Service struct {
	secret     []byte
	issuer     string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewService(secret, issuer string, accessTTL time.Duration) *Service {
	return &Service{secret: []byte(secret), issuer: issuer, accessTTL: accessTTL, refreshTTL: 30 * 24 * time.Hour}
}

func (s *Service) Generate(subject, tokenType string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Issuer:    s.issuer,
		Subject:   subject,
		ID:        tokenType,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.secret)
}

func (s *Service) Parse(tokenString, expectedType string) (*jwt.RegisteredClaims, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return s.secret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}), jwt.WithIssuer(s.issuer), jwt.WithExpirationRequired())
	if err != nil || !token.Valid || claims.ID != expectedType || claims.Subject == "" {
		return nil, ErrInvalidToken
	}
	return claims, nil
}

func (s *Service) AccessTTL() time.Duration  { return s.accessTTL }
func (s *Service) RefreshTTL() time.Duration { return s.refreshTTL }
