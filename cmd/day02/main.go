package main

import (
	"advent2023/util"
	"fmt"
	"log"
	"strconv"
	"strings"
)

func main() {
	lines, err := util.ReadInputLines(2)
	if err != nil {
		log.Fatalf("%s", err)
	}

	limits := map[string]int{"red": 12, "green": 13, "blue": 14}
	sum1 := 0
	sum2 := 0
	for id, game := range lines {
		possible := true
		_, game, _ = strings.Cut(game, ": ")
		// fmt.Printf("game %3d: %s\n        ", id+1, game)
		minCubes := map[string]int{}
		for _, round := range strings.Split(game, ";") {
			for _, draw := range strings.Split(round, ", ") {
				fields := strings.Fields(draw)
				numStr := fields[0]
				color := fields[1]
				num, _ := strconv.Atoi(numStr)
				// fmt.Printf(" num=%d|color='%s'", num, color)
				if num > limits[color] {
					// fmt.Printf("\n     %3d: not scientifically possible (%2d/%2d %s)\n", id+1, num, limits[color], color)
					possible = false
				}
				if num > minCubes[color] {
					minCubes[color] = num
				}
			}
			// fmt.Printf(";")
		}
		// fmt.Printf("\n     %3d: possible\n", id+1)
		if possible {
			sum1 += id + 1
		}
		sum2 += minCubes["red"] * minCubes["green"] * minCubes["blue"]
	}
	fmt.Println(sum1)
	fmt.Println(sum2)
}

// vim: set ts=2 sw=2:
