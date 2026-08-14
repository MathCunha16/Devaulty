package port

type Crypto interface {
	Encrypt(plainData, secretKey, aad []byte) (cipherText, iv, authTag []byte, err error)
	Decrypt(cipherText, iv, authTag, secretKey, aad []byte) (plainData []byte, err error)
}
