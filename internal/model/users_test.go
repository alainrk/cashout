package model

import (
	"testing"
)

func TestUserSessionValueAndScan(t *testing.T) {
	tests := []struct {
		name    string
		session UserSession
	}{
		{
			name: "normal state with body",
			session: UserSession{
				State: StateNormal,
				Body:  "test body content",
			},
		},
		{
			name: "inserting expense state",
			session: UserSession{
				State: StateInsertingExpense,
				Body:  `{"amount": 50.0}`,
			},
		},
		{
			name: "empty session",
			session: UserSession{
				State: "",
				Body:  "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test Value method
			value, err := tt.session.Value()
			if err != nil {
				t.Fatalf("UserSession.Value() error = %v", err)
			}

			// Test Scan method
			var scanned UserSession
			err = scanned.Scan(value)
			if err != nil {
				t.Fatalf("UserSession.Scan() error = %v", err)
			}

			// Compare original and scanned
			if scanned.State != tt.session.State {
				t.Errorf("Scanned State = %v, want %v", scanned.State, tt.session.State)
			}
			if scanned.Body != tt.session.Body {
				t.Errorf("Scanned Body = %v, want %v", scanned.Body, tt.session.Body)
			}
		})
	}

	// Test Scan with nil
	t.Run("scan nil value", func(t *testing.T) {
		var session UserSession
		err := session.Scan(nil)
		if err != nil {
			t.Errorf("UserSession.Scan(nil) error = %v", err)
		}
		if session.State != "" {
			t.Errorf("UserSession.Scan(nil) State = %v, want empty", session.State)
		}
	})

	// Test Scan with invalid type
	t.Run("scan invalid type", func(t *testing.T) {
		var session UserSession
		err := session.Scan("invalid")
		if err == nil {
			t.Error("UserSession.Scan() with invalid type should return error")
		}
	})
}
