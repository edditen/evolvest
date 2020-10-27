package cmdroot

import (
	"github.com/spf13/cobra"
)

const (
	CConfig    = "config"
	CPrint     = "print-config"
	CEnableEnv = "enable-env"
)

var (
	// command flag
	CmdConfig    string
	CmdPrint     bool
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
		if CmdPrint {
			return printConfigs(cmd)
		}
		// show version

		return cmd.Help()
	}

	pFlags := root.PersistentFlags()
	// only get form command line
	pFlags.StringVarP(&CmdConfig, CConfig, "c", "./conf/config.yaml", `Config file.")`)
	pFlags.BoolVarP(&CmdPrint, CPrint, "p", false, "Prints configs and exits")
	pFlags.BoolVarP(&cmdEnableEnv, CEnableEnv, "e", false, "Enable config from env")
	return root
}
