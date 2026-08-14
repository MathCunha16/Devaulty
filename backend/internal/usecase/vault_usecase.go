package usecase

import (
	"context"
	"devaulty-backend/internal/domain/port"
	"devaulty-backend/internal/dto"
	"encoding/base64"
	"errors"
	"fmt"
	"time"
)

const (
	MasterPasswordHashKey = "master_password_hash"
	MasterPasswordSaltKey = "master_password_salt"
	SaltLength            = 16 // 16-byte
	SessionTimeout        = 15 * time.Minute
)

var (
	ErrMasterPasswordAlreadyConfigured = errors.New("master password already configured")
	ErrMasterPasswordNotConfigured     = errors.New("master password not configured")
	ErrInvalidMasterPassword           = errors.New("invalid master password")
	ErrVaultLocked                     = errors.New("vault is locked")
)

type VaultUseCase struct {
	keyDeriver       port.KeyDeriver
	masterKeySession port.MasterKeySession
	appSettingRepo   port.AppSettingRepository
}

func NewVaultUseCase(keyDeriver port.KeyDeriver, masterKeySession port.MasterKeySession, appSettingRepo port.AppSettingRepository) *VaultUseCase {
	return &VaultUseCase{
		keyDeriver:       keyDeriver,
		masterKeySession: masterKeySession,
		appSettingRepo:   appSettingRepo,
	}
}

func (uc *VaultUseCase) SetupMasterPassword(ctx context.Context, password []byte) error {
	defer clear(password)

	saltBytes, err := uc.keyDeriver.GenerateSalt(SaltLength)
	if err != nil {
		return err
	}
	defer clear(saltBytes)

	hashedPassword, err := uc.keyDeriver.HashPassword(password, saltBytes)
	if err != nil {
		return err
	}

	hashBase64 := base64.StdEncoding.EncodeToString(hashedPassword)
	saltBase64 := base64.StdEncoding.EncodeToString(saltBytes)

	err = uc.appSettingRepo.SaveMasterPasswordSettings(ctx, hashBase64, saltBase64)
	if err != nil {
		return err
	}

	secretKey, err := uc.keyDeriver.DeriveKey(password, saltBytes)
	if err != nil {
		return err
	}
	defer clear(secretKey)

	uc.masterKeySession.SetKey(secretKey)
	return nil
}

func (uc *VaultUseCase) CheckIfMasterSetupIsRequired(ctx context.Context) (bool, error) {
	alreadyConfigured, err := uc.appSettingRepo.ExistsByKey(ctx, MasterPasswordHashKey)
	if err != nil {
		return false, fmt.Errorf("error checking if master password exists: %w", err)
	}
	return !alreadyConfigured, nil
}

func (uc *VaultUseCase) UnlockVault(ctx context.Context, password []byte) (bool, error) {
	defer clear(password)

	mapPassword, err := uc.appSettingRepo.FindByKey(ctx, MasterPasswordHashKey)
	if err != nil {
		return false, fmt.Errorf("error finding master password hash: %w", err)
	}
	if mapPassword == nil {
		return false, ErrMasterPasswordNotConfigured
	}
	hashedPassword, err := base64.StdEncoding.DecodeString(mapPassword.Value)
	if err != nil {
		return false, fmt.Errorf("error decoding master password hash: %w", err)
	}
	defer clear(hashedPassword)

	mapSalt, err := uc.appSettingRepo.FindByKey(ctx, MasterPasswordSaltKey)
	if err != nil {
		return false, fmt.Errorf("error finding master password salt: %w", err)
	}
	if mapSalt == nil {
		return false, fmt.Errorf("master password not configured")
	}
	saltBytes, err := base64.StdEncoding.DecodeString(mapSalt.Value)
	if err != nil {
		return false, fmt.Errorf("error decoding master password salt: %w", err)
	}
	defer clear(saltBytes)

	if !uc.keyDeriver.VerifyPassword(password, saltBytes, hashedPassword) {
		return false, ErrInvalidMasterPassword
	}

	secretKey, err := uc.keyDeriver.DeriveKey(password, saltBytes)
	if err != nil {
		return false, fmt.Errorf("error deriving key: %w", err)
	}
	defer clear(secretKey)

	uc.masterKeySession.SetKey(secretKey)

	return true, nil
}

func (uc *VaultUseCase) LockVault() {
	// This will purges key bytes completely from RAM
	uc.masterKeySession.Clear()
}

func (uc *VaultUseCase) GetSessionStatus() dto.VaultStatusView {
	isActive := uc.masterKeySession.HasKey()
	secondsLeft := uc.masterKeySession.GetSecondsRemaining(SessionTimeout)
	return dto.VaultStatusView{
		Active:      isActive,
		SecondsLeft: secondsLeft,
	}
}
