package app

import (
	"os"
	"os/exec"
)

// Wrapper function to run os command
func Run(command string, args ...string) error {
	program, err := exec.LookPath(command)

	if err != nil {
		return err
	}

	// execute the command
	cmd := exec.Command(program, args...)

	// grab stdout and stderr from the command and pipe them to the parent process
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// run the command
	err = cmd.Run()

	if err != nil {
		return err
	}

	return nil
}
