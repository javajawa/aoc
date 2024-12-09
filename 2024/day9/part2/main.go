package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"tea-cats.co.uk/aoc/2024"
	"time"
)

const debug = false

func main() {
	defer utils.TimeTrack(time.Now(), "main")

	buffer := loadData()

	defer utils.TimeTrack(time.Now(), "process")

	checksum := uint64(0)

	var currentWriteIndex, currentReadIndex uint64

	initialReadIndex := uint64(len(buffer)-1) & ^uint64(1)
	currentReadIndex = initialReadIndex

	// Counter for the number of blocks we have put back on the disk
	currentDiskBlock := uint64(0)

	for currentWriteIndex = 0; currentWriteIndex <= initialReadIndex; currentWriteIndex++ {
		// Even entries represent a file
		if currentWriteIndex%2 == 0 {
			fileLength := buffer[currentWriteIndex]

			// Skip files that have been moved in the defragmentation process
			if fileLength < 0 {
				if debug {
					fmt.Printf("%s", strings.Repeat("x", -fileLength))
				}
				currentDiskBlock += uint64(-fileLength)
				continue
			}

			updateChecksum(&checksum, &currentDiskBlock, currentWriteIndex>>1, fileLength)
		} else {
			spaceLength := buffer[currentWriteIndex]
			currentReadIndex = initialReadIndex

			for ; spaceLength > 0 && currentReadIndex > currentWriteIndex; currentReadIndex -= 2 {
				fileLength := buffer[currentReadIndex]

				// Do not move files that have already been moved, or won't fit in this space.
				if fileLength < 0 || fileLength > spaceLength {
					continue
				}

				// Write out the relocated file
				updateChecksum(&checksum, &currentDiskBlock, currentReadIndex>>1, fileLength)

				// Update our remaining free space, and mark the file as moved by making the space negative
				spaceLength -= fileLength
				buffer[currentReadIndex] = -fileLength
			}

			if debug {
				fmt.Printf("%s", strings.Repeat(".", spaceLength))
			}
			currentDiskBlock += uint64(spaceLength)
		}
	}

	fmt.Printf("\nChecksum: %d\n", checksum)
	fmt.Printf("Final locations: write=%d, reader=%d, block=%d\n", currentWriteIndex, currentReadIndex, currentDiskBlock)
}

func updateChecksum(checksum *uint64, currentBlock *uint64, fileId uint64, blocks int) {
	// Sum of ints is n(n+1)/2. We want end-start of that.
	// [(start+len)(start+len+1) - (start)(start + 1)]/2
	// [(s^2 + sl + s + l^2 + sl + l) - (s^2 + s)]/2
	// [(      sl     + l^2 + sl + l)]/2
	// (l^2 + 2sl + l)/2
	// l(l + 2s + 1)/2
	//
	// But, something-something (start+len-1) so we end up with an off by two error.
	if debug {
		fmt.Printf("%s", strings.Repeat(strconv.FormatUint(fileId%10, 10), blocks))
	}

	blocksUint := uint64(blocks)
	*checksum += fileId * blocksUint * ((*currentBlock << 1) + blocksUint - 1) / 2
	*currentBlock += blocksUint
}

func loadData() []int {
	defer utils.TimeTrack(time.Now(), "loadData")
	dataFile, err := os.Open("2024/input-9.txt")

	if err != nil {
		panic(err)
	}

	defer utils.CloseWithLog(dataFile)

	buffer, err := io.ReadAll(dataFile)
	if err != nil {
		panic(err)
	}

	result := make([]int, len(buffer))

	for i, c := range buffer {
		result[i] = int(c - '0')
	}

	return result
}
