package main

import (
	"fmt"
	"tea-cats.co.uk/aoc/2024"
	"tea-cats.co.uk/aoc/2024/day13"
	"time"
)

const debug = false

func main() {
	defer utils.TimeTrack(time.Now(), "main")

	requests := day13.LoadData()
	score := 0

	defer utils.TimeTrack(time.Now(), "process")

	for _, test := range requests {
		// Initial Matrix:
		// / test.ButtonA.X   test.ButtonB.X \ / A \  _ / test.Target.X \
		// \ test.ButtonA.Y   test.ButtonB.Y / \ B /  - \ test.Target.Y /

		// So, we need to invert that matrix
		determinate := test.ButtonA.X*test.ButtonB.Y - test.ButtonB.X*test.ButtonA.Y
		expectedA := test.ButtonB.Y*test.Target.X - test.ButtonB.X*test.Target.Y
		expectedB := test.ButtonA.X*test.Target.Y - test.ButtonA.Y*test.Target.X

		if debug {
			fmt.Println("determinate:", determinate, "expected:", expectedA, "expected:", expectedB)
		}

		if expectedA%determinate == 0 && expectedB%determinate == 0 {
			expectedA = expectedA / determinate
			expectedB = expectedB / determinate
			s := 3*expectedA + expectedB
			if debug {
				fmt.Println("a:", expectedA, "b:", expectedB, "score:", s)
			}
			score += s
		}
	}

	fmt.Println(score)
}
