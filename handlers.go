package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func handleGetProfiles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Parse query parameters
	filters := QueryFilters{
		Page:  1,
		Limit: 10,
		Order: "asc",
	}

	// Parse string filters
	if gender := r.URL.Query().Get("gender"); gender != "" {
		gender = strings.ToLower(strings.TrimSpace(gender))
		if gender != "male" && gender != "female" {
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(ErrorResponse{
				Status:  "error",
				Message: "Invalid gender value",
			})
			return
		}
		filters.Gender = gender
	}

	if ageGroup := r.URL.Query().Get("age_group"); ageGroup != "" {
		ageGroup = strings.ToLower(strings.TrimSpace(ageGroup))
		validGroups := []string{"child", "teenager", "adult", "senior"}
		valid := false
		for _, vg := range validGroups {
			if ageGroup == vg {
				valid = true
				break
			}
		}
		if !valid {
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(ErrorResponse{
				Status:  "error",
				Message: "Invalid age_group value",
			})
			return
		}
		filters.AgeGroup = ageGroup
	}

	if countryID := r.URL.Query().Get("country_id"); countryID != "" {
		filters.CountryID = strings.ToUpper(strings.TrimSpace(countryID))
	}

	// Parse numeric filters
	if minAge := r.URL.Query().Get("min_age"); minAge != "" {
		val, err := strconv.Atoi(minAge)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(ErrorResponse{
				Status:  "error",
				Message: "Invalid min_age value",
			})
			return
		}
		filters.MinAge = &val
	}

	if maxAge := r.URL.Query().Get("max_age"); maxAge != "" {
		val, err := strconv.Atoi(maxAge)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(ErrorResponse{
				Status:  "error",
				Message: "Invalid max_age value",
			})
			return
		}
		filters.MaxAge = &val
	}

	if minGenderProb := r.URL.Query().Get("min_gender_probability"); minGenderProb != "" {
		val, err := strconv.ParseFloat(minGenderProb, 64)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(ErrorResponse{
				Status:  "error",
				Message: "Invalid min_gender_probability value",
			})
			return
		}
		filters.MinGenderProbability = &val
	}

	if minCountryProb := r.URL.Query().Get("min_country_probability"); minCountryProb != "" {
		val, err := strconv.ParseFloat(minCountryProb, 64)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(ErrorResponse{
				Status:  "error",
				Message: "Invalid min_country_probability value",
			})
			return
		}
		filters.MinCountryProbability = &val
	}

	// Parse sorting
	if sortBy := r.URL.Query().Get("sort_by"); sortBy != "" {
		sortBy = strings.ToLower(strings.TrimSpace(sortBy))
		validSorts := []string{"age", "created_at", "gender_probability"}
		valid := false
		for _, vs := range validSorts {
			if sortBy == vs {
				valid = true
				break
			}
		}
		if !valid {
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(ErrorResponse{
				Status:  "error",
				Message: "Invalid sort_by value",
			})
			return
		}
		filters.SortBy = sortBy
	}

	if order := r.URL.Query().Get("order"); order != "" {
		order = strings.ToLower(strings.TrimSpace(order))
		if order != "asc" && order != "desc" {
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(ErrorResponse{
				Status:  "error",
				Message: "Invalid order value",
			})
			return
		}
		filters.Order = order
	}

	// Parse pagination
	if page := r.URL.Query().Get("page"); page != "" {
		val, err := strconv.Atoi(page)
		if err != nil || val < 1 {
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(ErrorResponse{
				Status:  "error",
				Message: "Invalid page value",
			})
			return
		}
		filters.Page = val
	}

	if limit := r.URL.Query().Get("limit"); limit != "" {
		val, err := strconv.Atoi(limit)
		if err != nil || val < 1 || val > 50 {
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(ErrorResponse{
				Status:  "error",
				Message: "Invalid limit value (must be 1-50)",
			})
			return
		}
		filters.Limit = val
	}

	// Build and execute query
	profiles, total, err := queryProfiles(db, filters)
	if err != nil {
		log.Println("Query error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Status:  "error",
			Message: "Internal server error",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ProfilesResponse{
		Status: "success",
		Page:   filters.Page,
		Limit:  filters.Limit,
		Total:  total,
		Data:   profiles,
	})
}

func handleSearchProfiles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	query := r.URL.Query().Get("q")
	if query == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Status:  "error",
			Message: "Missing or empty parameter",
		})
		return
	}

	// Parse natural language query
	filters, err := parseNaturalLanguageQuery(query)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Status:  "error",
			Message: "Unable to interpret query",
		})
		return
	}

	// Parse pagination parameters
	if page := r.URL.Query().Get("page"); page != "" {
		val, err := strconv.Atoi(page)
		if err != nil || val < 1 {
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(ErrorResponse{
				Status:  "error",
				Message: "Invalid page value",
			})
			return
		}
		filters.Page = val
	} else {
		filters.Page = 1
	}

	if limit := r.URL.Query().Get("limit"); limit != "" {
		val, err := strconv.Atoi(limit)
		if err != nil || val < 1 || val > 50 {
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(ErrorResponse{
				Status:  "error",
				Message: "Invalid limit value (must be 1-50)",
			})
			return
		}
		filters.Limit = val
	} else {
		filters.Limit = 10
	}

	// Build and execute query
	profiles, total, err := queryProfiles(db, *filters)
	if err != nil {
		log.Println("Query error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Status:  "error",
			Message: "Internal server error",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ProfilesResponse{
		Status: "success",
		Page:   filters.Page,
		Limit:  filters.Limit,
		Total:  total,
		Data:   profiles,
	})
}

func queryProfiles(db *sql.DB, filters QueryFilters) ([]Profile, int, error) {
	// Build WHERE clause
	whereClauses := []string{"1=1"}
	args := []interface{}{}

	if filters.Gender != "" {
		whereClauses = append(whereClauses, "gender = $"+fmt.Sprintf("%d", len(args)+1))
		args = append(args, filters.Gender)
	}

	if filters.AgeGroup != "" {
		whereClauses = append(whereClauses, "age_group = $"+fmt.Sprintf("%d", len(args)+1))
		args = append(args, filters.AgeGroup)
	}

	if filters.CountryID != "" {
		whereClauses = append(whereClauses, "country_id = $"+fmt.Sprintf("%d", len(args)+1))
		args = append(args, filters.CountryID)
	}

	if filters.MinAge != nil {
		whereClauses = append(whereClauses, "age >= $"+fmt.Sprintf("%d", len(args)+1))
		args = append(args, *filters.MinAge)
	}

	if filters.MaxAge != nil {
		whereClauses = append(whereClauses, "age <= $"+fmt.Sprintf("%d", len(args)+1))
		args = append(args, *filters.MaxAge)
	}

	if filters.MinGenderProbability != nil {
		whereClauses = append(whereClauses, "gender_probability >= $"+fmt.Sprintf("%d", len(args)+1))
		args = append(args, *filters.MinGenderProbability)
	}

	if filters.MinCountryProbability != nil {
		whereClauses = append(whereClauses, "country_probability >= $"+fmt.Sprintf("%d", len(args)+1))
		args = append(args, *filters.MinCountryProbability)
	}

	whereClause := strings.Join(whereClauses, " AND ")

	// Get total count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM profiles WHERE %s", whereClause)
	var total int
	err := db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Build ORDER BY clause
	orderClause := "created_at ASC"
	if filters.SortBy != "" {
		orderClause = filters.SortBy + " " + strings.ToUpper(filters.Order)
	}

	// Build pagination
	offset := (filters.Page - 1) * filters.Limit

	// Build final query
	query := fmt.Sprintf(`
		SELECT id, name, gender, gender_probability, age, age_group, 
		       country_id, country_name, country_probability, created_at
		FROM profiles
		WHERE %s
		ORDER BY %s
		LIMIT $%d OFFSET $%d
	`, whereClause, orderClause, len(args)+1, len(args)+2)

	args = append(args, filters.Limit, offset)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var profiles []Profile
	for rows.Next() {
		var p Profile
		err := rows.Scan(&p.ID, &p.Name, &p.Gender, &p.GenderProbability,
			&p.Age, &p.AgeGroup, &p.CountryID, &p.CountryName,
			&p.CountryProbability, &p.CreatedAt)
		if err != nil {
			return nil, 0, err
		}
		profiles = append(profiles, p)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return profiles, total, nil
}
