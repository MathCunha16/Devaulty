package scheduler

import (
	"devaulty-backend/internal/domain/port"
	"time"
)

var (
	SessionPurgeInterval = time.Minute * 1
)

type VaultAutoLock struct {
	masterKeySession port.MasterKeySession
}

func NewVaultAutoLock(masterKeySession port.MasterKeySession) *VaultAutoLock {
	return &VaultAutoLock{masterKeySession: masterKeySession}
}

func (s *VaultAutoLock) PurgeExpiredSession() {
	ticker := time.NewTicker(SessionPurgeInterval)
	defer ticker.Stop()
	for range ticker.C {
		// HasKey automatically checks expiration and purges expired key bytes from memory
		_ = s.masterKeySession.HasKey()
	}
}
