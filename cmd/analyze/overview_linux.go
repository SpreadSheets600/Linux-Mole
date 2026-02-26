//go:build linux

package main

import (
	"os"
	"strings"
)

func isWSLRuntime() bool {
	data, err := os.ReadFile("/proc/sys/kernel/osrelease")
	if err == nil && strings.Contains(strings.ToLower(string(data)), "microsoft") {
		return true
	}
	data, err = os.ReadFile("/proc/version")
	return err == nil && strings.Contains(strings.ToLower(string(data)), "microsoft")
}

func platformOverviewEntries() []dirEntry {
	if isWSLRuntime() {
		return []dirEntry{
			{Name: "Linux Runtime", Path: "/", IsDir: true, Size: -1},
			{Name: "User Local", Path: "/usr/local", IsDir: true, Size: -1},
			{Name: "Variable Data", Path: "/var", IsDir: true, Size: -1},
			{Name: "Opt", Path: "/opt", IsDir: true, Size: -1},
		}
	}

	return []dirEntry{
		{Name: "System", Path: "/", IsDir: true, Size: -1},
		{Name: "User Local", Path: "/usr/local", IsDir: true, Size: -1},
		{Name: "Variable Data", Path: "/var", IsDir: true, Size: -1},
		{Name: "Opt", Path: "/opt", IsDir: true, Size: -1},
	}
}
