package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ParseDate attempts to parse a date string in various formats.
// It returns a time.Time object and an error.
// Supported formats include:
// - dd mm yyyy (with various separators like space, -, /, or .)
// - dd mm yy (two-digit year, same separators)
// - dd mm (uses current year, same separators)
// - d m (single digit day/month, uses current year, same separators)
// - Any combination of the above (d-mm-yyyy, dd/m/yy, etc.)
func ParseDate(dateStr string) (time.Time, error) {
	// Trim spaces and normalize the string
	dateStr = strings.TrimSpace(dateStr)

	// If empty, return error
	if dateStr == "" {
		return time.Time{}, fmt.Errorf("empty date string")
	}

	// Define regex patterns for different date components
	// This covers formats with various separators (space, -, /, .)
	pattern := regexp.MustCompile(`^(\d{1,2})[\s\-\/\.](\d{1,2})(?:[\s\-\/\.](\d{2}|\d{4}))?$`)

	matches := pattern.FindStringSubmatch(dateStr)
	if matches == nil {
		return time.Time{}, fmt.Errorf("invalid date format: %s", dateStr)
	}

	// Parse day and month
	day, err := strconv.Atoi(matches[1])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid day: %s", matches[1])
	}

	month, err := strconv.Atoi(matches[2])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid month: %s", matches[2])
	}

	// Validate day and month
	if month < 1 || month > 12 {
		return time.Time{}, fmt.Errorf("invalid month: %d (must be between 1 and 12)", month)
	}

	// Check if year is present
	var year int
	currentYear := time.Now().Year()

	if len(matches) > 3 && matches[3] != "" {
		year, err = strconv.Atoi(matches[3])
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid year: %s", matches[3])
		}

		// Handle 2-digit year
		if year < 100 {
			// Assume 20xx for years 00-49, 19xx for years 50-99
			if year < 50 {
				year += 2000
			} else {
				year += 1900
			}
		}
	} else {
		// If year is not provided, use current year
		year = currentYear
	}

	// Validate day with respect to month and year
	maxDay := 31
	switch month {
	case 4, 6, 9, 11:
		maxDay = 30
	case 2:
		// Check for leap year
		if year%400 == 0 || (year%4 == 0 && year%100 != 0) {
			maxDay = 29
		} else {
			maxDay = 28
		}
	}

	if day < 1 || day > maxDay {
		return time.Time{}, fmt.Errorf("invalid day: %d for month %d (must be between 1 and %d)", day, month, maxDay)
	}

	// Create the time.Time object
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC), nil
}
