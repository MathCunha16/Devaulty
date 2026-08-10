package security_test

import (
	"testing"
	"time"

	"devaulty-backend/internal/adapter/out/security"

	"github.com/stretchr/testify/assert"
)

func TestMasterKeySessionHolder_Initialization(t *testing.T) {
	session := security.NewMasterKeySessionHolder()

	assert.False(t, session.HasKey())
	assert.Nil(t, session.GetKey())
	assert.Equal(t, int64(0), session.GetSecondsRemaining(15*time.Minute))
}

func TestMasterKeySessionHolder_SetAndGetKey(t *testing.T) {
	session := security.NewMasterKeySessionHolder()
	keyBytes := []byte("12345678901234567890123456789012") // 32 bytes

	session.SetKey(keyBytes)

	assert.True(t, session.HasKey())
	assert.Equal(t, keyBytes, session.GetKey())
	assert.Greater(t, session.GetSecondsRemaining(15*time.Minute), int64(0))
}

func TestMasterKeySessionHolder_DefensiveCopy(t *testing.T) {
	session := security.NewMasterKeySessionHolder()
	originalKey := []byte("12345678901234567890123456789012")

	session.SetKey(originalKey)

	// Retrieve key and mutate returned slice
	retrievedKey := session.GetKey()
	retrievedKey[0] = 'X'

	// Verify that internal session key was NOT mutated by the caller
	secondRetrieval := session.GetKey()
	assert.Equal(t, originalKey, secondRetrieval)
	assert.NotEqual(t, retrievedKey, secondRetrieval)
}

func TestMasterKeySessionHolder_Clear(t *testing.T) {
	session := security.NewMasterKeySessionHolder()
	keyBytes := []byte("12345678901234567890123456789012")

	session.SetKey(keyBytes)
	assert.True(t, session.HasKey())

	session.Clear()

	assert.False(t, session.HasKey())
	assert.Nil(t, session.GetKey())
	assert.Equal(t, int64(0), session.GetSecondsRemaining(15*time.Minute))
}

func TestMasterKeySessionHolder_SetEmptyKey(t *testing.T) {
	session := security.NewMasterKeySessionHolder()
	keyBytes := []byte("12345678901234567890123456789012")

	session.SetKey(keyBytes)
	assert.True(t, session.HasKey())

	// Setting empty key should clear the session
	session.SetKey([]byte{})
	assert.False(t, session.HasKey())
	assert.Nil(t, session.GetKey())
}

func TestMasterKeySessionHolder_Touch(t *testing.T) {
	session := security.NewMasterKeySessionHolder()
	keyBytes := []byte("12345678901234567890123456789012")

	session.SetKey(keyBytes)
	initialRemaining := session.GetSecondsRemaining(15 * time.Minute)

	session.Touch()
	touchedRemaining := session.GetSecondsRemaining(15 * time.Minute)

	assert.True(t, session.HasKey())
	assert.GreaterOrEqual(t, touchedRemaining, initialRemaining-2)
}
