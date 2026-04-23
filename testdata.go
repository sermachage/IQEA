package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
)

// GenerateTestProfiles creates a sample dataset for testing
func GenerateTestProfilesLegacy(count int) []RawProfile {
	names := []string{
		"Emmanuel", "Zainab", "Kwame", "Amara", "Joseph", "Fatima", "Kofi", "Khadija",
		"Ibrahim", "Aisha", "Jamal", "Leila", "Ahmed", "Noor", "Hassan", "Yasmin",
		"Ali", "Layla", "Omar", "Hana", "Mohammed", "Sarah", "Tariq", "Mariam",
	}

	surnames := []string{
		"Okafor", "Adeyemi", "Asante", "Kipchoge", "Mukwaya", "Mensah", "Kamau",
		"Traore", "Diallo", "Nguyen", "Aziz", "Hassan", "Ibrahim", "Amara",
	}

	countries := []struct {
		name string
		code string
	}{
		{"Nigeria", "NG"},
		{"Ghana", "GH"},
		{"Kenya", "KE"},
		{"Uganda", "UG"},
		{"Tanzania", "TZ"},
		{"Cameroon", "CM"},
		{"Senegal", "SN"},
		{"Angola", "AO"},
		{"Benin", "BJ"},
		{"Mali", "ML"},
	}

	profiles := make([]RawProfile, count)

	for i := 0; i < count; i++ {
		name := names[rand.Intn(len(names))] + " " + surnames[rand.Intn(len(surnames))]
		gender := "male"
		if rand.Float64() > 0.5 {
			gender = "female"
		}
		age := rand.Intn(80) + 8

		// Determine age group
		var ageGroup string
		if age <= 12 {
			ageGroup = "child"
		} else if age <= 19 {
			ageGroup = "teenager"
		} else if age <= 59 {
			ageGroup = "adult"
		} else {
			ageGroup = "senior"
		}

		country := countries[rand.Intn(len(countries))]

		profiles[i] = RawProfile{
			Name:                fmt.Sprintf("%s_%d", name, i),
			Gender:              gender,
			GenderProbability:   0.80 + rand.Float64()*0.19,
			Age:                 age,
			AgeGroup:            ageGroup,
			CountryID:           country.code,
			CountryName:         country.name,
			CountryProbability:  0.75 + rand.Float64()*0.24,
		}
	}

	return profiles
}

// GenerateAndSaveTestProfiles generates and saves test profiles to a file
func GenerateAndSaveTestProfilesLegacy(filename string, count int) error {
	profiles := GenerateTestProfilesLegacy(count)
	data, err := json.MarshalIndent(profiles, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal profiles: %v", err)
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	fmt.Printf("Generated and saved %d test profiles to %s\n", count, filename)
	return nil
}
