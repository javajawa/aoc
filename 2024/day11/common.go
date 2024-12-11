package day11

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"tea-cats.co.uk/aoc/2024"
	"time"
)

type StoneValue uint64
type BlinkCount int
type StoneCount uint

type BlinkResult []StoneCount

type ResultCache map[StoneValue]BlinkResult

var PowersOfTen = [...]StoneValue{1, 10, 100, 1e3, 1e4, 1e5, 1e6, 1e7, 1e8, 1e9, 1e10, 1e11, 1e12, 1e13, 1e14, 1e15, 1e16}
var CacheHits = uint64(0)
var CacheMisses = uint64(0)

const Debug = false

func (cache *ResultCache) GetCachedCount(stone StoneValue, blinks BlinkCount) StoneCount {
	if blinks <= 0 {
		return 1
	}

	cached := (*cache)[stone]
	iblinks := int(blinks)

	if iblinks < len(cached) {
		if cached[iblinks] > 0 {
			CacheHits++
			return cached[iblinks]
		}
	}

	if Debug {
		fmt.Printf("Evalauting %d with %d blinks\n", stone, blinks)
	}
	CacheMisses++

	newCache := make([]StoneCount, iblinks)
	copy(newCache, cached)
	newCache[0] = 1
	cached = newCache

	stones := StoneCount(1)

	if stone == 0 {
		if Debug {
			fmt.Printf("0 -> 1\n")
		}
		stones = cache.GetCachedCount(1, blinks-1)
		copy(cached[1:], (*cache)[1][:blinks-1])
	} else if stone < 10 {
		if Debug {
			fmt.Printf("x -> 2024*x\n")
		}
		target := stone * 2024
		stones = cache.GetCachedCount(target, blinks-1)
		copy(cached[1:], (*cache)[target][:blinks-1])
	} else {
		digits := 0
		for ; stone >= PowersOfTen[digits]; digits++ {
		}

		if Debug {
			fmt.Printf("%d has digits: %d\n", stone, digits)
		}

		if digits&1 == 1 {
			if Debug {
				fmt.Printf("xx -> 2024*xx\n")
			}
			target := stone * 2024
			stones = cache.GetCachedCount(target, blinks-1)
			copy(cached[1:], (*cache)[target][:blinks-1])
		} else {
			split := digits >> 1
			if Debug {
				fmt.Printf("splitting at %d\n", split)
			}

			right := stone % PowersOfTen[split]
			left := (stone - right) / PowersOfTen[split]

			if Debug {
				fmt.Printf("%d -> %d %d\n", stone, left, right)
			}

			stones = cache.GetCachedCount(left, blinks-1) + cache.GetCachedCount(right, blinks-1)

			lcache := (*cache)[left]
			rcache := (*cache)[right]

			for i := 0; i < iblinks-1; i++ {
				cached[i+1] = lcache[i] + rcache[i]
			}
		}
	}

	(*cache)[stone] = cached

	return stones
}

type Request struct {
	blinks BlinkCount
	stones []StoneValue
}

func NewRequest(blinks BlinkCount, stones []StoneValue) *Request {
	return &Request{blinks: blinks, stones: stones}
}

func (r *Request) Process(knownSequence *ResultCache) StoneCount {
	defer utils.TimeTrack(time.Now(), "process")

	count := StoneCount(0)

	for _, val := range r.stones {
		count += knownSequence.GetCachedCount(val, r.blinks)
	}

	return count
}

func CacheStats(knownSequence ResultCache) {
	fmt.Printf("\n=== Stats ===\n")
	cachedPoints := 0
	maxValue := StoneValue(0)
	for value, cached := range knownSequence {
		if value > maxValue {
			maxValue = value
		}
		cachedPoints = cachedPoints + len(cached)
	}
	fmt.Printf("Total stones seen:   %16d\n", len(knownSequence))
	fmt.Printf("Largest stone value: %16d\n", maxValue)
	fmt.Printf("Total points known:  %16d\n", cachedPoints)
	fmt.Printf("Mean cached points:  %16.1f\n", float64(cachedPoints)/float64(len(knownSequence)))
	fmt.Printf("Cache Hits:          %16d\n", CacheHits)
	fmt.Printf("Cache Misses:        %16d\n", CacheMisses)
	fmt.Printf("Cache Hit Rate      %16.1f%%\n", float64(100*CacheHits)/float64(CacheMisses+CacheHits))
	fmt.Printf("\n")
}

func LoadData() []StoneValue {
	defer utils.TimeTrack(time.Now(), "loadData")
	dataFile, err := os.Open("2024/input-11.txt")

	if err != nil {
		panic(err)
	}

	defer utils.CloseWithLog(dataFile)

	reader := bufio.NewReader(dataFile)
	line, err := reader.ReadString('\n')

	if err != nil {
		panic(err)
	}

	strData := strings.Fields(line)
	intData := make([]StoneValue, len(strData))

	for i, d := range strData {
		x, err := strconv.ParseUint(d, 10, 64)
		if err != nil {
			panic(err)
		}
		intData[i] = StoneValue(x)
	}

	return intData
}
