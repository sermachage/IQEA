package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type RawProfile struct {
	Name               string  `json:"name"`
	Gender             string  `json:"gender"`
	GenderProbability  float64 `json:"gender_probability"`
	Age                int     `json:"age"`
	AgeGroup           string  `json:"age_group"`
	CountryID          string  `json:"country_id"`
	CountryName        string  `json:"country_name"`
	CountryProbability float64 `json:"country_probability"`
}

type profileFile struct {
	Profiles []RawProfile `json:"profiles"`
}

func decodeProfilesJSON(data []byte) ([]RawProfile, error) {
	var profiles []RawProfile
	if err := json.Unmarshal(data, &profiles); err == nil {
		return profiles, nil
	}

	var wrapped profileFile
	if err := json.Unmarshal(data, &wrapped); err == nil && len(wrapped.Profiles) > 0 {
		return wrapped.Profiles, nil
	}

	return nil, fmt.Errorf("failed to parse JSON: expected an array of profiles or an object with a profiles field")
}

func seedDatabase(db *sql.DB, jsonFilePath string) error {
	log.Println("Starting database seeding...")

	// Read JSON file
	data, err := ioutil.ReadFile(jsonFilePath)
	if err != nil {
		// Try to download from URL if local file not found
		log.Println("Local file not found, attempting to download...")
		// This would be the provided link in the requirements
		return fmt.Errorf("JSON file not found at %s", jsonFilePath)
	}

	profiles, err := decodeProfilesJSON(data)
	if err != nil {
		return err
	}

	log.Printf("Loaded %d profiles from JSON\n", len(profiles))

	// Clear existing data
	_, err = db.Exec("TRUNCATE TABLE profiles RESTART IDENTITY CASCADE;")
	if err != nil {
		log.Printf("Could not truncate table (may not exist yet): %v\n", err)
	}

	// Insert profiles
	stmt, err := db.Prepare(`
		INSERT INTO profiles (name, gender, gender_probability, age, age_group, country_id, country_name, country_probability, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (name) DO NOTHING
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	inserted := 0
	skipped := 0

	for _, p := range profiles {
		result, err := stmt.Exec(
			p.Name,
			p.Gender,
			p.GenderProbability,
			p.Age,
			p.AgeGroup,
			p.CountryID,
			p.CountryName,
			p.CountryProbability,
			time.Now().UTC(),
		)

		if err != nil {
			log.Printf("Error inserting profile %s: %v\n", p.Name, err)
			skipped++
			continue
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			log.Printf("Could not get rows affected: %v\n", err)
			continue
		}

		if rowsAffected > 0 {
			inserted++
		} else {
			skipped++
		}
	}

	log.Printf("Seeding complete: %d inserted, %d skipped\n", inserted, skipped)

	// Verify count
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM profiles").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to verify count: %v", err)
	}

	log.Printf("Total profiles in database: %d\n", count)

	return nil
}

func downloadAndSeedProfiles(db *sql.DB, url string) error {
	log.Println("Downloading profiles from:", url)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download profiles: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download profiles: status %d", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	profiles, err := decodeProfilesJSON(data)
	if err != nil {
		return err
	}

	log.Printf("Downloaded %d profiles\n", len(profiles))

	// Check if profiles exist
	var existingCount int
	err = db.QueryRow("SELECT COUNT(*) FROM profiles").Scan(&existingCount)
	if err != nil {
		log.Printf("Could not check existing profiles: %v\n", err)
	}

	if existingCount == 0 {
		log.Println("Database is empty, seeding with downloaded profiles...")

		stmt, err := db.Prepare(`
			INSERT INTO profiles (name, gender, gender_probability, age, age_group, country_id, country_name, country_probability, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			ON CONFLICT (name) DO NOTHING
		`)
		if err != nil {
			return fmt.Errorf("failed to prepare statement: %v", err)
		}
		defer stmt.Close()

		inserted := 0

		for _, p := range profiles {
			_, err := stmt.Exec(
				p.Name,
				p.Gender,
				p.GenderProbability,
				p.Age,
				p.AgeGroup,
				p.CountryID,
				p.CountryName,
				p.CountryProbability,
				time.Now().UTC(),
			)

			if err != nil {
				log.Printf("Error inserting profile %s: %v\n", p.Name, err)
				continue
			}
			inserted++
		}

		log.Printf("Inserted %d profiles\n", inserted)
	} else {
		log.Printf("Database already has %d profiles, skipping seed\n", existingCount)
	}

	return nil
}
