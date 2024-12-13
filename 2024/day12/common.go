package day12

import (
	"bufio"
	"fmt"
	"image"
	"io"
	"os"
	utils "tea-cats.co.uk/aoc/2024"
	"time"
)

const Debug = false

type Region struct {
	Symbol   byte
	RegionId int
	Points   utils.Set[image.Point]
}

type GridPoint struct {
	symbol byte
	Region *Region
}

func LocateRegions(data *utils.Grid[GridPoint]) []Region {
	defer utils.TimeTrack(time.Now(), "LocateRegions")
	point := image.Point{}

	regions := make([]Region, 0)
	neighbours := [...]image.Point{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}

	for point.Y = 0; point.Y < data.Height; point.Y++ {
		for point.X = 0; point.X < data.Width; point.X++ {
			target := data.AtPoint(point)

			for _, neighbourVector := range neighbours {
				neighbor := data.AtPoint(point.Add(neighbourVector))

				if neighbor == nil {
					continue
				}
				if neighbor.symbol != target.symbol {
					continue
				}
				if neighbor.Region == nil {
					continue
				}

				if target.Region == nil {
					if Debug {
						fmt.Printf("Adding %v to Region %d(%s) due to neighbour %v\n", point, neighbor.Region.RegionId, string(target.symbol), point.Add(neighbourVector))
					}

					neighbor.Region.Points.Add(point)
					target.Region = neighbor.Region
					continue
				}

				if target.Region.RegionId == neighbor.Region.RegionId {
					continue
				}

				donorRegion := neighbor.Region
				targetRegion := target.Region

				if Debug {
					fmt.Printf("Merging %v into %v\n", donorRegion, targetRegion)
				}

				for p := range donorRegion.Points {
					targetRegion.Points.Add(p)
					data.AtPoint(p).Region = targetRegion
				}
				donorRegion.Points.Clear()

				if Debug {
					fmt.Printf(" -> %v ; %v\n", donorRegion, targetRegion)
				}
			}

			if target.Region == nil {
				if Debug {
					fmt.Printf("Creating new %s Region %d for %v\n", string(target.symbol), len(regions), point)
				}

				r := Region{Symbol: target.symbol, RegionId: len(regions), Points: utils.Set[image.Point]{point: struct{}{}}}
				regions = append(regions, r)
				target.Region = &r
			}
		}
	}
	return regions
}

func LoadData() utils.Grid[GridPoint] {
	defer utils.TimeTrack(time.Now(), "LoadData")
	dataFile, err := os.Open("2024/input-12.txt")

	if err != nil {
		panic(err)
	}

	defer utils.CloseWithLog(dataFile)

	reader := bufio.NewReader(dataFile)
	lines := make([]GridPoint, 0)
	width := 0

	for {
		line, err := reader.ReadSlice('\n')

		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		width = len(line) - 1
		for _, c := range line {
			if c == '\n' {
				break
			}
			lines = append(lines, GridPoint{c, nil})
		}
	}

	return utils.Grid[GridPoint]{
		Data:   lines,
		Width:  width,
		Height: len(lines) / width,
	}
}
