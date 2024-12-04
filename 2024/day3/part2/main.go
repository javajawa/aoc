package main

import (
	"fmt"
	"io"
	"os"
	"tea-cats.co.uk/aoc/2024"
	"time"
)

func main() {
	defer utils.TimeTrack(time.Now(), "main")

	data := loadData()
	result := parse(data)

	fmt.Printf("Answer: %d\n", result)
}

func loadData() []byte {
	defer utils.TimeTrack(time.Now(), "loadData")
	dataFile, err := os.Open("2024/input-3.txt")

	if err != nil {
		panic(err)
	}

	defer utils.CloseWithLog(dataFile)

	data, err := io.ReadAll(dataFile)

	if err != nil {
		panic(err)
	}

	return data
}

func parse(input []byte) int {
	defer utils.TimeTrack(time.Now(), "parse")

	acc := 0
	active := true

	// Minimum length of `mul(0,0)` is 8 bytes
	stop := len(input) - 8

outer:
	for i := 0; i < stop; i++ {
		if input[i] == 'd' {
			// Doing look aheads here so we don't swallow an `m`.
			if string(input[i:i+4]) == "do()" {
				active = true
				i += 3
			} else if string(input[i:i+7]) == "don't()" {
				active = false
				i += 6
			}
			continue
		}

		if input[i] != 'm' {
			continue
		}
		i++
		if input[i] != 'u' {
			continue
		}
		i++
		if input[i] != 'l' {
			continue
		}
		i++
		if input[i] != '(' {
			continue
		}

		leftOperand := 0
		leftDigits := 0

		for leftDigits = 0; leftDigits <= 3; leftDigits++ {
			i++
			c := input[i]

			if c == ',' {
				break
			}

			if c < '0' || c > '9' {
				continue outer
			}

			leftOperand = leftOperand*10 + int(c-'0')
		}

		if input[i] != ',' {
			continue
		}

		rightOperand := 0
		rightDigits := 0

		for rightDigits = 0; rightDigits <= 3; rightDigits++ {
			i++
			c := input[i]

			if c == ')' {
				break
			}

			if c < '0' || c > '9' {
				continue outer
			}

			rightOperand = rightOperand*10 + int(c-'0')
		}

		if input[i] != ')' {
			continue
		}

		if active {
			acc += leftOperand * rightOperand
		}
	}

	return acc
}
