package extractor

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

func extractDMG(filePath string, destination string, imageName string, contentName string) error {
	// Use hdiutil command to mount the .dmg file
	cmd := exec.Command("hdiutil", "attach", filePath)
	_, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to mount .dmg file")
	}

	volumePath := filepath.Join("Volumes", imageName)

	// Use cp command to copy the contents of the mounted .dmg to the destination
	cmd = exec.Command("cp", "-R", filepath.Join(volumePath, contentName), destination)
	_, err = cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to copy contents of .dmg file: %w", err)
	}

	// Use hdiutil command to detach the mounted .dmg
	cmd = exec.Command("hdiutil", "detach", volumePath)
	_, err = cmd.Output()
	if err != nil {
		return err
	}

	return nil
}
