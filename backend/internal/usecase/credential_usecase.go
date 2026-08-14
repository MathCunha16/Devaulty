package usecase

import (
	"bytes"
	"context"
	"devaulty-backend/internal/domain/model"
	"devaulty-backend/internal/domain/port"
	"devaulty-backend/internal/dto"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

var (
	ErrCredentialNotFound   = errors.New("credential not found")
	ErrInvalidSecretPayload = errors.New("invalid secret payload")
)

type CredentialUseCase struct {
	credentialRepo   port.CredentialRepository
	projectRepo      port.ProjectRepository
	itemTagRepo      port.ItemTagRepository
	crypto           port.Crypto
	masterKeySession port.MasterKeySession
	VaultUseCase
}

func NewCredentialUseCase(credentialRepo port.CredentialRepository, projectRepo port.ProjectRepository, itemTagRepo port.ItemTagRepository, crypto port.Crypto, masterKeySession port.MasterKeySession, vaultUseCase VaultUseCase) *CredentialUseCase {
	return &CredentialUseCase{
		credentialRepo:   credentialRepo,
		projectRepo:      projectRepo,
		itemTagRepo:      itemTagRepo,
		crypto:           crypto,
		masterKeySession: masterKeySession,
		VaultUseCase:     vaultUseCase,
	}
}

func (uc *CredentialUseCase) Create(ctx context.Context, cmd dto.CreateCredentialCommand) (*dto.CredentialView, error) {
	if err := ensureProjectExists(ctx, uc.projectRepo, cmd.ProjectID); err != nil {
		return nil, err
	}
	isRequired, err := uc.CheckIfMasterSetupIsRequired(ctx)
	if err != nil {
		return nil, fmt.Errorf("error checking if master setup is required: %w", err)
	}
	if isRequired {
		return nil, ErrMasterPasswordNotConfigured
	}
	secretKey := uc.masterKeySession.GetKey()
	if secretKey == nil {
		return nil, ErrVaultLocked
	}
	defer clear(secretKey)

	payloadBytes, err := buildSecretPayloadBytes(cmd.SecretType, cmd.Username, cmd.Password, cmd.APIKey, cmd.RawTextContent)
	if err != nil {
		return nil, err
	}
	defer clear(payloadBytes)

	credentialID := uuid.New()
	aad := uc.computeAad(cmd.ProjectID, credentialID)

	cipherText, iv, authTag, err := uc.crypto.Encrypt(payloadBytes, secretKey, aad)
	if err != nil {
		return nil, fmt.Errorf("error encrypting credential: %w", err)
	}

	credential := model.Credential{
		ID:                credentialID,
		ProjectID:         cmd.ProjectID,
		Title:             cmd.Title,
		SecretType:        cmd.SecretType,
		PayloadEncrypted:  cipherText,
		EncryptionIv:      iv,
		EncryptionAuthTag: authTag,
		Notes:             cmd.Notes,
		RelatedUrl:        cmd.RelatedURL,
		BaseEntity: model.BaseEntity{
			CreatedAt: time.Now(),
			UpdatedAt: nil,
		},
	}
	savedCredential, err := uc.credentialRepo.Save(ctx, &credential)
	if err != nil {
		return nil, fmt.Errorf("error saving credential: %w", err)
	}

	decryptedBytes, err := uc.crypto.Decrypt(savedCredential.PayloadEncrypted, savedCredential.EncryptionIv, savedCredential.EncryptionAuthTag, secretKey, aad)
	if err != nil {
		return nil, fmt.Errorf("error decrypting credential: %w", err)
	}
	defer clear(decryptedBytes)

	var payloadMap map[string]string
	if err := json.Unmarshal(decryptedBytes, &payloadMap); err != nil {
		return nil, fmt.Errorf("error unmarshalling decrypted credential payload: %w", err)
	}

	return &dto.CredentialView{
		ID:               savedCredential.ID,
		ProjectID:        savedCredential.ProjectID,
		Title:            savedCredential.Title,
		SecretType:       savedCredential.SecretType,
		DecryptedPayload: payloadMap,
		Notes:            savedCredential.Notes,
		RelatedURL:       savedCredential.RelatedUrl,
		Tags:             []dto.TagSummary{},
		CreatedAt:        savedCredential.CreatedAt,
		UpdatedAt:        savedCredential.UpdatedAt,
	}, nil
}

func (uc *CredentialUseCase) GetById(ctx context.Context, projectID, id uuid.UUID) (*dto.CredentialView, error) {
	isRequired, err := uc.CheckIfMasterSetupIsRequired(ctx)
	if err != nil {
		return nil, fmt.Errorf("error checking if master setup is required: %w", err)
	}
	if isRequired {
		return nil, ErrMasterPasswordNotConfigured
	}
	secretKey := uc.masterKeySession.GetKey()
	if secretKey == nil {
		return nil, ErrVaultLocked
	}
	defer clear(secretKey)
	if err := ensureProjectExists(ctx, uc.projectRepo, projectID); err != nil {
		return nil, err
	}

	credential, err := uc.credentialRepo.FindByIDAndProjectID(ctx, projectID, id)
	if err != nil {
		return nil, fmt.Errorf("error checking if credential exists: %w", err)
	}
	if credential == nil {
		return nil, ErrCredentialNotFound
	}

	tags, err := uc.itemTagRepo.FindTagsForItem(ctx, model.ItemTypeCredential, projectID, id)
	if err != nil {
		return nil, fmt.Errorf("error fetching tags for credential: %w", err)
	}

	aad := uc.computeAad(projectID, id)
	decryptedBytes, err := uc.crypto.Decrypt(credential.PayloadEncrypted, credential.EncryptionIv, credential.EncryptionAuthTag, secretKey, aad)
	if err != nil {
		return nil, fmt.Errorf("error decrypting credential: %w", err)
	}
	defer clear(decryptedBytes)

	var payloadMap map[string]string
	if err := json.Unmarshal(decryptedBytes, &payloadMap); err != nil {
		return nil, fmt.Errorf("error unmarshalling decrypted credential payload: %w", err)
	}

	return mapCredentialToView(credential, payloadMap, tags), nil
}

func (uc *CredentialUseCase) GetAllByProjectID(ctx context.Context, projectID uuid.UUID, page, size int) (model.Page[dto.CredentialSummaryView], error) {
	isRequired, err := uc.CheckIfMasterSetupIsRequired(ctx)
	if err != nil {
		return model.Page[dto.CredentialSummaryView]{}, fmt.Errorf("error checking if master setup is required: %w", err)
	}
	if isRequired {
		return model.Page[dto.CredentialSummaryView]{}, ErrMasterPasswordNotConfigured
	}
	secretKey := uc.masterKeySession.GetKey()
	if secretKey == nil {
		return model.Page[dto.CredentialSummaryView]{}, ErrVaultLocked
	}
	defer clear(secretKey)
	if err := ensureProjectExists(ctx, uc.projectRepo, projectID); err != nil {
		return model.Page[dto.CredentialSummaryView]{}, err
	}

	credentialPage, err := uc.credentialRepo.FindAllByProjectID(ctx, projectID, page, size)
	if err != nil {
		return model.Page[dto.CredentialSummaryView]{}, fmt.Errorf("error fetching credentials: %w", err)
	}
	if len(credentialPage.Content) == 0 {
		return model.NewPage([]dto.CredentialSummaryView{}, credentialPage.Number, credentialPage.Size, credentialPage.TotalElements), nil
	}

	credentialIDs := make([]uuid.UUID, len(credentialPage.Content))
	for i, c := range credentialPage.Content {
		credentialIDs[i] = c.ID
	}

	tagsMap, err := uc.itemTagRepo.FindTagsForItems(ctx, model.ItemTypeCredential, projectID, credentialIDs)
	if err != nil {
		return model.Page[dto.CredentialSummaryView]{}, fmt.Errorf("error fetching tags for credentials: %w", err)
	}

	summaries := make([]dto.CredentialSummaryView, len(credentialPage.Content))
	for i, c := range credentialPage.Content {
		tags := tagsMap[c.ID]
		tagSummaries := make([]dto.TagSummary, len(tags))
		for j, t := range tags {
			tagSummaries[j] = dto.TagSummary{
				ID:    t.ID,
				Name:  t.Name,
				Color: t.Color,
			}
		}
		summaries[i] = dto.CredentialSummaryView{
			ID:         c.ID,
			ProjectID:  c.ProjectID,
			Title:      c.Title,
			SecretType: c.SecretType,
			RelatedURL: c.RelatedUrl,
			Tags:       tagSummaries,
			CreatedAt:  c.CreatedAt,
			UpdatedAt:  c.UpdatedAt,
		}
	}

	return model.NewPage(summaries, credentialPage.Number, credentialPage.Size, credentialPage.TotalElements), nil
}

func (uc *CredentialUseCase) Update(ctx context.Context, cmd dto.UpdateCredentialCommand) (*dto.CredentialView, error) {
	isRequired, err := uc.CheckIfMasterSetupIsRequired(ctx)
	if err != nil {
		return nil, fmt.Errorf("error checking if master setup is required: %w", err)
	}
	if isRequired {
		return nil, ErrMasterPasswordNotConfigured
	}
	secretKey := uc.masterKeySession.GetKey()
	if secretKey == nil {
		return nil, ErrVaultLocked
	}
	defer clear(secretKey)
	if err := ensureProjectExists(ctx, uc.projectRepo, cmd.ProjectID); err != nil {
		return nil, err
	}

	credential, err := uc.credentialRepo.FindByIDAndProjectID(ctx, cmd.ProjectID, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("error checking if credential exists: %w", err)
	}
	if credential == nil {
		return nil, ErrCredentialNotFound
	}
	if cmd.Title != nil {
		credential.Title = *cmd.Title
	}
	if cmd.Notes != nil {
		credential.Notes = cmd.Notes
	}
	if cmd.RelatedURL != nil {
		credential.RelatedUrl = cmd.RelatedURL
	}
	aad := uc.computeAad(cmd.ProjectID, cmd.ID)

	secretTypeChanged := cmd.SecretType != nil && *cmd.SecretType != credential.SecretType
	secretFieldsProvided := cmd.APIKey != nil || cmd.Username != nil || cmd.Password != nil || cmd.RawTextContent != nil

	if secretTypeChanged || secretFieldsProvided {
		targetSecretType := credential.SecretType
		if cmd.SecretType != nil {
			targetSecretType = *cmd.SecretType
		}

		var mergedMap map[string]string
		if secretTypeChanged {
			mergedMap = make(map[string]string)
		} else {
			decryptedBytes, err := uc.crypto.Decrypt(credential.PayloadEncrypted, credential.EncryptionIv, credential.EncryptionAuthTag, secretKey, aad)
			if err != nil {
				return nil, fmt.Errorf("error decrypting existing credential for update: %w", err)
			}
			_ = json.Unmarshal(decryptedBytes, &mergedMap)
			clear(decryptedBytes)
			if mergedMap == nil {
				mergedMap = make(map[string]string)
			}
		}

		switch targetSecretType {
		case model.CredentialSecretTypeLogin:
			if len(cmd.Username) > 0 {
				mergedMap["username"] = string(cmd.Username)
			}
			if len(cmd.Password) > 0 {
				mergedMap["password"] = string(cmd.Password)
			}
			if len(bytes.TrimSpace([]byte(mergedMap["password"]))) == 0 {
				return nil, fmt.Errorf("%w: password is required for LOGIN secret type", ErrInvalidSecretPayload)
			}
		case model.CredentialSecretTypeApiKey:
			if len(cmd.APIKey) > 0 {
				mergedMap["apiKey"] = string(cmd.APIKey)
			}
			if len(bytes.TrimSpace([]byte(mergedMap["apiKey"]))) == 0 {
				return nil, fmt.Errorf("%w: apikey is required for APIKEY secret type", ErrInvalidSecretPayload)
			}
		case model.CredentialSecretTypeRawText:
			if len(cmd.RawTextContent) > 0 {
				mergedMap["rawText"] = string(cmd.RawTextContent)
			}
			if len(bytes.TrimSpace([]byte(mergedMap["rawText"]))) == 0 {
				return nil, fmt.Errorf("%w: raw text is required for RAWTEXT secret type", ErrInvalidSecretPayload)
			}
		default:
			return nil, ErrInvalidSecretPayload
		}

		payloadBytes, err := json.Marshal(mergedMap)
		if err != nil {
			return nil, fmt.Errorf("error marshalling updated credential payload: %w", err)
		}
		defer clear(payloadBytes)

		cipherText, iv, authTag, err := uc.crypto.Encrypt(payloadBytes, secretKey, aad)
		if err != nil {
			return nil, fmt.Errorf("error encrypting updated credential: %w", err)
		}
		credential.SecretType = targetSecretType
		credential.PayloadEncrypted = cipherText
		credential.EncryptionIv = iv
		credential.EncryptionAuthTag = authTag
	}

	now := time.Now()
	credential.UpdatedAt = &now
	saved, err := uc.credentialRepo.Save(ctx, credential)
	if err != nil {
		return nil, err
	}

	decryptedBytes, err := uc.crypto.Decrypt(saved.PayloadEncrypted, saved.EncryptionIv, saved.EncryptionAuthTag, secretKey, aad)
	if err != nil {
		return nil, fmt.Errorf("error decrypting credential: %w", err)
	}
	defer clear(decryptedBytes)

	var payloadMap map[string]string
	if err := json.Unmarshal(decryptedBytes, &payloadMap); err != nil {
		return nil, fmt.Errorf("error unmarshalling decrypted credential payload: %w", err)
	}

	tags, err := uc.itemTagRepo.FindTagsForItem(ctx, model.ItemTypeCredential, cmd.ProjectID, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("error fetching tags for credential: %w", err)
	}

	return mapCredentialToView(saved, payloadMap, tags), nil
}

func (uc *CredentialUseCase) Delete(ctx context.Context, projectID, id uuid.UUID) error {
	isRequired, err := uc.CheckIfMasterSetupIsRequired(ctx)
	if err != nil {
		return fmt.Errorf("error checking if master setup is required: %w", err)
	}
	if isRequired {
		return ErrMasterPasswordNotConfigured
	}
	secretKey := uc.masterKeySession.GetKey()
	if secretKey == nil {
		return ErrVaultLocked
	}
	defer clear(secretKey)

	if err := ensureProjectExists(ctx, uc.projectRepo, projectID); err != nil {
		return err
	}

	deleted, err := uc.credentialRepo.DeleteByIDAndProjectID(ctx, projectID, id)
	if err != nil {
		return fmt.Errorf("error deleting credential: %w", err)
	}
	if !deleted {
		return ErrCredentialNotFound
	}

	if err = uc.itemTagRepo.RemoveAllTagsFromItem(ctx, model.ItemTypeCredential, id); err != nil {
		log.Printf("warning: failed to remove tags from credential %s: %v", id, err)
	}
	return nil
}

// --- helper methods ---

func (uc *CredentialUseCase) computeAad(projectID, CredentialID uuid.UUID) []byte {
	aad := make([]byte, 0, 32)
	aad = append(aad, projectID[:]...)
	aad = append(aad, CredentialID[:]...)
	return aad
}

func buildSecretPayloadBytes(secretType model.CredentialSecretType, username, password, apikey, rawText []byte) ([]byte, error) {
	var secretMap map[string]string

	switch secretType {
	case model.CredentialSecretTypeLogin:
		if len(bytes.TrimSpace(password)) == 0 {
			return nil, fmt.Errorf("%w: password is required for LOGIN secret type", ErrInvalidSecretPayload)
		}
		secretMap = map[string]string{
			"username": string(username),
			"password": string(password),
		}
	case model.CredentialSecretTypeApiKey:
		if len(bytes.TrimSpace(apikey)) == 0 {
			return nil, fmt.Errorf("%w: apikey is required for APIKEY secret type", ErrInvalidSecretPayload)
		}
		secretMap = map[string]string{
			"apiKey": string(apikey),
		}
	case model.CredentialSecretTypeRawText:
		if len(bytes.TrimSpace(rawText)) == 0 {
			return nil, fmt.Errorf("%w: raw text is required for RAWTEXT secret type", ErrInvalidSecretPayload)
		}
		secretMap = map[string]string{
			"rawText": string(rawText),
		}
	default:
		return nil, ErrInvalidSecretPayload
	}

	payloadBytes, err := json.Marshal(secretMap)
	if err != nil {
		return nil, fmt.Errorf("error marshalling secret payload: %w", err)
	}

	return payloadBytes, nil
}

func mapCredentialToView(credential *model.Credential, decryptedPayload map[string]string, tags []model.Tag) *dto.CredentialView {
	tagSummaries := make([]dto.TagSummary, len(tags))
	for i, t := range tags {
		tagSummaries[i] = dto.TagSummary{
			ID:    t.ID,
			Name:  t.Name,
			Color: t.Color,
		}
	}
	return &dto.CredentialView{
		ID:               credential.ID,
		ProjectID:        credential.ProjectID,
		Title:            credential.Title,
		SecretType:       credential.SecretType,
		DecryptedPayload: decryptedPayload,
		Notes:            credential.Notes,
		RelatedURL:       credential.RelatedUrl,
		Tags:             tagSummaries,
		CreatedAt:        credential.CreatedAt,
		UpdatedAt:        credential.UpdatedAt,
	}
}
