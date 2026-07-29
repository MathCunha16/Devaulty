package model

import "github.com/google/uuid"

type CredentialSecretType string

const (
	CredentialSecretTypeLogin   CredentialSecretType = "LOGIN"
	CredentialSecretTypeApiKey  CredentialSecretType = "API_KEY"
	CredentialSecretTypeRawText CredentialSecretType = "RAW_TEXT"
)

type Credential struct {
	ID                uuid.UUID            `json:"id" db:"id"`
	ProjectID         uuid.UUID            `json:"projectId" db:"project_id"`
	Title             string               `json:"title" db:"title"`
	SecretType        CredentialSecretType `json:"secretType" db:"secret_type"`
	PayloadEncrypted  []byte               `json:"payloadEncrypted" db:"payload_encrypted"`
	EncryptionIv      []byte               `json:"encryptionIv" db:"encryption_iv"`
	EncryptionAuthTag []byte               `json:"encryptionAuthTag" db:"encryption_auth_tag"`
	Notes             *string              `json:"notes,omitempty" db:"notes"`
	RelatedUrl        *string              `json:"relatedUrl,omitempty" db:"related_url"`
	BaseEntity
}
