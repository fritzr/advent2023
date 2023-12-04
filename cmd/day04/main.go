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

func countWins(winners []int, present []int) int {
  count := 0
  for _, winner := range winners {
    for _, have := range present {
      if winner == have {
        count += 1
      }
    }
  }
  return count
}

func main() {
	lines, err := util.ReadInputLines(4)
	if err != nil {
		log.Fatalf("%s", err)
	}

  score := 0
  totalCopies := int64(0)
  gameCopies := make([]int64, len(lines))
  for gameIndex, line := range lines {
    _, numberSets, _ := strings.Cut(line, ": ")
    winningStr, presentStr, _ := strings.Cut(numberSets, "|")
    winning := parseNumbers(winningStr)
    present := parseNumbers(presentStr)
    wins := countWins(winning, present)
    if wins > 0 {
      score += 1 << (wins - 1)
    }
    thisCopies := gameCopies[gameIndex] + 1
    totalCopies += thisCopies
    for copyIndex := gameIndex + 1; copyIndex < len(gameCopies) && copyIndex <= gameIndex + wins; copyIndex++ {
      // fmt.Printf("won %d copies of %d\n", thisCopies, copyIndex + 1)
      gameCopies[copyIndex] += thisCopies
    }
  }

	fmt.Println(score)
  fmt.Println(totalCopies)
}

// vim: set ts=2 sw=2:
