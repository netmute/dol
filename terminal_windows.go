//go:build windows

package main

import (
	"fmt"
	"os"
	"syscall"
)

const (
	enableProcessedOutput           = 0x0001
	enableLineInput                 = 0x0002
	enableEchoInput                 = 0x0004
	enableVirtualTerminalInput      = 0x0200
	enableVirtualTerminalProcessing = 0x0004
)

var procSetConsoleMode = syscall.NewLazyDLL("kernel32.dll").NewProc("SetConsoleMode")

type terminalSession struct {
	in         *os.File
	out        *os.File
	oldInMode  uint32
	oldOutMode uint32
}

func openTerminal() (*terminalSession, error) {
	in, err := os.OpenFile("CONIN$", os.O_RDONLY, 0)
	if err != nil {
		return nil, fmt.Errorf("unable to open CONIN$: %w", err)
	}

	out, err := os.OpenFile("CONOUT$", os.O_WRONLY, 0)
	if err != nil {
		_ = in.Close()
		return nil, fmt.Errorf("unable to open CONOUT$: %w", err)
	}

	inHandle := syscall.Handle(in.Fd())
	outHandle := syscall.Handle(out.Fd())

	var inMode uint32
	if err := syscall.GetConsoleMode(inHandle, &inMode); err != nil {
		_ = in.Close()
		_ = out.Close()
		return nil, fmt.Errorf("unable to query input console mode: %w", err)
	}

	var outMode uint32
	if err := syscall.GetConsoleMode(outHandle, &outMode); err != nil {
		_ = in.Close()
		_ = out.Close()
		return nil, fmt.Errorf("unable to query output console mode: %w", err)
	}

	// Line/echo mode buffers input; VT input allows reading escape sequences.
	newInMode := (inMode &^ (enableLineInput | enableEchoInput)) | enableVirtualTerminalInput
	if err := setConsoleMode(inHandle, newInMode); err != nil {
		_ = in.Close()
		_ = out.Close()
		return nil, fmt.Errorf("could not switch terminal input mode: %w", err)
	}

	// VT processing is required so the terminal interprets ESC-based DSR queries.
	newOutMode := outMode | enableProcessedOutput | enableVirtualTerminalProcessing
	if err := setConsoleMode(outHandle, newOutMode); err != nil {
		_ = setConsoleMode(inHandle, inMode)
		_ = in.Close()
		_ = out.Close()
		return nil, fmt.Errorf("could not switch terminal output mode: %w", err)
	}

	return &terminalSession{
		in:         in,
		out:        out,
		oldInMode:  inMode,
		oldOutMode: outMode,
	}, nil
}

func (t *terminalSession) Read(p []byte) (int, error) {
	return t.in.Read(p)
}

func (t *terminalSession) Write(p []byte) (int, error) {
	return t.out.Write(p)
}

func (t *terminalSession) Close() error {
	_ = setConsoleMode(syscall.Handle(t.in.Fd()), t.oldInMode)
	_ = setConsoleMode(syscall.Handle(t.out.Fd()), t.oldOutMode)

	if err := t.in.Close(); err != nil {
		_ = t.out.Close()
		return err
	}
	return t.out.Close()
}

func setConsoleMode(handle syscall.Handle, mode uint32) error {
	r1, _, err := procSetConsoleMode.Call(uintptr(handle), uintptr(mode))
	if r1 == 0 {
		if err == syscall.Errno(0) {
			return syscall.EINVAL
		}
		return err
	}
	return nil
}
