package parse

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

var (
	ErrInvalidFormat = errors.New("invalid timestamp format")
	ErrOverflow      = errors.New("overflow")
)

// String attempts a fuzzy parse of a string into a time.Time object. The default precisions is
// seconds, but finer precisions are assumed for larger values. For a full description of the behavior,
// review the package tests.
//
// Return values are set with the default `Local` time zone.
func String(s string) (time.Time, error) {
	ticks, err := strconv.Atoi(s)
	if errors.Is(err, strconv.ErrRange) {
		return overflowString(s)
	}
	if err != nil {
		return time.Time{}, ErrInvalidFormat
	}

	// We already have the string, but we will convert it again in
	// Int to handle and leading 0s.
	return Int(int64(ticks))
}

// Int attempts a fuzzy parse of an int64 into a time.Time object. The default precisions is
// seconds, but finer precisions are assumed for larger values. For a full description of the behavior,
// review the package tests.
//
// Return values are set with the default `Local` time zone.
func Int(input int64) (time.Time, error) {
	var seconds, nanoseconds int64
	switch {
	// negative nanosecond
	case input <= -9_999_999_999_999_999:
		seconds = input / 1_000_000_000
		nanoseconds = input % 1_000_000_000
	// negative microseconds
	case input <= -100_000_000_000_000:
		seconds = input / 1_000_000
		nanoseconds = (input % 1_000_000) * 1_000 // convert microseconds to nanoseconds
	// negative milliseconds
	case input <= -30_000_000_000:
		seconds = input / 1_000
		nanoseconds = (input % 1_000) * 1_000_000 // convert milliseconds to nanoseconds
	// seconds
	case input <= 99_999_999_999:
		seconds = input
	// milliseconds
	case input <= 99_999_999_999_999:
		seconds = input / 1_000
		nanoseconds = (input % 1_000) * 1_000_000 // convert milliseconds to nanoseconds
	// microseconds
	case input <= 9_999_999_999_999_998:
		seconds = input / 1_000_000
		nanoseconds = (input % 1_000_000) * 1_000 // convert microseconds to nanoseconds
	// nanoseconds
	default:
		seconds = input / 1_000_000_000
		nanoseconds = input % 1_000_000_000
	}

	return time.Unix(seconds, nanoseconds), nil
}

// overflowString attempts to split a string that's larger than an int64 into nanosecond and
// second portions.
func overflowString(raw string) (time.Time, error) {
	sign := int64(1)
	input := raw
	if raw[0] == '-' {
		sign = -1
		input = raw[1:]
	}

	input = strings.TrimLeft(input, "0")

	nanoseconds, err := strconv.Atoi(input[len(input)-9:])
	if err != nil {
		return time.Time{}, ErrInvalidFormat
	}

	seconds, err := strconv.Atoi(input[:len(input)-9])
	if errors.Is(err, strconv.ErrRange) {
		return time.Time{}, ErrOverflow
	}
	if err != nil {
		return time.Time{}, ErrInvalidFormat
	}

	return time.Unix(sign*int64(seconds), sign*int64(nanoseconds)), nil
}
