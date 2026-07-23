package main

import (
	"flag"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func main() {
	direction := flag.String("direction", "up", "migration direction: up, down, or force")
	steps := flag.Int("steps", 1, "number of migrations for down")
	version := flag.Int("version", 0, "version used with force")
	flag.Parse()

	_ = godotenv.Load()
	viper.AutomaticEnv()
	databaseURL := viper.GetString("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	migrator, err := migrate.New("file://migrations", databaseURL)
	if err != nil {
		log.Fatalf("create migrator: %v", err)
	}
	defer migrator.Close()

	switch *direction {
	case "up":
		err = migrator.Up()
	case "down":
		err = migrator.Steps(-*steps)
	case "force":
		if *version < 0 {
			log.Fatal("version must be zero or greater")
		}
		err = migrator.Force(*version)
	default:
		log.Fatalf("unsupported direction %q", *direction)
	}
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("run migrations: %v", err)
	}
}
