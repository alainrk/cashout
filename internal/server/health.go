
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
			json.NewEncoder(w).Encode(map[string]string{"status": "error", "db": "down"})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok", "db": "up"})
	}
}
