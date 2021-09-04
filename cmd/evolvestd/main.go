package main

import (
	"github.com/edditen/evolvest/cmd/evolvestd/command"
	"github.com/edditen/evolvest/pkg/rootcmd"
	"github.com/spf13/cobra"
	"log"
)

func main() {
	cmd := rootcmd.InitRootCommand(
		"evolvestd server",
		"evolvestd server wiki: TODO.",
	)

	cmd.AddCommand(GetServeCommand())

	if err := cmd.Execute(); err != nil {
		log.Fatalf("cmd execute error: %+v\n", err)
	}
}

func GetServeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Run a service",
		Run: func(cmd *cobra.Command, args []string) {
			evolvestd := command.NewEvolvestd()
			if err := evolvestd.Init(); err != nil {
				log.Fatalf("Init error: %+v", err)
			}
			errC := make(chan error)
			evolvestd.Run(errC)
			evolvestd.WaitSignal(errC, func() {
				log.Println("prepare clean...")
				evolvestd.Shutdown()
				log.Println("clean finished")
			})
		},
	}

	return cmd
}
