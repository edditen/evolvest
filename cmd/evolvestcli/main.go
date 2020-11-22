package main

import (
	"flag"
	"fmt"
	"github.com/EdgarTeng/evolvest/cmd/evolvestcli/client"
	ecli "github.com/EdgarTeng/evolvest/cmd/evolvestcli/completer"
	"github.com/c-bata/go-prompt"
	"github.com/c-bata/go-prompt/completer"
	"log"
)

func main() {
	addr := parseArgs()

	startRpc(addr)

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

func parseArgs() (addr string) {
	flag.StringVar(&addr, "a", "127.0.0.1:8763", "address")
	flag.Parse()
	return
}

func startRpc(addr string) {
	if err := client.StartClient(addr); err != nil {
		log.Fatalf("connect '%s' error, %v\n", addr, err)
	}
}
