package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/DanStough/epok/parse"
)

// newParseCmd creates the parse subcommand.
func newParseCmd() *cobra.Command {
	parseCmd := &cobra.Command{
		Use:     "parse unix-timestamp",
		Aliases: []string{"p", "fuzzy-parse"},
		Short:   "fuzzy-parse unix epoch timestamps",
		Long: `Use the parse command to convert unix epoch timestamps into human-readable date-times or 
other formats. It can handle timestamps in various precisions and formats.`,
		GroupID: groupIDEpochCommands,
		Example: `epok parse 1751074598`,

		Args: cobra.ExactArgs(1), // TODO: this will need to change when stdin is supported

		PreRun: func(cmd *cobra.Command, args []string) {
			// This binds all of the flags to viper, but it won't inherit the required behavior
			// https://github.com/spf13/viper/issues/397
			f := cmd.Flags()
			normalizeFunc := f.GetNormalizeFunc()
			f.SetNormalizeFunc(func(fs *pflag.FlagSet, name string) pflag.NormalizedName {
				result := normalizeFunc(fs, name)
				name = strings.ReplaceAll(string(result), "-", "_")
				return pflag.NormalizedName(name)
			})
			viper.BindPFlags(f)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return parseArg(args[0])
		},
		SilenceUsage: true,
	}

	return parseCmd
}

func parseArg(s string) error {
	tLocal, err := parse.String(s)
	if err != nil {
		return err
	}

	tUTC := tLocal.In(time.UTC)
	now := time.Now()
	diff := now.Sub(tLocal)

	fmt.Println("Local: ", tLocal.Weekday(), tLocal.Month(), tLocal.Day(), tLocal.Year(), tLocal.Hour(), tLocal.Minute(), tLocal.Second())
	fmt.Println("UTC: ", tUTC.Weekday(), tUTC.Month(), tUTC.Day(), tUTC.Year(), tUTC.Hour(), tUTC.Minute(), tUTC.Second())

	// TODO (dans): format this with the optional year and day.
	// This will need to take into consideration leap days and seconds.
	fmt.Println("Relative: ", diff, "ago")
	return nil
}
