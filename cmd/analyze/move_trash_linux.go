//go:build linux

package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const trashTimeout = 30 * time.Second

func moveToTrashPlatform(path string) error {
	if cmd, ok := trashCommand(path); ok {
		ctx, cancel := context.WithTimeout(context.Background(), trashTimeout)
		defer cancel()
		cmd = exec.CommandContext(ctx, cmd.Args[0], cmd.Args[1:]...)
		out, err := cmd.CombinedOutput()
		if err == nil {
			return nil
		}
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("timeout moving to Trash")
		}
		if msg := strings.TrimSpace(string(out)); msg != "" {
			return fmt.Errorf("failed to move to Trash: %s", msg)
		}
	}

	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return fmt.Errorf("failed to move to Trash")
	}
	trashDir := filepath.Join(home, ".local", "share", "Trash", "files")
	if err := os.MkdirAll(trashDir, 0o755); err != nil {
		return fmt.Errorf("failed to move to Trash: %w", err)
	}

	base := filepath.Base(path)
	dest := filepath.Join(trashDir, base)
	for i := 1; ; i++ {
		if _, err := os.Lstat(dest); os.IsNotExist(err) {
			break
		}
		dest = filepath.Join(trashDir, fmt.Sprintf("%s.%d", base, i))
	}

	if err := os.Rename(path, dest); err != nil {
		return fmt.Errorf("failed to move to Trash: %w", err)
	}
	return nil
}

func trashCommand(path string) (*exec.Cmd, bool) {
	if _, err := exec.LookPath("gio"); err == nil {
		return exec.Command("gio", "trash", path), true
	}
	if _, err := exec.LookPath("trash-put"); err == nil {
		return exec.Command("trash-put", path), true
	}
	if _, err := exec.LookPath("kioclient5"); err == nil {
		return exec.Command("kioclient5", "move", path, "trash:/"), true
	}
	return nil, false
}
