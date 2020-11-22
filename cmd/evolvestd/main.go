package main

import (
	"flag"
	"github.com/EdgarTeng/evolvest/embed/rpc"
	"github.com/EdgarTeng/evolvest/embed/server"
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

}

func startServer() {

	// recover data from file
	store.Recover()

	syncPort := ":" + config.Config().SyncPort
	logger.Info("Sync server running, on listen %s", syncPort)
	go func() {
		if err := rpc.StartServer(syncPort); err != nil {
			logger.Fatal("init server failed, %v", err)
		}
	}()

	// start server
	serverPort := ":" + config.Config().ServerPort
	logger.Info("Server running, on listen %s", serverPort)
	go func() {
		if err := server.StartServer(serverPort); err != nil {
			logger.Fatal("init server failed, %v", err)
		}
	}()

	logger.Info("Server started!")
	utils.WaitSignal(store.Persistent)
	logger.Info("Bye!")

}
