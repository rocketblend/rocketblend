package extractor

import (
	"os/exec"
)

func extractDMG(filePath string, destination string, imageName string) error {
	// Use hdiutil command to mount the .dmg file
	cmd := exec.Command("hdiutil", "attach", filePath)
	_, err := cmd.Output()
	if err != nil {
		return err
	}

	// Use cp command to copy the contents of the mounted .dmg to the destination
	cmd = exec.Command("cp", "-R", "/Volumes/"+imageName+"/*", destination)
	_, err = cmd.Output()
	if err != nil {
		return err
	}

	// Use hdiutil command to detach the mounted .dmg
	cmd = exec.Command("hdiutil", "detach", "/Volumes/"+imageName)
	_, err = cmd.Output()
	if err != nil {
		return err
	}

	return nil
}
