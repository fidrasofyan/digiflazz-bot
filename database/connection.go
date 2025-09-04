package database

import (
	"context"
	"database/sql"
	"embed"
	"log"

	"github.com/fidrasofyan/digiflazz-bot/internal/config"
	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)

var DBConn *sql.DB
var Sqlc *Queries

//go:embed migrations/*.sql
var embeddedMigrations embed.FS

func MustLoadDatabase(ctx context.Context) {
	var err error

	// Open database
	DBConn, err = sql.Open("sqlite", config.Cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	// Ping
	if err := DBConn.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Set PRAGMA
	pragmas := []string{
		"PRAGMA journal_mode=WAL;",
		"PRAGMA synchronous=NORMAL;",
		"PRAGMA temp_store=MEMORY;",
	}
	for _, pragma := range pragmas {
		if _, err := DBConn.ExecContext(ctx, pragma); err != nil {
			log.Fatalf("Failed to set pragma: %v", err)
		}
	}

	// Initialize
	Sqlc, err = Prepare(ctx, DBConn)
	if err != nil {
		log.Fatalf("Failed to initialize SQLC: %v", err)
	}
	log.Printf("Database connected (%s)\n", config.Cfg.DatabaseURL)

	// Apply migrations
	if err := applyMigrations(DBConn); err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}
}

func applyMigrations(db *sql.DB) error {
	goose.SetBaseFS(embeddedMigrations)
	if err := goose.SetDialect("sqlite"); err != nil {
		return err
	}
	if err := goose.Up(db, "migrations"); err != nil {
		return err
	}
	return nil
}
