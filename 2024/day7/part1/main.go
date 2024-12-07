package main

import (
	"bufio"
	"fmt"
	"io"
	"math/bits"
	"os"
	"strconv"
	"strings"
	"tea-cats.co.uk/aoc/2024"
	"time"
)

type request struct {
	target   uint64
	operands []uint64
	length   uint16
}

func main() {
	defer utils.TimeTrack(time.Now(), "main")

	data := loadData()

	defer utils.TimeTrack(time.Now(), "process")

	const debug = false

	validOptions := uint64(0)
	variationsConsidered := uint64(0)
	variationsPossible := uint64(0)

nextNumber:
	for _, row := range data {
		totalPermutations := uint64(1) << (row.length - 1)
		variationsPossible += totalPermutations
	nextPermutation:
		for permutation := uint64(0); permutation < totalPermutations; permutation++ {
			variationsConsidered++
			rowAccumulator := row.operands[0]
			// We flip the order so that a known sequence e.g. 010101xxxxxx
			// in the permutations is processed as xxxxxx010101 by the binary processing logic.
			// This means that if we exceed the target value with the first set of operations,
			// we can easily prune all operations that start with that sequence by skipping the
			// rest of that block.
			runPermutation := permutation
			runPermutation = bits.Reverse64(permutation) >> (65 - row.length)

			if debug {
				fmt.Printf("%v variation %d\n", row, permutation)
				fmt.Printf("  x = %d\n", rowAccumulator)
			}

			for field := uint16(1); field < row.length; field++ {
				if (runPermutation & 1) != 0 {
					if debug {
						fmt.Printf("  x = %d + %d = %d\n", rowAccumulator, row.operands[field], rowAccumulator+row.operands[field])
					}
					rowAccumulator = rowAccumulator + row.operands[field]
				} else {
					if debug {
						fmt.Printf("  x = %d * %d = %d\n", rowAccumulator, row.operands[field], rowAccumulator*row.operands[field])
					}
					rowAccumulator = rowAccumulator * row.operands[field]
				}
				if rowAccumulator > row.target {
					// We're reading the binary string right -> left
					// But it is the inversion of the outer loop
					// So we are effectively reading `permutation` left -> right
					fieldInOriginalPermutation := row.length - field - 1

					// After `field` fields, we are out of bounds. Any combination of later fields
					// will also terminate here. So we can prune all those branches.
					// We achieve this by setting all the remaining bits high.
					permutation = permutation | ((1 << fieldInOriginalPermutation) - 1)

					// When iterating the loop, one more will be added to `permutation`, taking us
					// out of this branch
					continue nextPermutation
				}
				runPermutation = runPermutation >> 1
			}

			if rowAccumulator == row.target {
				if debug {
					fmt.Printf("Valid!\n")
				}
				validOptions += row.target
				continue nextNumber
			}
		}
	}

	fmt.Printf("Result: %d\n", validOptions)
	fmt.Printf("Considered %d variations (%.1f%% of possible variations)\n", variationsConsidered, 100*(float64(variationsConsidered)/float64(variationsPossible)))
}

func loadData() []request {
	defer utils.TimeTrack(time.Now(), "loadData")
	dataFile, err := os.Open("2024/input-7.txt")

	if err != nil {
		panic(err)
	}

	defer utils.CloseWithLog(dataFile)

	requests := make([]request, 0)

	scanner := bufio.NewReader(dataFile)

	for {
		buffer, err := scanner.ReadString('\n')

		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		fields := strings.Fields(buffer)
		length := uint16(len(fields) - 1)

		if length > 63 {
			panic("invalid input -- too many operands")
		}

		target, err := strconv.ParseUint(fields[0][:len(fields[0])-1], 10, 64)
		if err != nil {
			panic(err)
		}
		operands := make([]uint64, length)

		for fieldNo, operand := range fields[1:] {
			operands[fieldNo], err = strconv.ParseUint(operand, 10, 64)
			if err != nil {
				panic(err)
			}
		}
		//fmt.Println(requests[i])

		requests = append(requests, request{target: target, operands: operands, length: length})
	}
	return requests
}
