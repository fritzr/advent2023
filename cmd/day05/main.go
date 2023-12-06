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

func mapValue(valueMap util.RangeSet[int], value int) int {
	newValue := valueMap.Get(value)
	if newValue != nil {
		value += *newValue
	}
	return value
}

func mapMinValue(seeds []int, maps []util.RangeSet[int]) int {
	minValue := -1
	for _, value := range seeds {
		for _, valueMap := range maps {
			value = mapValue(valueMap, value)
		}
		if minValue < 0 || value < minValue {
			minValue = value
		}
	}
	return minValue
}

func mapMinRange(seeds []int, maps []util.RangeSet[int]) int {
	seedSet := util.RangeSet[int]{}
	for index := 0; index < len(seeds); index += 2 {
		seedSet.Add(util.Span{seeds[index], seeds[index] + seeds[index+1]}, 0)
	}
	for _, valueMap := range maps {
		seedSet = valueMap.Intersect(seedSet, func(_, _, _ util.Span, delta1, delta2 *int) int {
			// TODO
			return *delta1 + *delta2
		})
	}
	minValue := -1
	seedSet.Do(func(s util.Span, delta *int) bool {
		value := s[0] + *delta
		if minValue < 0 || value < minValue {
			minValue = value
		}
		return true
	})
	return minValue
}

func main() {
	lines, err := util.ReadInputLines(5)
	if err != nil {
		log.Fatalf("%s", err)
	}

	_, seedLine, _ := strings.Cut(lines[0], ": ")
	seeds := util.ParseNumberList(seedLine)
	maps := parseMaps(lines[2:])

	fmt.Println(mapMinValue(seeds, maps))
	fmt.Println(mapMinRange(seeds, maps))
}

// vim: set ts=2 sw=2:
