package config

import (
	"database/sql"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "carbon_calculator"
)

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", "user=postgres dbname=carbon_calculator sslmode=disable")
	if err != nil {
		return nil, err
	}

	// Créer la table users si elle n'existe pas
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return nil, err
	}

	// Créer la table results si elle n'existe pas
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS results (
			id UUID PRIMARY KEY,
			user_id UUID REFERENCES users(id),
			category VARCHAR(50) NOT NULL,
			value FLOAT NOT NULL,
			inputs JSONB,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func createTables(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id VARCHAR(36) PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS results (
			id VARCHAR(36) PRIMARY KEY,
			user_id VARCHAR(36) REFERENCES users(id),
			category VARCHAR(50) NOT NULL,
			value FLOAT NOT NULL,
			inputs JSONB NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			return err
		}
	}

	return nil
}
