package cmd

import (
	"bytes"
	"context"
	"io"
	"testing"
	"testing/synctest"
	"time"

	"github.com/stretchr/testify/require"
)

// Test_Parse covers basic command functionality and validation.
func Test_Parse(t *testing.T) {
	testCases := []struct {
		name string
		args []string
		in   string

		expectedError  string
		expectedOutput []string
	}{
		{
			name: "happy path - valid unix timestamp arg",
			args: []string{
				"1751770507\n",
			},
			expectedOutput: []string{
				"Local:  Saturday July 5 2025 22 55 7",
				"UTC:  Sunday July 6 2025 2 55 7",
			},
		},
		{
			name: "happy path - valid unix timestamp stdin",
			args: []string{},
			in:   "1751770507\n",
			expectedOutput: []string{
				"Local:  Saturday July 5 2025 22 55 7",
				"UTC:  Sunday July 6 2025 2 55 7",
			},
		},
		{
			name: "invalid argument",
			args: []string{
				"orange",
			},
			expectedError: "could not parse input: invalid timestamp format",
		},
		{
			name: "too many arguments",
			args: []string{
				"1751770507",
				"1751770508",
			},
			expectedError: "accepts at most 1 arg(s), received 2",
		},
	}

	// TODO (dans): refactor out the run method
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parseCmd := newParseCmd()
			parseCmd.SetArgs(tc.args)

			inStream := bytes.NewBufferString(tc.in)
			parseCmd.SetIn(inStream)

			outStream := bytes.NewBufferString("")
			parseCmd.SetOut(outStream)

			errorStream := bytes.NewBufferString("")
			parseCmd.SetErr(errorStream)

			err := parseCmd.Execute()
			if tc.expectedError != "" {
				require.EqualError(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)
			}

			output, err := io.ReadAll(outStream)
			require.NoError(t, err)

			// Usually the error returned from cmd.Execute is the same thing
			// return on STDERR.
			errors, err := io.ReadAll(errorStream)
			require.NoError(t, err)

			t.Logf("STDOUT:\n%s\n", output)
			t.Logf("STDERR:\n%s\n", errors)

			for _, expectedOutput := range tc.expectedOutput {
				require.Contains(t, string(output), expectedOutput)
			}
		})
	}
}

// Test_ParseCancellation makes sure that the command responds to context
// cancellation, like receiving SIGINT.
func Test_ParseCancellation(t *testing.T) {
	synctest.Run(func() {
		parseCmd := newParseCmd()
		parseCmd.SetArgs([]string{})

		// Open a reader and don't write to it to simulate
		// waiting on stdin.
		reader, _ := io.Pipe()
		defer reader.Close()
		parseCmd.SetIn(reader)

		outStream := bytes.NewBufferString("")
		parseCmd.SetOut(outStream)

		errorStream := bytes.NewBufferString("")
		parseCmd.SetErr(errorStream)

		// Create an expiring context to simulate an OS signal
		ctx, cancel := context.WithCancel(context.Background())
		parseCmd.SetContext(ctx)

		var done bool
		var err error
		go func() {
			err = parseCmd.Execute()
			done = true
		}()

		time.Sleep(1 * time.Second) // This time doesn't matter since we're in a synctest
		require.False(t, done)

		cancel()
		synctest.Wait()

		require.True(t, done)

		require.EqualError(t, err, "context canceled")
	})
}
