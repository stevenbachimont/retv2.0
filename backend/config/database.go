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
	db, err := sql.Open("postgres", "user=postgres password=postgres dbname=carbon_calculator sslmode=disable")
	if err != nil {
		return nil, err
	}

	// Créer la table users
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL
		)
	`)
	if err != nil {
		return nil, err
	}

	// Créer la table results avec le champ month
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS results (
			id UUID PRIMARY KEY,
			user_id UUID REFERENCES users(id),
			category VARCHAR(50) NOT NULL,
			value FLOAT NOT NULL,
			inputs JSONB,
			month DATE NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(user_id, category, month)
		)
	`)
	if err != nil {
		return nil, err
	}

	return db, nil
}
