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

type HandRanker struct {
	cardOrder     string
	cardRanks     map[rune]int
	handRankCache map[string]HandRank
	jokerCard     rune
}

func NewHandRanker(cardOrder string, joker rune) HandRanker {
	return HandRanker{
		cardOrder:     cardOrder,
		cardRanks:     generateCardRanks(cardOrder),
		handRankCache: map[string]HandRank{},
		jokerCard:     joker,
	}
}

// lowest to highest face value
var ranker = NewHandRanker("23456789TJQKA", ' ')
var jokerRanker = NewHandRanker("J23456789TQKA", 'J')

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

func (r HandRanker) LookupRank(h Hand) HandRank {
	cacheValue, ok := r.handRankCache[h.Hand]
	if ok {
		return cacheValue
	}
	rank := r.Rank(h)
	r.handRankCache[h.Hand] = rank
	return rank
}

func (r HandRanker) Rank(h Hand) HandRank {
	counts := map[rune]int{}
	for _, card := range h.Hand {
		counts[card]++
	}
	highCount := 0
	secondHighCount := 0
	for card, count := range counts {
		if card == r.jokerCard {
			continue
		}
		if count >= highCount {
			if highCount > secondHighCount {
				secondHighCount = highCount
			}
			highCount = count
		} else if count > secondHighCount {
			secondHighCount = count
		}
	}
	highCount += counts[r.jokerCard]
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

func (r HandRanker) Compare(h1, h2 Hand) int {
	rank1 := r.LookupRank(h1)
	rank2 := r.LookupRank(h2)
	//fmt.Printf("Rank(%s) -> %v\n", h1.Hand, rank1)
	//fmt.Printf("Rank(%s) -> %v\n", h2.Hand, rank2)
	if rank1 != rank2 {
		return int(rank1) - int(rank2)
	}
	index := 0
	for index < len(h1.Hand) {
		diff := r.cardRanks[rune(h1.Hand[index])] - r.cardRanks[rune(h2.Hand[index])]
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
	slices.SortFunc(hands, ranker.Compare)
	part1 := 0
	for index, hand := range hands {
		part1 += (index + 1) * hand.Bid
	}
	fmt.Println(part1)
	slices.SortFunc(hands, jokerRanker.Compare)
	part2 := 0
	for index, hand := range hands {
		part2 += (index + 1) * hand.Bid
	}
	fmt.Println(part2)
}
