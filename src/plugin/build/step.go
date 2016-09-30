package build

import (
	"io"
	"os/exec"
)

type BuildStep struct {
	Outer io.ReadCloser
	Error io.ReadCloser
}

func (b *BuildStep) Run() error {
	return nil
}

type ShellCommandStep struct {
	BuildStep
	env       []string
	workDir   string
	command   string
	arguments []string
}

/*
 * @param env
 * @param workDir
 * @param command
 * @param argument
 * @return ShellCommandStep
 */
func NewShellCommandStep(env []string, workDir string, command string, arguments []string) *ShellCommandStep {
	return &ShellCommandStep{env: env, workDir: workDir, command: command, arguments: arguments}
}

func (s *ShellCommandStep) Run() error {
	s.BuildStep.Run()
	cmd := exec.Command(s.command, s.arguments...)
	cmd.Dir = s.workDir
	cmd.Env = s.env

	var err error = nil
	s.Outer, err = cmd.StdoutPipe()
	if err != nil {
		return err
	}

	s.Error, err = cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}
