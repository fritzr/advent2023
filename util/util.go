package util

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
)

var inputName = flag.String("name", "input.txt", "input filename (relative to day directory)")
var inputPath = flag.String("path", "", "explicit input path (overrides -name)")

func InputPath(day int) string {
	if !flag.Parsed() {
		flag.Parse()
	}
	if inputPath != nil && *inputPath != "" {
		return *inputPath
	}
	return fmt.Sprintf(path.Join("cmd", "day%02d", "%s"), day, *inputName)
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
