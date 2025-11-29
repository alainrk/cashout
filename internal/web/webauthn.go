package web

import (
	"encoding/hex"
	"encoding/json"
	"net/http"
	"regexp"
	"strings"
	"time"

	"cashout/internal/client"
	"cashout/internal/repository"
)

// isSecureRequest checks if the request is over HTTPS, accounting for reverse proxies
func isSecureRequest(r *http.Request) bool {
	// Check X-Forwarded-Proto header (set by reverse proxies)
	if proto := r.Header.Get("X-Forwarded-Proto"); proto == "https" {
		return true
	}
	// Fall back to checking TLS directly
	return r.TLS != nil
}

// Session ID should be 64 hex characters (32 bytes encoded as hex)
var sessionIDPattern = regexp.MustCompile(`^[0-9a-f]{64}$`)

// isValidSessionID validates the format of a session ID to prevent injection attacks
func isValidSessionID(sessionID string) bool {
	return sessionIDPattern.MatchString(sessionID)
}

// handlePasskeyCheck checks if user has passkeys registered
func (s *Server) handlePasskeyCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendJSONError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	email := strings.TrimSpace(strings.ToLower(req.Email))
	if email == "" {
		s.sendJSONError(w, "Email is required", http.StatusBadRequest)
		return
	}

	// Get user
	user, exists, err := s.repositories.Users.GetByEmail(email)
	if err != nil || !exists {
		// Don't reveal if user exists for security
		s.sendJSONSuccess(w, map[string]any{
			"hasPasskey": false,
		})
		return
	}

	// Check for credentials
	creds, err := s.repositories.WebAuthn.GetUserCredentials(user.TgID)
	if err != nil {
		s.logger.Errorf("Failed to get credentials: %v", err)
		s.sendJSONSuccess(w, map[string]any{
			"hasPasskey": false,
		})
		return
	}

	s.sendJSONSuccess(w, map[string]any{
		"hasPasskey": len(creds) > 0,
	})
}

// handlePasskeyBeginLogin starts passkey authentication
func (s *Server) handlePasskeyBeginLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendJSONError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	email := strings.TrimSpace(strings.ToLower(req.Email))
	if email == "" {
		s.sendJSONError(w, "Email is required", http.StatusBadRequest)
		return
	}

	// Get user with credentials
	// Use generic error to prevent email enumeration attacks
	user, exists, err := s.repositories.Users.GetByEmail(email)
	if err != nil || !exists {
		// Don't reveal if user exists - same as handlePasskeyCheck
		s.sendJSONError(w, "Authentication failed", http.StatusUnauthorized)
		return
	}

	// Load credentials
	userWithCreds, err := s.repositories.WebAuthn.DB.GetUserWithWebAuthnCredentials(user.TgID)
	if err != nil {
		s.sendJSONError(w, "Failed to load user data", http.StatusInternalServerError)
		return
	}

	// Begin login
	assertion, sessionID, err := s.repositories.WebAuthn.BeginLogin(userWithCreds)
	if err != nil {
		if err == repository.ErrNoCredentials {
			s.sendJSONError(w, "No passkeys registered", http.StatusBadRequest)
			return
		}
		s.logger.Errorf("Failed to begin passkey login: %v", err)
		s.sendJSONError(w, "Failed to start passkey login", http.StatusInternalServerError)
		return
	}

	// Store session ID in cookie for finish step
	http.SetCookie(w, &http.Cookie{
		Name:     "webauthn_session",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   isSecureRequest(r),
		SameSite: http.SameSiteLaxMode,
		MaxAge:   300, // 5 minutes
	})

	s.sendJSONSuccess(w, map[string]any{
		"options": assertion,
	})
}

// handlePasskeyFinishLogin completes passkey authentication
func (s *Server) handlePasskeyFinishLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get session ID from cookie
	cookie, err := r.Cookie("webauthn_session")
	if err != nil {
		s.sendJSONError(w, "No active passkey session", http.StatusBadRequest)
		return
	}
	sessionID := cookie.Value

	// Validate session ID format to prevent injection attacks
	if !isValidSessionID(sessionID) {
		s.sendJSONError(w, "Invalid session", http.StatusBadRequest)
		return
	}

	// Get session to find user
	waSession, err := s.repositories.WebAuthn.DB.GetWebAuthnSession(sessionID)
	if err != nil {
		s.sendJSONError(w, "Invalid session", http.StatusBadRequest)
		return
	}

	// Get user with credentials
	user, err := s.repositories.WebAuthn.DB.GetUserWithWebAuthnCredentials(waSession.TgID)
	if err != nil {
		s.sendJSONError(w, "User not found", http.StatusNotFound)
		return
	}

	// Finish login using secure library method
	_, err = s.repositories.WebAuthn.FinishLogin(user, sessionID, r)
	if err != nil {
		s.logger.Errorf("Failed to finish passkey login: %v", err)
		s.sendJSONError(w, "Authentication failed", http.StatusUnauthorized)
		return
	}

	// Create web session (same as code auth)
	session, err := s.repositories.Auth.CreateWebSession(user.TgID)
	if err != nil {
		s.logger.Errorf("Failed to create session: %v", err)
		s.sendJSONError(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	// Clear webauthn session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "webauthn_session",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})

	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    session.ID,
		Path:     "/",
		HttpOnly: true,
		Secure:   isSecureRequest(r),
		SameSite: http.SameSiteLaxMode,
		MaxAge:   86400, // 24 hours
	})

	s.sendJSONSuccess(w, map[string]any{
		"message":  "Login successful",
		"redirect": basePath + "/dashboard",
	})
}

// handlePasskeyBeginRegister starts passkey registration
// (Protected endpoint - requires existing auth session)
func (s *Server) handlePasskeyBeginRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get current user from context (must be authenticated)
	user := client.GetUserFromContext(r.Context())
	if user == nil {
		s.sendJSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Verify user has email
	if user.Email == nil || *user.Email == "" {
		s.sendJSONError(w, "Email required for passkey registration", http.StatusBadRequest)
		return
	}

	// Load existing credentials
	userWithCreds, err := s.repositories.WebAuthn.DB.GetUserWithWebAuthnCredentials(user.TgID)
	if err != nil {
		s.sendJSONError(w, "Failed to load user data", http.StatusInternalServerError)
		return
	}

	// Begin registration
	creation, sessionID, err := s.repositories.WebAuthn.BeginRegistration(userWithCreds)
	if err != nil {
		s.logger.Errorf("Failed to begin passkey registration: %v", err)
		s.sendJSONError(w, "Failed to start passkey registration", http.StatusInternalServerError)
		return
	}

	// Store session ID in cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "webauthn_session",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   isSecureRequest(r),
		SameSite: http.SameSiteLaxMode,
		MaxAge:   300, // 5 minutes
	})

	s.sendJSONSuccess(w, map[string]any{
		"options": creation,
	})
}

// handlePasskeyFinishRegister completes passkey registration
func (s *Server) handlePasskeyFinishRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get current user from context
	user := client.GetUserFromContext(r.Context())
	if user == nil {
		s.sendJSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get webauthn session ID from cookie
	cookie, err := r.Cookie("webauthn_session")
	if err != nil {
		s.sendJSONError(w, "No active passkey session", http.StatusBadRequest)
		return
	}
	sessionID := cookie.Value

	// Validate session ID format to prevent injection attacks
	if !isValidSessionID(sessionID) {
		s.sendJSONError(w, "Invalid session", http.StatusBadRequest)
		return
	}

	// Extract and validate credential name from header
	var credName *string
	if name := r.Header.Get("X-Credential-Name"); name != "" {
		// Validate credential name length and sanitize
		name = strings.TrimSpace(name)
		if len(name) > 100 {
			s.sendJSONError(w, "Credential name too long (max 100 characters)", http.StatusBadRequest)
			return
		}
		if len(name) > 0 {
			credName = &name
		}
	}

	// Finish registration using secure verification method
	// The request body contains the credential directly
	err = s.repositories.WebAuthn.FinishRegistration(user, sessionID, credName, r)
	if err != nil {
		s.logger.Errorf("Failed to finish passkey registration: %v", err)
		s.sendJSONError(w, "Registration failed", http.StatusBadRequest)
		return
	}

	// Clear webauthn session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "webauthn_session",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})

	s.sendJSONSuccess(w, map[string]any{
		"message": "Passkey registered successfully",
	})
}

// handlePasskeyList lists user's passkeys
func (s *Server) handlePasskeyList(w http.ResponseWriter, r *http.Request) {
	user := client.GetUserFromContext(r.Context())
	if user == nil {
		s.sendJSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	creds, err := s.repositories.WebAuthn.GetUserCredentials(user.TgID)
	if err != nil {
		s.logger.Errorf("Failed to get credentials: %v", err)
		s.sendJSONError(w, "Failed to load passkeys", http.StatusInternalServerError)
		return
	}

	// Format for response
	type passkeyInfo struct {
		ID         string     `json:"id"`
		Name       string     `json:"name"`
		CreatedAt  time.Time  `json:"createdAt"`
		LastUsedAt *time.Time `json:"lastUsedAt"`
	}

	passkeys := make([]passkeyInfo, len(creds))
	for i, cred := range creds {
		name := "Unnamed Passkey"
		if cred.CredentialName != nil {
			name = *cred.CredentialName
		}

		passkeys[i] = passkeyInfo{
			ID:         hex.EncodeToString(cred.ID),
			Name:       name,
			CreatedAt:  cred.CreatedAt,
			LastUsedAt: cred.LastUsedAt,
		}
	}

	s.sendJSONSuccess(w, map[string]any{
		"passkeys": passkeys,
	})
}

// handlePasskeyDelete deletes a passkey
func (s *Server) handlePasskeyDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := client.GetUserFromContext(r.Context())
	if user == nil {
		s.sendJSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		CredentialID string `json:"credentialId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendJSONError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	credID, err := hex.DecodeString(req.CredentialID)
	if err != nil {
		s.sendJSONError(w, "Invalid credential ID", http.StatusBadRequest)
		return
	}

	// Verify ownership
	cred, err := s.repositories.WebAuthn.DB.GetWebAuthnCredential(credID)
	if err != nil {
		s.sendJSONError(w, "Credential not found", http.StatusNotFound)
		return
	}

	if cred.TgID != user.TgID {
		s.sendJSONError(w, "Unauthorized", http.StatusForbidden)
		return
	}

	// Delete
	if err := s.repositories.WebAuthn.DeleteCredential(credID); err != nil {
		s.logger.Errorf("Failed to delete credential: %v", err)
		s.sendJSONError(w, "Failed to delete passkey", http.StatusInternalServerError)
		return
	}

	s.sendJSONSuccess(w, map[string]any{
		"message": "Passkey deleted successfully",
	})
}
