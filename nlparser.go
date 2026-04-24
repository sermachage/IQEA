package main

import (
	"regexp"
	"sort"
	"strings"
)

// Country mapping
var countryMap = map[string]string{
	"nigeria":                      "NG",
	"ghana":                        "GH",
	"kenya":                        "KE",
	"uganda":                       "UG",
	"tanzania":                     "TZ",
	"cameroon":                     "CM",
	"senegal":                      "SN",
	"mali":                         "ML",
	"ivory coast":                  "CI",
	"benin":                        "BJ",
	"burkina faso":                 "BF",
	"burundi":                      "BI",
	"cape verde":                   "CV",
	"central african republic":     "CF",
	"chad":                         "TD",
	"comoros":                      "KM",
	"congo":                        "CG",
	"democratic republic of congo": "CD",
	"djibouti":                     "DJ",
	"egypt":                        "EG",
	"equatorial guinea":            "GQ",
	"eritrea":                      "ER",
	"eswatini":                     "SZ",
	"ethiopia":                     "ET",
	"gabon":                        "GA",
	"gambia":                       "GM",
	"guinea":                       "GN",
	"guinea-bissau":                "GW",
	"lesotho":                      "LS",
	"liberia":                      "LR",
	"libya":                        "LY",
	"madagascar":                   "MG",
	"malawi":                       "MW",
	"mauritania":                   "MR",
	"mauritius":                    "MU",
	"morocco":                      "MA",
	"mozambique":                   "MZ",
	"namibia":                      "NA",
	"niger":                        "NE",
	"rwanda":                       "RW",
	"são tomé and príncipe":        "ST",
	"seychelles":                   "SC",
	"sierra leone":                 "SL",
	"somalia":                      "SO",
	"south africa":                 "ZA",
	"south sudan":                  "SS",
	"sudan":                        "SD",
	"togo":                         "TG",
	"tunisia":                      "TN",
	"zambia":                       "ZM",
	"zimbabwe":                     "ZW",
	"angola":                       "AO",
}

// Age group classifications
var ageGroupRanges = map[string][2]int{
	"child":    {0, 12},
	"teenager": {13, 19},
	"adult":    {20, 59},
	"senior":   {60, 150},
}

// Age qualifiers
var ageQualifiers = map[string][2]int{
	"young":  {16, 24},
	"middle": {25, 50},
	"old":    {50, 150},
}

func parseNaturalLanguageQuery(query string) (*QueryFilters, error) {
	query = strings.ToLower(strings.TrimSpace(query))

	filters := &QueryFilters{
		Page:  1,
		Limit: 10,
		Order: "asc",
	}

	// Check for gender keywords
	hasMale := strings.Contains(query, "male") || strings.Contains(query, "man") || strings.Contains(query, "men")
	hasFemale := strings.Contains(query, "female") || strings.Contains(query, "woman") || strings.Contains(query, "women")

	// Handle gender logic
	if hasMale && hasFemale {
		// Both genders mentioned - check if "and" is used
		if strings.Contains(query, " and ") {
			// "male and female" - check order and context
			// If both are present with "and", we can't filter by gender
			// but we can still filter by other criteria
		} else {
			// No explicit "and", so prioritize female
			filters.Gender = "female"
		}
	} else if hasFemale {
		filters.Gender = "female"
	} else if hasMale {
		filters.Gender = "male"
	}

	// Extract age group
	for ageGroup := range ageGroupRanges {
		if strings.Contains(query, ageGroup) {
			filters.AgeGroup = ageGroup
			break
		}
	}

	// Extract age qualifiers (like "young", "middle", "old")
	for qualifier, ageRange := range ageQualifiers {
		if strings.Contains(query, qualifier) {
			filters.MinAge = &ageRange[0]
			filters.MaxAge = &ageRange[1]
			break
		}
	}

	// Extract age constraints
	if strings.Contains(query, "above") || strings.Contains(query, "over") {
		// Extract numeric value after "above" or "over"
		age := extractNumericAge(query, "above|over")
		if age > 0 {
			filters.MinAge = &age
		}
	} else if strings.Contains(query, "below") || strings.Contains(query, "under") {
		// Extract numeric value after "below" or "under"
		age := extractNumericAge(query, "below|under")
		if age > 0 {
			filters.MaxAge = &age
		}
	}

	// Extract country
	// Match longer country names first (e.g., "nigeria" before "niger")
	countryNames := make([]string, 0, len(countryMap))
	for countryName := range countryMap {
		countryNames = append(countryNames, countryName)
	}
	sort.Slice(countryNames, func(i, j int) bool {
		return len(countryNames[i]) > len(countryNames[j])
	})

	for _, countryName := range countryNames {
		pattern := `\b` + regexp.QuoteMeta(countryName) + `\b`
		if matched, _ := regexp.MatchString(pattern, query); matched {
			filters.CountryID = countryMap[countryName]
			break
		}
	}

	// Check if we parsed anything meaningful
	if filters.Gender == "" && filters.AgeGroup == "" && filters.CountryID == "" &&
		filters.MinAge == nil && filters.MaxAge == nil {
		return nil, ErrCannotParse
	}

	return filters, nil
}

func extractNumericAge(query string, pattern string) int {
	// Simple extraction: look for digits after age keywords
	parts := strings.Fields(query)
	for i, part := range parts {
		if strings.Contains(pattern, strings.ToLower(part)) && i+1 < len(parts) {
			// Try to parse the next word as a number
			nextPart := parts[i+1]
			// Remove non-numeric characters
			numStr := ""
			for _, ch := range nextPart {
				if ch >= '0' && ch <= '9' {
					numStr += string(ch)
				}
			}
			if numStr != "" {
				age := 0
				for _, ch := range numStr {
					age = age*10 + int(ch-'0')
				}
				if age > 0 && age <= 150 {
					return age
				}
			}
		}
	}
	return 0
}

var ErrCannotParse = &cannotParseError{}

type cannotParseError struct{}

func (e *cannotParseError) Error() string {
	return "cannot parse query"
}
