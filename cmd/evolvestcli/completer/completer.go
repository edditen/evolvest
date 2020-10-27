package completer

import (
	"fmt"
	"github.com/EdgarTeng/evolvest/cmd/evolvestcli/client"
	"github.com/c-bata/go-prompt"
	"strings"
)

type Completer struct {
}

func NewCompleter() *Completer {
	return &Completer{}
}

func (c *Completer) Complete(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "set", Description: "<key> <value> 'Set key with value'"},
		{Text: "get", Description: "<key> 'Get value of key'"},
		{Text: "del", Description: "<key> 'Del value of key'"},
		{Text: "exit", Description: "Exit the prompt"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func Executor(s string) {
	s = strings.TrimSpace(s)
	if s == "" {
		return
	}
	fn, err := client.ParseCommand(s)
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
		return
	}
	fn()
	return
}
