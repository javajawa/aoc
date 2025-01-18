package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"tea-cats.co.uk/aoc/2024"
	"time"
)

func main() {
	defer utils.TimeTrack(time.Now(), "main")

	secrets := loadData()

	basket := make(map[uint32]int)

	for _, secret := range secrets {
		for key, value := range mapSecretDelta(secret, 2000) {
			basket[key] = basket[key] + value
		}
	}

	keys := make([]uint32, 0, len(basket))

	for key := range basket {
		keys = append(keys, key)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return basket[keys[i]] > basket[keys[j]]
	})

	for i := 0; i < 10; i++ {
		fmt.Printf("%d => %d\n", sequenceDecode(keys[i]), basket[keys[i]])
	}

	utils.PrintMemUsage()
}

const diffHistory = 4
const bitsPerDiff = 5
const negativeDiffFlag = 1 << (bitsPerDiff - 1)
const diffMask = negativeDiffFlag - 1
const fullDiffMask = 1<<(bitsPerDiff*diffHistory) - 1

func binRep(x int) uint32 {
	if x >= 0 {
		return uint32(x & diffMask)
	}
	return uint32((-x)&diffMask) | negativeDiffFlag
}

func sequenceDecode(val uint32) []int {
	numbers := make([]int, 4)

	for i := 3; i >= 0; i-- {
		numbers[i] = int(val & diffMask)
		if val&negativeDiffFlag == negativeDiffFlag {
			numbers[i] = -numbers[i]
		}
		val >>= 5
	}
	return numbers
}

func mapSecretDelta(secret uint64, rounds int) map[uint32]int {
	var newDigit int
	currentSequence := uint32(fullDiffMask)
	previousDigit := int(secret % 10)
	memory := make(map[uint32]int)

	for i := 0; i < rounds; i++ {
		secret ^= secret << 6
		secret &= 0xFFFFFF
		secret ^= secret >> 5
		secret ^= secret << 11
		secret &= 0xFFFFFF

		newDigit = int(secret % 10)

		currentSequence <<= bitsPerDiff
		currentSequence += binRep(newDigit - previousDigit)
		currentSequence &= fullDiffMask

		_, exists := memory[currentSequence]
		if !exists {
			memory[currentSequence] = newDigit
		}
		previousDigit = newDigit
	}
	return memory
}

func loadData() []uint64 {
	defer utils.TimeTrack(time.Now(), "loadData")
	dataFile, err := os.Open("2024/input-22.txt")

	if err != nil {
		panic(err)
	}

	defer utils.CloseWithLog(dataFile)

	var next uint64
	secrets := make([]uint64, 0, 2200)

	for {
		_, err := fmt.Fscanf(dataFile, "%d\n", &next)

		if err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
				break
			}
			panic(err)
		}
		secrets = append(secrets, next)
	}
	return secrets
}
