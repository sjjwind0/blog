package shell

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

func RunShell(path string, name string, args ...string) (stdOutput string, stdError string, err error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = path

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", "", err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return "", "", err
	}

	if err := cmd.Start(); err != nil {
		return "", "", err
	}

	bytesErr, err := ioutil.ReadAll(stderr)
	if err != nil {
		return "", "", err
	}

	if len(bytesErr) != 0 {
		fmt.Printf("stderr is not nil: %s", bytesErr)
		return "", "", err
	}

	bytesOut, err := ioutil.ReadAll(stdout)
	if err != nil {
		return "", "", err
	}

	if err := cmd.Wait(); err != nil {
		return "", "", err
	}

	return string(bytesOut), string(bytesErr), nil
}

type ShellCompleteCallback func(stdOutput string, stdError string)

func RunShellAsync(callbak ShellCompleteCallback,
	path string, name string, args ...string) (*os.Process, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = path

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	go func() {
		bytesErr, err := ioutil.ReadAll(stderr)
		if err != nil {
			callbak("", "")
			return
		}

		if len(bytesErr) != 0 {
			fmt.Printf("stderr is not nil: %s", bytesErr)
			callbak("", "")
			return
		}

		bytesOut, err := ioutil.ReadAll(stdout)
		if err != nil {
			callbak("", "")
			return
		}

		if err := cmd.Wait(); err != nil {
			callbak("", "")
			return
		}
		callbak(string(bytesOut), string(bytesErr))
	}()
	return cmd.Process, nil
}
