package main

import (
	"bufio"
	"slices"
	"strings"
	"testing"
)

func TestReadInt(t *testing.T) {

	cases := []struct {
		name    string
		input   string
		prefix  string
		want    int
		wantErr bool
	}{
		{"valid input", "*3\r\n", "*", 3, false},
		{"zero", "*0\r\n", "*", 0, false},
		{"no prefix", "3\r\n", "*", 0, true},
		{"non numeric", "*abc\r\n", "*", 0, true},
		{"desync", "\ntest\r\n", "*", 0, true},
		{"no trailing crlf", "*3", "*", 0, true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := bufio.NewReader(strings.NewReader(c.input))
			got, err := readInt(r, c.prefix)
			if c.wantErr {
				if err == nil {
					t.Fatalf("expected an error")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != c.want {
				t.Errorf("got %d: want %d", got, c.want)
			}
		})
	}
}
func TestReadBulkString(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{"valid input", "$3\r\nhey\r\n", "hey", false},
		{"no $ prefix", "3\r\nhey\r\n", "", true},
		{"no trailing crlf", "$3\r\nhey", "", true},
		{"empty string", "$0\r\n\r\n", "", false},
		{"unterminated header", "$3hey", "", true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := bufio.NewReader(strings.NewReader(c.input))
			got, err := readBulkString(r)
			if c.wantErr {
				if err == nil {
					t.Fatalf("expected an error")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error")
			}

			if got != c.want {
				t.Errorf("got %s, want %s", got, c.want)
			}
		})
	}

}

func TestReadCommands(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		want    []string
		wantErr bool
	}{
		{"valid input", "*2\r\n$5\r\nHello\r\n$5\r\nWorld\r\n", []string{"Hello", "World"}, false},
		{"missing string", "*2\r\n$5\r\nHello\r\n", nil, true},
		{"no trailing clrf", "*2\r\n$5\r\nHello\r\n$5\r\nWorld", nil, true},
		{"0 size array", "*0\r\n", []string{}, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := bufio.NewReader(strings.NewReader(c.input))
			got, err := readCommands(r)
			if c.wantErr {
				if err == nil {
					t.Fatalf("expected an error")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected errror")
			}

			if !slices.Equal(got, c.want) {
				t.Errorf("got %q: want %q", got, c.want)
			}
		})
	}
}

func TestCrlf(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid input", "\r\n", false},
		{"dysnc 1", "\n$", true},
		{"dysnc 2", "n\r", true},
		{"dysnc 3", "ld", true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := bufio.NewReader(strings.NewReader(c.input))
			err := crlf(r)
			if c.wantErr {
				if err == nil {
					t.Fatalf("expected an error")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
