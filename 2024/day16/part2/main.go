package main

import (
	"bufio"
	"fmt"
	"image"
	"log"
	"math"
	"os"
	"sort"
	"tea-cats.co.uk/aoc/2024"
	"time"
)

const debug = false

// An Item is something we manage in a priority queue.
type Item struct {
	position     image.Point
	direction    image.Point
	costToArrive int
	priority     int // The priority of the item in the queue.
	path         []image.Point
}

// A PriorityQueue  and holds Items.
type PriorityQueue struct {
	length   int
	capacity int
	data     []*Item
}

func (pq *PriorityQueue) Len() int { return pq.length }

func (pq *PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the lowest, not highest, weight/priority so we use greater than here.
	itemI, itemJ := pq.data[i], pq.data[j]

	if itemI.priority == itemJ.priority {
		return itemI.costToArrive > itemJ.costToArrive
	}

	return itemI.priority > itemJ.priority
	//return pq.data[i].priority > pq.data[j].priority
}

func (pq *PriorityQueue) Swap(i, j int) {
	pq.data[i], pq.data[j] = pq.data[j], pq.data[i]
}

func (pq *PriorityQueue) pop() *Item {
	if pq.length == 0 {
		log.Printf("pop on empty queue\n")
		return nil
	}

	pq.length--

	item := pq.data[pq.length]
	pq.data[pq.length] = nil

	return item
}

func (pq *PriorityQueue) put(previous *Item, direction image.Point, cost int, dest image.Point) {

	position := previous.position.Add(direction)

	if debug {
		fmt.Printf(" * Queuing move of %v %s [current cost=%d", debugDirToCompass(direction), position, cost)
	}

	dist := dest.Sub(position)

	if debug {
		fmt.Printf(", distance=%v", dist)
	}

	// Total priority = min total cost
	//    = current cost + perfect run to target
	predictedCost := cost

	if dist.X != 0 {
		// Horizontal moves -- minimum of dist.X, plus 1000 if we're not facing that way
		predictedCost += utils.Abs(dist.X)
		if dist.X*direction.X <= 0 {
			if debug {
				fmt.Printf(", requires E/W turn")
			}
			predictedCost += 1000
		}
	}

	if dist.Y != 0 {
		// Horizontal moves -- minimum of dist.Y, plus 1000 if we're not facing that way
		predictedCost += utils.Abs(dist.Y)
		if dist.Y*direction.Y <= 0 {
			if debug {
				fmt.Printf(", requires N/S turn")
			}
			predictedCost += 1000
		}
	}

	if debug {
		fmt.Printf(", expectation=%d]\n", predictedCost)
	}

	if pq.length == pq.capacity {
		newData := make([]*Item, pq.length+128)
		copy(newData, pq.data)
		pq.data = newData
		pq.capacity = len(pq.data)
	}

	newPath := make([]image.Point, len(previous.path)+1)
	copy(newPath, previous.path)
	newPath[len(previous.path)] = position

	pq.data[pq.length] = &Item{position, direction, cost, predictedCost, newPath}
	pq.length++
}

func (pq *PriorityQueue) sort() {
	sort.Sort(pq)
}

type isWall bool
type dijkstraGrid struct {
	Data   []isWall
	Width  int
	Height int
}

var pointUp = image.Point{Y: -1}
var pointDown = image.Point{Y: 1}
var pointRight = image.Point{X: 1}
var pointLeft = image.Point{X: -1}

func (grid *dijkstraGrid) isWall(point image.Point) isWall {
	if point.X < 0 || point.Y < 0 || point.X >= grid.Width || point.Y >= grid.Height {
		return true
	}
	return grid.Data[point.Y*grid.Width+point.X]
}

type visit struct {
	position  image.Point
	direction image.Point
}

func (grid *dijkstraGrid) findRoute() int {
	defer utils.TimeTrack(time.Now(), "findRoute")

	start := image.Point{Y: grid.Height - 2}
	dest := image.Point{X: grid.Width - 2, Y: 1}
	preStart := &Item{
		position:     start,
		direction:    pointRight,
		costToArrive: 0,
		priority:     0,
		path:         []image.Point{},
	}

	if debug {
		fmt.Printf("Finding path from %v to %v\n", start, dest)
		fmt.Printf("==========================\n\n")
	}

	queue := PriorityQueue{length: 0, data: make([]*Item, 128)}
	queue.put(preStart, pointRight, 0, dest)

	minPathLength := math.MaxInt

	visited := map[visit]int{}
	seats := map[image.Point]bool{}

	var current *Item

	limit := math.MaxInt
	if debug {
		limit = 2000
	}

	for x := 0; x < limit; x++ {
		queue.sort()

		if debug {
			fmt.Printf("\n--------------------------------------\n\n")
			fmt.Println("Locating next node")
			for i := 0; i < queue.length; i++ {
				n := queue.data[i]
				fmt.Printf(" * %3d: %v %s [%d -> %d]\n", i, n.position, debugDirToCompass(n.direction), n.costToArrive, n.priority)
			}
		}

		for current = queue.pop(); current != nil; current = queue.pop() {
			previousCost, ok := visited[visit{current.position, current.direction}]

			if ok && previousCost < current.costToArrive {
				if debug {
					fmt.Printf(" - Skipping visited %v:%s\n", current.position, debugDirToCompass(current.direction))
				}
				continue
			}
			break
		}

		if current == nil {
			break
		}

		visited[visit{current.position, current.direction}] = current.costToArrive

		if debug {
			fmt.Printf("\nExpanding node %v (facing %s, cost %d)\n", current.position, debugDirToCompass(current.direction), current.costToArrive)
		}

		if current.priority > minPathLength {
			if debug {
				fmt.Printf("Discaring too long path %v\n", *current)
			}
			continue
		}

		if current.position == dest {
			minPathLength = current.costToArrive
			for _, n := range current.path {
				seats[n] = true
			}
			continue
		}

		nextPoint := current.position.Add(current.direction)
		if !grid.isWall(nextPoint) {
			queue.put(current, current.direction, current.costToArrive+1, dest)
		}

		if current.direction.X == 0 {
			// Current N/S, can turn to E/W
			nextPoint = current.position.Add(pointRight)
			if !grid.isWall(nextPoint) {
				queue.put(current, pointRight, current.costToArrive+1001, dest)
			}

			nextPoint = current.position.Add(pointLeft)
			if !grid.isWall(nextPoint) {
				queue.put(current, pointLeft, current.costToArrive+1001, dest)
			}
		} else {
			nextPoint = current.position.Add(pointUp)
			if !grid.isWall(nextPoint) {
				queue.put(current, pointUp, current.costToArrive+1001, dest)
			}

			nextPoint = current.position.Add(pointDown)
			if !grid.isWall(nextPoint) {
				queue.put(current, pointDown, current.costToArrive+1001, dest)
			}
		}
	}

	point := image.Point{X: 0, Y: 0}

	for point.Y = 0; point.Y < grid.Height; point.Y++ {
		for point.X = 0; point.X < grid.Width; point.X++ {
			if grid.isWall(point) {
				fmt.Printf("#")
			} else if seats[point] {
				fmt.Printf("O")
			} else {
				fmt.Printf(".")
			}
		}
		fmt.Printf("\n")
	}

	return len(seats)
}

func main() {
	defer utils.TimeTrack(time.Now(), "main")

	grid := loadData()
	cost := grid.findRoute()
	fmt.Printf("Cost: %d\n", cost)
	utils.PrintMemUsage()
}

func loadData() dijkstraGrid {
	defer utils.TimeTrack(time.Now(), "loadData")
	dataFile, err := os.Open("2024/input-16.txt")

	if err != nil {
		panic(err)
	}

	defer utils.CloseWithLog(dataFile)

	scanner := bufio.NewReader(dataFile)

	var width, lines int
	var data []isWall

	for y := 0; true; y++ {
		line, err := scanner.ReadSlice('\n')

		if len(line) == 0 || line[0] == '\n' {
			break
		}
		if err != nil {
			panic(err)
		}

		for _, c := range line {
			switch c {
			case '\n':
				break
			case '#':
				data = append(data, true)
			default:
				data = append(data, false)
			}
		}

		width = len(line) - 1
		lines++
	}

	return dijkstraGrid{
		Data:   data,
		Width:  width,
		Height: lines,
	}
}

func debugDirToCompass(d image.Point) string {
	if d.X > 0 {
		return "Right"
	}
	if d.X < 0 {
		return "Left"
	}
	if d.Y < 0 {
		return "Up"
	}
	return "Down"
}
