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
	seconds := 100
	center := grid.Max.Sub(grid.Min).Sub(image.Point{1, 1}).Div(2)

	fmt.Printf("Grid=%v, Center=%v\n", center)

	counts := make(map[image.Point]int)
	quadrants := [4]int{0, 0, 0, 0}

	for i, robot := range robots {
		final := robot.finalPosition(grid, seconds)
		if debug {
			fmt.Printf("Robot %d ends at %v\n", i, final)
			counts[final]++
		}

		if final.X == center.X || final.Y == center.Y {
			continue
		}

		quadrant := 0

		if final.X < center.X {
			quadrant = 1
		}
		if final.Y < center.Y {
			quadrant += 2
		}
		quadrants[quadrant]++
	}

	if debug {
		point := image.Point{}
		for point.Y = grid.Min.Y; point.Y < grid.Max.Y; point.Y++ {
			for point.X = grid.Min.X; point.X < grid.Max.X; point.X++ {
				count, ok := counts[point]
				if point.X == center.X {
					fmt.Printf(" ")
				} else if point.Y == center.Y {
					fmt.Printf(" ")
				} else if !ok || count == 0 {
					fmt.Printf(".")
				} else if count > 9 {
					fmt.Printf("+")
				} else {
					fmt.Printf("%d", count)
				}
			}
			fmt.Printf("\n")
		}
	}

	fmt.Printf("Quadrants=%v, Safety Factor: %d\n", quadrants, quadrants[0]*quadrants[1]*quadrants[2]*quadrants[3])
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
