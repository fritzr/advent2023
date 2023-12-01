package main

import (
	"advent2023/util"
	"fmt"
	"log"
	"strings"
)

const digits = "0123456789"
var digitWords = []string{
	"zero", "0o",
	"one", "o1e",
	"two", "t2o",
	"three", "t3e",
	"four", "4",
	"five", "5e",
	"six", "6",
	"seven", "7",
	"eight", "e8t",
	"nine", "n9e",
}

func lineSum(line string) int {
	d1 := strings.IndexAny(line, digits)
	d2 := strings.LastIndexAny(line, digits)
	if d1 < 0 { return 0 }
	return 10 * int(line[d1] - '0') + int(line[d2] - '0')
}

func main() {
	lines, err := util.ReadInputLines(1)
	if err != nil { log.Fatalf("%s", err) }

	sum1 := 0
	sum2 := 0
	r := strings.NewReplacer(digitWords...)
	for _, line := range lines {
		sum1 += lineSum(line)

		newLine := r.Replace(line)
		for newLine != line {
			line = newLine
			newLine = r.Replace(line)
		}
		sum2 += lineSum(newLine)
	}
	fmt.Println(sum1)
	fmt.Println(sum2)
}

// vim: set ts=2 sw=2:
