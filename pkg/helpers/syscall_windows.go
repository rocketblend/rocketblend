//go:build windows

package helpers

import (
	"os/exec"
	"syscall"
)

// SetupSysProcAttr sets up the SysProcAttr for the given command to hide the window on Windows.
func SetupSysProcAttr(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: 0x08000000,
	}
}
