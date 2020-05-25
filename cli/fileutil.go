package cli

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func writeFile(destination, content string) error {
	file, err := os.OpenFile(destination, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0664)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}
	return nil
}

// Deletes all files and directories that were created for the command to execute.
func cleanup(files []string) {
	errors := []string{}
	defer func() { notifyCleanupResult(errors) }()

	for _, file := range files {
		toBeRemoved := file
		isInsideDir := filepath.Dir(file) != "."
		if isInsideDir {
			toBeRemoved = filepath.Dir(file)
		}

		if err := os.RemoveAll(toBeRemoved); err != nil {
			log.Print(err)
			errors = append(errors, toBeRemoved)
		}
	}
}

func notifyCleanupResult(errors []string) {
	if len(errors) == 0 {
		return
	}

	fmt.Println("Some files failed to remove")
	for _, file := range errors {
		fmt.Println(file)
	}
}
