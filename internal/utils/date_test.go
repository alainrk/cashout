package utils

import (
	"testing"
	"time"
)

func TestParseDate(t *testing.T) {
	currentYear := time.Now().Year()

	tests := []struct {
		name    string
		input   string
		want    time.Time
		wantErr bool
	}{
		// Valid formats with different separators
		{
			name:  "dd-mm-yyyy with dash",
			input: "15-03-2024",
			want:  time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			name:  "dd/mm/yyyy with slash",
			input: "25/12/2023",
			want:  time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC),
		},
		{
			name:  "dd.mm.yyyy with dot",
			input: "01.01.2025",
			want:  time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:  "dd mm yyyy with space",
			input: "31 05 2024",
			want:  time.Date(2024, 5, 31, 0, 0, 0, 0, time.UTC),
		},
		// Two-digit year formats
		{
			name:  "dd-mm-yy (00-49 -> 20xx)",
			input: "15-03-24",
			want:  time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			name:  "dd-mm-yy (50-99 -> 19xx)",
			input: "15-03-99",
			want:  time.Date(1999, 3, 15, 0, 0, 0, 0, time.UTC),
		},
		// No year (uses current year)
		{
			name:  "dd-mm without year",
			input: "15-03",
			want:  time.Date(currentYear, 3, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			name:  "dd/mm without year",
			input: "25/12",
			want:  time.Date(currentYear, 12, 25, 0, 0, 0, 0, time.UTC),
		},
		// Single digit day/month
		{
			name:  "d-m-yyyy single digits",
			input: "5-3-2024",
			want:  time.Date(2024, 3, 5, 0, 0, 0, 0, time.UTC),
		},
		{
			name:  "d-m without year",
			input: "5-3",
			want:  time.Date(currentYear, 3, 5, 0, 0, 0, 0, time.UTC),
		},
		// February leap year tests
		{
			name:  "29 Feb leap year",
			input: "29-02-2024",
			want:  time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC),
		},
		// Invalid formats
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "invalid format",
			input:   "2024-03-15",
			wantErr: true,
		},
		{
			name:    "invalid month",
			input:   "15-13-2024",
			wantErr: true,
		},
		{
			name:    "invalid day for month",
			input:   "31-02-2024",
			wantErr: true,
		},
		{
			name:    "29 Feb non-leap year",
			input:   "29-02-2023",
			wantErr: true,
		},
		{
			name:    "invalid day",
			input:   "32-01-2024",
			wantErr: true,
		},
		{
			name:    "zero day",
			input:   "0-01-2024",
			wantErr: true,
		},
		{
			name:    "zero month",
			input:   "15-0-2024",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDate(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !got.Equal(tt.want) {
				t.Errorf("ParseDate() = %v, want %v", got, tt.want)
			}
		})
	}
}
