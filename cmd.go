package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
)

type cmdInstruction string

type withOptsFunc func(cmd *exec.Cmd)

func (c cmdInstruction) String() string {
	return string(c)
}

func (c cmdInstruction) AppendOption(opt string) cmdInstruction {
	return c + " " + cmdInstruction(opt)
}

func makeCmdInstruction(format string, a ...interface{}) cmdInstruction {
	return cmdInstruction(fmt.Sprintf(format, a...))
}

// doExec execute exec cmd and print as stdout、stderr
func (c cmdInstruction) doExec(opts ...withOptsFunc) error {
	cmd := exec.Command("sh", "-c", c.String())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	for _, opt := range opts {
		opt(cmd)
	}

	colorPrint(blue, cmd.String())

	// block until cmd execute completed
	return cmd.Run()
}

// doExecInto execute exec cmd and unmarshal output into given v
func (c cmdInstruction) doExecInto(v interface{}, opts ...withOptsFunc) error {
	cmd := exec.Command("sh", "-c", c.String())

	for _, opt := range opts {
		opt(cmd)
	}

	colorPrint(blue, cmd.String())

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err = cmd.Start(); err != nil {
		return err
	}

	if err = json.NewDecoder(stdout).Decode(v); err != nil {
		return err
	}

	slurp, err := io.ReadAll(stderr)
	if err != nil {
		return err
	}

	if len(slurp) > 0 {
		return fmt.Errorf("%s", slurp)
	}

	if err = cmd.Wait(); err != nil {
		return err
	}

	return nil
}

// doExecOutPut execute exec cmd and return output as []byte
func (c cmdInstruction) doExecOutPut(opts ...withOptsFunc) ([]byte, error) {
	cmd := exec.Command("sh", "-c", c.String())

	for _, opt := range opts {
		opt(cmd)
	}

	colorPrint(blue, cmd.String())

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return out, nil
}
