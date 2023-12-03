package main

import (
	"advent2023/util"
	"fmt"
	"log"
	"strings"
  "strconv"
)

func main() {
	lines, err := util.ReadInputLines(2)
	if err != nil { log.Fatalf("%s", err) }

	limits := map[string]int{"red": 12, "green": 13, "blue": 14}
	sum := 0
NewGame:
	for id, game := range lines {
		_, game, _ = strings.Cut(game, ": ")
    // fmt.Printf("game %3d: %s\n        ", id+1, game)
		for _, round := range strings.Split(game, ";") {
			for _, draw := range strings.Split(round, ", ") {
				fields := strings.Fields(draw)
        numStr := fields[0]
        color := fields[1]
        num, _ := strconv.Atoi(numStr)
        // fmt.Printf(" num=%d|color='%s'", num, color)
				if num > limits[color] {
          // fmt.Printf("\n     %3d: not scientifically possible (%2d/%2d %s)\n", id+1, num, limits[color], color)
					continue NewGame
				}
			}
      // fmt.Printf(";")
		}
    // fmt.Printf("\n     %3d: possible\n", id+1)
		sum += id + 1
	}
	fmt.Println(sum)
}

// vim: set ts=2 sw=2:
