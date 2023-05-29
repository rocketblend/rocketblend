package extractor

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

func (e *extractor) extractDMGWithContext(ctx context.Context, filePath string, destination string) error {
	logContext := map[string]interface{}{
		"filePath":    filePath,
		"destination": destination,
	}

	e.logger.Debug("Starting DMG extraction", logContext)

	// Mount the DMG file
	cmd := exec.CommandContext(ctx, "hdiutil", "attach", "-nobrowse", filePath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logContext["error"] = err.Error()
		e.logger.Error("Could not mount DMG file", logContext)
		return fmt.Errorf("could not mount DMG file: %s", err)
	}

	e.logger.Debug("DMG file mounted", logContext)

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
		e.logger.Error("Could not determine image name", logContext)
		return fmt.Errorf("could not determine image name")
	}

	defer func() {
		// Use hdiutil command to detach the mounted .dmg
		cmd = exec.Command("hdiutil", "detach", imageName)
		cmd.Run()

		e.logger.Debug("DMG file unmounted", logContext)
	}()

	// copy the files from the mounted image to the destination path
	appFiles, err := filepath.Glob(filepath.Join(imageName, "*.app"))
	if err != nil {
		logContext["error"] = err.Error()
		e.logger.Error("Could not search for app files", logContext)
		return fmt.Errorf("could not search for app files: %s", err)
	}
	if len(appFiles) == 0 {
		e.logger.Error("No app files found in the DMG", logContext)
		return fmt.Errorf("no app files found in the DMG")
	}

	logContext["appFiles"] = appFiles
	e.logger.Debug("Found app files", logContext)

	for _, appFile := range appFiles {
		// Use cp command to copy the .app file to the destination
		cmd = exec.CommandContext(ctx, "cp", "-R", appFile, destination)
		err = cmd.Run()
		if err != nil {
			logContext["error"] = err.Error()
			e.logger.Error("Could not copy app files", logContext)
			return fmt.Errorf("could not copy app files: %s", err)
		}
	}

	e.logger.Info("DMG extraction complete", logContext)

	return nil
}
