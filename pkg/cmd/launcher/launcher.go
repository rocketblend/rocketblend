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

	// Start the command
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("error running command '%v': %s", command, err)
	}

	// // Capture the output of the command
	// output, err := cmd.CombinedOutput()
	// if err != nil {
	// 	return fmt.Errorf("Error running command '%v': %s\nOutput: %s", command, err, output)
	// }

	return nil
}
