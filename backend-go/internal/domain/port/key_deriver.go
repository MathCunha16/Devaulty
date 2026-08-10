package port

type KeyDeriver interface {
	GenerateSalt(size int) ([]byte, error)
	DeriveKey(password, salt []byte) ([]byte, error)
	HashPassword(password, salt []byte) ([]byte, error)
	VerifyPassword(password, salt, expectedHash []byte) bool
}
