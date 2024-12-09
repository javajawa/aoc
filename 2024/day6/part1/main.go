package main

import (
	"fmt"
	"os"
	"tea-cats.co.uk/aoc/2024"
	"time"
)

const debug = false
const size = 130

type Direction uint8

const (
	North Direction = iota
	East  Direction = iota
	South Direction = iota
	West  Direction = iota
)

type CellState uint8

const (
	Clear       CellState = iota
	Obstruction CellState = iota
	Visited     CellState = iota
)

type Point struct {
	x, y uint8
}

type Maze struct {
	area      [size][size]CellState
	guard     Point
	direction Direction
	visited   uint16
}

func main() {
	defer utils.TimeTrack(time.Now(), "main")

	maze := loadData()

	for maze.move() {
	}

	if debug {
		maze.Print()
	}
	fmt.Printf("Visited: %d\n", maze.visited)
}

func loadData() Maze {
	defer utils.TimeTrack(time.Now(), "loadData")
	dataFile, err := os.Open("2024/input-6.txt")

	if err != nil {
		panic(err)
	}

	defer utils.CloseWithLog(dataFile)

	result := Maze{
		guard:     Point{0, 0},
		direction: North,
	}

	raw := make([]byte, size*size+size)

	_, err = dataFile.Read(raw)
	if err != nil {
		panic(err)
	}

	row := uint8(0)
	col := int16(0)

	for _, char := range raw {
		switch char {
		case '#':
			result.area[row][col] = Obstruction
		case '^':
			result.area[row][col] = Visited
			result.guard.y = row
			result.guard.x = uint8(col)
			fmt.Println("Starting from ", row, ":", uint8(col))
			result.visited = 1
		case '\n':
			row += 1
			col = -1
		}
		col += 1
	}

	return result
}

func (maze *Maze) move() bool {
	var next Point

	switch maze.direction {
	case North:
		next = Point{maze.guard.x, maze.guard.y - 1}
	case East:
		next = Point{maze.guard.x + 1, maze.guard.y}
	case South:
		next = Point{maze.guard.x, maze.guard.y + 1}
	case West:
		next = Point{maze.guard.x - 1, maze.guard.y}
	}

	if next.x >= size || next.y >= size {
		if debug {
			fmt.Printf("Escaping at %v\n", next)
		}
		return false
	}

	if maze.area[next.y][next.x] == Obstruction {
		if maze.direction == West {
			maze.direction = North
		} else {
			maze.direction += 1
		}
		if debug {
			fmt.Printf("Encountered obstruction at %v, turning to %v\n", next, maze.direction)
		}
		return true
	}

	//fmt.Printf("Moving %v to %v\n", maze.direction, next)
	maze.guard = next
	if maze.area[next.y][next.x] != Visited {
		maze.area[next.y][next.x] = Visited
		maze.visited++
	}
	return true
}

func (maze *Maze) Print() {
	buffer := make([]byte, size*size+size)

	for row := 0; row < size; row++ {
		for col := 0; col < size; col++ {
			switch maze.area[row][col] {
			case Obstruction:
				buffer[row*(size+1)+col] = '#'
			case Visited:
				buffer[row*(size+1)+col] = 'X'
			case Clear:
				buffer[row*(size+1)+col] = ' '
			}
		}
		buffer[row*(size+1)+size] = '\n'
	}

	fmt.Println(string(buffer))
}
