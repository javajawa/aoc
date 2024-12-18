package main

import (
	"bufio"
	"fmt"
	"image"
	"log"
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

func (pq *PriorityQueue) put(position image.Point, direction image.Point, cost int, dest image.Point) {

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

	pq.data[pq.length] = &Item{position, direction, cost, predictedCost}
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
type history map[visit]struct{}

func (h history) has(i *Item) bool {
	_, ok := h[visit{i.position, i.direction}]

	return ok
}

func (grid *dijkstraGrid) findRoute() int {
	defer utils.TimeTrack(time.Now(), "findRoute")

	start := image.Point{X: 1, Y: grid.Height - 2}
	dest := image.Point{X: grid.Width - 2, Y: 1}

	if debug {
		fmt.Printf("Finding path from %v to %v\n", start, dest)
		fmt.Printf("==========================\n\n")
	}

	queue := PriorityQueue{length: 0, data: make([]*Item, 128)}
	queue.put(start, image.Point{X: 1}, 0, dest)
	visited := history{}

	var current *Item

	limit := len(grid.Data)
	if debug {
		limit = 20
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

		for current = queue.pop(); visited.has(current); current = queue.pop() {
			if debug {
				fmt.Printf(" - Skipping visited %v:%s\n", current.position, debugDirToCompass(current.direction))
			}
		}
		visited[visit{current.position, current.direction}] = struct{}{}

		if debug {
			fmt.Printf("\nExpanding node %v (facing %s, cost %d)\n", current.position, debugDirToCompass(current.direction), current.costToArrive)
		}

		if current.position == dest {
			return current.costToArrive
		}

		nextPoint := current.position.Add(current.direction)
		if !grid.isWall(nextPoint) {
			queue.put(nextPoint, current.direction, current.costToArrive+1, dest)
		}

		if current.direction.X == 0 {
			// Current N/S, can turn to E/W
			nextPoint = current.position.Add(pointRight)
			if !grid.isWall(nextPoint) {
				queue.put(nextPoint, pointRight, current.costToArrive+1001, dest)
			}

			nextPoint = current.position.Add(pointLeft)
			if !grid.isWall(nextPoint) {
				queue.put(nextPoint, pointLeft, current.costToArrive+1001, dest)
			}
		} else {
			nextPoint = current.position.Add(pointUp)
			if !grid.isWall(nextPoint) {
				queue.put(nextPoint, pointUp, current.costToArrive+1001, dest)
			}

			nextPoint = current.position.Add(pointDown)
			if !grid.isWall(nextPoint) {
				queue.put(nextPoint, pointDown, current.costToArrive+1001, dest)
			}
		}
	}

	return 0
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
