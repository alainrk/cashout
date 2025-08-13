package server

import (
	"cashout/internal/db"
	"encoding/json"
	"net/http"
	"os"
)

func HandleHealthCheck(db *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := os.Getenv("HEALTH_CHECK_TOKEN")
		if token != "" {
			headerToken := r.Header.Get("Authorization")
			if headerToken != "Bearer "+token {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
		}

		err := db.Ping()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			err = json.NewEncoder(w).Encode(map[string]string{"status": "error", "db": "down"})
			if err != nil {
				http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(map[string]string{"status": "ok", "db": "up"})
		if err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}
