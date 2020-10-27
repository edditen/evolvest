package completer

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"os"
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
	} else if s == "quit" || s == "exit" {
		fmt.Println("Bye!")
		os.Exit(0)
		return
	}

	fmt.Println("executing: ", s)
	return
}
