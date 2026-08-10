package port

import "time"

type MasterKeySession interface {
	SetKey(key []byte)
	GetKey() []byte
	HasKey() bool
	Clear()
	Touch()
	GetSecondsRemaining(timeout time.Duration) int64
}
