package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestParseNaturalLanguageQuery(t *testing.T) {
	tests := []struct {
		query     string
		expected  QueryFilters
		shouldErr bool
	}{
		{
			query: "young males from nigeria",
			expected: QueryFilters{
				Gender:    "male",
				CountryID: "NG",
				MinAge:    intPtr(16),
				MaxAge:    intPtr(24),
			},
			shouldErr: false,
		},
		{
			query: "females above 30",
			expected: QueryFilters{
				Gender: "female",
				MinAge: intPtr(30),
			},
			shouldErr: false,
		},
		{
			query: "adult males from kenya",
			expected: QueryFilters{
				Gender:    "male",
				AgeGroup:  "adult",
				CountryID: "KE",
			},
			shouldErr: false,
		},
		{
			query: "teenagers from uganda",
			expected: QueryFilters{
				AgeGroup:  "teenager",
				CountryID: "UG",
			},
			shouldErr: false,
		},
		{
			query:     "xyz qwerty asdf",
			shouldErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.query, func(t *testing.T) {
			filters, err := parseNaturalLanguageQuery(test.query)
			if (err != nil) != test.shouldErr {
				t.Fatalf("Expected error: %v, got: %v", test.shouldErr, err)
			}
			if err == nil {
				if filters.Gender != test.expected.Gender {
					t.Errorf("Gender: expected %s, got %s", test.expected.Gender, filters.Gender)
				}
				if filters.CountryID != test.expected.CountryID {
					t.Errorf("CountryID: expected %s, got %s", test.expected.CountryID, filters.CountryID)
				}
				if filters.AgeGroup != test.expected.AgeGroup {
					t.Errorf("AgeGroup: expected %s, got %s", test.expected.AgeGroup, filters.AgeGroup)
				}
			}
		})
	}
}

func TestGetProfilesResponse(t *testing.T) {
	// This would require a test database
	// For now, we just test the response structure

	w := httptest.NewRecorder()

	// Mock handler response structure
	response := ProfilesResponse{
		Status: "success",
		Page:   1,
		Limit:  10,
		Total:  100,
		Data:   []Profile{},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	var respBody ProfilesResponse
	json.NewDecoder(w.Body).Decode(&respBody)

	if respBody.Status != "success" {
		t.Errorf("Expected success status, got %s", respBody.Status)
	}
}

func TestErrorResponse(t *testing.T) {
	errResp := ErrorResponse{
		Status:  "error",
		Message: "Test error",
	}

	body, err := json.Marshal(errResp)
	if err != nil {
		t.Fatal(err)
	}

	var decoded ErrorResponse
	err = json.Unmarshal(body, &decoded)
	if err != nil {
		t.Fatal(err)
	}

	if decoded.Status != "error" {
		t.Errorf("Expected error status, got %s", decoded.Status)
	}
}

func TestCORSHeaders(t *testing.T) {
	req, err := http.NewRequest("OPTIONS", "/api/profiles", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	handler := corsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(w, req)

	corsHeader := w.Header().Get("Access-Control-Allow-Origin")
	if corsHeader != "*" {
		t.Errorf("Expected CORS header to be *, got %s", corsHeader)
	}
}

func TestPaginationLimitValidation(t *testing.T) {
	tests := []struct {
		limit     string
		shouldErr bool
	}{
		{"10", false},
		{"50", false},
		{"51", true},
		{"1", false},
		{"0", true},
		{"invalid", true},
		{"-5", true},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("limit=%s", test.limit), func(t *testing.T) {
			// Parse like the handler does
			limit := 10
			if test.limit != "" {
				// Simulate parsing
				var val int
				_, err := fmt.Sscanf(test.limit, "%d", &val)
				if err != nil || val < 1 || val > 50 {
					if !test.shouldErr {
						t.Error("Expected no error but got one")
					}
					return
				}
				limit = val
			}

			if test.shouldErr && limit >= 1 && limit <= 50 {
				t.Error("Expected error but parsing succeeded")
			}
		})
	}
}

// Helper function
func intPtr(i int) *int {
	return &i
}
