package client

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	CmdKeys = "keys"
	CmdPul  = "pull"
	CmdPush = "push"
)

var commands = []string{CmdKeys, CmdPul, CmdPush}

type Command interface {
	Execute(args ...string) (string, error)
}

func NewCommand(cmd string) (Command, error) {
	switch cmd {
	case CmdKeys:
		return &KeysCommand{baseCommand{
			client: GetEvolvestClient(),
		}}, nil

	case CmdPul:
		return &PullCommand{baseCommand{
			client: GetEvolvestClient(),
		}}, nil
	case CmdPush:
		return &PushCommand{baseCommand{
			client: GetEvolvestClient(),
		}}, nil
	}

	return nil, fmt.Errorf("cmd %s not support", cmd)

}

func ParseCommand(line string) (func(), error) {
	worlds := strings.Fields(line)
	cmd := worlds[0]

	if contains(commands, cmd) {
		c, err := NewCommand(worlds[0])
		if err != nil {
			return nil, err
		}
		return func() {

			if ret, err := c.Execute(worlds[1:]...); err != nil {
				fmt.Printf("err: %s\n", err.Error())
				return
			} else {
				fmt.Println(ret)
			}
		}, nil
	} else if contains([]string{"exit", "quit"}, cmd) {
		return func() {
			fmt.Println("Bye!")
			os.Exit(1)
		}, nil
	}

	return nil, fmt.Errorf("command '%s' is not supported", cmd)
}

func contains(array []string, e string) bool {
	for _, item := range array {
		if item == e {
			return true
		}
	}
	return false
}

type baseCommand struct {
	client *EvolvestClient
}

type KeysCommand struct {
	baseCommand
}

func (c *KeysCommand) Execute(args ...string) (string, error) {
	var pattern string
	if len(args) == 0 {
		pattern = ".*"
	} else if len(args) == 1 {
		pattern = args[0]
	} else {
		return "", fmt.Errorf("wrong format, have multiple keys")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return c.client.Keys(ctx, pattern)
}

type PullCommand struct {
	baseCommand
}

func (c *PullCommand) Execute(args ...string) (string, error) {
	if len(args) != 0 {
		return "", fmt.Errorf("wrong format, no required parameters")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return c.client.Pull(ctx)
}

type PushCommand struct {
	baseCommand
}

func (c *PushCommand) Execute(args ...string) (string, error) {
	if len(args) < 4 {
		return "", fmt.Errorf("wrong format, missing required parameters")
	}
	text := fmt.Sprintf("%s %s %s %s %s",
		args[0], args[1], args[2], args[3], args[4])

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return c.client.Push(ctx, text)
}
