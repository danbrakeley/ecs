package test

import (
	"fmt"
	"strings"
)

func didPanic(fn func()) (didPanic bool) {
	defer func() {
		if r := recover(); r != nil {
			didPanic = true
		}
	}()
	didPanic = false
	fn()
	return
}

func formatStringDiff(expected, actual string) []string {
	i := firstDifference(expected, actual)
	if i != -1 {
		return []string{
			fmt.Sprintf("expected: %s", expected),
			fmt.Sprintf("  actual: %s", actual),
			fmt.Sprintf("          %s^", strings.Repeat(" ", i)),
		}
	}
	return nil
}

func firstDifference(a, b string) int {
	minLen := len(a)
	if len(b) < minLen {
		minLen = len(b)
	}

	for i := 0; i < minLen; i++ {
		if a[i] != b[i] {
			return i
		}
	}
	if len(a) != len(b) {
		return minLen
	}
	// no difference
	return -1
}
