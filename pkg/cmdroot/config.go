package cmdroot

import (
	"encoding/json"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cmdViper = viper.New()

// GetViper return Viper bind command/env/config
func GetViper() *viper.Viper {
	return cmdViper
}

// printConfigs print configs from json file
func printConfigs(cmd *cobra.Command) error {
	cmd.Printf("%s's configs[env=%t]:\n", cmd.Use, cmdEnableEnv)
	keys := cmdViper.AllSettings()
	bs, err := json.MarshalIndent(keys, "", "  ")
	if err != nil {
		cmd.Println(err)
		return err
	}
	cmd.Println(string(bs))
	return nil
}
