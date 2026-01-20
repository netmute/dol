//go:build darwin || freebsd || netbsd || openbsd || dragonfly

package main

import "syscall"

const (
	ioctlReadTermios  = syscall.TIOCGETA
	ioctlWriteTermios = syscall.TIOCSETA
)
