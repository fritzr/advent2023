package util

import (
	"bufio"
	"cmp"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"strings"
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

type Set[T comparable] map[T]bool

func ParseNumbers(numberFields string, f func(int)) {
	for _, field := range strings.Fields(numberFields) {
		value, _ := strconv.Atoi(field)
		f(value)
	}
}

func ParseNumberList(numberFields string) []int {
	result := make([]int, 0)
	ParseNumbers(numberFields, func(value int) { result = append(result, value) })
	return result
}

func ParseNumberSet(numberFields string) Set[int] {
	result := make(Set[int])
	ParseNumbers(numberFields, func(value int) { result[value] = true })
	return result
}

func Min[T cmp.Ordered](x, y T) T {
	if x < y {
		return x
	}
	return y
}

func FindMin[T cmp.Ordered](slice []T) *T {
	if len(slice) == 0 {
		return nil
	}
	min := &slice[0]
	for index := range slice {
		if slice[index] < *min {
			min = &slice[index]
		}
	}
	return min
}

func FindMax[T cmp.Ordered](slice []T) *T {
	if len(slice) == 0 {
		return nil
	}
	max := &slice[0]
	for index := range slice {
		if slice[index] > *max {
			max = &slice[index]
		}
	}
	return max
}
