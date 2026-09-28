package executor

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
)

// Executor defines an interface for executing shell commands securely.
// This allows us to mock system calls during testing.
type Executor interface {
	Execute(ctx context.Context, command string, args ...string) (string, error)
	LookPath(file string) (string, error)
}

type DefaultExecutor struct{}

func NewDefaultExecutor() *DefaultExecutor {
	return &DefaultExecutor{}
}

// Execute securely runs a system command and captures its standard output and error.
// SECURITY NOTE: This uses os/exec.CommandContext which passes arguments directly to the OS
// executable without invoking a shell (like sh or cmd.exe). This structurally prevents
// command injection vulnerabilities (e.g. chaining with `&&` or `;`), making it safe to
// pass LLM-generated arguments directly.
func (e *DefaultExecutor) Execute(ctx context.Context, command string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, command, args...)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("command failed: %w\nStderr: %s", err, errBuf.String())
	}
	return outBuf.String(), nil
}

func (e *DefaultExecutor) LookPath(file string) (string, error) {
	return exec.LookPath(file)
}
