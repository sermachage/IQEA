package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	// load .env file
	_ = godotenv.Load()

	// Initialize database
	var err error
	db, err = initializeDatabase()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Handle CLI commands first
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "seed":
			handleSeed()
			os.Exit(0)
		case "generate":
			generateCommand()
			os.Exit(0)
		}
	}

	// Set up routes (only runs if no CLI command is given)
	http.HandleFunc("GET /api/profiles", handleGetProfiles)
	http.HandleFunc("GET /api/profiles/search", handleSearchProfiles)
	http.HandleFunc("GET /health", handleHealth)
	http.HandleFunc("OPTIONS /", handleOptions)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, corsMiddleware(http.DefaultServeMux)))
}

func initializeDatabase() (*sql.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=iqea sslmode=disable"
	}

	dbConn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	if err := dbConn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	dbConn.SetMaxOpenConns(25)
	dbConn.SetMaxIdleConns(5)

	if err := setupDatabase(dbConn); err != nil {
		log.Printf("Warning: database setup error (may already exist): %v\n", err)
	}

	log.Println("Database connection successful")
	return dbConn, nil
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func handleOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.WriteHeader(http.StatusOK)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprintf(w, `{"status":"ok"}`)
}
