package security_test

import (
	"bytes"
	"crypto/rand"
	"testing"

	"devaulty-backend/internal/adapter/out/security"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAESGCMCryptoAdapter_EncryptAndDecrypt(t *testing.T) {
	crypto := security.NewAESGCMCrypto()

	key := make([]byte, 32)
	_, _ = rand.Read(key)

	aad := []byte("32BytesLongAADProjectIDAndCredID")
	plainData := []byte(`{"username":"admin","password":"mySecretPassword123"}`)

	t.Run("Encrypt and decrypt success", func(t *testing.T) {
		cipherText, iv, authTag, err := crypto.Encrypt(plainData, key, aad)
		require.NoError(t, err)
		assert.NotEmpty(t, cipherText)
		assert.Len(t, iv, 12)
		assert.Len(t, authTag, 16)

		decrypted, err := crypto.Decrypt(cipherText, iv, authTag, key, aad)
		require.NoError(t, err)
		assert.True(t, bytes.Equal(plainData, decrypted))
	})

	t.Run("Decrypt failure with wrong key", func(t *testing.T) {
		cipherText, iv, authTag, err := crypto.Encrypt(plainData, key, aad)
		require.NoError(t, err)

		wrongKey := make([]byte, 32)
		_, _ = rand.Read(wrongKey)

		_, err = crypto.Decrypt(cipherText, iv, authTag, wrongKey, aad)
		assert.Error(t, err)
	})

	t.Run("Decrypt failure with wrong AAD", func(t *testing.T) {
		cipherText, iv, authTag, err := crypto.Encrypt(plainData, key, aad)
		require.NoError(t, err)

		wrongAAD := []byte("Wrong32BytesLongAADString1234567")

		_, err = crypto.Decrypt(cipherText, iv, authTag, key, wrongAAD)
		assert.Error(t, err)
	})

	t.Run("Decrypt failure with tampered authTag", func(t *testing.T) {
		cipherText, iv, authTag, err := crypto.Encrypt(plainData, key, aad)
		require.NoError(t, err)

		tamperedAuthTag := make([]byte, len(authTag))
		copy(tamperedAuthTag, authTag)
		tamperedAuthTag[0] ^= 0xFF

		_, err = crypto.Decrypt(cipherText, iv, tamperedAuthTag, key, aad)
		assert.Error(t, err)
	})
}
