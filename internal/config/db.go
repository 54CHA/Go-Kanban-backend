package config

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var DB *pgxpool.Pool

func InitDB() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		log.Fatalf("Error parsing database config: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	// Verify connection
	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("Unable to ping database: %v", err)
	}

	DB = pool

	// Initialize tables
	if err := createTables(); err != nil {
		log.Fatalf("Error creating tables: %v", err)
	}

	log.Println("Successfully connected to database")
}

func createTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS tasks (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			description TEXT,
			priority TEXT NOT NULL,
			status TEXT NOT NULL,
			start_date TEXT NOT NULL,
			due_date TEXT,
			due_time TEXT,
			parent_id TEXT REFERENCES tasks(id) ON DELETE CASCADE
		)`,
	}

	for _, query := range queries {
		_, err := DB.Exec(context.Background(), query)
		if err != nil {
			return err
		}
	}

	return nil
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
} 