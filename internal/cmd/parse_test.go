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

const (

	// This is 1751770507 relative to a "now" of 2000-01-01 0:00
	parseOutputBefore = `UTC       Sunday, July 6, 2025      02:55:07Z

Relative: 223634h55m7s from now`

	// This is 946080000 relative to a "now" of 2000-01-01 0:00
	parseOutputAfter = `
UTC       Saturday, December 25, 1999    00:00:00Z

Relative: 168h0m0s ago`
)

// Test_Parse covers basic command functionality and validation.
func Test_Parse(t *testing.T) {
	testCases := []testCase{
		{
			name: "happy path - valid unix timestamp arg",
			args: []string{
				"parse",
				"1751770507\n",
			},
			expectedOutput: []string{
				"LOCALE    DATE                      TIME",
				parseOutputBefore,
			},
		},
		{
			name: "happy path - valid unix timestamp stdin",
			args: []string{
				"parse",
			},
			in: "1751770507\n",
			expectedOutput: []string{
				"LOCALE    DATE                      TIME",
				parseOutputBefore,
			},
		},
		{
			name: "happy path - timestamp is before 'now'",
			args: []string{
				"parse",
			},
			in: "946080000\n",
			expectedOutput: []string{
				"LOCALE    DATE                           TIME",
				parseOutputAfter,
			},
		},
		{
			name: "happy path - JSON output",
			args: []string{
				"parse",
				"-ojson",
			},
			in: "1751770507\n",
			expectedOutput: []string{
				"{\"Epoch\":\"1751770507\"",
				",{\"Name\":\"UTC\",\"Time\":\"2025-07-06T02:55:07Z\"}],\"Now\":\"1999-12-31T19:00:00-05:00\"}",
			},
		},
		{
			name: "invalid argument",
			args: []string{
				"parse",
				"orange",
			},
			expectedError: "could not parse input: invalid timestamp format",
		},
		{
			name: "too many arguments",
			args: []string{
				"parse",
				"1751770507",
				"1751770508",
			},
			expectedError: "accepts at most 1 arg(s), received 2",
		},
	}

	for _, tc := range testCases {
		testCommand(t, tc)
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
