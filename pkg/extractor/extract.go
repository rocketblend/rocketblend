package extractor

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// TODO: Better cleanup on context cancel
func (e *extractor) extractDMGWithContext(ctx context.Context, filePath string, destination string) error {
	logContext := map[string]interface{}{
		"filePath":    filePath,
		"destination": destination,
	}

	e.logger.Debug("starting .dmg extraction", logContext)

	// Mount the DMG file
	cmd := exec.CommandContext(ctx, "hdiutil", "attach", "-nobrowse", filePath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return err
		}

		logContext["error"] = err.Error()
		e.logger.Error("could not mount .dmg file", logContext)
		return fmt.Errorf("could not mount .dmg file: %s", err)
	}

	e.logger.Debug(".dmg mounted", logContext)

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
		e.logger.Error("could not determine image name", logContext)
		return fmt.Errorf("could not determine image name")
	}

	defer func() {
		// Use hdiutil command to detach the mounted .dmg
		cmd = exec.Command("hdiutil", "detach", imageName)
		cmd.Run()

		e.logger.Debug(".dmg unmounted", logContext)
	}()

	// copy the files from the mounted image to the destination path
	appFiles, err := filepath.Glob(filepath.Join(imageName, "*.app"))
	if err != nil {
		logContext["error"] = err.Error()
		e.logger.Error("could not search for app files", logContext)
		return fmt.Errorf("could not search for app files: %s", err)
	}
	if len(appFiles) == 0 {
		e.logger.Error("no app files found in the .dmg", logContext)
		return fmt.Errorf("no app files found in the .dmg")
	}

	logContext["appFiles"] = appFiles
	e.logger.Debug("found app files", logContext)

	for _, appFile := range appFiles {
		if err := ctx.Err(); err != nil {
			return err
		}

		// Use cp command to copy the .app file to the destination
		cmd = exec.CommandContext(ctx, "cp", "-R", appFile, destination)
		err = cmd.Run()
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return err
			}

			logContext["error"] = err.Error()
			e.logger.Error("could not copy app files", logContext)
			return fmt.Errorf("could not copy app files: %s", err)
		}
	}

	e.logger.Info("extraction of .dmg is complete", logContext)

	return nil
}
