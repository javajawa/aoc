package main

import (
	"bufio"
	"errors"
	"fmt"
	"image"
	"io"
	"os"
	"tea-cats.co.uk/aoc/2024"
	"time"
)

const debug = false

type instruction byte

const (
	up    instruction = '^'
	right instruction = '>'
	left  instruction = '<'
	down  instruction = 'v'
)

type cell byte

const (
	cellEmpty cell = '.'
	cellBox   cell = 'O'
	cellWall  cell = '#'
	cellRobot cell = '@'
)

type grid struct {
	utils.Grid[cell]
	robot image.Point
}

func (g *grid) at(point image.Point) cell {
	if point.X < 0 || point.Y < 0 || point.X >= g.Width || point.Y >= g.Height {
		return cellWall
	}
	return g.Data[point.Y*g.Width+point.X]
}

func (g *grid) set(point image.Point, c cell) {
	if point.X < 0 || point.Y < 0 || point.X >= g.Width || point.Y >= g.Height {
		panic(fmt.Sprintf("Cannot go to %v", point))
	}
	g.Data[point.Y*g.Width+point.X] = c
}

func (g *grid) shift(dir image.Point) {
	if debug {
		fmt.Printf("Trying to move from %v by %v\n", g.robot, dir)
	}
	if g.tryShift(g.robot, dir) {
		g.robot = g.robot.Add(dir)
		if debug {
			fmt.Printf("  - Robot now at %v\n", g.robot)
		}
	}
}

func (g *grid) tryShift(point image.Point, dir image.Point) bool {
	next := point.Add(dir)
	self := g.at(point)
	target := g.at(next)

	if debug {
		fmt.Printf("  - Checking move of %s from %v to %v (%s): ", string(self), point, next, string(target))
	}

	if target == cellWall {
		if debug {
			fmt.Printf("target is wall, not moving\n")
		}
		return false
	}

	if target == cellBox {
		if debug {
			fmt.Printf("is box, checking inside:\n")
		}
		if !g.tryShift(next, dir) {
			return false
		}
	}

	// Fall through case: next cell was cell_empty _or_
	// We shifted boxes to make it cell_empty.
	if debug {
		fmt.Printf("performing swap\n")
	}

	g.set(next, self)
	g.set(point, cellEmpty)

	return true
}

func main() {
	defer utils.TimeTrack(time.Now(), "main")

	grid, instructions := loadData()

	processInstructions(&grid, instructions)
	total := sumValue(grid)

	//printGrid(grid)
	fmt.Printf("Result: %d\n", total)
}

func processInstructions(g *grid, instructions []instruction) {
	defer utils.TimeTrack(time.Now(), "processInstructions")

	directions := map[instruction]image.Point{
		up:    {0, -1},
		right: {1, 0},
		down:  {0, 1},
		left:  {-1, 0},
	}

	for _, dir := range instructions {
		if debug {
			printGrid(*g)
		}
		g.shift(directions[dir])
	}
}

func sumValue(g grid) int {
	defer utils.TimeTrack(time.Now(), "sumValue")

	row := -1
	acc := 0
	for i, c := range g.Data {
		x := i % g.Width
		if x == 0 {
			row++
		}
		if c == cellBox {
			acc += row*100 + x
		}
	}
	return acc
}

func loadData() (grid, []instruction) {
	defer utils.TimeTrack(time.Now(), "loadData")
	dataFile, err := os.Open("2024/input-15.txt")

	if err != nil {
		panic(err)
	}

	defer utils.CloseWithLog(dataFile)

	scanner := bufio.NewReader(dataFile)

	var width, lines int
	var data []cell
	var robot image.Point

	for y := 0; true; y++ {
		line, err := scanner.ReadSlice('\n')

		if line[0] == '\n' {
			break
		}
		if err != nil {
			panic(err)
		}

		for x, c := range line {
			if c == '\n' {
				continue
			}
			cellType := cell(c)
			if cellType == cellRobot {
				robot = image.Point{X: x, Y: lines}
			}

			data = append(data, cellType)
		}

		width = len(line) - 1
		lines++
	}

	fmt.Printf("Found %d points, width=%d, height=%d\n", len(data), width, lines)

	instructions := make([]instruction, 0)

	for {
		line, err := scanner.ReadSlice('\n')

		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			panic(err)
		}

		for _, c := range line {
			if c == '\n' {
				continue
			}
			instructions = append(instructions, instruction(c))
		}
	}

	fmt.Printf("Read %d instructions\n", len(instructions))

	return grid{
		Grid: utils.Grid[cell]{
			Data:   data,
			Width:  width,
			Height: lines,
		},
		robot: robot,
	}, instructions
}

func printGrid(g grid) {
	for row := 0; row < g.Height; row++ {
		start := row * g.Width
		end := start + g.Width
		fmt.Println(string(g.Data[start:end]))
	}
}
