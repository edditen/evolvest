package completer

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/edditen/evolvest/cmd/evolvestcli/client"
	"strings"
)

type Completer struct {
}

func NewCompleter() *Completer {
	return &Completer{}
}

func (c *Completer) Complete(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "keys", Description: "<pattern> 'Keys of pattern'"},
		{Text: "pull", Description: "'Pull values'"},
		{Text: "push", Description: "<txid> <flag> <cmd> <key> [val] 'Push Command'"},
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
