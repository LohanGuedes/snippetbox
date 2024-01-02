package assert

import (
	"strings"
	"testing"
)

func StringContains(t *testing.T, expectedSubstring, actual string) {
	t.Helper()

	if !strings.Contains(actual, expectedSubstring) {
		t.Errorf("Got: %v; Expected to contain: %v", actual, expectedSubstring)
	}
}

func Equal[T comparable](t *testing.T, expected, actual T) {
	t.Helper()

	if actual != expected {
		t.Errorf("Expected %v; Got %v", expected, actual)
	}
}
