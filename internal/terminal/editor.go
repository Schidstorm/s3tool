package terminal

import (
	"os"
	"os/exec"
)

var editCommand = []string{"editor"}
var viewCommand = []string{"less"}

func EditFile(c Context, filePath string) error {

	return openFile(c, filePath, false)
}

func ShowFile(c Context, filePath string) error {
	return openFile(c, filePath, true)
}

func openFile(c Context, filePath string, readonly bool) error {
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
	c.SuspendApp(func() {
		err = cmd.Run()
	})
	return err
}
