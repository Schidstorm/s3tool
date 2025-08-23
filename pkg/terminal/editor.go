package terminal

import (
	"errors"
	"os"
	"os/exec"
)

var editCommand = []string{"editor"}
var viewCommand = []string{"less"}

func EditFile(filePath string) error {
	return openFile(filePath, false)
}

func ShowFile(filePath string) error {
	return openFile(filePath, true)
}

func openFile(filePath string, readonly bool) error {
	args := []string{filePath}

	var command []string
	if readonly {
		command = append(command, viewCommand...)
	} else {
		command = append(command, editCommand...)
	}
	command = append(command, args...)

	cmd := exec.Command(command[0], command[1:]...)
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
