package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/lipgloss/v2/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/DanStough/epok/internal/styles"
	"github.com/DanStough/epok/parse"
)

const (
	RFC3339NanoTime = "15:04:05.999999999Z07:00"
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
pbpaste | epoch parse

# Override displayed timezones
epok parse 1751074598 -z Big\ Ben=Europe/London,Tokyo\ SkyTree=Asia/Tokyo,Empire\ State=America/New_York
`,

		Args: cobra.MaximumNArgs(1),

		PreRun: func(cmd *cobra.Command, _ []string) {
			bindFlags(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runParse(cmd, args)
		},
		SilenceUsage: true,
	}

	defaultLocales := map[string]string{
		"Local": "Local",
		"UTC":   "UTC",
	}

	parseCmd.Flags().StringToStringP("timezone", "z", defaultLocales,
		"override the map of locales:timezones. "+
			"Use 'Local' for system time.")
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

	mode, err := getOutput()
	if err != nil {
		return err
	}

	timezones := viper.GetStringMapString("timezone")
	if len(timezones) == 0 {
		return errors.New("must specify at least one locale timezone")
	}

	locales := make(map[string]*time.Location, len(timezones))
	for name, timezone := range timezones {
		loc, err := time.LoadLocation(timezone)
		if err != nil {
			return fmt.Errorf("invalid timezone %s for locale %s: %w", timezone, name, err)
		}
		locales[name] = loc
	}

	input = strings.TrimSpace(input)
	timestamp, err := parse.String(input)
	if err != nil {
		return fmt.Errorf("could not parse input: %w", err)
	}

	out := newParseOutput(input, timestamp, locales)

	switch mode {
	case outputModePretty:
		return out.writePretty(cmd.OutOrStdout())
	case outputModeSimple:
		return out.writeSimple(cmd.OutOrStdout())
	case outputModeJson:
		return out.writeJson(cmd.OutOrStdout())
	default:
		return fmt.Errorf("unexpected output format: %s", mode)
	}
}

func readFromStdin(cmd *cobra.Command) (string, error) {
	inputChan := make(chan string, 1)
	// We don't want to block on the error, so we use a buffered channel to allow cleanup.
	errChan := make(chan error, 1)

	// We run in a separate goroutine so we can still respond to
	// context cancellations.
	go func() {
		defer close(inputChan)
		defer close(errChan)

		reader := cmd.InOrStdin()
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
	case <-ctx.Done():
		return "", ctx.Err()
	case err := <-errChan:
		return "", err
	case input = <-inputChan:
		// Fallthrough with assignment to input
	}

	return input, nil
}

type parseOutput struct {
	Epoch   string
	Locales []Locale

	// Derived
	Now time.Time

	relative time.Duration // TODO: make this public so it can be rendered in JSON once we decide on a format
}

type Locale struct {
	Name string
	Time time.Time
}

func newParseOutput(input string, localTime time.Time, localesByTz map[string]*time.Location) *parseOutput {
	now := time.Now()

	localesByTime := make([]Locale, 0, len(localesByTz))
	for name, loc := range localesByTz {
		l := Locale{
			Name: name,
			Time: localTime.In(loc),
		}
		localesByTime = append(localesByTime, l)
	}
	sort.Slice(localesByTime, func(i, j int) bool { return localesByTime[i].Name < localesByTime[j].Name })

	return &parseOutput{
		Epoch:   input,
		Now:     now,
		Locales: localesByTime,

		relative: now.Sub(localTime),
	}
}

func (o *parseOutput) writeSimple(w io.Writer) error {
	var errs error

	tw := tabwriter.NewWriter(w, 0, 0, 4, ' ', tabwriter.TabIndent)
	_, err := fmt.Fprintf(tw, "%s\t%s\t%s\n", "LOCALE", "DATE", "TIME")
	errs = errors.Join(errs, err)

	for _, locale := range o.Locales {
		_, err = fmt.Fprintf(tw, "%s\t%s\t%s\n", locale.Name, formatDate(locale.Time), locale.Time.Format(RFC3339NanoTime))
		errs = errors.Join(errs, err)
	}

	duration, label := formatLocalDiff(o.relative)
	_, err = fmt.Fprintf(tw, "\nRelative: %s %s\n", duration, label)
	errs = errors.Join(errs, err)

	return errs
}

func (o *parseOutput) writePretty(w io.Writer) error {
	sheet := styles.NewEpokTheme().Sheet()

	rows := make([][]string, 0, len(o.Locales))
	for _, locale := range o.Locales {
		rows = append(rows, []string{locale.Name, formatDate(locale.Time), locale.Time.Format(RFC3339NanoTime)})
	}

	t := table.New().
		Border(sheet.Table.BorderThickness).
		BorderStyle(sheet.Table.Border).
		StyleFunc(func(row, col int) lipgloss.Style {
			var style lipgloss.Style

			switch {
			case row == table.HeaderRow:
				return sheet.Table.Header
			case row%2 == 0:
				style = sheet.Table.EvenRow
			default:
				style = sheet.Table.OddRow
			}

			localeWidth := 8
			for _, locale := range o.Locales {
				if len(locale.Name) > localeWidth {
					localeWidth = len(locale.Name)
				}
			}

			switch col {
			case 0:
				style = style.Width(localeWidth + 2) // include padding
			case 1:
				style = style.Width(32)
			case 2:
				style = style.Width(28)
			}

			return style
		}).
		// TODO (dans): timezone could be a separate field
		Headers("Locale", "Date", "Time").
		Rows(rows...)

	var errs error
	_, err := lipgloss.Fprintln(w, t)
	errs = errors.Join(errs, err)

	duration, label := formatLocalDiff(o.relative)
	_, err = fmt.Fprintln(w, sheet.Keyword.Render("Relative:"), sheet.Text.Render(duration), sheet.TextSubdued.Italic(true).Render(label))
	return errors.Join(errs, err)
}

func (o *parseOutput) writeJson(w io.Writer) error {
	bytes, err := json.Marshal(o)
	if err != nil {
		return fmt.Errorf("could not marshal parse output JSON: %w", err)
	}
	_, err = w.Write(bytes)
	if err != nil {
		return fmt.Errorf("could not write parse output: %w", err)
	}
	return nil
}

func formatDate(t time.Time) string {
	return fmt.Sprintf("%s, %s %d, %d", t.Weekday(), t.Month(), t.Day(), t.Year())
}

// TODO: break this down by year and day. The highest unit is currently hour.
// This will need to take into consideration leap days and seconds.
func formatLocalDiff(diff time.Duration) (string, string) {
	label := "ago"
	if diff < 0 {
		diff = diff.Abs()
		label = "from now"
	}

	return diff.String(), label
}
