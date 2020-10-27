package client

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	CmdGet = "get"
	CmdSet = "set"
	CmdDel = "del"
)

var commands = []string{CmdGet + " ", CmdSet + " ", CmdDel + " "}

type Command interface {
	Execute(args ...string) (string, error)
}

func NewCommand(cmd string) (Command, error) {
	switch cmd {
	case CmdGet:
		return &GetCommand{baseCommand{
			client: GetEvolvestClient(),
		}}, nil
	case CmdSet:
		return &SetCommand{baseCommand{
			client: GetEvolvestClient(),
		}}, nil
	case CmdDel:
		return &DelCommand{baseCommand{
			client: GetEvolvestClient(),
		}}, nil
	}
	return nil, fmt.Errorf("cmd %s not support", cmd)

}

func ParseCommand(line string) (func(), error) {
	worlds := strings.Fields(line)

	switch {
	case hasPrefixes(line, commands...):
		cmd, err := NewCommand(worlds[0])
		if err != nil {
			return nil, err
		}
		return func() {
			rest, err := cmd.Execute(worlds[1:]...)
			if err != nil {
				fmt.Printf("err: %s\n", err.Error())
				return
			}
			fmt.Println(rest)
		}, nil
	case worlds[0] == "exit" || worlds[0] == "quit":
		return func() {
			fmt.Println("Bye!")
			os.Exit(1)
		}, nil
	}

	return nil, fmt.Errorf("command '%s' cant not be executed", line)
}

func hasPrefixes(s string, prefixes ...string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(s, prefix) {
			return true
		}
	}
	return false
}

type baseCommand struct {
	client *EvolvestClient
}

type GetCommand struct {
	baseCommand
}

func (c *GetCommand) Execute(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("wrong format, missing key")
	}
	if len(args) > 1 {
		return "", fmt.Errorf("wrong format, have multiple keys")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return c.client.Get(ctx, args[0])
}

type SetCommand struct {
	baseCommand
}

func (c *SetCommand) Execute(args ...string) (string, error) {
	if len(args) != 2 {
		return "", fmt.Errorf("format wrong, must be as <key> <value>")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return "ok", c.client.Set(ctx, args[0], args[1])
}

type DelCommand struct {
	baseCommand
}

func (c *DelCommand) Execute(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("wrong format, missing key")
	}
	if len(args) > 1 {
		return "", fmt.Errorf("wrong format, have multiple keys")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return "ok", c.client.Del(ctx, args[0])
}
