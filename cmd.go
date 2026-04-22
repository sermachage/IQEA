package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"
)

func generateCommand() {
	fs := flag.NewFlagSet("generate", flag.ExitOnError)
	output := fs.String("output", "profiles.json", "Output file path")
	count := fs.Int("count", 2026, "Number of profiles to generate")
	fs.Parse(os.Args[2:])

	if err := GenerateAndSaveTestProfiles(*output, *count); err != nil {
		log.Fatalf("Failed to generate profiles: %v\n", err)
	}
}

// ComprehensiveTestProfiles generates realistic demographic profiles
func ComprehensiveTestProfiles(count int) []RawProfile {
	maleFirstNames := []string{
		"Emmanuel", "Kwame", "Joseph", "Ahmed", "Hassan", "Ibrahim", "Ali", "Omar",
		"Jamal", "Tariq", "Mohammed", "Kofi", "Yusuf", "Amadou", "Moussa",
		"Ismail", "Rashid", "Karim", "Samir", "Khalid", "Majid", "Nabil", "Salim",
		"Adil", "Aziz", "Hakim", "Kamal", "Malik", "Nasir", "Rasheed",
		"Samuel", "David", "John", "Peter", "Mark", "Paul", "James", "Michael",
		"Charles", "Robert", "William", "Daniel", "Richard", "Thomas",
	}

	femaleFirstNames := []string{
		"Zainab", "Amara", "Fatima", "Khadija", "Aisha", "Leila", "Hana", "Yasmin",
		"Mariam", "Noor", "Layla", "Sarah", "Maryam", "Alia", "Lina", "Dina",
		"Nadia", "Rania", "Sana", "Samira", "Tania", "Vicky", "Wafa", "Yara",
		"Zara", "Amina", "Baida", "Citra", "Dalia", "Elia", "Farah", "Gina",
		"Hannah", "Iris", "Jessica", "Karen", "Linda", "Monica", "Nancy", "Olivia",
		"Patricia", "Rachel", "Sandra", "Teresa", "Ursula", "Victoria", "Wendy",
	}

	surnames := []string{
		"Okafor", "Adeyemi", "Asante", "Kipchoge", "Mukwaya", "Mensah", "Kamau",
		"Traore", "Diallo", "Toure", "Sow", "Ba", "Diop", "N'Diaye", "Kone",
		"Cisse", "Dabo", "Deme", "Sall", "Fall", "Bah", "Bel",
		"Camara", "Conde", "Conte", "Dibba", "Dibs",
		"Hassan", "Ahmed", "Ali", "Mohammed", "Fatima", "Omar",
		"Rashid", "Karim", "Aziz", "Malik", "Nasir", "Salim", "Samir",
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
		{"Ivory Coast", "CI"},
		{"Burkina Faso", "BF"},
		{"Rwanda", "RW"},
		{"Democratic Republic of Congo", "CD"},
		{"South Africa", "ZA"},
		{"Ethiopia", "ET"},
		{"Egypt", "EG"},
		{"Sudan", "SD"},
		{"Morocco", "MA"},
		{"Tunisia", "TN"},
		{"Zambia", "ZM"},
		{"Zimbabwe", "ZW"},
		{"Namibia", "NA"},
	}

	profiles := make([]RawProfile, count)
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < count; i++ {
		gender := "male"
		firstName := maleFirstNames[rand.Intn(len(maleFirstNames))]
		if rand.Float64() > 0.5 {
			gender = "female"
			firstName = femaleFirstNames[rand.Intn(len(femaleFirstNames))]
		}

		surname := surnames[rand.Intn(len(surnames))]
		name := fmt.Sprintf("%s %s", firstName, surname)

		age := rand.Intn(80) + 8

		// Determine age group based on age
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

		// Gender confidence is usually high
		genderProb := 0.80 + rand.Float64()*0.19

		// Country confidence varies
		countryProb := 0.75 + rand.Float64()*0.24

		profiles[i] = RawProfile{
			Name:               name,
			Gender:             gender,
			GenderProbability:  genderProb,
			Age:                age,
			AgeGroup:           ageGroup,
			CountryID:          country.code,
			CountryName:        country.name,
			CountryProbability: countryProb,
		}
	}

	return profiles
}

// GenerateTestProfiles creates a sample dataset for testing
func GenerateTestProfiles(count int) []RawProfile {
	return ComprehensiveTestProfiles(count)
}

// GenerateAndSaveTestProfiles generates and saves test profiles to a file
func GenerateAndSaveTestProfiles(filename string, count int) error {
	log.Printf("Generating %d test profiles...\n", count)
	profiles := GenerateTestProfiles(count)

	log.Printf("Marshaling to JSON...\n")
	data, err := json.MarshalIndent(profiles, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal profiles: %v", err)
	}

	log.Printf("Writing to file: %s\n", filename)
	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	log.Printf("✓ Generated and saved %d test profiles to %s\n", count, filename)
	log.Printf("  File size: %.2f MB\n", float64(len(data))/1024/1024)

	return nil
}
