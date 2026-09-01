//go:build windows

package main

import (
	"fmt"
	"os"
	"syscall"
)

func acquireSingleInstanceLock(dataDir string) (*os.File, error) {
	lockPath := dataDir + "\\devaulty.lock"
	file, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, fmt.Errorf("open lock file: %w", err)
	}

	if err := lockFileEx(file); err != nil {
		_ = file.Close()
		return nil, fmt.Errorf("another Devaulty instance is already running: %w", err)
	}

	return file, nil
}

func releaseSingleInstanceLock(file *os.File) {
	if file == nil {
		return
	}
	_ = unlockFileEx(file)
	_ = file.Close()
}

func lockFileEx(file *os.File) error {
	hFile := syscall.Handle(file.Fd())
	var ov syscall.Overlapped
	return syscall.LockFileEx(hFile, syscall.LOCKFILE_EXCLUSIVE_LOCK, 0, 1, 0, &ov)
}

func unlockFileEx(file *os.File) error {
	hFile := syscall.Handle(file.Fd())
	var ov syscall.Overlapped
	return syscall.UnlockFileEx(hFile, 0, 1, 0, &ov)
}
