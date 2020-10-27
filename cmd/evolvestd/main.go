package main

import (
	"github.com/EdgarTeng/evolvest/embed/rpc"
	"github.com/EdgarTeng/evolvest/pkg/kit"
)

func main() {
	port := ":8762"
	go rpc.StartServer(port)
	kit.WaitSignal()
}
