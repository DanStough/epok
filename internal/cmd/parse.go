package cmd

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/spf13/cobra"

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
		Example: `# fuzzy-parse timestamp
epok parse 1751074598

# Read from stdin
pbpaste | epoch parse`,

		Args: cobra.MaximumNArgs(1),

		PreRun: func(cmd *cobra.Command, _ []string) {
			bindFlags(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runParse(cmd, args)
		},
		SilenceUsage: true,
	}

	return parseCmd
}

func runParse(cmd *cobra.Command, args []string) error {
	var input string
	var err error
	if len(args) == 0 {
		input, err = readFromStdin(cmd)
		if err != nil {
			return err
		}
	} else {
		input = args[0]
	}

	timestamp, err := parse.String(strings.TrimSpace(input))
	if err != nil {
		return fmt.Errorf("could not parse input: %w", err)
	}

	return writeEpoch(cmd, timestamp)
}

func readFromStdin(cmd *cobra.Command) (string, error) {
	reader := cmd.InOrStdin()

	inputChan := make(chan string)
	// We don't want to block on the error, so we use a buffered channel to allow cleanup.
	errChan := make(chan error, 1)

	// We run in a separate goroutine so we can still respond to
	// context cancellations.
	go func() {
		defer close(inputChan)
		defer close(errChan)

		data, err := io.ReadAll(reader)
		if err != nil {
			// Since this is buffered, we can send the error without blocking
			// and clean up the goroutine. This is useful for testing.
			errChan <- err
			return
		}
		inputChan <- string(data)
	}()

	ctx := cmd.Context()
	var input string
	select {
	case err := <-errChan:
		return "", err
	case <-ctx.Done():
		return "", ctx.Err()
	case input = <-inputChan:
		// Fallthrough with assignment to input
	}

	return input, nil
}

func writeEpoch(cmd *cobra.Command, tsLocal time.Time) error {
	tUTC := tsLocal.In(time.UTC)
	now := time.Now()
	diff := now.Sub(tsLocal)

	var errs error

	_, err := fmt.Fprintln(cmd.OutOrStdout(), "Local: ", tsLocal.Weekday(), tsLocal.Month(), tsLocal.Day(), tsLocal.Year(), tsLocal.Hour(), tsLocal.Minute(), tsLocal.Second())
	errs = errors.Join(errs, err)

	_, err = fmt.Fprintln(cmd.OutOrStdout(), "UTC: ", tUTC.Weekday(), tUTC.Month(), tUTC.Day(), tUTC.Year(), tUTC.Hour(), tUTC.Minute(), tUTC.Second())
	errs = errors.Join(errs, err)

	// TODO (dans): format this with the optional year and day.
	// This will need to take into consideration leap days and seconds.
	_, err = fmt.Fprintln(cmd.OutOrStdout(), "Relative: ", diff, "ago")
	return errors.Join(errs, err)
}
