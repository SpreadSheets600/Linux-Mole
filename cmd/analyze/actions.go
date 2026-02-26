package main

import (
	"context"
	"os/exec"
	"path/filepath"
	"runtime"
)

func revealTargetName() string {
	if runtime.GOOS == "darwin" {
		return "Finder"
	}
	return "file manager"
}

func openPath(path string) error {
	ctx, cancel := context.WithTimeout(context.Background(), openCommandTimeout)
	defer cancel()

	cmd, ok := platformOpenCommand(ctx, path)
	if !ok {
		return nil
	}
	return cmd.Run()
}

func revealPath(path string) error {
	ctx, cancel := context.WithTimeout(context.Background(), openCommandTimeout)
	defer cancel()

	cmd, ok := platformRevealCommand(ctx, path)
	if !ok {
		return nil
	}
	return cmd.Run()
}

func platformOpenCommand(ctx context.Context, path string) (*exec.Cmd, bool) {
	switch runtime.GOOS {
	case "darwin":
		return exec.CommandContext(ctx, "open", path), true
	case "linux":
		if _, err := exec.LookPath("xdg-open"); err == nil {
			return exec.CommandContext(ctx, "xdg-open", path), true
		}
		if _, err := exec.LookPath("gio"); err == nil {
			return exec.CommandContext(ctx, "gio", "open", path), true
		}
	}
	return nil, false
}

func platformRevealCommand(ctx context.Context, path string) (*exec.Cmd, bool) {
	switch runtime.GOOS {
	case "darwin":
		return exec.CommandContext(ctx, "open", "-R", path), true
	case "linux":
		parent := filepath.Dir(path)
		return platformOpenCommand(ctx, parent)
	}
	return nil, false
}
