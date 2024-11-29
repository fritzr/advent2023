package main

import (
	"advent2023/util"
	"fmt"
	"log"
	"slices"
	"strconv"
	"strings"
)

// generate reverse mapping of card face value to rank (low to high)
func generateCardRanks(cardOrder string) map[rune]int {
	ranks := make(map[rune]int, len(cardOrder))
	rank := 1
	for _, value := range cardOrder {
		ranks[value] = rank
		rank++
	}
	return ranks
}

// lowest to highest face value
const cardOrder = "23456789TJQKA"

// map face value to numeric rank
var cardRanks = generateCardRanks(cardOrder)

type HandRank int

const RankHighCard HandRank = 0
const RankPair HandRank = 1
const RankTwoPair HandRank = 2
const RankTriple HandRank = 3
const RankFullHouse HandRank = 4
const RankQuad HandRank = 5
const RankPenta HandRank = 6

type Hand struct {
	Hand string
	Bid  int
}

var rankCache = map[string]HandRank{}

func (h Hand) Rank() HandRank {
	cacheValue, ok := rankCache[h.Hand]
	if ok {
		return cacheValue
	}
	counts := map[rune]int{}
	for _, value := range h.Hand {
		counts[value]++
	}
	highCount := 0
	secondHighCount := 0
	for _, count := range counts {
		if count >= highCount {
			if highCount > secondHighCount {
				secondHighCount = highCount
			}
			highCount = count
		} else if count > secondHighCount {
			secondHighCount = count
		}
	}
	switch highCount {
	case 2:
		if secondHighCount == 2 {
			return RankTwoPair
		}
		return RankPair
	case 3:
		if secondHighCount == 2 {
			return RankFullHouse
		}
		return RankTriple
	case 4:
		return RankQuad
	case 5:
		return RankPenta
	case 1:
		fallthrough
	default:
		return RankHighCard
	}
}

func (h1 Hand) Compare(h2 Hand) int {
	rank1 := h1.Rank()
	rank2 := h2.Rank()
	//fmt.Printf("Rank(%s) -> %v\n", h1.Hand, rank1)
	//fmt.Printf("Rank(%s) -> %v\n", h2.Hand, rank2)
	if rank1 != rank2 {
		return int(rank1) - int(rank2)
	}
	index := 0
	for index < len(h1.Hand) {
		diff := cardRanks[rune(h1.Hand[index])] - cardRanks[rune(h2.Hand[index])]
		if diff != 0 {
			//fmt.Printf("Compare(%s, %s) -> %d\n", h1.Hand, h2.Hand, diff)
			return diff
		}
		index++
	}
	// panic(fmt.Sprintf("equal hands %s %s", h1.Hand, h2.Hand))
	return 0
}

func readHands() ([]Hand, error) {
	lines, err := util.ReadInputLines(7)
	if err != nil {
		return nil, fmt.Errorf("opening input file: %w", err)
	}

	hands := make([]Hand, 0, len(lines))
	for _, line := range lines {
		hand, bid, ok := strings.Cut(line, " ")
		if ok {
			bidValue, err := strconv.Atoi(bid)
			if err != nil {
				continue
			}
			hands = append(hands, Hand{Hand: hand, Bid: bidValue})
		}
	}
	return hands, nil
}

func main() {
	hands, err := readHands()
	if err != nil {
		log.Fatalf("parsing hands: %s", err)
	}
	slices.SortFunc(hands, Hand.Compare)
	part1 := int64(0)
	for index, hand := range hands {
		part1 += int64(index+1) * int64(hand.Bid)
	}
	fmt.Println(part1)
}
