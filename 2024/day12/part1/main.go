package main

import (
	"fmt"
	"image"
	"tea-cats.co.uk/aoc/2024"
	"tea-cats.co.uk/aoc/2024/day12"
	"time"
)

const debug = day12.Debug

func main() {
	defer utils.TimeTrack(time.Now(), "main")

	data := day12.LoadData()
	regions := day12.LocateRegions(&data)

	if debug {
		point := image.Point{}
		for point.Y = 0; point.Y < data.Height; point.Y++ {
			for point.X = 0; point.X < data.Width; point.X++ {
				fmt.Printf("%2d ", data.AtPoint(point).Region.RegionId)
			}
			fmt.Println()
		}
	}

	cost := calculateCost(regions, data)

	fmt.Println(cost)
}

func calculateCost(regions []day12.Region, data utils.Grid[day12.GridPoint]) int {
	defer utils.TimeTrack(time.Now(), "calculateCost")
	cost := 0

	for _, plot := range regions {
		perim := 0

		for point := range plot.Points {
			self := data.AtPoint(point)

			above := data.AtPoint(point.Add(image.Point{Y: -1}))
			if above == nil || above.Region != self.Region {
				perim++
			}

			left := data.AtPoint(point.Add(image.Point{X: -1}))
			if left == nil || left.Region != self.Region {
				perim++
			}

			below := data.AtPoint(point.Add(image.Point{Y: 1}))
			if below == nil || below.Region != self.Region {
				perim++
			}

			right := data.AtPoint(point.Add(image.Point{X: 1}))
			if right == nil || right.Region != self.Region {
				perim++
			}
		}

		if debug {
			fmt.Println(plot)
			fmt.Printf("Region %d (%s): size=%d, perim=%d\n", plot.RegionId, string(plot.Symbol), len(plot.Points), perim)
		}

		cost += len(plot.Points) * perim
	}
	return cost
}
