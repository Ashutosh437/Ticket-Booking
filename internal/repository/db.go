package repository

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

func InitDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Create tables
	userTableQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		created_at DATETIME NOT NULL
	);`

	ticketTableQuery := `
	CREATE TABLE IF NOT EXISTS tickets (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		description TEXT NOT NULL,
		status TEXT NOT NULL,
		user_id TEXT NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`

	if _, err := db.Exec(userTableQuery); err != nil {
		return nil, fmt.Errorf("failed to create users table: %w", err)
	}

	if _, err := db.Exec(ticketTableQuery); err != nil {
		return nil, fmt.Errorf("failed to create tickets table: %w", err)
	}

	log.Println("Database initialized successfully at:", dbPath)
	return db, nil
}
