package util

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
)

func InputPath(day int) string {
	return path.Join(fmt.Sprintf("cmd/day%02d/input.txt", day))
}

func OpenInput(day int) (*os.File, error) {
	return os.Open(InputPath(day))
}

func ReadInput(day int) (string, error) {
	f, err := OpenInput(day)
	if err != nil {
		return "", err
	}
	defer f.Close()
	raw, err := io.ReadAll(f)
	if err == nil {
		return string(raw), nil
	}
	return "", err
}

func ReadInputLines(day int) ([]string, error) {
	f, err := OpenInput(day)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	lines := make([]string, 0, 128)
	for s.Scan() {
		lines = append(lines, s.Text())
	}
	return lines, s.Err()
}
