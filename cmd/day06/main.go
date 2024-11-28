package main

import (
	"advent2023/util"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"unicode"
)

type Race struct {
	Time   int
	Record int
}

func parseRaces(lines []string) []Race {
	races := make([]Race, 0, 4)
	mode := ' '
	for _, line := range lines {
		index := 0
		for _, field := range strings.Fields(line) {
			value := 0
			if unicode.IsDigit(rune(field[0])) {
				value, _ = strconv.Atoi(field)
			} else {
				mode = rune(field[0])
				continue
			}
			if index >= len(races) {
				races = append(races, Race{})
			}
			if mode == 'T' {
				races[index].Time = value
			}
			if mode == 'D' {
				races[index].Record = value
			}
			index++
		}
	}
	return races
}

func parseOneRace(lines []string) Race {
	race := Race{}
	mode := ' '
	for _, line := range lines {
		index := 0
		modeStr, valueStr, _ := strings.Cut(line, " ")
		mode = rune(modeStr[0])
		valueStr = strings.ReplaceAll(valueStr, " ", "")
		value, _ := strconv.Atoi(valueStr)
		if mode == 'T' {
			race.Time = value
		}
		if mode == 'D' {
			race.Record = value
		}
		index++
	}
	return race
}

func numRecordBreakers(race Race) int {
	// solve distance (record) = t(T - t)
	// T = race time, t = speed = time held
	a := float64(race.Time) / 2
	b := math.Sqrt(float64(race.Time*race.Time-4*race.Record)) / 2
	lo := a - b
	hi := a + b
	loClamp := math.Ceil(lo)
	if loClamp == lo {
		loClamp++
	}
	hiClamp := math.Floor(hi)
	if hiClamp == hi {
		hiClamp--
	}
	sum := int(hiClamp) - int(loClamp) + 1
	//fmt.Printf("T=%d, D=%d, roots = %v, %v, sum=%d\n", race.Time, race.Record, lo, hi, sum)
	return sum
}

func main() {
	lines, err := util.ReadInputLines(6)
	if err != nil {
		log.Fatalf("%s", err)
	}
	races := parseRaces(lines)

	part1 := 1
	for _, race := range races {
		part1 = part1 * numRecordBreakers(race)
	}
	fmt.Println(part1)
	fmt.Println(numRecordBreakers(parseOneRace(lines)))
}

// vim: set ts=2 sw=2:
