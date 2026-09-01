//go:build windows

package main

import (
	"fmt"
	"os"

	"golang.org/x/sys/windows"
)

const (
	lockFileExclusiveLock   = 0x00000002
	lockFileFailImmediately = 0x00000001
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
	hFile := windows.Handle(file.Fd())
	var ov windows.Overlapped
	return windows.LockFileEx(
		hFile,
		lockFileExclusiveLock|lockFileFailImmediately,
		0,
		1,
		0,
		&ov,
	)
}

func unlockFileEx(file *os.File) error {
	hFile := windows.Handle(file.Fd())
	var ov windows.Overlapped
	return windows.UnlockFileEx(hFile, 0, 1, 0, &ov)
}
