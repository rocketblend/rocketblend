package extractor

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

func extractDMG(filePath string, destination string, imageName string, contentName string) error {
	// Use hdiutil command to mount the .dmg file
	cmd := exec.Command("hdiutil", "attach", filePath)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to mount .dmg file")
	}

	volumePath := filepath.Join("/Volumes", imageName)

	// Use cp command to copy the contents of the mounted .dmg to the destination
	cmd = exec.Command("cp", "-R", filepath.Join(volumePath, contentName), destination)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to copy contents of .dmg file")
	}

	// Use hdiutil command to detach the mounted .dmg
	cmd = exec.Command("hdiutil", "detach", volumePath)
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
