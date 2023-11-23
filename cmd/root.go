package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const defaultCfgFileName = ".smuggler"
const defaultCfgFileType = "yaml"

var cfgFile string

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType(defaultCfgFileType)
		viper.SetConfigName(defaultCfgFileName)
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func NewCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "smuggler",
		Short: "Showing things that are not very easy to get",
		Long:  `Showing things that are not very easy to get`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.PersistentFlags().StringVar(
		&cfgFile,
		"config",
		"",
		fmt.Sprintf(
			"config file (default is $HOME/%s.%s)",
			defaultCfgFileName,
			defaultCfgFileType,
		),
	)

	cmd.AddCommand(NewOauth2Cmd())

	return cmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := NewCmd().Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}
