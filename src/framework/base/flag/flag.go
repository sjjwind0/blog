package flag

import (
	"errors"
	"fmt"
	"strings"
)

type Command struct {
	Command string
	Args    []string
}

func (c *Command) show() {
	fmt.Println("command: ", c.Command)
	for _, arg := range c.Args {
		fmt.Println(" - ", arg)
	}
}

func preHandleCommand(command string) string {
	var prevChar rune = ' '
	var handleCommand []rune = nil
	for _, c := range command {
		if (c == ' ' && prevChar != ' ') || c != ' ' {
			handleCommand = append(handleCommand, c)
		}
		prevChar = c
	}
	lastIndex := len(handleCommand) - 1
	for i := lastIndex; i >= 0; i-- {
		if handleCommand[i] == ' ' {
			lastIndex = i - 1
			continue
		}
		break
	}
	handleCommand = handleCommand[0 : lastIndex+1]
	return string(handleCommand)
}

func Parse(command string) (*Command, error) {
	cc := preHandleCommand(command)
	if len(cc) == 0 {
		return nil, errors.New("command is nil")
	}
	l := strings.Split(cc, " ")
	if len(l) == 0 {
		return nil, errors.New("command length can not be zero")
	}
	c := &Command{}
	c.Command = l[0]
	c.Args = l[1:]
	return c, nil
}
