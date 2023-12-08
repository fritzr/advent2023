package main

import (
	"advent2023/util"
	"fmt"
	"log"
	"strconv"
	"strings"
)

type Set[T comparable] map[T]bool

func parseNumbers(numberFields string, f func(int)) {
	for _, field := range strings.Fields(numberFields) {
		value, _ := strconv.Atoi(field)
		f(value)
	}
}

func parseNumberList(numberFields string) []int {
	result := make([]int, 0)
	parseNumbers(numberFields, func(value int) { result = append(result, value) })
	return result
}

func parseNumberSet(numberFields string) Set[int] {
	result := make(Set[int])
	parseNumbers(numberFields, func(value int) { result[value] = true })
	return result
}

func countWins(winners Set[int], numbers []int) int {
	count := 0
	for _, number := range numbers {
		if winners[number] {
			count += 1
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
		winning := parseNumberSet(winningStr)
		present := parseNumberList(presentStr)
		wins := countWins(winning, present)
		if wins > 0 {
			score += 1 << (wins - 1)
		}
		thisCopies := gameCopies[gameIndex] + 1
		totalCopies += thisCopies
		for copyIndex := gameIndex + 1; copyIndex < len(gameCopies) && copyIndex <= gameIndex+wins; copyIndex++ {
			// fmt.Printf("after game %d, total = %d, copies = %v\n", gameIndex + 1, totalCopies, gameCopies[gameIndex:])
			gameCopies[copyIndex] += thisCopies
		}
	}

	fmt.Println(score)
	fmt.Println(totalCopies)
}

// vim: set ts=2 sw=2:
