package cmdroot

import (
	"github.com/spf13/cobra"
)

const (
	CConfig    = "config"
	CShow      = "show-config-only"
	CEnableEnv = "enable-env-variable"
)

var (
	// command flag
	CmdConfig    string
	CmdShow      bool
	cmdEnableEnv bool
)

// return root cobra command
func getRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:   programName,
		Short: "A brief description of command were not set",
		Long:  "A longer description were not set",

		// children of this command will inherit and execute.
		PersistentPreRunE: nil,
	}

	// root help function
	root.RunE = func(cmd *cobra.Command, args []string) error {
		// show config
		if CmdShow {
			return printConfigs(cmd)
		}
		// show version

		return cmd.Help()
	}

	pFlags := root.PersistentFlags()
	// only get form command line
	pFlags.StringVar(&CmdConfig, CConfig, "./conf/config.yaml", `config file.")`)
	pFlags.BoolVarP(&CmdShow, CShow, "S", false, "Prints configs and exits")
	pFlags.BoolVar(&cmdEnableEnv, CEnableEnv, false, "enable config from env")
	return root
}
