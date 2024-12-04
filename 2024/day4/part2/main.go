package main

import (
	"bufio"
	"fmt"
	"os"
	"tea-cats.co.uk/aoc/2024"
	"time"
)

func main() {
	defer utils.TimeTrack(time.Now(), "main")

	fmt.Println("Answer:", loadData())
}

const size = 140
const forward = uint32('M'<<16) + uint32('A'<<8) + uint32('S')
const backwards = uint32('S'<<16) + uint32('A'<<8) + uint32('M')

func loadData() int {
	defer utils.TimeTrack(time.Now(), "loadData")
	dataFile, err := os.Open("2024/input-4.txt")

	if err != nil {
		panic(err)
	}

	defer utils.CloseWithLog(dataFile)

	var diagToRight [size + size - 1]uint32
	var diagToLeft [size + size - 1]uint32

	reader := bufio.NewScanner(dataFile)
	row := 0
	matches := 0

	for reader.Scan() {
		line := reader.Bytes()

		for col, char := range line {
			c := uint32(char)

			// Counting diagonals from the top.
			// The ones which go from top-right to bottom-left are indexed by row+col
			// The ones which go from top-left to bottom-right are indexed by rol-col, with the first one
			// being the middle
			diagLeftId := row + col
			diagRightId := (size - 1) + (row - col)

			diagToLeft[diagLeftId] = (diagToLeft[diagLeftId] << 8) + c
			diagToRight[diagRightId] = (diagToRight[diagRightId] << 8) + c

			// We can't get a full X in the first two columns
			if col < 2 {
				continue
			}

			toRight := diagToRight[diagRightId] & 0x00FFFFFF
			// The to-left is two letter to the left of the to-right
			toLeft := diagToLeft[diagLeftId-2] & 0x00FFFFFF

			if toRight == backwards || toRight == forward {
				if toLeft == backwards || toLeft == forward {
					matches++
				}
			}
		}

		row++
	}

	return matches
}
