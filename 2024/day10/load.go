package day10

import (
	"bufio"
	"fmt"
	"image"
	"io"
	"os"
	utils "tea-cats.co.uk/aoc/2024"
	"time"
)

type PointHeight int16

type Grid struct {
	data   []PointHeight
	width  int
	height int
}

func (grid *Grid) At(point image.Point) PointHeight {
	if point.X < 0 || point.Y < 0 || point.X >= grid.width || point.Y >= grid.height {
		return -1
	}
	return grid.data[point.Y*grid.width+point.X]
}

func LoadData() (Grid, [10][]image.Point) {
	defer utils.TimeTrack(time.Now(), "loadData")
	dataFile, err := os.Open("2024/input-10.txt")

	if err != nil {
		panic(err)
	}

	defer utils.CloseWithLog(dataFile)

	scanner := bufio.NewReader(dataFile)

	var width, lines int
	var data []PointHeight

	points := [10][]image.Point{}

	for y := 0; true; y++ {
		line, err := scanner.ReadSlice('\n')

		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		width = len(line) - 1
		lines++

		for x, c := range line {
			if c == '\n' {
				continue
			}

			data = append(data, PointHeight(c-'0'))

			points[c-'0'] = append(points[c-'0'], image.Point{X: x, Y: y})
		}
	}

	fmt.Printf("Found %d points, width=%d, height=%d\n", len(data), width, lines)

	return Grid{
		data:   data,
		width:  width,
		height: lines,
	}, points
}
