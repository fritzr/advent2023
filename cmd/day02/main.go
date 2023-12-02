package main

import (
	"advent2023/util"
	"fmt"
	"log"
	"strings"
)

func main() {
	lines, err := util.ReadInputLines(2)
	if err != nil { log.Fatalf("%s", err) }

	limits := map[string]int{"red": 12, "green": 13, "blue": 14}
	sum := 0
NewGame:
	for id, game := range lines {
		_, game, _ = strings.Cut(game, ":")
		for _, round := range strings.Split(game, ";") {
			for _, draw := range strings.Split(round, ", ") {
				num, color, _ := strings.Cut(draw, " ")
				num := strconv.Atoi(num)
				if num > limits[color] {
					continue NewGame
				}
			}
		}
		sum += id + 1
	}
	fmt.Println(sum)
}

// vim: set ts=2 sw=2:
