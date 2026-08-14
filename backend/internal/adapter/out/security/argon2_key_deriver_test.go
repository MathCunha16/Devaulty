package security_test

import (
	"bytes"
	"testing"

	"devaulty-backend/internal/adapter/out/security"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestArgon2KeyDeriver_GenerateSalt(t *testing.T) {
	deriver := security.NewArgon2KeyDeriver()

	t.Run("Generates salt of specified length", func(t *testing.T) {
		salt, err := deriver.GenerateSalt(16)
		require.NoError(t, err)
		assert.Len(t, salt, 16)
	})

	t.Run("Generates distinct random salts", func(t *testing.T) {
		salt1, err1 := deriver.GenerateSalt(16)
		salt2, err2 := deriver.GenerateSalt(16)
		require.NoError(t, err1)
		require.NoError(t, err2)
		assert.False(t, bytes.Equal(salt1, salt2))
	})
}

func TestArgon2KeyDeriver_DeriveKey(t *testing.T) {
	deriver := security.NewArgon2KeyDeriver()
	password := []byte("mySecretMasterPassword123")
	salt := []byte("0123456789abcdef")

	t.Run("Derives valid 32-byte key", func(t *testing.T) {
		key, err := deriver.DeriveKey(password, salt)
		require.NoError(t, err)
		assert.Len(t, key, 32)
	})

	t.Run("Derives identical key for identical inputs", func(t *testing.T) {
		key1, err1 := deriver.DeriveKey(password, salt)
		key2, err2 := deriver.DeriveKey(password, salt)
		require.NoError(t, err1)
		require.NoError(t, err2)
		assert.True(t, bytes.Equal(key1, key2))
	})

	t.Run("Derives different key for different passwords", func(t *testing.T) {
		differentPassword := []byte("differentMasterPassword123")
		key1, _ := deriver.DeriveKey(password, salt)
		key2, _ := deriver.DeriveKey(differentPassword, salt)
		assert.False(t, bytes.Equal(key1, key2))
	})

	t.Run("Derives different key for different salts", func(t *testing.T) {
		differentSalt := []byte("fedcba9876543210")
		key1, _ := deriver.DeriveKey(password, salt)
		key2, _ := deriver.DeriveKey(password, differentSalt)
		assert.False(t, bytes.Equal(key1, key2))
	})
}

func TestArgon2KeyDeriver_HashAndPasswordVerification(t *testing.T) {
	deriver := security.NewArgon2KeyDeriver()
	password := []byte("mySecretMasterPassword123")
	wrongPassword := []byte("wrongMasterPassword123")

	salt, err := deriver.GenerateSalt(16)
	require.NoError(t, err)

	t.Run("HashPassword and VerifyPassword success", func(t *testing.T) {
		hash, err := deriver.HashPassword(password, salt)
		require.NoError(t, err)
		assert.Len(t, hash, 32)

		// Verification with correct password
		match := deriver.VerifyPassword(password, salt, hash)
		assert.True(t, match)
	})

	t.Run("VerifyPassword failure with wrong password", func(t *testing.T) {
		hash, err := deriver.HashPassword(password, salt)
		require.NoError(t, err)

		match := deriver.VerifyPassword(wrongPassword, salt, hash)
		assert.False(t, match)
	})
}
