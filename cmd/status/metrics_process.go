package main

import (
	"bufio"
	"context"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func collectTopProcesses() []ProcessInfo {
	if !commandExists("ps") {
		return nil
	}
	commands, ok := processPSCommands(runtime.GOOS)
	if !ok {
		return nil
	}

	for _, cmdArgs := range commands {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		out, err := runCmd(ctx, "ps", cmdArgs.args...)
		cancel()
		if err != nil {
			continue
		}
		procs := parseTopProcesses(out, 5, cmdArgs.hasHeader)
		if len(procs) > 0 {
			return procs
		}
	}
	return nil
}

type psCommand struct {
	args      []string
	hasHeader bool
}

func processPSCommands(goos string) ([]psCommand, bool) {
	switch goos {
	case "darwin":
		return []psCommand{
			{args: []string{"-Aceo", "pcpu,pmem,comm", "-r"}, hasHeader: true},
		}, true
	case "linux":
		return []psCommand{
			// GNU procps on Linux/WSL. Headers removed via "=" and sorted by CPU.
			{args: []string{"-eo", "pcpu=,pmem=,comm=", "--sort=-pcpu"}, hasHeader: false},
			// Fallback for restricted ps implementations without --sort.
			{args: []string{"-eo", "pcpu,pmem,comm"}, hasHeader: true},
		}, true
	default:
		return nil, false
	}
}

func parseTopProcesses(out string, limit int, hasHeader bool) []ProcessInfo {
	scanner := bufio.NewScanner(strings.NewReader(out))
	var procs []ProcessInfo
	lineIndex := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		lineIndex++
		if hasHeader && lineIndex == 1 {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}

		cpuVal, cpuErr := strconv.ParseFloat(fields[0], 64)
		memVal, memErr := strconv.ParseFloat(fields[1], 64)
		if cpuErr != nil || memErr != nil {
			continue
		}
		name := fields[len(fields)-1]
		// Strip path from command name.
		if idx := strings.LastIndex(name, "/"); idx >= 0 {
			name = name[idx+1:]
		}
		if name == "" {
			continue
		}
		procs = append(procs, ProcessInfo{
			Name:   name,
			CPU:    cpuVal,
			Memory: memVal,
		})
		if len(procs) >= limit {
			break
		}
	}
	if len(procs) == 0 {
		return nil
	}
	return procs
}
