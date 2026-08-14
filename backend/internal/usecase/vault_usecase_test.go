package usecase_test

import (
	"context"
	"encoding/base64"
	"errors"
	"testing"
	"time"

	"devaulty-backend/internal/domain/model"
	"devaulty-backend/internal/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- MOCK DEFINITIONS ---

type MockKeyDeriver struct {
	mock.Mock
}

func (m *MockKeyDeriver) GenerateSalt(size int) ([]byte, error) {
	args := m.Called(size)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockKeyDeriver) DeriveKey(password, salt []byte) ([]byte, error) {
	args := m.Called(password, salt)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockKeyDeriver) HashPassword(password, salt []byte) ([]byte, error) {
	args := m.Called(password, salt)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockKeyDeriver) VerifyPassword(password, salt, expectedHash []byte) bool {
	args := m.Called(password, salt, expectedHash)
	return args.Bool(0)
}

type MockMasterKeySession struct {
	mock.Mock
}

func (m *MockMasterKeySession) SetKey(key []byte) {
	m.Called(key)
}

func (m *MockMasterKeySession) GetKey() []byte {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).([]byte)
}

func (m *MockMasterKeySession) HasKey() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockMasterKeySession) Clear() {
	m.Called()
}

func (m *MockMasterKeySession) Touch() {
	m.Called()
}

func (m *MockMasterKeySession) GetSecondsRemaining(timeout time.Duration) int64 {
	args := m.Called(timeout)
	return args.Get(0).(int64)
}

type MockAppSettingRepository struct {
	mock.Mock
}

func (m *MockAppSettingRepository) Save(ctx context.Context, setting *model.AppSetting) (*model.AppSetting, error) {
	args := m.Called(ctx, setting)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.AppSetting), args.Error(1)
}

func (m *MockAppSettingRepository) FindByKey(ctx context.Context, key string) (*model.AppSetting, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.AppSetting), args.Error(1)
}

func (m *MockAppSettingRepository) ExistsByKey(ctx context.Context, key string) (bool, error) {
	args := m.Called(ctx, key)
	return args.Bool(0), args.Error(1)
}

func (m *MockAppSettingRepository) SaveMasterPasswordSettings(ctx context.Context, hashValue, saltValue string) error {
	args := m.Called(ctx, hashValue, saltValue)
	return args.Error(0)
}

// --- UNIT TESTS ---

func TestVaultUseCase_SetupMasterPassword(t *testing.T) {
	ctx := context.Background()

	t.Run("SetupMasterPassword_Success", func(t *testing.T) {
		mockKeyDeriver := new(MockKeyDeriver)
		mockSession := new(MockMasterKeySession)
		mockAppSettingRepo := new(MockAppSettingRepository)
		uc := usecase.NewVaultUseCase(mockKeyDeriver, mockSession, mockAppSettingRepo)

		password := []byte("MySuperSecretPassword123")
		saltBytes := []byte("0123456789abcdef")
		hashedPassword := []byte("32BytesHashedPasswordOutputBytes")
		secretKey := []byte("32BytesDerivedKeyAES256SecretBytes")

		hashBase64 := base64.StdEncoding.EncodeToString(hashedPassword)
		saltBase64 := base64.StdEncoding.EncodeToString(saltBytes)

		mockKeyDeriver.On("GenerateSalt", usecase.SaltLength).Return(saltBytes, nil)
		mockKeyDeriver.On("HashPassword", mock.Anything, mock.Anything).Return(hashedPassword, nil)
		mockAppSettingRepo.On("SaveMasterPasswordSettings", ctx, hashBase64, saltBase64).Return(nil)
		mockKeyDeriver.On("DeriveKey", mock.Anything, mock.Anything).Return(secretKey, nil)
		mockSession.On("SetKey", secretKey).Return()

		err := uc.SetupMasterPassword(ctx, password)

		assert.NoError(t, err)
		mockAppSettingRepo.AssertExpectations(t)
		mockKeyDeriver.AssertExpectations(t)
		mockSession.AssertExpectations(t)
	})

	t.Run("SetupMasterPassword_AlreadyConfigured", func(t *testing.T) {
		mockKeyDeriver := new(MockKeyDeriver)
		mockSession := new(MockMasterKeySession)
		mockAppSettingRepo := new(MockAppSettingRepository)
		uc := usecase.NewVaultUseCase(mockKeyDeriver, mockSession, mockAppSettingRepo)

		password := []byte("MySuperSecretPassword123")
		saltBytes := []byte("0123456789abcdef")
		hashedPassword := []byte("32BytesHashedPasswordOutputBytes")

		hashBase64 := base64.StdEncoding.EncodeToString(hashedPassword)
		saltBase64 := base64.StdEncoding.EncodeToString(saltBytes)

		mockKeyDeriver.On("GenerateSalt", usecase.SaltLength).Return(saltBytes, nil)
		mockKeyDeriver.On("HashPassword", mock.Anything, mock.Anything).Return(hashedPassword, nil)
		mockAppSettingRepo.On("SaveMasterPasswordSettings", ctx, hashBase64, saltBase64).Return(usecase.ErrMasterPasswordAlreadyConfigured)

		err := uc.SetupMasterPassword(ctx, password)

		assert.ErrorIs(t, err, usecase.ErrMasterPasswordAlreadyConfigured)
		mockAppSettingRepo.AssertExpectations(t)
		mockKeyDeriver.AssertNotCalled(t, "DeriveKey", mock.Anything, mock.Anything)
	})

	t.Run("SetupMasterPassword_GenerateSaltError", func(t *testing.T) {
		mockKeyDeriver := new(MockKeyDeriver)
		mockSession := new(MockMasterKeySession)
		mockAppSettingRepo := new(MockAppSettingRepository)
		uc := usecase.NewVaultUseCase(mockKeyDeriver, mockSession, mockAppSettingRepo)

		password := []byte("MySuperSecretPassword123")
		mockKeyDeriver.On("GenerateSalt", usecase.SaltLength).Return(nil, errors.New("entropy error"))

		err := uc.SetupMasterPassword(ctx, password)

		assert.Error(t, err)
		assert.Equal(t, "entropy error", err.Error())
	})

	t.Run("SetupMasterPassword_DBSaveError", func(t *testing.T) {
		mockKeyDeriver := new(MockKeyDeriver)
		mockSession := new(MockMasterKeySession)
		mockAppSettingRepo := new(MockAppSettingRepository)
		uc := usecase.NewVaultUseCase(mockKeyDeriver, mockSession, mockAppSettingRepo)

		password := []byte("MySuperSecretPassword123")
		saltBytes := []byte("0123456789abcdef")
		hashedPassword := []byte("32BytesHashedPasswordOutputBytes")

		hashBase64 := base64.StdEncoding.EncodeToString(hashedPassword)
		saltBase64 := base64.StdEncoding.EncodeToString(saltBytes)

		mockKeyDeriver.On("GenerateSalt", usecase.SaltLength).Return(saltBytes, nil)
		mockKeyDeriver.On("HashPassword", mock.Anything, mock.Anything).Return(hashedPassword, nil)
		mockAppSettingRepo.On("SaveMasterPasswordSettings", ctx, hashBase64, saltBase64).Return(errors.New("db error"))

		err := uc.SetupMasterPassword(ctx, password)

		assert.Error(t, err)
		mockSession.AssertNotCalled(t, "SetKey", mock.Anything)
	})
}

func TestVaultUseCase_CheckIfMasterSetupIsRequired(t *testing.T) {
	ctx := context.Background()

	t.Run("CheckIfMasterSetupIsRequired_True", func(t *testing.T) {
		mockAppSettingRepo := new(MockAppSettingRepository)
		uc := usecase.NewVaultUseCase(new(MockKeyDeriver), new(MockMasterKeySession), mockAppSettingRepo)

		mockAppSettingRepo.On("ExistsByKey", ctx, usecase.MasterPasswordHashKey).Return(false, nil)

		required, err := uc.CheckIfMasterSetupIsRequired(ctx)

		assert.NoError(t, err)
		assert.True(t, required)
	})

	t.Run("CheckIfMasterSetupIsRequired_False", func(t *testing.T) {
		mockAppSettingRepo := new(MockAppSettingRepository)
		uc := usecase.NewVaultUseCase(new(MockKeyDeriver), new(MockMasterKeySession), mockAppSettingRepo)

		mockAppSettingRepo.On("ExistsByKey", ctx, usecase.MasterPasswordHashKey).Return(true, nil)

		required, err := uc.CheckIfMasterSetupIsRequired(ctx)

		assert.NoError(t, err)
		assert.False(t, required)
	})

	t.Run("CheckIfMasterSetupIsRequired_DBError", func(t *testing.T) {
		mockAppSettingRepo := new(MockAppSettingRepository)
		uc := usecase.NewVaultUseCase(new(MockKeyDeriver), new(MockMasterKeySession), mockAppSettingRepo)

		mockAppSettingRepo.On("ExistsByKey", ctx, usecase.MasterPasswordHashKey).Return(false, errors.New("db error"))

		required, err := uc.CheckIfMasterSetupIsRequired(ctx)

		assert.Error(t, err)
		assert.False(t, required)
	})
}

func TestVaultUseCase_UnlockVault(t *testing.T) {
	ctx := context.Background()

	password := []byte("MySuperSecretPassword123")
	saltBytes := []byte("0123456789abcdef")
	hashedPassword := []byte("32BytesHashedPasswordOutputBytes")
	secretKey := []byte("32BytesDerivedKeyAES256SecretBytes")

	encodedHash := base64.StdEncoding.EncodeToString(hashedPassword)
	encodedSalt := base64.StdEncoding.EncodeToString(saltBytes)

	t.Run("UnlockVault_Success", func(t *testing.T) {
		mockKeyDeriver := new(MockKeyDeriver)
		mockSession := new(MockMasterKeySession)
		mockAppSettingRepo := new(MockAppSettingRepository)
		uc := usecase.NewVaultUseCase(mockKeyDeriver, mockSession, mockAppSettingRepo)

		mockAppSettingRepo.On("FindByKey", ctx, usecase.MasterPasswordHashKey).Return(&model.AppSetting{Key: usecase.MasterPasswordHashKey, Value: encodedHash}, nil)
		mockAppSettingRepo.On("FindByKey", ctx, usecase.MasterPasswordSaltKey).Return(&model.AppSetting{Key: usecase.MasterPasswordSaltKey, Value: encodedSalt}, nil)
		mockKeyDeriver.On("VerifyPassword", mock.Anything, saltBytes, hashedPassword).Return(true)
		mockKeyDeriver.On("DeriveKey", mock.Anything, saltBytes).Return(secretKey, nil)
		mockSession.On("SetKey", secretKey).Return()

		unlocked, err := uc.UnlockVault(ctx, password)

		assert.NoError(t, err)
		assert.True(t, unlocked)
		mockAppSettingRepo.AssertExpectations(t)
		mockKeyDeriver.AssertExpectations(t)
		mockSession.AssertExpectations(t)
	})

	t.Run("UnlockVault_NotConfigured_MissingHash", func(t *testing.T) {
		mockKeyDeriver := new(MockKeyDeriver)
		mockSession := new(MockMasterKeySession)
		mockAppSettingRepo := new(MockAppSettingRepository)
		uc := usecase.NewVaultUseCase(mockKeyDeriver, mockSession, mockAppSettingRepo)

		mockAppSettingRepo.On("FindByKey", ctx, usecase.MasterPasswordHashKey).Return(nil, nil)

		unlocked, err := uc.UnlockVault(ctx, password)

		assert.ErrorIs(t, err, usecase.ErrMasterPasswordNotConfigured)
		assert.False(t, unlocked)
	})

	t.Run("UnlockVault_InvalidPassword", func(t *testing.T) {
		mockKeyDeriver := new(MockKeyDeriver)
		mockSession := new(MockMasterKeySession)
		mockAppSettingRepo := new(MockAppSettingRepository)
		uc := usecase.NewVaultUseCase(mockKeyDeriver, mockSession, mockAppSettingRepo)

		mockAppSettingRepo.On("FindByKey", ctx, usecase.MasterPasswordHashKey).Return(&model.AppSetting{Key: usecase.MasterPasswordHashKey, Value: encodedHash}, nil)
		mockAppSettingRepo.On("FindByKey", ctx, usecase.MasterPasswordSaltKey).Return(&model.AppSetting{Key: usecase.MasterPasswordSaltKey, Value: encodedSalt}, nil)
		mockKeyDeriver.On("VerifyPassword", mock.Anything, saltBytes, hashedPassword).Return(false)

		unlocked, err := uc.UnlockVault(ctx, password)

		assert.ErrorIs(t, err, usecase.ErrInvalidMasterPassword)
		assert.False(t, unlocked)
		mockSession.AssertNotCalled(t, "SetKey", mock.Anything)
	})

	t.Run("UnlockVault_DBError", func(t *testing.T) {
		mockKeyDeriver := new(MockKeyDeriver)
		mockSession := new(MockMasterKeySession)
		mockAppSettingRepo := new(MockAppSettingRepository)
		uc := usecase.NewVaultUseCase(mockKeyDeriver, mockSession, mockAppSettingRepo)

		mockAppSettingRepo.On("FindByKey", ctx, usecase.MasterPasswordHashKey).Return(nil, errors.New("db connection failure"))

		unlocked, err := uc.UnlockVault(ctx, password)

		assert.Error(t, err)
		assert.False(t, unlocked)
	})
}

func TestVaultUseCase_LockVault(t *testing.T) {
	mockSession := new(MockMasterKeySession)
	uc := usecase.NewVaultUseCase(new(MockKeyDeriver), mockSession, new(MockAppSettingRepository))

	mockSession.On("Clear").Return()

	uc.LockVault()

	mockSession.AssertExpectations(t)
}

func TestVaultUseCase_GetSessionStatus(t *testing.T) {
	mockSession := new(MockMasterKeySession)
	uc := usecase.NewVaultUseCase(new(MockKeyDeriver), mockSession, new(MockAppSettingRepository))

	mockSession.On("HasKey").Return(true)
	mockSession.On("GetSecondsRemaining", usecase.SessionTimeout).Return(int64(895))

	status := uc.GetSessionStatus()

	assert.True(t, status.Active)
	assert.Equal(t, int64(895), status.SecondsLeft)
	mockSession.AssertExpectations(t)
}
