package main

import (
	"bufio"
	"strings"
	"testing"
)

func TestReadInt(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("*3\r\n"))
	n, err := readInt(r, "*")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if n != 3 {
		t.Errorf("expected 3, got %d", n)
	}

	r2 := bufio.NewReader(strings.NewReader(""))
	n, err = readInt(r2, "*")
	if err == nil {
		t.Fatalf("expected an error")
	}

	r3 := bufio.NewReader(strings.NewReader("noprefix\r\n"))
	n, err = readInt(r3, "*")

	if err == nil {
		t.Fatalf("expected an error")
	}
}
