package parse

import (
	"errors"
	"testing"
	"time"
)

// Test covers both String and Int functions, as they share the same logic for parsing timestamps.
func Test(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected time.Time
		err      error
	}{
		{
			name:     "valid seconds",
			input:    "1633072800",
			expected: time.Unix(1633072800, 0),
		},
		{
			name:     "valid seconds: 0",
			input:    "0",
			expected: time.Unix(0, 0),
		},
		{
			name:     "valid seconds: 1",
			input:    "1",
			expected: time.Unix(1, 0),
		},
		{
			name:     "valid second: 99999999999",
			input:    "99999999999",
			expected: time.Unix(99999999999, 0),
		},
		{
			name:     "valid second with leading 0's: 099999999999",
			input:    "099999999999",
			expected: time.Unix(99999999999, 0),
		},
		{
			name:     "valid millisecond: 100000000000",
			input:    "100000000000",
			expected: time.Unix(0, 100_000_000_000_000_000),
		},
		{
			name:     "valid millisecond: 99999999999999",
			input:    "99999999999999",
			expected: time.Unix(99999999999, 999000000),
		},
		{
			name:     "valid microseconds: 100000000000000",
			input:    "100000000000000",
			expected: time.Unix(100000000, 0),
		},
		{
			// this is the last valid microsecond value on epochconverter.com
			name:     "valid microseconds: 9999999999999998",
			input:    "9999999999999998",
			expected: time.Unix(9999999999, 999998000),
		},
		{
			name:     "valid nanoseconds: 9999999999999999",
			input:    "9999999999999999",
			expected: time.Unix(9999999, 999999999),
		},
		{
			name:     "valid nanoseconds: int64 overflow",
			input:    "9999999999999999999", // This is larger than int64 can handle
			expected: time.Unix(9999999999, 999999999),
		},
		{
			name:     "valid negative seconds: -1",
			input:    "-1",
			expected: time.Unix(-1, 0),
		},
		{
			name:     "valid negative seconds: -1",
			input:    "-29999999999",
			expected: time.Unix(-29999999999, 0),
		},
		{
			name:     "valid negative millisecond: -30000000000",
			input:    "-30000000000",
			expected: time.Unix(-30_000_000, 0),
		},
		{
			name:     "valid negative millisecond: -99999999999999",
			input:    "-99999999999999",
			expected: time.Unix(-99_999_999_999, -999_000_000),
		},
		{
			name:     "valid negative microseconds: -100000000000000",
			input:    "-100000000000000",
			expected: time.Unix(-100_000_000, 0),
		},
		{
			// this is the last valid microsecond value on epochconverter.com
			name:     "valid negative microseconds: -9999999999999998",
			input:    "-9999999999999998",
			expected: time.Unix(-9_999_999_999, -999_998_000),
		},
		{
			name:     "valid negative nanoseconds: -9999999999999999",
			input:    "-9999999999999999",
			expected: time.Unix(-9_999_999, -999_999_999),
		},
		{
			name:     "valid negative nanoseconds: int64 overflow",
			input:    "-09999999999999999999", // This is larger than int64 can handle
			expected: time.Unix(-9999999999, -999999999),
		},
		{
			name:     "invalid format: not a number",
			input:    "invalid",
			expected: time.Time{},
			err:      ErrInvalidFormat,
		},
		{
			name:     "invalid format: empty string",
			input:    "",
			expected: time.Time{},
			err:      ErrInvalidFormat,
		},
		{
			name:     "invalid format: decimal",
			input:    "1.0",
			expected: time.Time{},
			err:      ErrInvalidFormat,
		},
		{
			name:     "invalid format: overflow",
			input:    "9999999999999999999999999999", // This is larger than int64 can handle for seconds
			expected: time.Time{},
			err:      ErrOverflow,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := String(tt.input)
			if !errors.Is(err, tt.err) {
				t.Errorf("expected error %v, got %v", tt.err, err)
			}
			if !result.Equal(tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
