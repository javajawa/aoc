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
const forward uint32 = ('X' << 24) + ('M' << 16) + ('A' << 8) + ('S')
const backwards uint32 = ('S' << 24) + ('A' << 16) + ('M' << 8) + ('X')

func loadData() int {
	defer utils.TimeTrack(time.Now(), "loadData")
	dataFile, err := os.Open("2024/input-4.txt")

	if err != nil {
		panic(err)
	}

	defer utils.CloseWithLog(dataFile)

	// We only need to track the current row for row-based matching
	var currentRow uint32
	// We have to keep track of each column for column matching
	var cols [size]uint32
	// There are 2n-1 diagonals in each direction -- you can visualise this
	// as there being one full length diagonal from corner to corner,
	// and moving outwards from there counting down to 0.
	//
	// For this code, "to right" is "diagonals going from top-left to bottom-right",
	// and "to left" is "diagonals going from top-right to bottom-left".
	var diagToRight [size + size - 1]uint32
	var diagToLeft [size + size - 1]uint32

	reader := bufio.NewScanner(dataFile)

	// Total number of XMASes in the word search
	matches := 0
	// Current row ID
	row := 0

	for reader.Scan() {
		// Read the next line from the file
		line := reader.Bytes()

		// Reset the current row state
		currentRow = 0

		for col, character := range line {
			if check(&currentRow, character) {
				matches++
			}

			if check(&cols[col], character) {
				matches++
			}

			// The top-right to bottom-left diagonals are
			// indexed by how far away from the top left they are
			//  1 2 3 4
			//  2 3 4 5
			//  3 4 5 6
			//  4 5 6 7
			if check(&diagToLeft[row+col], character) {
				matches++
			}

			// The top-left to bottom-right diagonals are
			// indexed by how far away from the center line they are
			// 4 3 2 1
			// 5 4 3 2
			// 6 5 4 3
			// 7 6 5 4
			if check(&diagToRight[(size-1)+(row-col)], character) {
				matches++
			}
		}

		row++
	}

	return matches
}

func check(val *uint32, new byte) bool {
	*val = (*val << 8) + uint32(new)

	return *val == backwards || *val == forward
}
