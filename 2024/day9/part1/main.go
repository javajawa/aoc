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

	// Find the last file in the input (has to be an even index)
	currentReadIndex = uint64(len(buffer)-1) & ^uint64(1)

	// Counter for the number of blocks we have put back on the disk
	currentDiskBlock := uint64(0)

	for currentWriteIndex = 0; currentWriteIndex <= currentReadIndex; currentWriteIndex++ {
		// Even entries represent a file, which we will not move
		if currentWriteIndex%2 == 0 {
			updateChecksum(&checksum, &currentDiskBlock, currentWriteIndex>>1, buffer[currentWriteIndex])
			continue
		}

		// Odd entries are spaces we are fragmenting files into
		for spaceLength := buffer[currentWriteIndex]; spaceLength > 0 && currentReadIndex > currentWriteIndex; currentReadIndex -= 2 {
			fileLength := buffer[currentReadIndex]

			// Fragment the file if needed
			writeLength := fileLength
			if fileLength > spaceLength {
				writeLength = spaceLength
			}

			// Write this chunk
			updateChecksum(&checksum, &currentDiskBlock, currentReadIndex>>1, writeLength)
			spaceLength -= writeLength

			// Keep state of the partially fragmented file
			if fileLength > writeLength {
				buffer[currentReadIndex] -= writeLength
				break
			}
		}
	}

	fmt.Printf("\nChecksum: %d\n", checksum)
	fmt.Printf("Final locations: write=%d, reader=%d, block=%d\n", currentWriteIndex, currentReadIndex, currentDiskBlock)
}

func updateChecksum(checksum *uint64, currentBlock *uint64, fileId uint64, blocks byte) {
	// Sum of ints is n(n+1)/2. We want end-start of that.
	// [(start+len)(start+len+1) - (start)(start + 1)]/2
	// [(s^2 + sl + s + l^2 + sl + l) - (s^2 + s)]/2
	// [(      sl     + l^2 + sl + l)]/2
	// (l^2 + 2sl + l)/2
	// l(l + 2s + 1)/2
	//
	// But, something-something (start+len-1) so we end up with an off by two error.
	if debug {
		fmt.Printf("%s", strings.Repeat(strconv.FormatUint(fileId, 10), int(blocks)))
	}

	blocksUint := uint64(blocks)
	*checksum += fileId * blocksUint * ((*currentBlock << 1) + blocksUint - 1) / 2
	*currentBlock += blocksUint
}

func loadData() []byte {
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

	for i, c := range buffer {
		buffer[i] = c - '0'
	}

	return buffer
}
