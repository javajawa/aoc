package main

import (
	"fmt"
	"image"
	"tea-cats.co.uk/aoc/2024"
	"tea-cats.co.uk/aoc/2024/day10"
	"time"
)

const debug = false

func main() {
	defer utils.TimeTrack(time.Now(), "main")

	grid, heightMap := day10.LoadData()

	defer utils.TimeTrack(time.Now(), "process")

	totalTrails := 0
	neighbours := [4]image.Point{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}

	for _, start := range heightMap[0] {
		locationsAtThisHeight := map[image.Point]struct{}{start: {}}

		for height := day10.PointHeight(1); height < 10; height++ {
			nextLocations := make(map[image.Point]struct{})

			for previous := range locationsAtThisHeight {
				if debug {
					fmt.Printf("Finding neighbours of %v with height %d\n", previous, height)
				}
				for _, neighbour := range neighbours {
					location := previous.Add(neighbour)
					if grid.At(location) == height {
						if debug {
							fmt.Printf("  %v (height %d) is adjacent to %v (height %d)\n", location, height, previous, height-1)
						}
						nextLocations[location] = struct{}{}
					}
				}
			}

			locationsAtThisHeight = nextLocations
		}

		if debug {
			fmt.Printf("Trails from %v: %d\n", start, len(locationsAtThisHeight))
		}
		totalTrails += len(locationsAtThisHeight)
	}

	fmt.Printf("Total Trails: %d\n", totalTrails)
}
