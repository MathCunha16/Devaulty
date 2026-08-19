package updater

import (
	"context"
	"devaulty-backend/internal/domain/port"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type nativeUpdater struct{}

func NewNativeUpdater() port.AppUpdater {
	return &nativeUpdater{}
}

func (n *nativeUpdater) InstallUpdate(ctx context.Context, tempDownloadFilePath string) error {
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("error trying to identify binary path: %w", err)
	}

	oldPath := execPath + ".old"
	_ = os.Remove(oldPath)

	if err := os.Rename(execPath, oldPath); err != nil {
		return fmt.Errorf("error trying to rename binary: %w", err)
	}

	if err := copyFile(tempDownloadFilePath, execPath); err != nil {
		_ = os.Rename(oldPath, execPath) // rollback
		return fmt.Errorf("error trying to install new binary: %w", err)
	}

	if runtime.GOOS != "windows" {
		_ = os.Chmod(execPath, 0755)
	}

	_ = os.Remove(tempDownloadFilePath)
	return nil
}

func (n *nativeUpdater) RestartApp() error {
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("error trying to identify binary path to restart: %w", err)
	}

	time.Sleep(time.Second * 2)
	cmd := exec.Command(execPath, os.Args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error trying to restart app: %w", err)
	}

	os.Exit(0)
	return nil
}

func (n *nativeUpdater) CleanupResidualFiles() error {
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("error trying to identify binary path: %w", err)
	}

	execDir := filepath.Dir(execPath)

	files, err := os.ReadDir(execDir)
	if err != nil {
		return fmt.Errorf("error trying to read directory: %w", err)
	}

	for _, file := range files {
		if !file.IsDir() &&
			(strings.HasSuffix(file.Name(), ".old") ||
				strings.HasSuffix(file.Name(), ".update.tmp")) {
			if err := os.Remove(filepath.Join(execDir, file.Name())); err != nil {
				log.Printf("error trying to remove residual file %s: %v", filepath.Join(execDir, file.Name()), err)
			}
		}
	}

	tempPattern := filepath.Join(os.TempDir(), "devaulty-update-*")
	matches, _ := filepath.Glob(tempPattern)

	for _, m := range matches {
		err = os.Remove(m)
		if err != nil {
			log.Printf("error trying to remove temp file %s: %v", m, err)
		}
	}

	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}
