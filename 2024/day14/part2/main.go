package main

import (
	"errors"
	"fmt"
	"image"
	"io"
	"os"
	"tea-cats.co.uk/aoc/2024"
	"time"
)

type robot struct {
	initial  image.Point
	movement image.Point
}

const debug = false

func (r *robot) finalPosition(grid image.Rectangle, seconds int) image.Point {
	return r.initial.Add(r.movement.Mul(seconds)).Mod(grid)
}

func main() {
	defer utils.TimeTrack(time.Now(), "main")

	robots := loadData()
	grid := image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: 101, Y: 103}}
	center := grid.Max.Sub(grid.Min).Sub(image.Point{1, 1}).Div(2)

	fmt.Printf("Grid=%v, Center=%v\n", grid, center)

nextSecond:
	for second := 0; second < 101*103; second++ {
		counts := make(map[image.Point]int)
		centerCol := 0
		lefts := 0
		rights := 0

		for i, robot := range robots {
			final := robot.finalPosition(grid, second)
			if debug {
				fmt.Printf("Robot %d ends at %v\n", i, final)
			}

			old, ok := counts[final]
			if ok {
				if debug {
					fmt.Printf("Time %d: Robot %d (%v) is standing on %d (%v)!\n", second, i, robot, old, robots[old])
				}
				continue nextSecond
			}
			counts[final] = i

			if final.X == center.X {
				centerCol++
			} else if final.X < center.X {
				lefts++
			} else {
				rights++
			}
		}

		fmt.Printf("No stacking at t=%d\n", second)
	}
}

func loadData() []robot {
	defer utils.TimeTrack(time.Now(), "loadData")
	dataFile, err := os.Open("2024/input-14.txt")

	if err != nil {
		panic(err)
	}

	defer utils.CloseWithLog(dataFile)

	robots := make([]robot, 0)

	for {
		r := robot{}

		count, err := fmt.Fscanf(dataFile, "p=%d,%d v=%d,%d\n", &r.initial.X, &r.initial.Y, &r.movement.X, &r.movement.Y)

		if err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
				break
			}
			panic(err)
		}
		if count != 4 {
			panic("invalid input")
		}

		robots = append(robots, r)
	}

	return robots
}
