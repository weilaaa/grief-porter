package main

import (
	"os"
	"os/exec"
)

// doExec execute exec cmd return nil if exit 0
// otherwise return error
func doExec(cmd, dir string) error {
	_cmd := exec.Command("sh", "-c", cmd)
	_cmd.Dir = dir
	_cmd.Stdout = os.Stdout
	_cmd.Stderr = os.Stderr

	colorPrint(blue, _cmd.String())

	// block until cmd execute completed
	return _cmd.Run()
}
