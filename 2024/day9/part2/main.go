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

type fileSpec struct {
	fileId uint64
	size   int
}

type priorityQueue struct {
	sourceBuffer []int
	readHead     uint64
	queues       [9][]fileSpec
}

func newPriorityQueue(buffer []int) *priorityQueue {
	var queues [9][]fileSpec

	for i, _ := range queues {
		queues[i] = make([]fileSpec, 0)
	}

	return &priorityQueue{
		sourceBuffer: buffer,
		readHead:     uint64(len(buffer)-1) & ^uint64(1),
		queues:       queues,
	}
}

func (queues *priorityQueue) get(spaceSize int, after uint64) fileSpec {
	matched := fileSpec{
		fileId: 0,
		size:   0,
	}

	// Find the highest ID file that fits in the space
	for i := 0; i < spaceSize; i++ {
		// Skip empty queues
		if len(queues.queues[i]) == 0 {
			continue
		}

		potentialFileId := queues.queues[i][0].fileId
		// If the first item in a queue is from before our current position,
		// we can drop that entire queue.
		if potentialFileId < (after >> 1) {
			queues.queues[i] = []fileSpec{}
		} else if potentialFileId > matched.fileId {
			// Otherwise, keep it if it's file from further on in the disk
			matched = fileSpec{fileId: potentialFileId, size: i + 1}
		}
	}

	// If the matched file is from further on in the disk than we have de-fragged...
	if matched.fileId > (after >> 1) {
		// Mark the file as moved
		queues.sourceBuffer[matched.fileId<<1] = -matched.size
		// Remove it from the queue
		queues.queues[matched.size-1] = queues.queues[matched.size-1][1:]
		// And give it back to the moving function
		return matched
	}

	// Otherwise, scan the filesystem from right to left, filling queues
	// until we either reach the left-to-right process or a file we can move.
	for ; queues.readHead > after; queues.readHead -= 2 {
		size := queues.sourceBuffer[queues.readHead]
		matched = fileSpec{fileId: queues.readHead >> 1, size: size}

		if size <= spaceSize {
			queues.sourceBuffer[queues.readHead] = -queues.sourceBuffer[queues.readHead]
			queues.readHead -= 2
			return matched
		}

		queues.queues[size-1] = append(queues.queues[size-1], matched)
	}

	return fileSpec{fileId: 0}
}

func main() {
	defer utils.TimeTrack(time.Now(), "main")

	buffer := loadData()

	defer utils.TimeTrack(time.Now(), "process")

	checksum := uint64(0)

	queue := newPriorityQueue(buffer)
	initialReadIndex := queue.readHead

	// Counter for the number of blocks we have put back on the disk
	currentDiskBlock := uint64(0)

	for scanPosition := uint64(0); scanPosition <= initialReadIndex; scanPosition++ {
		// Even entries represent a file
		if scanPosition%2 == 0 {
			fileLength := buffer[scanPosition]

			// Skip files that have been moved in the defragmentation process
			if fileLength < 0 {
				if debug {
					fmt.Printf("%s", strings.Repeat("x", -fileLength))
				}
				currentDiskBlock += uint64(-fileLength)
				continue
			}

			// Write out the file in the same position as it originally was
			updateChecksum(&checksum, &currentDiskBlock, fileSpec{fileId: scanPosition >> 1, size: fileLength})
		} else {
			spaceLength := buffer[scanPosition]

			for spaceLength > 0 {
				fileToMove := queue.get(spaceLength, scanPosition)

				if fileToMove.fileId == 0 {
					break
				}

				// Write out the relocated file
				updateChecksum(&checksum, &currentDiskBlock, fileToMove)
				// Update our remaining free space
				spaceLength -= fileToMove.size
			}

			if debug {
				fmt.Printf("%s", strings.Repeat(".", spaceLength))
			}
			currentDiskBlock += uint64(spaceLength)
		}
	}

	fmt.Printf("\nChecksum: %d\n", checksum)
	fmt.Printf("Final locations: block=%d\n", currentDiskBlock)
}

func updateChecksum(checksum *uint64, currentBlock *uint64, file fileSpec) {
	// Sum of ints is n(n+1)/2. We want end-start of that.
	// [(start+len)(start+len+1) - (start)(start + 1)]/2
	// [(s^2 + sl + s + l^2 + sl + l) - (s^2 + s)]/2
	// [(      sl     + l^2 + sl + l)]/2
	// (l^2 + 2sl + l)/2
	// l(l + 2s + 1)/2
	//
	// But, something-something (start+len-1) so we end up with an off by two error.
	if debug {
		fmt.Printf("%s", strings.Repeat(strconv.FormatUint(file.fileId%10, 10), file.size))
	}

	blocksUint := uint64(file.size)
	*checksum += file.fileId * blocksUint * ((*currentBlock << 1) + blocksUint - 1) / 2
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
