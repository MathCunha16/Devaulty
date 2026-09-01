//go:build !windows

package main

import (
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

func acquireSingleInstanceLock(dataDir string) (*os.File, error) {
	lockPath := dataDir + "/devaulty.lock"
	file, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, fmt.Errorf("open lock file: %w", err)
	}

	if err := unix.Flock(int(file.Fd()), unix.LOCK_EX|unix.LOCK_NB); err != nil {
		_ = file.Close()
		return nil, fmt.Errorf("another Devaulty instance is already running: %w", err)
	}

	return file, nil
}

func releaseSingleInstanceLock(file *os.File) {
	if file == nil {
		return
	}
	_ = unix.Flock(int(file.Fd()), unix.LOCK_UN)
	_ = file.Close()
}
