package auth

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/eduardoquea3/finance-go/internal/platform/token"
	"github.com/eduardoquea3/finance-go/internal/user"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type Service struct {
	users  user.Repository
	tokens *token.Service
}

func NewService(users user.Repository, tokens *token.Service) *Service {
	return &Service{users: users, tokens: tokens}
}

func (s *Service) Login(ctx context.Context, email, password string) (TokenResponse, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	account, err := s.users.FindByEmail(ctx, email)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(account.PasswordHash), []byte(password)) != nil {
		return TokenResponse{}, ErrInvalidCredentials
	}
	return s.issue(subjectString(account.ID))
}

func (s *Service) Refresh(refreshToken string) (TokenResponse, error) {
	claims, err := s.tokens.Parse(refreshToken, "refresh")
	if err != nil {
		return TokenResponse{}, ErrInvalidCredentials
	}
	return s.issue(claims.Subject)
}

func (s *Service) issue(subject string) (TokenResponse, error) {
	access, err := s.tokens.Generate(subject, "access", s.tokens.AccessTTL())
	if err != nil {
		return TokenResponse{}, err
	}
	refresh, err := s.tokens.Generate(subject, "refresh", s.tokens.RefreshTTL())
	if err != nil {
		return TokenResponse{}, err
	}
	return TokenResponse{AccessToken: access, RefreshToken: refresh, TokenType: "Bearer", ExpiresIn: int64(s.tokens.AccessTTL().Seconds())}, nil
}

func subjectString(id uint) string { return strconv.FormatUint(uint64(id), 10) }
