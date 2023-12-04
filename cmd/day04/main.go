package main

import (
	"advent2023/util"
  "fmt"
  "strings"
  "strconv"
  "log"
)

func parseNumbers(numberFields string) []int {
  result := make([]int, 0)
  for _, field := range strings.Fields(numberFields) {
    value, _ := strconv.Atoi(field)
    result = append(result, value)
  }
  return result
}

func main() {
	lines, err := util.ReadInputLines(4)
	if err != nil {
		log.Fatalf("%s", err)
	}

  score := 0
  for _, line := range lines {
    _, numberSets, _ := strings.Cut(line, ": ")
    winningStr, presentStr, _ := strings.Cut(numberSets, "|")
    winning := parseNumbers(winningStr)
    present := parseNumbers(presentStr)

    count := 0
    for _, winner := range winning {
      for _, have := range present {
        if winner == have {
          count += 1
        }
      }
    }
    if count > 0 {
      score += 1 << (count - 1)
    }
  }

	fmt.Println(score)
}

// vim: set ts=2 sw=2:
