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

var pointUp = image.Point{Y: -1}
var pointDown = image.Point{Y: 1}
var pointRight = image.Point{X: 1}
var pointLeft = image.Point{X: -1}

type cell byte

const (
	cellEmpty    cell = '.'
	cellBoxLeft  cell = '['
	cellBoxRight cell = ']'
	cellWall     cell = '#'
	cellRobot    cell = '@'
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

	if dir.Y == 0 {
		if g.tryShiftHorizontal(g.robot, dir) {
			g.robot = g.robot.Add(dir)
			if debug {
				fmt.Printf("  - Robot now at %v\n", g.robot)
			}
		}
		return
	}

	if g.tryShiftVertical(g.robot, dir) {
		g.robot = g.robot.Add(dir)
		if debug {
			fmt.Printf("  - Robot now at %v\n", g.robot)
		}
	}
}

func (g *grid) tryShiftHorizontal(point image.Point, dir image.Point) bool {
	next := point.Add(dir)
	self := g.at(point)
	target := g.at(next)

	if debug {
		fmt.Printf("  - Trying move of %s from %v to %v (%s): ", string(self), point, next, string(target))
	}

	if target == cellWall {
		if debug {
			fmt.Printf("target is wall, not moving\n")
		}
		return false
	}

	if target == cellBoxLeft || target == cellBoxRight {
		if debug {
			fmt.Printf("is box, checking next:\n")
		}
		if !g.tryShiftHorizontal(next, dir) {
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

func (g *grid) tryShiftVertical(point image.Point, dir image.Point) bool {
	next := point.Add(dir)
	self := g.at(point)
	target := g.at(next)

	if debug {
		fmt.Printf("  - Trying move of %s from %v to %v (%s): ", string(self), point, next, string(target))
	}

	if target == cellWall {
		if debug {
			fmt.Printf("target is wall, not moving\n")
		}
		return false
	}

	if target == cellBoxLeft {
		if debug {
			fmt.Printf("is left of box, checking left cna move:\n")
		}
		if !g.checkShiftVertical(next, dir) {
			if debug {
				fmt.Printf("Not moving box %v as left side can't move\n", next)
			}
			return false
		}
		if debug {
			fmt.Printf("Now checking if the right of %v can move\n", point)
		}
		if !g.tryShiftVertical(next.Add(pointRight), dir) {
			if debug {
				fmt.Printf("Not moving box %v as right side can't move\n", next)
			}
			return false
		}
		if debug {
			fmt.Printf("Successfully checked left and moved right, moving left of %v\n", next)
		}
		g.tryShiftVertical(next, dir)
	}

	if target == cellBoxRight {
		if debug {
			fmt.Printf("is right of box, checking is right can move:\n")
		}
		if !g.checkShiftVertical(next, dir) {
			if debug {
				fmt.Printf("Not moving box %v as RIGHT side can't move\n", next)
			}
			return false
		}
		if debug {
			fmt.Printf("Now checking if the left of %v can move\n", next)
		}
		if !g.tryShiftVertical(next.Add(pointLeft), dir) {
			if debug {
				fmt.Printf("Not moving box %v as left side can't move\n", next)
			}
			return false
		}
		if debug {
			fmt.Printf("Successfully checked right and moved left, moving right of %v\n", next)
		}
		g.tryShiftVertical(next, dir)
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

func (g *grid) checkShiftVertical(point image.Point, dir image.Point) bool {
	next := point.Add(dir)
	self := g.at(point)
	target := g.at(next)

	if debug {
		fmt.Printf("    - Checking move of %s from %v to %v (%s): ", string(self), point, next, string(target))
	}

	switch target {
	case cellWall:
		if debug {
			fmt.Printf("target is wall, not approving move\n")
		}
		return false

	case cellEmpty:
		if debug {
			fmt.Printf("target is empty, approving move\n")
		}
		return true

	case cellBoxLeft:
		if debug {
			fmt.Printf("is left of box, checking both sides can move:\n")
		}
		return g.checkShiftVertical(next, dir) && g.checkShiftVertical(next.Add(pointRight), dir)

	case cellBoxRight:
		if debug {
			fmt.Printf("is right of box, checking both sides can move:\n")
		}
		return g.checkShiftVertical(next, dir) && g.checkShiftVertical(next.Add(pointLeft), dir)

	default:
		fmt.Printf("UKNOWN SYMBOL %s\n", string(target))
		return false
	}
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
		up:    pointUp,
		right: pointRight,
		down:  pointDown,
		left:  pointLeft,
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
		if c == cellBoxLeft {
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
			switch c {
			case '\n':
				break
			case '@':
				robot = image.Point{X: x * 2, Y: lines}
				data = append(data, cellRobot)
				data = append(data, cellEmpty)
			case '#':
				data = append(data, cellWall)
				data = append(data, cellWall)
			case 'O':
				data = append(data, cellBoxLeft)
				data = append(data, cellBoxRight)
			case '.':
				data = append(data, cellEmpty)
				data = append(data, cellEmpty)
			}
		}

		width = 2 * (len(line) - 1)
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
