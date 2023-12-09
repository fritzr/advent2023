package main

import (
	"advent2023/util"
	"fmt"
	"log"
	"strings"
)

func parseMaps(lines []string) []util.RangeMap {
	maps := make([]util.RangeMap, 0)
	var newMap *util.RangeMap
	for _, line := range lines {
		if line == "" {
			continue
		}
		if strings.HasSuffix(line, ":") { // new map:
			maps = append(maps, util.RangeMap{})
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

func mapMinValue(seeds []int, seedMap util.RangeMap) int {
	minValue := seedMap.Map(seeds[0])
	for _, seedValue := range seeds[1:] {
		mappedValue := seedMap.Map(seedValue)
		if mappedValue < minValue {
			minValue = mappedValue
		}
	}
	return minValue
}

func dumpMap(seeds []int, seedMap util.RangeMap) []int {
	fmt.Printf("test:\n")
	for _, seed := range seeds {
		value := seedMap.Map(seed)
		fmt.Printf("  %d=>%d", seed, value)
	}
	fmt.Printf("\n")
	return seeds
}

func Reduce(seeds []int, maps []util.RangeMap) util.RangeMap {
	result := maps[0]
	seeds = dumpMap(seeds, result)
	for _, rangeMap := range maps[1:] {
		fmt.Printf("map(\n     %s,\n     %s\n", result, rangeMap)
		newMap := result.CombineMap(rangeMap)
		fmt.Printf("  => %s\n)\n", newMap)
		seeds = dumpMap(seeds, newMap)
		result = newMap
	}
	return result
}

func mapMinRange(seeds []int, seedMap util.RangeMap) int {
	seedSet := util.RangeMap{}
	for index := 0; index < len(seeds); index += 2 {
		seedSet.Add(util.Span{seeds[index], seeds[index] + seeds[index+1]}, 0)
	}
	fmt.Println("seeds:")
	dumpMap(seeds, seedSet)
	seedMap = seedSet.CombineMap(seedMap)
	minValue := 0
	index := 0
	util.RangeSet[int](seedMap).Do(func(s util.Span, delta *int) bool {
		value := s[0] + *delta
		fmt.Printf("  %s,%+d => %d\n", s, *delta, value)
		if index == 0 || (seedSet.Maps(s[0]) && value < minValue) {
			minValue = value
		}
		index++
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

	seedMap := Reduce(seeds, maps)
	fmt.Println(mapMinValue(seeds, seedMap))
	fmt.Println(mapMinRange(seeds, seedMap))
}

// vim: set ts=2 sw=2:
