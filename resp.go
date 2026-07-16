package main

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"
)

func readInt(r *bufio.Reader, prefix string) (int, error) {
	s, err := r.ReadString('\n')
	if err != nil {
		return 0, errors.New("Cannot read string")
	}

	//strip sufix
	s = strings.TrimPrefix(s, prefix)
	//strip newline
	s = strings.TrimSuffix(s, "\r\n")

	n, err := strconv.Atoi(s)

	if err != nil {
		return 0, errors.New("Cannot convert string to int")
	}
	return n, nil
}

func readCommand(r *bufio.Reader) (string, error) {
	byteCount, err := readInt(r, "$")

	if err != nil {
		return "", errors.New("Cannot read int")
	}

	buf := make([]byte, byteCount)
	_, err = io.ReadFull(r, buf)
	if err != nil {
		return "", errors.New("Cannot read command")
	}

	clrf(r)
	return string(buf), nil
}

func clrf(r *bufio.Reader) error {
	buf := make([]byte, 2)
	_, err := io.ReadFull(r, buf)
	if err != nil {
		errors.New("clrf")
	}
	return nil
}
