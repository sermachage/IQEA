package main

import (
	"database/sql"
	_ "embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

//go:embed profiles.json
var embeddedProfiles []byte

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

// seedFromEmbedded seeds using the compile-time embedded profiles.json.
// Safe to call at server startup — does NOT call os.Exit.
func seedFromEmbedded(db *sql.DB) error {
	profiles, err := decodeProfilesJSON(embeddedProfiles)
	if err != nil {
		return fmt.Errorf("failed to decode embedded profiles: %v", err)
	}
	return insertProfiles(db, profiles)
}

// insertProfiles does the actual bulk insert shared by all seed paths.
func insertProfiles(db *sql.DB, profiles []RawProfile) error {
	log.Printf("Inserting %d profiles...\n", len(profiles))

	stmt, err := db.Prepare(`
		INSERT INTO profiles (name, gender, gender_probability, age, age_group, country_id, country_name, country_probability, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (name) DO NOTHING
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	inserted, skipped := 0, 0
	for _, p := range profiles {
		result, err := stmt.Exec(
			p.Name, p.Gender, p.GenderProbability, p.Age,
			p.AgeGroup, p.CountryID, p.CountryName, p.CountryProbability,
			time.Now().UTC(),
		)
		if err != nil {
			log.Printf("Error inserting profile %s: %v\n", p.Name, err)
			skipped++
			continue
		}
		rows, _ := result.RowsAffected()
		if rows > 0 {
			inserted++
		} else {
			skipped++
		}
	}

	log.Printf("Seed complete: %d inserted, %d skipped\n", inserted, skipped)

	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM profiles").Scan(&count); err != nil {
		return fmt.Errorf("failed to verify count: %v", err)
	}
	log.Printf("Total profiles in database: %d\n", count)
	return nil
}

func seedDatabase(db *sql.DB, jsonFilePath string) error {
	log.Println("Starting database seeding from file:", jsonFilePath)

	data, err := ioutil.ReadFile(jsonFilePath)
	if err != nil {
		log.Printf("JSON file not found at %s, falling back to embedded profiles\n", jsonFilePath)
		return seedFromEmbedded(db)
	}

	profiles, err := decodeProfilesJSON(data)
	if err != nil {
		return err
	}

	_, err = db.Exec("TRUNCATE TABLE profiles RESTART IDENTITY CASCADE;")
	if err != nil {
		log.Printf("Could not truncate table: %v\n", err)
	}

	return insertProfiles(db, profiles)
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

	var existingCount int
	if err := db.QueryRow("SELECT COUNT(*) FROM profiles").Scan(&existingCount); err != nil {
		log.Printf("Could not check existing profiles: %v\n", err)
	}

	if existingCount > 0 {
		log.Printf("Database already has %d profiles, skipping seed\n", existingCount)
		return nil
	}

	return insertProfiles(db, profiles)
}
