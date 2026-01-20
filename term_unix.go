//go:build darwin || linux || freebsd || netbsd || openbsd || dragonfly

package main

import (
	"syscall"
	"unsafe"
)

type termState struct {
	termios syscall.Termios
}

func makeRaw(fd int) (*termState, error) {
	termios, err := ioctlGetTermios(fd, ioctlReadTermios)
	if err != nil {
		return nil, err
	}

	// Based on "stty raw": disable canonical mode and echo, and disable signals.
	raw := *termios
	raw.Iflag &^= syscall.ICRNL | syscall.INLCR | syscall.IGNCR | syscall.IXON | syscall.IXOFF
	raw.Lflag &^= syscall.ICANON | syscall.ECHO | syscall.ECHOE | syscall.ISIG | syscall.IEXTEN
	raw.Cflag &^= syscall.CSIZE | syscall.PARENB
	raw.Cflag |= syscall.CS8
	raw.Cc[syscall.VMIN] = 1
	raw.Cc[syscall.VTIME] = 0

	if err := ioctlSetTermios(fd, ioctlWriteTermios, &raw); err != nil {
		return nil, err
	}

	return &termState{termios: *termios}, nil
}

func restore(fd int, state *termState) error {
	return ioctlSetTermios(fd, ioctlWriteTermios, &state.termios)
}

func ioctlGetTermios(fd int, req uintptr) (*syscall.Termios, error) {
	var termios syscall.Termios
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), req, uintptr(unsafe.Pointer(&termios)))
	if errno != 0 {
		return nil, errno
	}
	return &termios, nil
}

func ioctlSetTermios(fd int, req uintptr, termios *syscall.Termios) error {
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), req, uintptr(unsafe.Pointer(termios)))
	if errno != 0 {
		return errno
	}
	return nil
}
