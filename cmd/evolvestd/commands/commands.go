package commands

import (
	"github.com/EdgarTeng/evolvest/embed/rpc"
	"github.com/EdgarTeng/evolvest/pkg/cmdroot"
	"github.com/EdgarTeng/evolvest/pkg/common/config"
	"github.com/spf13/cobra"
	"log"
)

// Execute adds all child commands to the root command
func Execute() {
	cmdroot.InitCommand(
		"evolvestd",
		`evolvest service`,
		cmdroot.WithReport(), cmdroot.WithMonitor())
	cmdroot.AddCommand(newServer())
	cmdroot.Execute()
}

func newServer() *cobra.Command {
	serverCmd := &cobra.Command{
		Use:   "start",
		Short: "Start an evolvestd ",
		RunE:  startServer,
	}

	return serverCmd
}

func initBeforeStart() {
	initConfig()
	initGrpc()
}

func initGrpc() {
	if err := rpc.StartServer(":" + config.Config().ServerPort); err != nil {
		log.Fatalf("init grpc server failed, %v\n", err)
	}

}

func initConfig() {
	if err := config.InitConfig(cmdroot.CmdConfig); err != nil {
		log.Fatalf("init config failed, %v\n", err)
	}
}

func startServer(cmd *cobra.Command, args []string) error {
	initBeforeStart()

	cmd.Println("Server started!")

	cmdroot.WaitSignal()
	cmd.Println("Server stopping ...")

	return nil
}
