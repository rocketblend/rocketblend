package extractor

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

func extractDMG(filePath string, destination string) error {
	// Mount the DMG file
	cmd := exec.Command("hdiutil", "attach", "-nobrowse", filePath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("could not mount DMG file: %s", err)
	}

	// Extract the image name from the output of the hdiutil attach command
	imageName := ""
	output := string(out)
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "/dev/") {
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				imageName = parts[2]
				break
			}
		}
	}

	if imageName == "" {
		return fmt.Errorf("could not determine image name")
	}

	defer func() {
		// Use hdiutil command to detach the mounted .dmg
		cmd = exec.Command("hdiutil", "detach", imageName)
		cmd.Run()
	}()

	// copy the files from the mounted image to the destination path
	appFiles, err := filepath.Glob(filepath.Join(imageName, "*.app"))
	if err != nil {
		return fmt.Errorf("could not search for app files: %s", err)
	}
	if len(appFiles) == 0 {
		return fmt.Errorf("no app files found in the DMG")
	}

	for _, appFile := range appFiles {
		// Use cp command to copy the .app file to the destination
		cmd = exec.Command("cp", "-R", appFile, destination)
		err = cmd.Run()
		if err != nil {
			return fmt.Errorf("could not copy app files: %s", err)
		}
	}

	return nil
}
