package rootcmd

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// command flag
var (
	rootCmd           *cobra.Command
	flagConfigFile    string
	flagEnv           string
	flagLogConfigFile string
)

// InitRootCommand set command description
//
// Input:
//    short: short description
//    long:  longer description
func InitRootCommand(short, long string) *cobra.Command {
	file, _ := exec.LookPath(os.Args[0])
	_, programName := filepath.Split(file)

	rootCmd = &cobra.Command{
		Use:   programName,
		Short: short,
		Long:  long,

		// children of this command will inherit and execute.
		PersistentPreRunE: nil,
	}

	// root help function
	rootCmd.RunE = func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	}

	pFlags := rootCmd.PersistentFlags()

	pFlags.StringVarP(&flagConfigFile, "config", "c", "conf/conf.yaml", "Set config file")
	if err := viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config")); err != nil {
		log.Println("[InitRootCommand],", err)
		return nil
	}

	pFlags.StringVarP(&flagEnv, "env", "e", "test", `Set serve env`)
	if err := viper.BindPFlag("env", rootCmd.PersistentFlags().Lookup("env")); err != nil {
		log.Println("[InitRootCommand],", err)
		return nil
	}

	pFlags.StringVarP(&flagLogConfigFile, "log-config", "l", "conf/log.yaml", "Set log config file")
	if err := viper.BindPFlag("log-config", rootCmd.PersistentFlags().Lookup("log-config")); err != nil {
		log.Println("[InitRootCommand],", err)
		return nil
	}

	return rootCmd
}
