package config

import (
	"errors"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	HTTPAddress string        `mapstructure:"http_address"`
	DatabaseURL string        `mapstructure:"database_url"`
	JWTSecret   string        `mapstructure:"jwt_secret"`
	JWTIssuer   string        `mapstructure:"jwt_issuer"`
	JWTTTL      time.Duration `mapstructure:"jwt_ttl"`
}

func Load() (Config, error) {
	_ = godotenv.Load()

	viper.SetDefault("http_address", ":8080")
	viper.SetDefault("jwt_issuer", "finance-go")
	viper.SetDefault("jwt_ttl", 24*time.Hour)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	for key, env := range map[string]string{
		"http_address": "HTTP_ADDRESS",
		"database_url": "DATABASE_URL",
		"jwt_secret":   "JWT_SECRET",
		"jwt_issuer":   "JWT_ISSUER",
		"jwt_ttl":      "JWT_TTL",
	} {
		if err := viper.BindEnv(key, env); err != nil {
			return Config{}, err
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return Config{}, err
	}
	if cfg.DatabaseURL == "" {
		return Config{}, errors.New("DATABASE_URL is required")
	}
	if len(cfg.JWTSecret) < 32 {
		return Config{}, errors.New("JWT_SECRET must contain at least 32 characters")
	}
	return cfg, nil
}
