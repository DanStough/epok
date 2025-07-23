package cmd

import (
	"bytes"
	"io"
	"testing"
	"testing/synctest"

	"github.com/stretchr/testify/require"
)

// testCase is the structure of almost all command line tests.
type testCase struct {
	name string
	args []string
	in   string

	expectedError  string
	expectedOutput []string
}

// testCommand is a shared harness for testing all subcommands.
func testCommand(t *testing.T, tc testCase) {
	// Using synctest here means the relative time in the output is fixed
	synctest.Run(func() {
		t.Run(tc.name, func(t *testing.T) {
			cmd := NewRootCMD()
			cmd.SetArgs(tc.args)

			inStream := bytes.NewBufferString(tc.in)
			cmd.SetIn(inStream)

			outStream := bytes.NewBufferString("")
			cmd.SetOut(outStream)

			errorStream := bytes.NewBufferString("")
			cmd.SetErr(errorStream)

			err := cmd.Execute()
			if tc.expectedError != "" {
				require.EqualError(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)
			}

			actual, err := io.ReadAll(outStream)
			require.NoError(t, err)

			// Usually the error returned from cmd.Execute is the same thing
			// return on STDERR.
			errors, err := io.ReadAll(errorStream)
			require.NoError(t, err)

			t.Logf("STDOUT:\n%s\n", actual)
			t.Logf("STDERR:\n%s\n", errors)

			for _, expected := range tc.expectedOutput {
				require.Contains(t, string(actual), expected)
			}
		})
	})
}
