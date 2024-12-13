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
		sides := 0

		if debug {
			fmt.Printf("Region %d '%s'\n", plot.RegionId, string(plot.Symbol))
		}

		for point := range plot.Points {
			self := data.AtPoint(point)

			above := data.AtPoint(point.Add(image.Point{Y: -1}))
			edgeAbove := above == nil || above.Region != self.Region

			left := data.AtPoint(point.Add(image.Point{X: -1}))
			edgeLeft := left == nil || left.Region != self.Region

			below := data.AtPoint(point.Add(image.Point{Y: 1}))
			edgeBelow := below == nil || below.Region != self.Region

			right := data.AtPoint(point.Add(image.Point{X: 1}))
			edgeRight := right == nil || right.Region != self.Region

			if debug {
				fmt.Printf(" -> checking %v for corners: ", point)
			}

			if edgeAbove && edgeRight {
				if debug {
					fmt.Printf(" has convex ┐")
				}
				sides++
			}
			if !edgeBelow && !edgeLeft {
				diag := data.AtPoint(point.Add(image.Point{X: -1, Y: 1}))
				if diag.Region != self.Region {
					if debug {
						fmt.Printf(" has concave ┐")
					}
					sides++
				}
			}

			if edgeRight && edgeBelow {
				if debug {
					fmt.Printf(" has convex ┘")
				}
				sides++
			}
			if !edgeLeft && !edgeAbove {
				diag := data.AtPoint(point.Add(image.Point{X: -1, Y: -1}))
				if diag.Region != self.Region {
					if debug {
						fmt.Printf(" has concave ┘")
					}
					sides++
				}
			}

			if edgeBelow && edgeLeft {
				if debug {
					fmt.Printf(" has convex └")
				}
				sides++
			}
			if !edgeAbove && !edgeRight {
				diag := data.AtPoint(point.Add(image.Point{X: 1, Y: -1}))
				if diag.Region != self.Region {
					if debug {
						fmt.Printf(" has concave └")
					}
					sides++
				}
			}

			if edgeLeft && edgeAbove {
				if debug {
					fmt.Printf(" has convex ┌")
				}
				sides++
			}
			if !edgeRight && !edgeBelow {
				diag := data.AtPoint(point.Add(image.Point{X: 1, Y: 1}))
				if diag.Region != self.Region {
					if debug {
						fmt.Printf(" has concave ┌")
					}
					sides++
				}
			}

			if debug {
				fmt.Printf("\n")
			}
		}

		if debug {
			fmt.Println(plot)
			fmt.Printf("= size=%d, sides=%d, cost=%d\n", len(plot.Points), sides, len(plot.Points)*sides)
		}

		cost += len(plot.Points) * sides
	}
	return cost
}
