package launcher

import (
	"fmt"
	"os"
	"os/exec"
)

func Launch() error {
	command := []string{"open", "-p", ""}
	if len((os.Args)) > 1 {
		command[2] = os.Args[1]
	}

	cmd := exec.Command("rocketblend", command...)

	// Run the command
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error running command '%v': %s", command, err)
	}

	return nil
}
