package main

import (
	"bytes"
	"strings"
	"testing"
)

func Noop(in []string) []string {
	return in
}

func checkProcess(t *testing.T, input string, expected string) bool {
	outputWriter := new(bytes.Buffer)
	inputReader := strings.NewReader(input)

	Process(inputReader, outputWriter, Noop)

	if actual := outputWriter.String(); actual != expected {
		t.Errorf("expected '%s' got '%s'", expected, actual)
		return false
	}
	return true
}

func TestProcess(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"", ""},
		{"1\t2", "1\t2\n"},
		{"1\t2\n3\t4", "1\t2\n3\t4\n"},
	}
	for _, c := range cases {
		checkProcess(t, c.in, c.want)
	}
}
