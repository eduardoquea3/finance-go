package main

import (
	"log"
	"net/http"

	_ "github.com/eduardoquea3/finance-go/docs"
	"github.com/eduardoquea3/finance-go/internal/auth"
	"github.com/eduardoquea3/finance-go/internal/config"
	httpapi "github.com/eduardoquea3/finance-go/internal/http"
	"github.com/eduardoquea3/finance-go/internal/platform/database"
	"github.com/eduardoquea3/finance-go/internal/platform/token"
	"github.com/eduardoquea3/finance-go/internal/user"
)

// @title Finance Go API
// @version 1.0
// @description API financiera construida con Go y Gin.
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	db, err := database.Open(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}

	tokenService := token.NewService(cfg.JWTSecret, cfg.JWTIssuer, cfg.JWTTTL)
	users := user.NewPostgresRepository(db)
	authHandler := auth.NewHTTPHandler(auth.NewService(users, tokenService))
	router := httpapi.NewRouter(db, authHandler, tokenService)
	server := &http.Server{Addr: cfg.HTTPAddress, Handler: router}

	log.Printf("API listening on %s", cfg.HTTPAddress)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("run server: %v", err)
	}
}
