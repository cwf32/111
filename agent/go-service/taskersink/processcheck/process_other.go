//go:build !windows

package processcheck

import "errors"

// enumerateProcessNames is not supported on non-Windows platforms
func enumerateProcessNames() ([]string, error) {
	return nil, errors.New("process enumeration is only supported on Windows")
}
