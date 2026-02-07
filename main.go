package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

const programName = "dol"

var version = "dev"

func main() {
	os.Exit(run())
}

func run() int {
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Printf("%s %s\n", programName, version)
		return 0
	}

	tty, err := openTerminal()
	if err != nil {
		fmt.Fprintln(os.Stderr, "dol:", err)
		return 1
	}
	defer tty.Close()

	// Query: CSI ? 996 n
	fmt.Fprint(tty, "\x1b[?996n")

	// Read reply: CSI ? 997 ; 1 n  (dark)  or  CSI ? 997 ; 2 n (light)
	r := bufio.NewReader(tty)

	type result struct {
		state int
		ok    bool
	}
	done := make(chan result, 1)

	go func() {
		var (
			b     byte
			err   error
			state int
			ok    bool
		)
		for {
			b, err = r.ReadByte()
			if err != nil {
				break
			}
			if b != 0x1b { // ESC
				continue
			}
			b, err = r.ReadByte()
			if err != nil || b != '[' {
				continue
			}
			b, err = r.ReadByte()
			if err != nil || b != '?' {
				continue
			}

			// parse "997;Xn"
			var n int
			for {
				b, err = r.ReadByte()
				if err != nil {
					break
				}
				if b < '0' || b > '9' {
					break
				}
				n = n*10 + int(b-'0')
			}
			if n != 997 || b != ';' {
				continue
			}

			b, err = r.ReadByte()
			if err != nil {
				break
			}
			switch b {
			case '1':
				state = 1
			case '2':
				state = 2
			default:
				continue
			}

			b, err = r.ReadByte()
			if err != nil || b != 'n' {
				continue
			}

			ok = true
			break
		}
		done <- result{state: state, ok: ok}
	}()

	select {
	case res := <-done:
		if !res.ok {
			fmt.Fprintln(os.Stderr, "dol: no response from terminal (may be unsupported)")
			return 1
		}

		if res.state == 1 {
			fmt.Println("dark")
		} else {
			fmt.Println("light")
		}
	case <-time.After(500 * time.Millisecond):
		// Some terminals never reply to this DSR; time out to avoid hanging forever.
		fmt.Fprintln(os.Stderr, "dol: no response from terminal (may be unsupported)")
		return 1
	}
	return 0
}
