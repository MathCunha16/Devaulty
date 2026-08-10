package security

import (
	"devaulty-backend/internal/domain/port"
	"sync"
	"time"
)

const (
	DefaultTimeout = 15 * time.Minute
)

type MasterKeySessionHolderAdapter struct {
	mu             sync.RWMutex
	masterKey      []byte
	lastActivityAt *time.Time
}

func NewMasterKeySessionHolder() port.MasterKeySession {
	return &MasterKeySessionHolderAdapter{}
}

func (m *MasterKeySessionHolderAdapter) SetKey(key []byte) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.clearLocked()

	if len(key) == 0 {
		return
	}

	m.masterKey = make([]byte, len(key))
	copy(m.masterKey, key)

	now := time.Now()
	m.lastActivityAt = &now
}

func (m *MasterKeySessionHolderAdapter) GetKey() []byte {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.hasKeyLocked() {
		if m.isExpired(DefaultTimeout) {
			m.clearLocked()
			return nil
		}
		m.touchLocked()
	}
	if m.masterKey == nil {
		return nil
	}
	keyCopy := make([]byte, len(m.masterKey))
	copy(keyCopy, m.masterKey)
	return keyCopy
}

func (m *MasterKeySessionHolderAdapter) HasKey() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if !m.hasKeyLocked() {
		return false
	}
	if m.isExpired(DefaultTimeout) {
		m.clearLocked()
		return false
	}
	return true
}

func (m *MasterKeySessionHolderAdapter) hasKeyLocked() bool {
	return m.masterKey != nil
}

func (m *MasterKeySessionHolderAdapter) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.clearLocked()
}

func (m *MasterKeySessionHolderAdapter) clearLocked() {
	clear(m.masterKey)
	m.masterKey = nil
	m.lastActivityAt = nil
}

func (m *MasterKeySessionHolderAdapter) Touch() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.touchLocked()
}

func (m *MasterKeySessionHolderAdapter) touchLocked() {
	if m.hasKeyLocked() {
		if m.isExpired(DefaultTimeout) {
			m.clearLocked()
			return
		}
		now := time.Now()
		m.lastActivityAt = &now
	}
}

func (m *MasterKeySessionHolderAdapter) GetSecondsRemaining(timeout time.Duration) int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.masterKey == nil || m.lastActivityAt == nil {
		return 0
	}
	elapsedSeconds := time.Since(*m.lastActivityAt).Seconds()
	timeoutSeconds := timeout.Seconds()
	remaining := timeoutSeconds - elapsedSeconds
	if remaining < 0 {
		return 0
	}
	return int64(remaining)
}

// helper methods
func (m *MasterKeySessionHolderAdapter) isExpired(timeout time.Duration) bool {
	if (m.lastActivityAt == nil) || (timeout == 0) {
		return true
	}
	elapsedSeconds := time.Since(*m.lastActivityAt).Seconds()
	return elapsedSeconds > timeout.Seconds()
}
