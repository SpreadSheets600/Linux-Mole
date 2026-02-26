//go:build unix

package main

import (
	"io/fs"
	"syscall"
	"time"
)

func getStatBlocks(info fs.FileInfo) int64 {
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return 0
	}
	return stat.Blocks
}

func getAccessTime(info fs.FileInfo) time.Time {
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return time.Time{}
	}
	return time.Unix(statAtimeSec(stat), statAtimeNsec(stat))
}
