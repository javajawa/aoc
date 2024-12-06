package main

import (
	"fmt"
	"os"
	"slices"
	"tea-cats.co.uk/aoc/2024"
	"time"
)

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

type PointFromDirection struct {
	x, y uint8
	dir  Direction
}

type Maze struct {
	area         [size][size]CellState
	guard        Point
	direction    Direction
	obstructions []PointFromDirection
}

func main() {
	defer utils.TimeTrack(time.Now(), "main")

	maze := loadData()

	mazeWithoutExtraObstruction := Maze{
		area:      maze.area,
		guard:     maze.guard,
		direction: North,
	}

	// Find all cells the guard will naturally visit
	if mazeWithoutExtraObstruction.checkLoop() {
		panic("Already a loop?")
	}

	var i, j uint8
	possibleLoop := 0

	for i = 0; i < size; i++ {
		for j = 0; j < size; j++ {
			if mazeWithoutExtraObstruction.area[i][j] != Visited {
				continue
			}

			testMaze := Maze{
				area:      maze.area,
				guard:     maze.guard,
				direction: North,
			}
			testMaze.area[i][j] = Obstruction
			if testMaze.checkLoop() {
				possibleLoop++
				//fmt.Printf("Maze contains loop with extra obstruction %v\n", Point{x: j, y: i})
			}
		}
	}

	fmt.Printf("Loops found: %d\n", possibleLoop)
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
		case '\n':
			row += 1
			col = -1
		}
		col += 1
	}

	return result
}

func (maze *Maze) checkLoop() bool {
	var next Point

	for {
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
			//fmt.Printf("Escaping at %v\n", next)
			return false
		}

		if maze.area[next.y][next.x] != Obstruction {
			//fmt.Printf("Moving %v to %v\n", maze.direction, next)
			maze.guard = next
			maze.area[next.y][next.x] = Visited
			continue
		}

		marker := PointFromDirection{
			x:   next.x,
			y:   next.y,
			dir: maze.direction,
		}

		if slices.Contains(maze.obstructions, marker) {
			//fmt.Println("Loop found")
			return true
		}

		maze.obstructions = append(maze.obstructions, marker)

		if maze.direction == West {
			maze.direction = North
		} else {
			maze.direction += 1
		}
		//fmt.Printf("Encountered obstruction at %v, turning to %v\n", next, maze.direction)
	}
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
