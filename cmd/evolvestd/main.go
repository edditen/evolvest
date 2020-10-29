package main

import (
	"flag"
	"github.com/EdgarTeng/evolvest/embed/rpc"
	"github.com/EdgarTeng/evolvest/pkg/common/config"
	"github.com/EdgarTeng/evolvest/pkg/common/log"
	"github.com/EdgarTeng/evolvest/pkg/common/utils"
	"github.com/EdgarTeng/evolvest/pkg/store"
)

var (
	configFile string
	verbose    bool
)

func main() {
	parseArgs()
	prepare()
	startServer()
}

func parseArgs() {
	flag.StringVar(&configFile, "c", "./conf/config.yaml", "config file")
	flag.BoolVar(&verbose, "v", false, "verbose")
	flag.Parse()
	return
}

func prepare() {
	log.SetVerbose(verbose)
	//init config
	if err := config.InitConfig(configFile); err != nil {
		log.Fatal("init config failed, %v", err)
	}
	if cfg, err := config.PrintConfig(); err != nil {
		log.Warn("config info error, %v", err)
	} else {
		log.Info("show config, %s", cfg)
	}

	// init grpc
	port := ":" + config.Config().ServerPort
	log.Info("Server running, on listen %s\n", port)
	if err := rpc.StartServer(port); err != nil {
		log.Fatal("init grpc server failed, %v", err)
	}
}

func startServer() {

	// recover data from file
	store.Recover()

	log.Info("Server started!")
	utils.WaitSignal(store.Persistent)
	log.Info("Bye!")

}
