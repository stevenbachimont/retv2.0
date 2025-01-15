package config

import (
	"database/sql"
	"log"

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

	// Supprimer les tables existantes
	_, err = db.Exec(`
		DROP TABLE IF EXISTS results CASCADE;
		DROP TABLE IF EXISTS users CASCADE;
	`)
	if err != nil {
		log.Printf("Erreur lors de la suppression des tables: %v", err)
	}

	// Créer la table users avec la nouvelle structure
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
			username VARCHAR(50) NOT NULL,
			password VARCHAR(255) NOT NULL
		)
	`)
	if err != nil {
		return nil, err
	}

	// Créer la table results
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
