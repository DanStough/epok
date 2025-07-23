package cmd

import (
	"testing"
)

func Test_Now(t *testing.T) {
	testCases := []struct {
		name string
		args []string
		in   string

		expectedError  string
		expectedOutput []string
	}{
		{
			name: "happy path - seconds",
			args: []string{
				"now",
			},
			expectedOutput: []string{
				"946684800\n",
			},
		},
		{
			name: "happy path - milliseconds",
			args: []string{
				"now",
				"--precision",
				"milliseconds",
			},
			expectedOutput: []string{
				"946684800000\n",
			},
		},
		{
			name: "happy path - microseconds",
			args: []string{
				"now",
				"-p",
				"us",
			},
			expectedOutput: []string{
				"946684800000000\n",
			},
		},
		{
			name: "happy path - nanoseconds",
			args: []string{
				"now",
				"-pnanos",
			},
			expectedOutput: []string{
				"946684800000000000\n",
			},
		},
		{
			name: "happy path - json output",
			args: []string{
				"now",
				"-ojson",
			},
			expectedOutput: []string{
				"{\"Epoch\":\"946684800\",\"Now\":\"1999-12-31T19:00:00-05:00\"}\n",
			},
		},
		{
			name: "no arguments allowed",
			args: []string{
				"now",
				"banana",
			},
			expectedError: "accepts at most 0 arg(s), received 1",
		},
	}

	for _, tc := range testCases {
		testCommand(t, tc)
	}
}
