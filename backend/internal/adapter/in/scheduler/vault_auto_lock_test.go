package scheduler_test

import (
	"devaulty-backend/internal/adapter/in/scheduler"
	"devaulty-backend/internal/adapter/out/security"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMasterKeySession struct {
	mock.Mock
}

func (m *MockMasterKeySession) SetKey(key []byte) {
	m.Called(key)
}

func (m *MockMasterKeySession) GetKey() []byte {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).([]byte)
}

func (m *MockMasterKeySession) HasKey() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockMasterKeySession) Clear() {
	m.Called()
}

func (m *MockMasterKeySession) Touch() {
	m.Called()
}

func (m *MockMasterKeySession) GetSecondsRemaining(timeout time.Duration) int64 {
	args := m.Called(timeout)
	return int64(args.Int(0))
}

func TestVaultAutoLock_PurgeExpiredSession(t *testing.T) {
	origInterval := scheduler.SessionPurgeInterval
	scheduler.SessionPurgeInterval = 10 * time.Millisecond
	defer func() { scheduler.SessionPurgeInterval = origInterval }()

	t.Run("PurgeExpiredSession invokes HasKey on ticker tick", func(t *testing.T) {
		mockSession := new(MockMasterKeySession)
		mockSession.On("HasKey").Return(false)

		autoLock := scheduler.NewVaultAutoLock(mockSession)

		go autoLock.PurgeExpiredSession()

		// Allow ticker loop to execute
		time.Sleep(50 * time.Millisecond)

		mockSession.AssertCalled(t, "HasKey")
	})

	t.Run("PurgeExpiredSession with real MasterKeySessionHolderAdapter", func(t *testing.T) {
		sessionHolder := security.NewMasterKeySessionHolder()
		autoLock := scheduler.NewVaultAutoLock(sessionHolder)

		assert.False(t, sessionHolder.HasKey())

		go autoLock.PurgeExpiredSession()

		time.Sleep(50 * time.Millisecond)

		assert.False(t, sessionHolder.HasKey())
	})
}
