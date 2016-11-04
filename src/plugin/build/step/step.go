package step

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
)

type BuildStep interface {
	Description() string
	Run() (string, string, error)
}

type ShellCommandStep struct {
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

func (s *ShellCommandStep) Description() string {
	return s.command + " " + strings.Join(s.arguments, " ")
}

func (s *ShellCommandStep) Run() (string, string, error) {
	cmd := exec.Command(s.command, s.arguments...)
	cmd.Dir = s.workDir
	cmd.Env = s.env

	var err error = nil
	outReaderCloser, err := cmd.StdoutPipe()
	if err != nil {
		return "", "", err
	}

	errReaderCloser, err := cmd.StderrPipe()
	if err != nil {
		return "", "", err
	}

	if err := cmd.Start(); err != nil {
		fmt.Println("start err: ", err)
		return "", "", err
	}

	outBytes, err := ioutil.ReadAll(outReaderCloser)
	if err != nil {
		return "", "", err
	}

	errBytes, err := ioutil.ReadAll(errReaderCloser)
	if err != nil {
		return "", "", err
	}

	if err := cmd.Wait(); err != nil {
		return "", "", err
	}

	return string(outBytes), string(errBytes), err
}
