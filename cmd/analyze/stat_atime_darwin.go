//go:build darwin

package main

import "syscall"

func statAtimeSec(stat *syscall.Stat_t) int64 {
	return stat.Atimespec.Sec
}

func statAtimeNsec(stat *syscall.Stat_t) int64 {
	return stat.Atimespec.Nsec
}
