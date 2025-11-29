package repository

import (
	"bytes"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"

	"cashout/internal/model"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

var (
	// ErrNoCredentials is returned when attempting to authenticate a user without any registered credentials
	ErrNoCredentials = errors.New("user has no webauthn credentials")

	// ErrInvalidSession is returned when the WebAuthn session is invalid, expired, or doesn't match the user
	ErrInvalidSession = errors.New("invalid or expired webauthn session")

	// ErrTooManyCredentials is returned when a user attempts to register more credentials than allowed
	ErrTooManyCredentials = errors.New("maximum number of credentials reached")
)

// maxCredentialsPerUser limits the number of passkeys a single user can register.
// This prevents credential spam and DoS attacks. 10 is sufficient for most users
// (multiple devices, backup keys, etc.) while preventing abuse.
const maxCredentialsPerUser = 10

type WebAuthn struct {
	Repository
	webAuthn *webauthn.WebAuthn
}

// NewWebAuthn creates a new WebAuthn repository with configured instance.
//
// Required environment variables:
//   - WEBAUTHN_RP_ID: Relying Party ID (e.g., "example.com")
//   - WEBAUTHN_RP_ORIGIN: Relying Party origin (e.g., "https://example.com")
//
// Configuration details:
//   - Timeouts: 60 seconds for both registration and login (enforced)
//   - Attestation: PreferNoAttestation (standard for passkeys, no attestation verification)
//   - User Verification: Required (biometric or PIN)
//   - Resident Keys: Preferred (enables discoverable credentials)
func NewWebAuthn(repo Repository) (*WebAuthn, error) {
	rpID := os.Getenv("WEBAUTHN_RP_ID")
	if rpID == "" {
		return nil, errors.New("WEBAUTHN_RP_ID environment variable is required")
	}

	rpOrigin := os.Getenv("WEBAUTHN_RP_ORIGIN")
	if rpOrigin == "" {
		return nil, errors.New("WEBAUTHN_RP_ORIGIN environment variable is required")
	}

	wconfig := &webauthn.Config{
		RPID:          rpID,               // Domain for passkey (e.g., "example.com")
		RPDisplayName: "Cashout",          // Display name shown to users
		RPOrigins:     []string{rpOrigin}, // Allowed origins (must match exactly)
		Timeouts: webauthn.TimeoutsConfig{
			Login: webauthn.TimeoutConfig{
				Enforce:    true,                     // Strictly enforce timeout
				Timeout:    60000 * time.Millisecond, // 60 seconds for user to complete
				TimeoutUVD: 60000 * time.Millisecond, // 60 seconds with user verification
			},
			Registration: webauthn.TimeoutConfig{
				Enforce:    true,                     // Strictly enforce timeout
				Timeout:    60000 * time.Millisecond, // 60 seconds for user to complete
				TimeoutUVD: 60000 * time.Millisecond, // 60 seconds with user verification
			},
		},
		// Accept "none" attestation which is standard for passkeys.
		// Passkeys don't require attestation statement verification.
		AttestationPreference: protocol.PreferNoAttestation,
	}

	webAuthnInstance, err := webauthn.New(wconfig)
	if err != nil {
		return nil, err
	}

	return &WebAuthn{
		Repository: repo,
		webAuthn:   webAuthnInstance,
	}, nil
}

// BeginRegistration starts passkey registration ceremony
func (r *WebAuthn) BeginRegistration(user *model.User) (*protocol.CredentialCreation, string, error) {
	// Generate session ID
	sessionID, err := generateSessionID()
	if err != nil {
		return nil, "", err
	}

	// Start registration
	creation, sessionData, err := r.webAuthn.BeginRegistration(user)
	if err != nil {
		return nil, "", err
	}

	// Store session data
	session := &model.WebAuthnSession{
		ID:               sessionID,
		TgID:             user.TgID,
		Challenge:        sessionData.Challenge,
		UserVerification: string(sessionData.UserVerification),
		CeremonyType:     "registration",
		ExpiresAt:        time.Now().UTC().Add(5 * time.Minute),
	}

	if err := r.DB.CreateWebAuthnSession(session); err != nil {
		return nil, "", err
	}

	return creation, sessionID, nil
}

// FinishRegistration completes passkey registration
func (r *WebAuthn) FinishRegistration(user *model.User, sessionID string, credentialName *string, response *http.Request) error {
	// Get and validate session
	session, err := r.DB.GetWebAuthnSession(sessionID)
	if err != nil {
		return err
	}

	if !session.IsValid() || session.CeremonyType != "registration" || session.TgID != user.TgID {
		return ErrInvalidSession
	}

	// Parse the credential creation response
	var ccr protocol.CredentialCreationResponse
	if err := json.NewDecoder(response.Body).Decode(&ccr); err != nil {
		return errors.New("failed to decode credential response")
	}

	parsedResponse, err := ccr.Parse()
	if err != nil {
		return errors.New("failed to parse credential response")
	}

	// Verify the challenge using constant-time comparison to prevent timing attacks
	receivedChallenge := []byte(parsedResponse.Response.CollectedClientData.Challenge)
	expectedChallenge := []byte(session.Challenge)
	if subtle.ConstantTimeCompare(receivedChallenge, expectedChallenge) != 1 {
		r.Logger.Warnf("Registration challenge mismatch for user %d", user.TgID)
		return errors.New("challenge mismatch")
	}

	// Verify the origin
	rpOrigin := os.Getenv("WEBAUTHN_RP_ORIGIN")
	if parsedResponse.Response.CollectedClientData.Origin != rpOrigin {
		r.Logger.Warnf("Registration origin mismatch for user %d: expected %s, got %s",
			user.TgID, rpOrigin, parsedResponse.Response.CollectedClientData.Origin)
		return errors.New("origin mismatch")
	}

	// Verify the type
	if parsedResponse.Response.CollectedClientData.Type != protocol.CreateCeremony {
		r.Logger.Warnf("Registration ceremony type mismatch for user %d", user.TgID)
		return errors.New("ceremony type mismatch")
	}

	// Verify RP ID hash
	rpID := os.Getenv("WEBAUTHN_RP_ID")
	rpIDHash := sha256.Sum256([]byte(rpID))
	if !bytes.Equal(parsedResponse.Response.AttestationObject.AuthData.RPIDHash, rpIDHash[:]) {
		r.Logger.Warnf("Registration RP ID hash mismatch for user %d", user.TgID)
		return errors.New("RP ID hash mismatch")
	}

	// Verify flags
	if !parsedResponse.Response.AttestationObject.AuthData.Flags.HasAttestedCredentialData() {
		r.Logger.Warnf("Registration missing attested credential data for user %d", user.TgID)
		return errors.New("credential data missing")
	}
	if !parsedResponse.Response.AttestationObject.AuthData.Flags.HasUserPresent() {
		r.Logger.Warnf("Registration user not present for user %d", user.TgID)
		return errors.New("user not present")
	}

	// Check maximum credentials limit
	userCreds, err := r.DB.GetWebAuthnCredentialsByUser(user.TgID)
	if err != nil {
		return err
	}
	if len(userCreds) >= maxCredentialsPerUser {
		r.Logger.Warnf("User %d attempted to register credential but already has %d credentials (max: %d)",
			user.TgID, len(userCreds), maxCredentialsPerUser)
		return ErrTooManyCredentials
	}

	// Create credential object
	credentialID := parsedResponse.Response.AttestationObject.AuthData.AttData.CredentialID
	credentialPublicKey := parsedResponse.Response.AttestationObject.AuthData.AttData.CredentialPublicKey

	// Validate credential ID and public key are present
	if len(credentialID) == 0 {
		r.Logger.Warnf("Registration failed: empty credential ID for user %d", user.TgID)
		return errors.New("credential ID is empty")
	}
	if len(credentialPublicKey) == 0 {
		r.Logger.Warnf("Registration failed: empty public key for user %d", user.TgID)
		return errors.New("credential public key is empty")
	}
	// COSE public keys should be at least 32 bytes (minimum for EC2 keys)
	if len(credentialPublicKey) < 32 {
		r.Logger.Warnf("Registration failed: public key too short (%d bytes) for user %d", len(credentialPublicKey), user.TgID)
		return errors.New("credential public key too short")
	}

	// CRITICAL: Check if credential ID already exists to prevent duplicate registration
	existingCred, err := r.DB.GetWebAuthnCredential(credentialID)
	if err == nil && existingCred != nil {
		// Credential ID already exists
		if existingCred.TgID == user.TgID {
			return errors.New("credential already registered for this user")
		}
		// Even worse - credential registered to different user
		r.Logger.Warnf("Attempted to register credential %x already belonging to user %d by user %d",
			credentialID, existingCred.TgID, user.TgID)
		return errors.New("credential already exists")
	}

	credential := &webauthn.Credential{
		ID:              credentialID,
		PublicKey:       credentialPublicKey,
		AttestationType: parsedResponse.Response.AttestationObject.Format,
		Transport:       []protocol.AuthenticatorTransport{}, // Transports are optional
		Flags: webauthn.CredentialFlags{
			UserPresent:    parsedResponse.Response.AttestationObject.AuthData.Flags.HasUserPresent(),
			UserVerified:   parsedResponse.Response.AttestationObject.AuthData.Flags.HasUserVerified(),
			BackupEligible: parsedResponse.Response.AttestationObject.AuthData.Flags.HasBackupEligible(),
			BackupState:    parsedResponse.Response.AttestationObject.AuthData.Flags.HasBackupState(),
		},
		Authenticator: webauthn.Authenticator{
			AAGUID:    parsedResponse.Response.AttestationObject.AuthData.AttData.AAGUID,
			SignCount: parsedResponse.Response.AttestationObject.AuthData.Counter,
		},
	}

	// Store credential and cleanup session in a transaction
	dbCredential := model.FromWebAuthnCredential(user.TgID, credential, credentialName)
	if err := r.DB.CreateWebAuthnCredentialWithSessionCleanup(dbCredential, sessionID); err != nil {
		return err
	}

	credName := "unnamed"
	if credentialName != nil {
		credName = *credentialName
	}
	r.Logger.Infof("WebAuthn credential '%s' registered successfully for user %d (UV: %v, BE: %v, BS: %v)",
		credName, user.TgID,
		credential.Flags.UserVerified,
		credential.Flags.BackupEligible,
		credential.Flags.BackupState)
	return nil
}

// BeginLogin starts passkey authentication ceremony
func (r *WebAuthn) BeginLogin(user *model.User) (*protocol.CredentialAssertion, string, error) {
	// Check user has credentials
	if len(user.WebAuthnCredentials()) == 0 {
		return nil, "", ErrNoCredentials
	}

	// Generate session ID
	sessionID, err := generateSessionID()
	if err != nil {
		return nil, "", err
	}

	// Start login
	assertion, sessionData, err := r.webAuthn.BeginLogin(user)
	if err != nil {
		return nil, "", err
	}

	// Store session data
	session := &model.WebAuthnSession{
		ID:               sessionID,
		TgID:             user.TgID,
		Challenge:        sessionData.Challenge,
		UserVerification: string(sessionData.UserVerification),
		CeremonyType:     "authentication",
		ExpiresAt:        time.Now().UTC().Add(5 * time.Minute),
	}

	if err := r.DB.CreateWebAuthnSession(session); err != nil {
		return nil, "", err
	}

	return assertion, sessionID, nil
}

// FinishLogin completes passkey authentication
func (r *WebAuthn) FinishLogin(user *model.User, sessionID string, response *http.Request) (*model.WebAuthnCredential, error) {
	// Get and validate session
	session, err := r.DB.GetWebAuthnSession(sessionID)
	if err != nil {
		return nil, err
	}

	if !session.IsValid() || session.CeremonyType != "authentication" || session.TgID != user.TgID {
		return nil, ErrInvalidSession
	}

	// Reconstruct session data
	sessionData := webauthn.SessionData{
		Challenge:        session.Challenge,
		UserID:           user.WebAuthnID(),
		UserVerification: protocol.UserVerificationRequirement(session.UserVerification),
	}

	// Parse and validate assertion
	credential, err := r.webAuthn.FinishLogin(user, sessionData, response)
	if err != nil {
		r.Logger.Warnf("Authentication failed for user %d: %v", user.TgID, err)
		return nil, err
	}

	// Get the credential from database
	dbCred, err := r.DB.GetWebAuthnCredential(credential.ID)
	if err != nil {
		return nil, err
	}

	// CRITICAL: Verify the credential belongs to the user attempting to authenticate
	if dbCred.TgID != user.TgID {
		r.Logger.Warnf("Credential ownership mismatch: credential belongs to user %d but user %d attempted to use it",
			dbCred.TgID, user.TgID)
		return nil, errors.New("credential does not belong to user")
	}

	// Check for cloning (sign count should always increase)
	// Do this BEFORE updating the sign count
	if credential.Authenticator.SignCount > 0 && credential.Authenticator.SignCount <= dbCred.SignCount {
		dbCred.CloneWarning = true
		r.Logger.Warnf("Possible credential cloning detected for user %d (old count: %d, new count: %d)",
			user.TgID, dbCred.SignCount, credential.Authenticator.SignCount)
	}

	// Update credential (sign count, backup state, last used)
	now := time.Now().UTC()
	dbCred.SignCount = credential.Authenticator.SignCount
	dbCred.CloneWarning = credential.Authenticator.CloneWarning
	dbCred.FlagsBackupState = credential.Flags.BackupState
	dbCred.LastUsedAt = &now

	// Update credential and cleanup session in a transaction
	if err := r.DB.UpdateWebAuthnCredentialWithSessionCleanup(dbCred, sessionID); err != nil {
		return nil, err
	}

	credName := "unnamed"
	if dbCred.CredentialName != nil {
		credName = *dbCred.CredentialName
	}
	r.Logger.Infof("WebAuthn authentication successful for user %d using credential '%s' (sign count: %d, clone warning: %v)",
		user.TgID, credName, credential.Authenticator.SignCount, dbCred.CloneWarning)
	return dbCred, nil
}

// GetUserCredentials retrieves all credentials for a user
func (r *WebAuthn) GetUserCredentials(tgID int64) ([]model.WebAuthnCredential, error) {
	return r.DB.GetWebAuthnCredentialsByUser(tgID)
}

// DeleteCredential removes a credential
func (r *WebAuthn) DeleteCredential(credentialID []byte) error {
	// Get credential details before deletion for logging
	cred, err := r.DB.GetWebAuthnCredential(credentialID)
	if err != nil {
		return err
	}

	if err := r.DB.DeleteWebAuthnCredential(credentialID); err != nil {
		return err
	}

	credName := "unnamed"
	if cred.CredentialName != nil {
		credName = *cred.CredentialName
	}
	r.Logger.Infof("WebAuthn credential '%s' deleted for user %d", credName, cred.TgID)
	return nil
}

// CleanupExpiredSessions removes expired challenge sessions
func (r *WebAuthn) CleanupExpiredSessions() error {
	return r.DB.CleanupExpiredWebAuthnSessions()
}
