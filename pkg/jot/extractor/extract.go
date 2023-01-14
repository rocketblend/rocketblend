package extractor

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

func extractDMG(filePath string, destination string) error {
	// Open the .dmg file
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a new reader for the .dmg file
	dmgReader, err := NewDMGReader(file)
	if err != nil {
		return err
	}
	defer dmgReader.Close()

	// Iterate through the entries in the .dmg file
	for _, entry := range dmgReader.Entries {
		// Get the entry's name and path
		name := entry.Name()
		path := destination + "/" + name

		// Check if the entry is a directory or a file
		if entry.IsDir() {
			// Create the directory if it doesn't exist
			if _, err := os.Stat(path); os.IsNotExist(err) {
				err = os.MkdirAll(path, os.ModePerm)
				if err != nil {
					return err
				}
			}
		} else {
			// Create the directories leading up to the file if they don't exist
			if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
				return err
			}

			// Open the file for writing
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, entry.Mode())
			if err != nil {
				return err
			}
			defer f.Close()

			// Copy the contents of the entry to the file
			_, err = io.Copy(f, entry)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func extractDMGCMD(filePath string, destination string, imageName string) error {
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
