package cmd

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

const (
	groupIDEpochCommands = "epoch-manipulation"
)

// NewRootCMD creates the root command for the epok CLI application.
func NewRootCMD() *cobra.Command {
	cobra.OnInitialize(initConfig)

	rootCmd := &cobra.Command{
		Use:   "epok",
		Short: "a tool for working with unix epoch timestamps",
		Long: `epok is a command line tool for viewing 👀 and generating unix epoch timestamps.

Some things you can do with epok:
  - fuzzy-parse timestamps from multiple precisions into human readable date-times.
  - generate timestamps from multiple formats and expressions. (TBA)
  - view and search system timezone information. (TBA)

See the GitHub repository for more information: https://github.com/DanStough/epok`,
		Example: `# Get the human readable version of an epoch timestampe
epok parse 1751074598`,
	}

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.epok.yaml)")

	groups := []*cobra.Group{
		{
			ID:    groupIDEpochCommands,
			Title: "Epoch Manipulation",
		},
	}
	rootCmd.AddGroup(groups...)
	// Add subcommands here
	rootCmd.AddCommand(newParseCmd())

	return rootCmd
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".backend" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".epok")
		viper.EnvKeyReplacer(strings.NewReplacer("_", "-"))
		viper.SetEnvPrefix("epok")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Print("Using config file:", viper.ConfigFileUsed())
	}
}
