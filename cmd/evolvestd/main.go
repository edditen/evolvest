package main

import (
	"flag"
	"github.com/EdgarTeng/evolvest/embed/rpc"
	"github.com/EdgarTeng/evolvest/pkg/common/config"
	"github.com/EdgarTeng/evolvest/pkg/common/utils"
	"github.com/EdgarTeng/evolvest/pkg/store"
	"log"
)

var (
	configFile string
)

func main() {
	parseArgs()
	startServer()
}

func parseArgs() {
	flag.StringVar(&configFile, "c", "./conf/config.yaml", "config file")
	flag.Parse()
	startServer()
	return
}

func startServer() {
	//init config
	if err := config.InitConfig(configFile); err != nil {
		log.Fatalf("init config failed, %v\n", err)
	}

	// init grpc
	if err := rpc.StartServer(":" + config.Config().ServerPort); err != nil {
		log.Fatalf("init grpc server failed, %v\n", err)
	}

	// recover data from file
	store.Recover()

	log.Println("Server started!")
	utils.WaitSignal(store.Persistent)
	log.Println("Bye!")

}
