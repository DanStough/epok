package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"golang.org/x/term"

	"github.com/DanStough/epok/internal/styles"
)

const (

	//https://www.asciiart.eu/text-to-ascii-art
	asciiName = `
 /$$$$$$$$                     /$$      
| $$_____/                    | $$      
| $$        /$$$$$$   /$$$$$$ | $$   /$$
| $$$$$    /$$__  $$ /$$__  $$| $$  /$$/
| $$__/   | $$  \ $$| $$  \ $$| $$$$$$/ 
| $$      | $$  | $$| $$  | $$| $$_  $$ 
| $$$$$$$$| $$$$$$$/|  $$$$$$/| $$ \  $$
|________/| $$____/  \______/ |__/  \__/
          | $$                          
          | $$                          
          |__/                          
`
)

var cfgFile string

const (
	groupIDEpochCommands = "epoch-manipulation"
)

// NewRootCMD creates the root command for the epok CLI application.
func NewRootCMD() *cobra.Command {
	cobra.OnInitialize(initConfig)

	isInteractive := term.IsTerminal(int(os.Stdout.Fd()))

	banner := ""
	if isInteractive {
		sheet := styles.NewEpokTheme().Sheet()
		banner = sheet.Banner.Render(asciiName)
	}

	rootCmd := &cobra.Command{
		Use:   "epok",
		Short: "a tool for working with unix epoch timestamps",
		Long: banner +
			`
epok is a command line tool for viewing 👀 and generating unix epoch timestamps.

Some things you can do with epok:
  - fuzzy-parse timestamps from multiple precisions into human readable date-times.
  - generate timestamps from multiple formats and expressions. (TBA)
  - view and search system timezone information. (TBA)

See the GitHub repository for more information: https://github.com/DanStough/epok`,
		Example: `# Get the human readable version of an epoch timestamp
epok parse 1751074598

# read a timestamp from stdin and print in "simple" mode 
pbpaste | epok parse --output=simple

# create a timestamp with nanosecond precision and save to the clipboard
epok now --precision=nanosecond --output=simple | pbcopy
`,
	}

	// Persistent Flags
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "",
		"config file (default is $HOME/.epok.yaml)")
	rootCmd.PersistentFlags().StringP("output", "o", "pretty",
		"output format. Non-interactive outputs will automatically be downgraded to "+
			"\"simple\" Valid options are: simple, json, pretty")

	// Groups
	groups := []*cobra.Group{
		{
			ID:    groupIDEpochCommands,
			Title: "Epoch Manipulation",
		},
	}
	rootCmd.AddGroup(groups...)

	// Subcommands
	rootCmd.AddCommand(newNowCmd())
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

// bindFlags all the flags to viper, but it won't inherit the required behavior
// https://github.com/spf13/viper/issues/397. This function should be called as part
// of the preExecution phase.
func bindFlags(cmd *cobra.Command) {
	f := cmd.Flags()
	normalizeFunc := f.GetNormalizeFunc()
	f.SetNormalizeFunc(func(fs *pflag.FlagSet, name string) pflag.NormalizedName {
		result := normalizeFunc(fs, name)
		name = strings.ReplaceAll(string(result), "-", "_")
		return pflag.NormalizedName(name)
	})
	if err := viper.BindPFlags(f); err != nil {
		fmt.Println("Error binding flags:", err)
	}
}

type outputMode string

const (
	outputModeJson   outputMode = "json"
	outputModePretty outputMode = "pretty"
	outputModeSimple outputMode = "simple"
)

func getOutput() (outputMode, error) {
	str := viper.GetString("output")
	output := outputMode(str)

	switch output {
	// support shorthands
	case outputModePretty, "p":
		output = outputModePretty
	case outputModeSimple, "s":
		output = outputModeSimple
	case outputModeJson, "j":
		output = outputModeJson
	default:
		return "", fmt.Errorf("invalid output flag: %s", output)
	}

	isInteractive := term.IsTerminal(int(os.Stdout.Fd()))

	// Downgrade styling for non-interactive terminals
	if !isInteractive && output == outputModePretty {
		output = outputModeSimple
	}
	return output, nil
}
