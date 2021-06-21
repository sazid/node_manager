package app

import (
	"bufio"
	"io"
	"os"
	"strconv"
)

const (
	PidFilename = "pid.txt"

	// PidSentinelValue A sentinel PID value that does not really point
	// to any real process.
	//
	// PIDs are non-negative in either unix (linux/mac) or windows systems.
	// See https://stackoverflow.com/a/10019054 (unix)
	// and https://stackoverflow.com/a/46058651 (windows)
	PidSentinelValue = -99999999
)

// ReadNodePID reads the PID from the `io.Reader` and returns
// it as an `int`. The value will be set to `pidSentinelValue`
// in the event of an error.
func ReadNodePID(r io.Reader) (pid int, err error) {
	scan := bufio.NewScanner(r)
	scan.Scan()
	pid, err = strconv.Atoi(scan.Text())
	return
}

// KillProcess finds the process with the given PID and kills it.
// This ignores pid with value `pidSentinelValue`.
func KillProcess(pid int) error {
	if pid == PidSentinelValue {
		return nil
	}

	proc, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	if err = proc.Kill(); err != nil {
		return err
	}
	return nil
}
