package model

import (
	"encoding/binary"
	"time"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/lib/pq"
)

// WebAuthnCredential represents a stored passkey credential
type WebAuthnCredential struct {
	ID                  []byte         `gorm:"column:id;primaryKey"`
	TgID                int64          `gorm:"column:tg_id;not null;index"`
	PublicKey           []byte         `gorm:"column:public_key;not null"`
	AttestationType     string         `gorm:"column:attestation_type;not null"`
	AAGUID              []byte         `gorm:"column:aaguid;not null"`
	SignCount           uint32         `gorm:"column:sign_count;not null;default:0"`
	CloneWarning        bool           `gorm:"column:clone_warning;not null;default:false"`
	FlagsUserPresent    bool           `gorm:"column:flags_user_present;not null"`
	FlagsUserVerified   bool           `gorm:"column:flags_user_verified;not null"`
	FlagsBackupEligible bool           `gorm:"column:flags_backup_eligible;not null"`
	FlagsBackupState    bool           `gorm:"column:flags_backup_state;not null"`
	Transport           pq.StringArray `gorm:"column:transport;type:text[];not null;default:'{}'"`
	CredentialName      *string        `gorm:"column:credential_name"`
	LastUsedAt          *time.Time     `gorm:"column:last_used_at"`
	CreatedAt           time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt           time.Time      `gorm:"column:updated_at;autoUpdateTime"`

	// Association
	User *User `gorm:"foreignKey:TgID;references:TgID"`
}

// TableName overrides the table name
func (WebAuthnCredential) TableName() string {
	return "webauthn_credentials"
}

// ToWebAuthnCredential converts DB model to go-webauthn library credential
func (c *WebAuthnCredential) ToWebAuthnCredential() webauthn.Credential {
	// Convert transport strings to protocol.AuthenticatorTransport
	transports := make([]protocol.AuthenticatorTransport, len(c.Transport))
	for i, t := range c.Transport {
		transports[i] = protocol.AuthenticatorTransport(t)
	}

	// Build flags from stored boolean values
	flags := webauthn.CredentialFlags{
		UserPresent:    c.FlagsUserPresent,
		UserVerified:   c.FlagsUserVerified,
		BackupEligible: c.FlagsBackupEligible,
		BackupState:    c.FlagsBackupState,
	}

	return webauthn.Credential{
		ID:              c.ID,
		PublicKey:       c.PublicKey,
		AttestationType: c.AttestationType,
		Authenticator: webauthn.Authenticator{
			AAGUID:       c.AAGUID,
			SignCount:    c.SignCount,
			CloneWarning: c.CloneWarning,
		},
		Flags:     flags,
		Transport: transports,
	}
}

// FromWebAuthnCredential creates DB model from library credential
func FromWebAuthnCredential(tgID int64, cred *webauthn.Credential, name *string) *WebAuthnCredential {
	// Convert transports to pq.StringArray (properly handles empty arrays in PostgreSQL)
	transports := make(pq.StringArray, 0, len(cred.Transport))
	for _, t := range cred.Transport {
		transports = append(transports, string(t))
	}

	return &WebAuthnCredential{
		ID:                  cred.ID,
		TgID:                tgID,
		PublicKey:           cred.PublicKey,
		AttestationType:     cred.AttestationType,
		AAGUID:              cred.Authenticator.AAGUID,
		SignCount:           cred.Authenticator.SignCount,
		CloneWarning:        cred.Authenticator.CloneWarning,
		FlagsUserPresent:    cred.Flags.UserPresent,
		FlagsUserVerified:   cred.Flags.UserVerified,
		FlagsBackupEligible: cred.Flags.BackupEligible,
		FlagsBackupState:    cred.Flags.BackupState,
		Transport:           transports,
		CredentialName:      name,
	}
}

// WebAuthnSession stores challenge data during ceremonies
type WebAuthnSession struct {
	ID               string    `gorm:"column:id;primaryKey"`
	TgID             int64     `gorm:"column:tg_id;not null;index"`
	Challenge        string    `gorm:"column:challenge;not null"`
	UserVerification string    `gorm:"column:user_verification;not null"`
	CeremonyType     string    `gorm:"column:ceremony_type;not null"` // "registration" or "authentication"
	ExpiresAt        time.Time `gorm:"column:expires_at;not null;index"`
	CreatedAt        time.Time `gorm:"column:created_at;autoCreateTime"`

	// Association
	User *User `gorm:"foreignKey:TgID;references:TgID"`
}

// TableName overrides the table name
func (WebAuthnSession) TableName() string {
	return "webauthn_sessions"
}

// IsValid checks if the session is still valid
func (s *WebAuthnSession) IsValid() bool {
	return time.Now().UTC().Before(s.ExpiresAt)
}

// WebAuthn interface implementation for User model
// These methods must be added to the User type in users.go

// WebAuthnID returns the user's ID as bytes
func (u *User) WebAuthnID() []byte {
	// Use TgID as unique identifier (convert int64 to []byte)
	id := make([]byte, 8)
	binary.BigEndian.PutUint64(id, uint64(u.TgID))
	return id
}

// WebAuthnName returns the user's name for WebAuthn (uses email)
func (u *User) WebAuthnName() string {
	// Use email as account name for passkeys
	if u.Email != nil && *u.Email != "" {
		return *u.Email
	}
	// Fallback to username
	return u.TgUsername
}

// WebAuthnDisplayName returns the user's display name
func (u *User) WebAuthnDisplayName() string {
	// Use friendly name or username
	if u.Name != "" {
		return u.Name
	}
	return u.TgUsername
}

// WebAuthnCredentials returns the user's credentials for WebAuthn
func (u *User) WebAuthnCredentials() []webauthn.Credential {
	// This will be populated by repository layer
	// when credentials are loaded from database
	if u.Credentials == nil {
		return []webauthn.Credential{}
	}

	creds := make([]webauthn.Credential, len(u.Credentials))
	for i, c := range u.Credentials {
		creds[i] = c.ToWebAuthnCredential()
	}
	return creds
}

// WebAuthnIcon returns the user's icon URL (optional)
func (u *User) WebAuthnIcon() string {
	// No icon for now
	return ""
}
