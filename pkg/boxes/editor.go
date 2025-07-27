package boxes

import (
	"errors"
	"os"
	"os/exec"
)

func EditFile(filePath string) error {
	return openFile(filePath, false)
}

func ShowFile(filePath string) error {
	return openFile(filePath, true)
}

func openFile(filePath string, readonly bool) error {
	args := []string{filePath}
	if readonly {
		args = append(args, "-R")
	}

	cmd := exec.Command("editor", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	var err error
	wasSuspended := activeApp.Application.Suspend(func() {
		err = cmd.Run()
	})

	if !wasSuspended {
		return errors.New("application was already suspended")
	}

	return err
}
