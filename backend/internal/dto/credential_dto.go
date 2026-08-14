package dto

import (
	"devaulty-backend/internal/domain/model"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type SecretBytes []byte

type CreateCredentialCommand struct {
	ProjectID      uuid.UUID                  `json:"-"`
	Title          string                     `json:"title" binding:"required,min=2,max=255"`
	SecretType     model.CredentialSecretType `json:"secretType" binding:"required"`
	Username       SecretBytes                `json:"username,omitempty"`
	Password       SecretBytes                `json:"password,omitempty"`
	APIKey         SecretBytes                `json:"apiKey,omitempty"`
	RawTextContent SecretBytes                `json:"rawTextContent,omitempty"`
	Notes          *string                    `json:"notes,omitempty"`
	RelatedURL     *string                    `json:"relatedUrl,omitempty"`
}

type UpdateCredentialCommand struct {
	ID             uuid.UUID                   `json:"-"`
	ProjectID      uuid.UUID                   `json:"-"`
	Title          *string                     `json:"title,omitempty" binding:"omitempty,min=2,max=255"`
	SecretType     *model.CredentialSecretType `json:"secretType,omitempty"`
	Username       SecretBytes                 `json:"username,omitempty"`
	Password       SecretBytes                 `json:"password,omitempty"`
	APIKey         SecretBytes                 `json:"apiKey,omitempty"`
	RawTextContent SecretBytes                 `json:"rawTextContent,omitempty"`
	Notes          *string                     `json:"notes,omitempty"`
	RelatedURL     *string                     `json:"relatedUrl,omitempty"`
}

type CredentialSummaryView struct {
	ID         uuid.UUID                  `json:"id"`
	ProjectID  uuid.UUID                  `json:"projectId"`
	Title      string                     `json:"title"`
	SecretType model.CredentialSecretType `json:"secretType"`
	RelatedURL *string                    `json:"relatedUrl"`
	Tags       []TagSummary               `json:"tags"`
	CreatedAt  time.Time                  `json:"createdAt,omitempty"`
	UpdatedAt  *time.Time                 `json:"updatedAt,omitempty"`
}

type CredentialView struct {
	ID               uuid.UUID                  `json:"id"`
	ProjectID        uuid.UUID                  `json:"projectId"`
	Title            string                     `json:"title"`
	SecretType       model.CredentialSecretType `json:"secretType"`
	DecryptedPayload map[string]string          `json:"decryptedPayload"`
	Notes            *string                    `json:"notes"`
	RelatedURL       *string                    `json:"relatedUrl"`
	Tags             []TagSummary               `json:"tags"`
	CreatedAt        time.Time                  `json:"createdAt,omitempty"`
	UpdatedAt        *time.Time                 `json:"updatedAt,omitempty"`
}

func (sb *SecretBytes) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	*sb = SecretBytes(s)
	return nil
}
