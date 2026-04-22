package main

import (
	"database/sql"
	"log"
)

func setupDatabase(db *sql.DB) error {
	schema := `
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

	CREATE TABLE IF NOT EXISTS profiles (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		name VARCHAR(255) NOT NULL UNIQUE,
		gender VARCHAR(20) NOT NULL,
		gender_probability FLOAT NOT NULL,
		age INT NOT NULL,
		age_group VARCHAR(50) NOT NULL,
		country_id VARCHAR(2) NOT NULL,
		country_name VARCHAR(255) NOT NULL,
		country_probability FLOAT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		
		CHECK (gender IN ('male', 'female')),
		CHECK (gender_probability >= 0 AND gender_probability <= 1),
		CHECK (age >= 0 AND age <= 150),
		CHECK (age_group IN ('child', 'teenager', 'adult', 'senior')),
		CHECK (country_probability >= 0 AND country_probability <= 1)
	);

	CREATE INDEX IF NOT EXISTS idx_profiles_gender ON profiles(gender);
	CREATE INDEX IF NOT EXISTS idx_profiles_age ON profiles(age);
	CREATE INDEX IF NOT EXISTS idx_profiles_age_group ON profiles(age_group);
	CREATE INDEX IF NOT EXISTS idx_profiles_country_id ON profiles(country_id);
	CREATE INDEX IF NOT EXISTS idx_profiles_gender_probability ON profiles(gender_probability);
	CREATE INDEX IF NOT EXISTS idx_profiles_country_probability ON profiles(country_probability);
	CREATE INDEX IF NOT EXISTS idx_profiles_created_at ON profiles(created_at);
	`

	_, err := db.Exec(schema)
	if err != nil {
		log.Printf("Schema creation error (may already exist): %v", err)
		// Continue anyway as this might just mean the table already exists
	}

	return nil
}
