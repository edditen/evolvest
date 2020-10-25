package main

import (
	"github.com/EdgarTeng/evolvest/cmd/evolvestcli/client"
	"github.com/EdgarTeng/evolvest/pkg/kit"
)

func main() {
	port := ":8762"
	client.StartClient(port)
	kit.WaitSignal()
}
