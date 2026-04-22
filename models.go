package main

import "time"

type Profile struct {
	ID                   string    `json:"id"`
	Name                 string    `json:"name"`
	Gender               string    `json:"gender"`
	GenderProbability    float64   `json:"gender_probability"`
	Age                  int       `json:"age"`
	AgeGroup             string    `json:"age_group"`
	CountryID            string    `json:"country_id"`
	CountryName          string    `json:"country_name"`
	CountryProbability   float64   `json:"country_probability"`
	CreatedAt            time.Time `json:"created_at"`
}

type ProfilesResponse struct {
	Status string     `json:"status"`
	Page   int        `json:"page"`
	Limit  int        `json:"limit"`
	Total  int        `json:"total"`
	Data   []Profile  `json:"data"`
}

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type QueryFilters struct {
	Gender                  string
	AgeGroup                string
	CountryID               string
	MinAge                  *int
	MaxAge                  *int
	MinGenderProbability    *float64
	MinCountryProbability   *float64
	SortBy                  string
	Order                   string
	Page                    int
	Limit                   int
}
