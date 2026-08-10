package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"devaulty-backend/internal/domain/port"
	"devaulty-backend/internal/usecase"
	"fmt"
)

const (
	TagBitLength = 128 // 16 bytes
	IvByteLength = 12  // 12 bytes
)

type AESGCMCryptoAdapter struct{}

func NewAESGCMCrypto() port.Crypto {
	return &AESGCMCryptoAdapter{}
}

func (A *AESGCMCryptoAdapter) Encrypt(plainData, secretKey, aad []byte) (cipherText, iv, authTag []byte, err error) {
	iv = make([]byte, IvByteLength)
	if _, err := rand.Read(iv); err != nil {
		return nil, nil, nil, err
	}

	cipherBlock, err := aes.NewCipher(secretKey)
	if err != nil {
		return nil, nil, nil, err
	}

	gcm, err := cipher.NewGCM(cipherBlock)
	if err != nil {
		return nil, nil, nil, err
	}

	sealed := gcm.Seal(nil, iv, plainData, aad)

	tagLen := TagBitLength / 8
	cipherTextLen := len(sealed) - tagLen

	cipherText = sealed[:cipherTextLen]
	authTag = sealed[cipherTextLen:]

	return cipherText, iv, authTag, nil
}

func (A *AESGCMCryptoAdapter) Decrypt(cipherText, iv, authTag, secretKey, aad []byte) (plainData []byte, err error) {
	cipherBlock, err := aes.NewCipher(secretKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(cipherBlock)
	if err != nil {
		return nil, err
	}

	sealed := make([]byte, 0, len(cipherText)+len(authTag))
	sealed = append(sealed, cipherText...)
	sealed = append(sealed, authTag...)

	plainData, err = gcm.Open(nil, iv, sealed, aad)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", usecase.ErrInvalidMasterPassword, err)
	}

	return plainData, nil
}
