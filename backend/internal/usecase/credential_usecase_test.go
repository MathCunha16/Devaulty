package usecase_test

import (
	"bytes"
	"context"
	"devaulty-backend/internal/domain/model"
	"devaulty-backend/internal/dto"
	"devaulty-backend/internal/usecase"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCredentialRepository struct {
	mock.Mock
}

func (m *MockCredentialRepository) Save(ctx context.Context, credential *model.Credential) (*model.Credential, error) {
	args := m.Called(ctx, credential)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Credential), args.Error(1)
}

func (m *MockCredentialRepository) FindByIDAndProjectID(ctx context.Context, projectID, id uuid.UUID) (*model.Credential, error) {
	args := m.Called(ctx, projectID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Credential), args.Error(1)
}

func (m *MockCredentialRepository) FindAllByProjectID(ctx context.Context, projectID uuid.UUID, page, size int) (model.Page[model.Credential], error) {
	args := m.Called(ctx, projectID, page, size)
	return args.Get(0).(model.Page[model.Credential]), args.Error(1)
}

func (m *MockCredentialRepository) DeleteByIDAndProjectID(ctx context.Context, projectID, id uuid.UUID) (bool, error) {
	args := m.Called(ctx, projectID, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockCredentialRepository) ExistsByIDAndProjectID(ctx context.Context, id, projectID uuid.UUID) (bool, error) {
	args := m.Called(ctx, id, projectID)
	return args.Bool(0), args.Error(1)
}

func (m *MockCredentialRepository) FindExistingIDsByProjectID(ctx context.Context, ids []uuid.UUID, projectID uuid.UUID) ([]uuid.UUID, error) {
	args := m.Called(ctx, ids, projectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]uuid.UUID), args.Error(1)
}

type MockCrypto struct {
	mock.Mock
}

func (m *MockCrypto) Encrypt(plainData, secretKey, aad []byte) (cipherText, iv, authTag []byte, err error) {
	args := m.Called(plainData, secretKey, aad)
	if args.Get(0) == nil {
		return nil, nil, nil, args.Error(3)
	}
	return args.Get(0).([]byte), args.Get(1).([]byte), args.Get(2).([]byte), args.Error(3)
}

func (m *MockCrypto) Decrypt(cipherText, iv, authTag, secretKey, aad []byte) (plainData []byte, err error) {
	args := m.Called(cipherText, iv, authTag, secretKey, aad)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func setupCredentialUseCaseTest() (
	*usecase.CredentialUseCase,
	*MockCredentialRepository,
	*MockProjectRepository,
	*MockItemTagRepository,
	*MockCrypto,
	*MockMasterKeySession,
	*MockAppSettingRepository,
) {
	credRepo := new(MockCredentialRepository)
	projRepo := new(MockProjectRepository)
	tagRepo := new(MockItemTagRepository)
	crypto := new(MockCrypto)
	session := new(MockMasterKeySession)
	appSettingRepo := new(MockAppSettingRepository)

	vaultUC := usecase.NewVaultUseCase(new(MockKeyDeriver), session, appSettingRepo)
	uc := usecase.NewCredentialUseCase(credRepo, projRepo, tagRepo, crypto, session, *vaultUC)

	return uc, credRepo, projRepo, tagRepo, crypto, session, appSettingRepo
}

func TestCredentialUseCase_Create(t *testing.T) {
	ctx := context.Background()
	projectID := uuid.New()

	t.Run("Create LOGIN Credential Success", func(t *testing.T) {
		uc, credRepo, projRepo, _, crypto, session, appSettingRepo := setupCredentialUseCaseTest()
		subtestMasterKey := []byte("12345678901234567890123456789012")

		cmd := dto.CreateCredentialCommand{
			ProjectID:  projectID,
			Title:      "My AWS Login",
			SecretType: model.CredentialSecretTypeLogin,
			Username:   []byte("admin"),
			Password:   []byte("secretPassword123"),
		}

		cipherText := []byte("encryptedCipherBytes")
		iv := []byte("12BytesIvData")
		authTag := []byte("16BytesAuthTag")

		payloadMap := map[string]string{"username": "admin", "password": "secretPassword123"}
		decryptedPayloadBytes, _ := json.Marshal(payloadMap)

		projRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		appSettingRepo.On("ExistsByKey", ctx, usecase.MasterPasswordHashKey).Return(true, nil)
		session.On("GetKey").Return(subtestMasterKey)

		crypto.On("Encrypt", mock.Anything, subtestMasterKey, mock.Anything).Return(cipherText, iv, authTag, nil)
		credRepo.On("Save", ctx, mock.Anything).Return(&model.Credential{
			ID:                uuid.New(),
			ProjectID:         projectID,
			Title:             cmd.Title,
			SecretType:        cmd.SecretType,
			PayloadEncrypted:  cipherText,
			EncryptionIv:      iv,
			EncryptionAuthTag: authTag,
			BaseEntity:        model.BaseEntity{CreatedAt: time.Now()},
		}, nil)
		crypto.On("Decrypt", cipherText, iv, authTag, subtestMasterKey, mock.Anything).Return(decryptedPayloadBytes, nil)

		view, err := uc.Create(ctx, cmd)

		assert.NoError(t, err)
		assert.NotNil(t, view)
		assert.Equal(t, "My AWS Login", view.Title)
		assert.Equal(t, model.CredentialSecretTypeLogin, view.SecretType)
		assert.Equal(t, "admin", view.DecryptedPayload["username"])
		assert.Equal(t, "secretPassword123", view.DecryptedPayload["password"])

		// Assert subtest key was zeroed in RAM after operation
		assert.True(t, bytes.Equal(subtestMasterKey, make([]byte, 32)))

		projRepo.AssertExpectations(t)
		crypto.AssertExpectations(t)
		credRepo.AssertExpectations(t)
	})

	t.Run("Create API_KEY Credential Success", func(t *testing.T) {
		uc, credRepo, projRepo, _, crypto, session, appSettingRepo := setupCredentialUseCaseTest()
		subtestMasterKey := []byte("12345678901234567890123456789012")

		cmd := dto.CreateCredentialCommand{
			ProjectID:  projectID,
			Title:      "Stripe API Key",
			SecretType: model.CredentialSecretTypeApiKey,
			APIKey:     []byte("sk_test_123456789"),
		}

		cipherText := []byte("encryptedCipherBytes")
		iv := []byte("12BytesIvData")
		authTag := []byte("16BytesAuthTag")

		payloadMap := map[string]string{"apiKey": "sk_test_123456789"}
		decryptedPayloadBytes, _ := json.Marshal(payloadMap)

		projRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		appSettingRepo.On("ExistsByKey", ctx, usecase.MasterPasswordHashKey).Return(true, nil)
		session.On("GetKey").Return(subtestMasterKey)

		crypto.On("Encrypt", mock.Anything, subtestMasterKey, mock.Anything).Return(cipherText, iv, authTag, nil)
		credRepo.On("Save", ctx, mock.Anything).Return(&model.Credential{
			ID:                uuid.New(),
			ProjectID:         projectID,
			Title:             cmd.Title,
			SecretType:        cmd.SecretType,
			PayloadEncrypted:  cipherText,
			EncryptionIv:      iv,
			EncryptionAuthTag: authTag,
			BaseEntity:        model.BaseEntity{CreatedAt: time.Now()},
		}, nil)
		crypto.On("Decrypt", cipherText, iv, authTag, subtestMasterKey, mock.Anything).Return(decryptedPayloadBytes, nil)

		view, err := uc.Create(ctx, cmd)

		assert.NoError(t, err)
		assert.NotNil(t, view)
		assert.Equal(t, "sk_test_123456789", view.DecryptedPayload["apiKey"])

		// Assert subtest key was zeroed in RAM after operation
		assert.True(t, bytes.Equal(subtestMasterKey, make([]byte, 32)))
	})

	t.Run("Create Failure - Project Not Found", func(t *testing.T) {
		uc, _, projRepo, _, _, _, _ := setupCredentialUseCaseTest()
		cmd := dto.CreateCredentialCommand{ProjectID: projectID, Title: "Test"}

		projRepo.On("ExistsByID", ctx, projectID).Return(false, nil)

		view, err := uc.Create(ctx, cmd)
		assert.ErrorIs(t, err, usecase.ErrProjectNotFound)
		assert.Nil(t, view)
	})

	t.Run("Create Failure - Master Password Not Configured", func(t *testing.T) {
		uc, _, projRepo, _, _, _, appSettingRepo := setupCredentialUseCaseTest()
		cmd := dto.CreateCredentialCommand{ProjectID: projectID, Title: "Test"}

		projRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		appSettingRepo.On("ExistsByKey", ctx, usecase.MasterPasswordHashKey).Return(false, nil)

		view, err := uc.Create(ctx, cmd)
		assert.ErrorIs(t, err, usecase.ErrMasterPasswordNotConfigured)
		assert.Nil(t, view)
	})

	t.Run("Create Failure - Vault Locked", func(t *testing.T) {
		uc, _, projRepo, _, _, session, appSettingRepo := setupCredentialUseCaseTest()
		cmd := dto.CreateCredentialCommand{ProjectID: projectID, Title: "Test"}

		projRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		appSettingRepo.On("ExistsByKey", ctx, usecase.MasterPasswordHashKey).Return(true, nil)
		session.On("GetKey").Return(nil)

		view, err := uc.Create(ctx, cmd)
		assert.ErrorIs(t, err, usecase.ErrVaultLocked)
		assert.Nil(t, view)
	})

	t.Run("Create Failure - Missing Password for LOGIN", func(t *testing.T) {
		uc, _, projRepo, _, _, session, appSettingRepo := setupCredentialUseCaseTest()
		subtestMasterKey := []byte("12345678901234567890123456789012")

		cmd := dto.CreateCredentialCommand{
			ProjectID:  projectID,
			Title:      "Missing Password",
			SecretType: model.CredentialSecretTypeLogin,
			Username:   []byte("admin"),
			Password:   []byte(""),
		}

		projRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		appSettingRepo.On("ExistsByKey", ctx, usecase.MasterPasswordHashKey).Return(true, nil)
		session.On("GetKey").Return(subtestMasterKey)

		view, err := uc.Create(ctx, cmd)
		assert.Error(t, err)
		assert.Nil(t, view)

		assert.True(t, bytes.Equal(subtestMasterKey, make([]byte, 32)))
	})
}

func TestCredentialUseCase_GetById(t *testing.T) {
	ctx := context.Background()
	projectID := uuid.New()
	credentialID := uuid.New()

	t.Run("GetById Success", func(t *testing.T) {
		uc, credRepo, projRepo, tagRepo, crypto, session, appSettingRepo := setupCredentialUseCaseTest()
		subtestMasterKey := []byte("12345678901234567890123456789012")

		cipherText := []byte("encryptedCipherBytes")
		iv := []byte("12BytesIvData")
		authTag := []byte("16BytesAuthTag")

		payloadMap := map[string]string{"username": "user1", "password": "pass1"}
		decryptedPayloadBytes, _ := json.Marshal(payloadMap)

		credModel := &model.Credential{
			ID:                credentialID,
			ProjectID:         projectID,
			Title:             "DB Credentials",
			SecretType:        model.CredentialSecretTypeLogin,
			PayloadEncrypted:  cipherText,
			EncryptionIv:      iv,
			EncryptionAuthTag: authTag,
			BaseEntity:        model.BaseEntity{CreatedAt: time.Now()},
		}

		tagColor := "#ff0000"
		tags := []model.Tag{{ID: uuid.New(), Name: "prod", Color: &tagColor}}

		appSettingRepo.On("ExistsByKey", ctx, usecase.MasterPasswordHashKey).Return(true, nil)
		session.On("GetKey").Return(subtestMasterKey)
		projRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		credRepo.On("FindByIDAndProjectID", ctx, projectID, credentialID).Return(credModel, nil)
		tagRepo.On("FindTagsForItem", ctx, model.ItemTypeCredential, projectID, credentialID).Return(tags, nil)
		crypto.On("Decrypt", cipherText, iv, authTag, subtestMasterKey, mock.Anything).Return(decryptedPayloadBytes, nil)

		view, err := uc.GetById(ctx, projectID, credentialID)

		assert.NoError(t, err)
		assert.NotNil(t, view)
		assert.Equal(t, credentialID, view.ID)
		assert.Equal(t, "DB Credentials", view.Title)
		assert.Len(t, view.Tags, 1)
		assert.Equal(t, "prod", view.Tags[0].Name)
		assert.Equal(t, "user1", view.DecryptedPayload["username"])

		assert.True(t, bytes.Equal(subtestMasterKey, make([]byte, 32)))
	})

	t.Run("GetById Failure - Credential Not Found", func(t *testing.T) {
		uc, credRepo, projRepo, _, _, session, appSettingRepo := setupCredentialUseCaseTest()
		subtestMasterKey := []byte("12345678901234567890123456789012")

		appSettingRepo.On("ExistsByKey", ctx, usecase.MasterPasswordHashKey).Return(true, nil)
		session.On("GetKey").Return(subtestMasterKey)
		projRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		credRepo.On("FindByIDAndProjectID", ctx, projectID, credentialID).Return(nil, nil)

		view, err := uc.GetById(ctx, projectID, credentialID)

		assert.ErrorIs(t, err, usecase.ErrCredentialNotFound)
		assert.Nil(t, view)

		assert.True(t, bytes.Equal(subtestMasterKey, make([]byte, 32)))
	})
}

func TestCredentialUseCase_GetAllByProjectID(t *testing.T) {
	ctx := context.Background()
	projectID := uuid.New()

	t.Run("GetAllByProjectID Success", func(t *testing.T) {
		uc, credRepo, projRepo, tagRepo, _, session, appSettingRepo := setupCredentialUseCaseTest()
		subtestMasterKey := []byte("12345678901234567890123456789012")

		cred1 := model.Credential{ID: uuid.New(), ProjectID: projectID, Title: "Cred 1", SecretType: model.CredentialSecretTypeLogin}
		cred2 := model.Credential{ID: uuid.New(), ProjectID: projectID, Title: "Cred 2", SecretType: model.CredentialSecretTypeApiKey}

		pageModel := model.NewPage([]model.Credential{cred1, cred2}, 1, 10, 2)

		appSettingRepo.On("ExistsByKey", ctx, usecase.MasterPasswordHashKey).Return(true, nil)
		session.On("GetKey").Return(subtestMasterKey)
		projRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		credRepo.On("FindAllByProjectID", ctx, projectID, 1, 10).Return(pageModel, nil)
		tagRepo.On("FindTagsForItems", ctx, model.ItemTypeCredential, projectID, []uuid.UUID{cred1.ID, cred2.ID}).
			Return(map[uuid.UUID][]model.Tag{}, nil)

		pageResult, err := uc.GetAllByProjectID(ctx, projectID, 1, 10)

		assert.NoError(t, err)
		assert.Equal(t, int64(2), pageResult.TotalElements)
		assert.Len(t, pageResult.Content, 2)
		assert.Equal(t, "Cred 1", pageResult.Content[0].Title)
		assert.Equal(t, "Cred 2", pageResult.Content[1].Title)

		assert.True(t, bytes.Equal(subtestMasterKey, make([]byte, 32)))
	})
}

func TestCredentialUseCase_Update(t *testing.T) {
	ctx := context.Background()
	projectID := uuid.New()
	credentialID := uuid.New()

	t.Run("Update Title Only Success", func(t *testing.T) {
		uc, credRepo, projRepo, tagRepo, crypto, session, appSettingRepo := setupCredentialUseCaseTest()
		subtestMasterKey := []byte("12345678901234567890123456789012")

		newTitle := "Updated AWS Title"
		cmd := dto.UpdateCredentialCommand{
			ID:        credentialID,
			ProjectID: projectID,
			Title:     &newTitle,
		}

		cipherText := []byte("existingCipherText")
		iv := []byte("12BytesIvData")
		authTag := []byte("16BytesAuthTag")

		existingCred := &model.Credential{
			ID:                credentialID,
			ProjectID:         projectID,
			Title:             "Old Title",
			SecretType:        model.CredentialSecretTypeLogin,
			PayloadEncrypted:  cipherText,
			EncryptionIv:      iv,
			EncryptionAuthTag: authTag,
		}

		payloadMap := map[string]string{"username": "admin", "password": "secretPassword123"}
		decryptedPayloadBytes, _ := json.Marshal(payloadMap)

		appSettingRepo.On("ExistsByKey", ctx, usecase.MasterPasswordHashKey).Return(true, nil)
		session.On("GetKey").Return(subtestMasterKey)
		projRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		credRepo.On("FindByIDAndProjectID", ctx, projectID, credentialID).Return(existingCred, nil)
		credRepo.On("Save", ctx, mock.Anything).Return(existingCred, nil)
		tagRepo.On("FindTagsForItem", ctx, model.ItemTypeCredential, projectID, credentialID).Return([]model.Tag{}, nil)
		crypto.On("Decrypt", cipherText, iv, authTag, subtestMasterKey, mock.Anything).Return(decryptedPayloadBytes, nil)

		view, err := uc.Update(ctx, cmd)

		assert.NoError(t, err)
		assert.NotNil(t, view)
		assert.Equal(t, "Updated AWS Title", view.Title)

		assert.True(t, bytes.Equal(subtestMasterKey, make([]byte, 32)))
	})
}

func TestCredentialUseCase_Delete(t *testing.T) {
	ctx := context.Background()
	projectID := uuid.New()
	credentialID := uuid.New()

	t.Run("Delete Success", func(t *testing.T) {
		uc, credRepo, projRepo, tagRepo, _, session, appSettingRepo := setupCredentialUseCaseTest()
		subtestMasterKey := []byte("12345678901234567890123456789012")

		appSettingRepo.On("ExistsByKey", ctx, usecase.MasterPasswordHashKey).Return(true, nil)
		session.On("GetKey").Return(subtestMasterKey)
		projRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		credRepo.On("DeleteByIDAndProjectID", ctx, projectID, credentialID).Return(true, nil)
		tagRepo.On("RemoveAllTagsFromItem", ctx, model.ItemTypeCredential, credentialID).Return(nil)

		err := uc.Delete(ctx, projectID, credentialID)

		assert.NoError(t, err)
		credRepo.AssertExpectations(t)
		tagRepo.AssertExpectations(t)

		assert.True(t, bytes.Equal(subtestMasterKey, make([]byte, 32)))
	})

	t.Run("Delete Failure - Credential Not Found", func(t *testing.T) {
		uc, credRepo, projRepo, _, _, session, appSettingRepo := setupCredentialUseCaseTest()
		subtestMasterKey := []byte("12345678901234567890123456789012")

		appSettingRepo.On("ExistsByKey", ctx, usecase.MasterPasswordHashKey).Return(true, nil)
		session.On("GetKey").Return(subtestMasterKey)
		projRepo.On("ExistsByID", ctx, projectID).Return(true, nil)
		credRepo.On("DeleteByIDAndProjectID", ctx, projectID, credentialID).Return(false, nil)

		err := uc.Delete(ctx, projectID, credentialID)

		assert.ErrorIs(t, err, usecase.ErrCredentialNotFound)

		assert.True(t, bytes.Equal(subtestMasterKey, make([]byte, 32)))
	})
}
