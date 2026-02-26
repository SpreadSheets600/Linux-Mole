//go:build darwin

package main

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const trashTimeout = 30 * time.Second

func moveToTrashPlatform(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to resolve path: %w", err)
	}

	escapedPath := strings.ReplaceAll(absPath, "\\", "\\\\")
	escapedPath = strings.ReplaceAll(escapedPath, "\"", "\\\"")
	script := fmt.Sprintf(`tell application "Finder" to delete POSIX file "%s"`, escapedPath)

	ctx, cancel := context.WithTimeout(context.Background(), trashTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "osascript", "-e", script)
	output, err := cmd.CombinedOutput()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("timeout moving to Trash")
		}
		return fmt.Errorf("failed to move to Trash: %s", strings.TrimSpace(string(output)))
	}

	return nil
}
