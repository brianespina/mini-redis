package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func readInt(r *bufio.Reader, prefix string) (int, error) {
	s, err := r.ReadString('\n')
	if err != nil {
		return 0, err
	}

	if !strings.HasPrefix(s, prefix) {
		return 0, fmt.Errorf("expected prefix %q, got %q", prefix, s)
	}

	//strip sufix
	s = strings.TrimPrefix(s, prefix)
	//strip newline
	s = strings.TrimSuffix(s, "\r\n")

	n, err := strconv.Atoi(s)

	if err != nil {
		return 0, err
	}
	return n, nil
}

func readCommand(r *bufio.Reader) (string, error) {
	byteCount, err := readInt(r, "$")

	if err != nil {
		return "", err
	}

	buf := make([]byte, byteCount)
	_, err = io.ReadFull(r, buf)
	if err != nil {
		return "", err
	}

	err = crlf(r)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

func crlf(r *bufio.Reader) error {
	buf := make([]byte, 2)
	_, err := io.ReadFull(r, buf)
	if err != nil {
		return err
	}

	if buf[0] != '\r' || buf[1] != '\n' {
		return fmt.Errorf("expected crlf, got %q", buf)
	}

	return nil
}
