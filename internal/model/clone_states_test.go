package model

import (
	"testing"
)

func TestCloneStateConstants(t *testing.T) {
	tests := []struct {
		name  string
		state StateType
		want  string
	}{
		{
			name:  "selecting clone transaction",
			state: StateSelectingCloneTransaction,
			want:  "selecting_clone_transaction",
		},
		{
			name:  "selecting clone search category",
			state: StateSelectingCloneSearchCategory,
			want:  "selecting_clone_search_category",
		},
		{
			name:  "entering clone search query",
			state: StateEnteringCloneSearchQuery,
			want:  "entering_clone_search_query",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.state) != tt.want {
				t.Errorf("state = %q, want %q", tt.state, tt.want)
			}
		})
	}
}

func TestCloneStateSessionSerialization(t *testing.T) {
	tests := []struct {
		name    string
		session UserSession
	}{
		{
			name: "selecting clone transaction",
			session: UserSession{
				State: StateSelectingCloneTransaction,
				Body:  "",
			},
		},
		{
			name: "entering clone search query with category",
			session: UserSession{
				State: StateEnteringCloneSearchQuery,
				Body:  "Grocery",
			},
		},
		{
			name: "selecting clone search category with type",
			session: UserSession{
				State: StateSelectingCloneSearchCategory,
				Body:  "expense",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Serialize
			value, err := tt.session.Value()
			if err != nil {
				t.Fatalf("Value() error = %v", err)
			}

			// Deserialize
			var scanned UserSession
			err = scanned.Scan(value)
			if err != nil {
				t.Fatalf("Scan() error = %v", err)
			}

			if scanned.State != tt.session.State {
				t.Errorf("State = %q, want %q", scanned.State, tt.session.State)
			}
			if scanned.Body != tt.session.Body {
				t.Errorf("Body = %q, want %q", scanned.Body, tt.session.Body)
			}
		})
	}
}
