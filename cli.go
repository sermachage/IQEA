package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
)

func handleSeed() {
	flagset := flag.NewFlagSet("seed", flag.ExitOnError)
	file := flagset.String("file", "profiles.json", "Path to the JSON file with profiles")
	url := flagset.String("url", "", "URL to download profiles from")

	flagset.Parse(os.Args[2:])

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=iqea sslmode=disable"
	}

	seedDB, err := initDB(dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer seedDB.Close()

	if *url != "" {
		err = downloadAndSeedProfiles(seedDB, *url)
	} else {
		err = seedDatabase(seedDB, *file)
	}

	if err != nil {
		log.Fatal("Seeding failed:", err)
	}

	log.Println("Seeding completed successfully")
}

func initDB(dsn string) (*sql.DB, error) {
	dbConn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	if err := dbConn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	return dbConn, nil
}
