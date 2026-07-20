package auth

import (
	"testing"
	"time"
)

func TestTokenService(t *testing.T) {
	tokens := NewTokenService("01234567890123456789012345678901", "test", time.Hour)

	t.Run("generates and parses a token", func(t *testing.T) {
		token, err := tokens.Generate("user-123")
		if err != nil {
			t.Fatalf("generate token: %v", err)
		}
		claims, err := tokens.Parse(token)
		if err != nil {
			t.Fatalf("parse token: %v", err)
		}
		if claims.Subject != "user-123" {
			t.Fatalf("subject = %q, want %q", claims.Subject, "user-123")
		}
	})

	t.Run("rejects a token signed with another secret", func(t *testing.T) {
		other := NewTokenService("12345678901234567890123456789012", "test", time.Hour)
		token, err := other.Generate("user-123")
		if err != nil {
			t.Fatalf("generate token: %v", err)
		}
		if _, err := tokens.Parse(token); err == nil {
			t.Fatal("expected invalid token error")
		}
	})
}
