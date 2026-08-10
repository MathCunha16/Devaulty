package security

import (
	"crypto/rand"
	"crypto/subtle"
	"devaulty-backend/internal/domain/port"
	"fmt"

	"golang.org/x/crypto/argon2"
)

const (
	ArgonIterations  uint32 = 3
	ArgonMemory      uint32 = 65536 // this is 64MB, A.K.A. 64 * 1024 * 1024
	ArgonParallelism uint8  = 2
	KeyLength        uint32 = 32
)

type Argon2KeyDeriverAdapter struct{}

func NewArgon2KeyDeriver() port.KeyDeriver {
	return &Argon2KeyDeriverAdapter{}
}

func (a *Argon2KeyDeriverAdapter) GenerateSalt(size int) ([]byte, error) {
	if size <= 0 {
		return nil, fmt.Errorf("invalid salt size: %d", size)
	}
	salt := make([]byte, size)
	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("error generating salt: %w", err)
	}
	return salt, nil
}

func (a *Argon2KeyDeriverAdapter) DeriveKey(password, salt []byte) ([]byte, error) {
	if password == nil || salt == nil {
		return nil, fmt.Errorf("password and salt must not be nil")
	}
	key := argon2.IDKey(password, salt, ArgonIterations, ArgonMemory, ArgonParallelism, KeyLength)
	return key, nil
}

func (a *Argon2KeyDeriverAdapter) HashPassword(password, salt []byte) ([]byte, error) {
	hashedPassword, err := a.DeriveKey(password, salt)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %w", err)
	}
	return hashedPassword, nil
}

func (a *Argon2KeyDeriverAdapter) VerifyPassword(password, salt, expectedHash []byte) bool {
	hashedPassword, err := a.DeriveKey(password, salt)
	if err != nil {
		return false
	}
	return subtle.ConstantTimeCompare(hashedPassword, expectedHash) == 1
}
