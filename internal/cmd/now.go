package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/DanStough/epok/internal/styles"
)

// newNowCmd creates the now subcommand.
func newNowCmd() *cobra.Command {
	nowCmd := &cobra.Command{
		Use:   "now",
		Short: "create unix timestamp for current instant",
		Long: `Use the now command to create a unix epoch timestamp for the current instant. 
The precision can be adjusted using that flag.`,
		GroupID: groupIDEpochCommands,
		Example: `# generate unix timestamp in seconds
epok now

# generate unix timestamp in nanoseconds
epok now -p ns`,

		Args: cobra.MaximumNArgs(0),

		PreRun: func(cmd *cobra.Command, _ []string) {
			bindFlags(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runNow(cmd)
		},
		SilenceUsage: true,
	}

	nowCmd.Flags().StringP("precision", "p", "seconds", "precision for unix timestamp. valid units are seconds [s,secs], milliseconds [ms, millis], microseconds [us, micros], and nanoseconds [ns, nanos]")

	return nowCmd
}

type precision string

const (
	precisionSeconds      precision = "seconds"
	precisionMilliseconds precision = "milliseconds"
	precisionMicroseconds precision = "microseconds"
	precisionNanoseconds  precision = "nanoseconds"
)

func runNow(cmd *cobra.Command) error {
	str := viper.GetString("precision")
	prec := precision(str)

	switch prec {
	// support shorthands
	case precisionSeconds, "s", "secs":
		prec = precisionSeconds
	case precisionMilliseconds, "ms", "millis":
		prec = precisionMilliseconds
	case precisionMicroseconds, "us", "micros":
		prec = precisionMicroseconds
	case precisionNanoseconds, "ns", "nanos":
		prec = precisionNanoseconds
	default:
		return fmt.Errorf("invalid precision flag: %s", prec)
	}

	out := &NowOutput{
		Now:       time.Now(),
		precision: prec,
	}

	mode, err := getOutput()
	if err != nil {
		return err
	}

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

var _ json.Marshaler = (*NowOutput)(nil)

// NowOutput is the data needed to render the result of the now command.
// It will serialize the time to variable precision.
type NowOutput struct {
	Now time.Time

	precision precision
}

func (o *NowOutput) MarshalJSON() ([]byte, error) {
	ts, err := o.getEpochWithPrecision()
	if err != nil {
		return nil, err
	}

	str := fmt.Sprintf("{ \"Epoch\": \"%s\", \"Now\": \"%s\" }", ts, o.Now.Format(time.RFC3339))
	return []byte(str), nil
}

func (o *NowOutput) writeSimple(w io.Writer) error {
	ts, err := o.getEpochWithPrecision()
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "%s\n", ts)
	return err
}

func (o *NowOutput) writePretty(w io.Writer) error {
	sheet := styles.NewEpokTheme().Sheet()

	epoch, err := o.getEpochWithPrecision()
	if err != nil {
		return err
	}

	caser := cases.Title(language.English)
	label := fmt.Sprintf("%s Epoch:", caser.String(string(o.precision)))
	_, err = fmt.Fprintln(w, sheet.Keyword.Render(label), sheet.Text.Render(epoch))
	return err
}

func (o *NowOutput) writeJson(w io.Writer) error {
	bytes, err := json.Marshal(o)
	if err != nil {
		return fmt.Errorf("could not marshal parse output JSON: %w", err)
	}
	bytes = append(bytes, '\n')
	_, err = w.Write(bytes)
	if err != nil {
		return fmt.Errorf("could not write parse output: %w", err)
	}
	return nil
}

func (o *NowOutput) getEpochWithPrecision() (string, error) {
	var ts string
	switch o.precision {
	case precisionSeconds:
		ts = fmt.Sprintf("%d", o.Now.Unix())
	case precisionMilliseconds:
		ts = fmt.Sprintf("%d", o.Now.UnixNano()/int64(time.Millisecond))
	case precisionMicroseconds:
		ts = fmt.Sprintf("%d", o.Now.UnixNano()/int64(time.Microsecond))
	case precisionNanoseconds:
		ts = fmt.Sprintf("%d", o.Now.UnixNano())
	default:
		return "", fmt.Errorf("unexpected precision: %s", o.precision)
	}
	return ts, nil
}
