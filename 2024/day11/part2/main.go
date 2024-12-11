package main

import (
	"fmt"
	"tea-cats.co.uk/aoc/2024"
	"tea-cats.co.uk/aoc/2024/day11"
	"time"
)

func main() {
	defer utils.TimeTrack(time.Now(), "main")

	knownSequence := day11.ResultCache{}
	data := day11.NewRequest(75, day11.LoadData())
	fmt.Printf("Result: %d\n", data.Process(&knownSequence))

	day11.CacheStats(knownSequence)
}
