package main

import (
	"log"
	"net/http"

	_ "github.com/eduardoquea3/finance-go/docs"
	"github.com/eduardoquea3/finance-go/internal/auth"
	"github.com/eduardoquea3/finance-go/internal/config"
	"github.com/eduardoquea3/finance-go/internal/database"
	"github.com/eduardoquea3/finance-go/internal/httpapi"
)

// @title Finance Go API
// @version 1.0
// @description API financiera construida con Go y Gin.
// @host localhost:8080
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

	tokenService := auth.NewTokenService(cfg.JWTSecret, cfg.JWTIssuer, cfg.JWTTTL)
	router := httpapi.NewRouter(db, tokenService)
	server := &http.Server{Addr: cfg.HTTPAddress, Handler: router}

	log.Printf("API listening on %s", cfg.HTTPAddress)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("run server: %v", err)
	}
}
