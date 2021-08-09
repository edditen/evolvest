package main

import (
	"github.com/EdgarTeng/evolvest/cmd/evolvestd/command"
	"github.com/EdgarTeng/evolvest/pkg/rootcmd"
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
			server := command.NewEvolvestd()
			if err := server.Init(); err != nil {
				log.Fatalf("Init error: %+v\n", err)
			}
			if err := server.Run(); err != nil {
				log.Fatalf("Run error: %+v\n", err)
			}
		},
	}

	return cmd
}
