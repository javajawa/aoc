package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"tea-cats.co.uk/aoc/2024"
	"time"
)

type request struct {
	target   uint64
	operands []uint64
}

func main() {
	defer utils.TimeTrack(time.Now(), "main")

	data := loadData()

	defer utils.TimeTrack(time.Now(), "process")

	validOptions := uint64(0)

	// Tracking information
	potentialTrees := 0 // How many operations could have happened
	trees := 0          // How many operations we evaluated
	maxValues := 0      // How many values we were tracking at once
	maxRowNum := 0      // The row on which the max values was reached

	for rowNum, row := range data {
		shiftFactors := make([]uint64, len(row.operands))
		potentialTrees += 3 * IntPow(3, len(row.operands)-1) / 2

		// Calculate the value for concatenating the values together.
		// This is always the largest operator:
		//   x || 9 = 10 * x + 9
		//   x * 9 < 10x + 9
		//   x + 9 < 10x + 9
		conAcc := uint64(0)
		for i, operand := range row.operands {
			// Calculate the multiplication factor for the || operator
			shiftFactors[i] = 1
			for buf := operand; buf > 0; buf /= 10 {
				shiftFactors[i] *= 10
			}

			conAcc = conAcc*shiftFactors[i] + operand
		}

		// Hey, if we get exactly the answer from concatenation, that's a free result
		//  (My data set includes 0 of these)
		if conAcc == row.target {
			validOptions += row.target
			continue
		}

		// And if we didn't make it to the target, that's a free negative result
		if conAcc < row.target {
			continue
		}

		minimums := make([]uint64, len(row.operands)-1)
		maximums := make([]uint64, len(row.operands)-1)
		minTarget := row.target
		maxTarget := row.target

		// Calculate the minimum and maximum values at each
		// position that could in theory still reach an answer
		minimums[len(row.operands)-2] = minTarget
		maximums[len(row.operands)-2] = maxTarget
		for i := len(row.operands) - 1; i > 1; i-- {
			minTarget /= shiftFactors[i]
			minimums[i-2] = minTarget

			// Special case: x*1 < x+1
			// In the event of a 1, we keep the same maximum as the next position to the right.
			if row.operands[i] > 1 {
				maxTarget -= row.operands[i]
			}
			maximums[i-2] = maxTarget
		}

		tracked := []uint64{row.operands[0]}

		for i, operand := range row.operands[1:] {
			toCheck := 3 * len(tracked)

			// Keep track of which row had the most allocations.
			trees += toCheck
			if toCheck > maxValues {
				maxValues = toCheck
				maxRowNum = rowNum
			}

			// Pre-allocate the slice for the values we find at this step
			out := make([]uint64, toCheck)
			accepted := 0

			// Get the minimum / maximum allocations
			minimum := minimums[i]
			maximum := maximums[i]

			for _, previous := range tracked {
				next := previous + operand
				if minimum <= next && next <= maximum {
					out[accepted] = next
					accepted++
				}

				next = previous * operand
				if minimum <= next && next <= maximum {
					out[accepted] = next
					accepted++
				}

				next = previous*shiftFactors[i+1] + operand
				if minimum <= next && next <= maximum {
					out[accepted] = next
					accepted++
				}
			}

			tracked = out[0:accepted]
		}

		if len(tracked) > 0 {
			validOptions += row.target
		}
	}

	fmt.Printf("Result: %d\n", validOptions)
	fmt.Printf("Evaluated %.1f%% of %d possible operations, with row %d having %d operations evaluated in one step\n", float64(trees)/float64(potentialTrees)*100, potentialTrees, maxRowNum, maxValues)
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
		length := len(fields) - 1

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

		requests = append(requests, request{target: target, operands: operands})
	}
	return requests
}

func IntPow(base, exp int) int {
	result := 1
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
