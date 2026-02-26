//go:build linux

package main

import "syscall"

func statAtimeSec(stat *syscall.Stat_t) int64 {
	return stat.Atim.Sec
}

func statAtimeNsec(stat *syscall.Stat_t) int64 {
	return stat.Atim.Nsec
}
