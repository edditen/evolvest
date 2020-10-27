package main

import (
	"fmt"
	"github.com/EdgarTeng/evolvest/cmd/evolvestcli/client"
	ecli "github.com/EdgarTeng/evolvest/cmd/evolvestcli/completer"
	"github.com/c-bata/go-prompt"
	"github.com/c-bata/go-prompt/completer"
)

func main() {
	startRpc()

	c := ecli.NewCompleter()
	fmt.Printf("evolve-prompt\n")
	fmt.Println("Please use `exit` to exit this program.")
	defer fmt.Println("Bye!")
	p := prompt.New(
		ecli.Executor,
		c.Complete,
		prompt.OptionTitle("evolve-prompt: interactive client"),
		prompt.OptionPrefix(">>> "),
		prompt.OptionInputTextColor(prompt.Yellow),
		prompt.OptionCompletionWordSeparator(completer.FilePathCompletionSeparator),
	)
	p.Run()
}

func startRpc() {
	port := "127.0.0.1:8762"
	client.StartClient(port)
}
