package main

import (
	"fmt"
	"math"
	"tea-cats.co.uk/aoc/2024"
	"time"
)

type stoneValue uint64
type blinkCount int
type stoneCount uint

type blinkResult []stoneCount

type resultCache map[stoneValue]blinkResult

const debug = false

func (cache *resultCache) getCachedCount(stone stoneValue, blinks blinkCount) stoneCount {
	if debug {
		fmt.Printf("Evalauting %d with %d blinks\n", stone, blinks)
	}

	if blinks <= 0 {
		return 1
	}

	cached := (*cache)[stone]
	iblinks := int(blinks)

	if iblinks < len(cached) {
		if cached[iblinks] > 0 {
			return cached[iblinks]
		}
	}

	newCache := make([]stoneCount, iblinks)
	copy(newCache, cached)
	newCache[0] = 1
	cached = newCache

	stones := stoneCount(1)

	if stone == 0 {
		if debug {
			fmt.Printf("0 -> 1\n")
		}
		stones = cache.getCachedCount(1, blinks-1)
		//copy(cached[1:], (*cache)[1][:blinks-1])
	} else if stone < 10 {
		if debug {
			fmt.Printf("x -> 2024*x\n")
		}
		target := stone * 2024
		stones = cache.getCachedCount(target, blinks-1)
		//copy(cached[1:], (*cache)[target][:blinks-1])
	} else {
		digits := int(math.Ceil(math.Log10(float64(stone))))

		if digits&1 == 1 {
			if debug {
				fmt.Printf("xx -> 2024*xx\n")
			}
			target := stone * 2024
			stones = cache.getCachedCount(target, blinks-1)
			//copy(cached[1:], (*cache)[target][:blinks-1])
		} else {
			split := digits >> 1

			right := stone % intPow(split)
			left := (stone - right) / intPow(split)

			if debug {
				fmt.Printf("%d -> %d %d\n", stone, left, right)
			}

			stones = cache.getCachedCount(left, blinks-1) + cache.getCachedCount(right, blinks-1)

			//lcache := (*cache)[left]
			//rcache := (*cache)[right]
			//
			//for i := 0; i < iblinks-1; i++ {
			//	cached[i+1] = lcache[i] + rcache[i]
			//}
		}
	}

	(*cache)[stone] = cached

	return stones
}

func main() {
	defer utils.TimeTrack(time.Now(), "main")

	const blinks = 25

	knownSequence := resultCache{}
	//data := []stoneValue{125, 17}
	data := []stoneValue{112, 1110, 163902, 0, 7656027, 83039, 9, 74}
	count := stoneCount(0)

	//slices.Sort(data)

	for _, val := range data {
		count += knownSequence.getCachedCount(val, blinks)
	}

	fmt.Println(count)
}

func intPow(exp int) stoneValue {
	base := stoneValue(10)
	result := stoneValue(1)
	for {
		if exp&1 == 1 {
			result *= base
		}
		exp >>= 1
		if exp == 0 {
			break
		}
		base *= base
	}

	return result
}
