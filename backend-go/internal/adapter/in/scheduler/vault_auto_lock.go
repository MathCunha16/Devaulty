package scheduler

import (
	"devaulty-backend/internal/domain/port"
	"time"
)

const (
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
		_ = s.masterKeySession.HasKey()
	}
}
