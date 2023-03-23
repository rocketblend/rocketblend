package generator

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func CommandDocs(cmd *cobra.Command, path string) error {
	err := doc.GenMarkdownTree(cmd, path)
	if err != nil {
		return err
	}

	err = os.Rename(filepath.Join(path, cmd.Name()+".md"), filepath.Join(path, "README.md"))
	if err != nil {
		return err
	}

	err = RemoveTextFromFilenames(path, cmd.Name()+"_")
	if err != nil {
		return err
	}

	err = ReplaceStringInFiles(path, "## "+cmd.Name(), "##"+strings.ToUpper(cmd.Name()))
	if err != nil {
		return err
	}

	return nil
}

func RemoveTextFromFilenames(path string, textToRemove string) error {
	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, file := range files {
		oldName := file.Name()
		newName := strings.ReplaceAll(oldName, textToRemove, "")
		if oldName != newName {
			newPath := filepath.Join(path, newName)
			if _, err := os.Stat(newPath); err == nil {
				// File already exists
				continue
			}
			err = os.Rename(filepath.Join(path, oldName), newPath)
			if err != nil {
				return fmt.Errorf("failed to rename file %s to %s: %v", oldName, newName, err)
			}
		}
	}

	return nil
}

func ReplaceStringInFiles(path string, oldStr string, newStr string) error {
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) != ".md" {
			return nil
		}

		fmt.Println("Replacing string in file: " + info.Name())

		// Open file for reading
		file, err := os.Open(info.Name())
		if err != nil {
			return err
		}
		defer file.Close()

		// Create a temporary file for writing
		tmpfile, err := os.CreateTemp(path, "tmp_*")
		if err != nil {
			return err
		}
		defer tmpfile.Close()

		// Replace string in file and write to temporary file
		fileBytes, err := io.ReadAll(file)
		if err != nil {
			return err
		}
		fileContents := string(fileBytes)
		newContents := strings.ReplaceAll(fileContents, oldStr, newStr)
		err = os.WriteFile(tmpfile.Name(), []byte(newContents), 0666)
		if err != nil {
			return err
		}

		fmt.Println("Renaming file: " + info.Name())

		// Rename original file with temporary file
		err = os.Rename(tmpfile.Name(), info.Name())
		if err != nil {
			return err
		}

		err = os.Remove(tmpfile.Name())
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
