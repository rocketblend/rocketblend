//go:build !windows

package helpers

import "os/exec"

// SetupSysProcAttr is a no-op on non-Windows platforms.
func SetupSysProcAttr(cmd *exec.Cmd) {
	// No-op for non-Windows platforms
}
