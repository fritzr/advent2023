package main

import (
	"advent2023/util"
	"fmt"
	"log"
	"strings"
)

func parseMaps(lines []string) []util.RangeSet[int] {
	maps := make([]util.RangeSet[int], 0)
	var newMap *util.RangeSet[int]
	for _, line := range lines {
		if line == "" {
			continue
		}
		if strings.HasSuffix(line, ":") { // new map:
			maps = append(maps, util.RangeSet[int]{})
			newMap = &maps[len(maps)-1]
			continue
		}
		//
		rangeNums := util.ParseNumberList(line)
		delta := rangeNums[0] - rangeNums[1]
		sourceSpan := util.Span{rangeNums[1], rangeNums[1] + rangeNums[2]}
		newMap.Add(sourceSpan, delta)
	}
	return maps
}

func main() {
	lines, err := util.ReadInputLines(5)
	if err != nil {
		log.Fatalf("%s", err)
	}

	_, seedLine, _ := strings.Cut(lines[0], ": ")
	seeds := util.ParseNumberList(seedLine)
	maps := parseMaps(lines[2:])

	minValue := -1
	for _, value := range seeds {
		for _, valueMap := range maps {
			newValue := valueMap.Get(value)
			if newValue != nil {
				value += *newValue
			}
		}
		if minValue < 0 || value < minValue {
			minValue = value
		}
	}
	fmt.Println(minValue)
}

// vim: set ts=2 sw=2:
