package main

import (
	"flag"
	"github.com/EdgarTeng/evolvest/embed/rpc"
	"github.com/EdgarTeng/evolvest/pkg/common/config"
	"github.com/EdgarTeng/evolvest/pkg/common/logger"
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
	logger.SetVerbose(verbose)
	//init config
	if err := config.InitConfig(configFile); err != nil {
		logger.Fatal("init config failed, %v", err)
	}
	if cfg, err := config.PrintConfig(); err != nil {
		logger.Warn("config info error, %v", err)
	} else {
		logger.Info("show config, %s", cfg)
	}

	// init grpc
	port := ":" + config.Config().ServerPort
	logger.Info("Server running, on listen %s", port)
	if err := rpc.StartServer(port); err != nil {
		logger.Fatal("init grpc server failed, %v", err)
	}
}

func startServer() {

	// recover data from file
	store.Recover()

	logger.Info("Server started!")
	utils.WaitSignal(store.Persistent)
	logger.Info("Bye!")

}
