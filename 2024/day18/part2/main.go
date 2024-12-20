package main

import (
	"errors"
	"fmt"
	"image"
	"io"
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
		return nil
	}

	pq.length--

	item := pq.data[pq.length]
	pq.data[pq.length] = nil

	return item
}

func (pq *PriorityQueue) put(position image.Point, cost int, dest image.Point) {

	if debug {
		fmt.Printf(" * Queuing move to %v [current cost=%d", position, cost)
	}

	dist := dest.Sub(position)

	if debug {
		fmt.Printf(", distance=%v", dist)
	}

	// Total priority = min total cost
	//    = current cost + perfect run to target
	predictedCost := cost + utils.Abs(dist.X) + utils.Abs(dist.Y)

	if debug {
		fmt.Printf(", expectation=%d]\n", predictedCost)
	}

	if pq.length == pq.capacity {
		newData := make([]*Item, pq.length+128)
		copy(newData, pq.data)
		pq.data = newData
		pq.capacity = len(pq.data)
	}

	pq.data[pq.length] = &Item{position, cost, predictedCost}
	pq.length++
}

func (pq *PriorityQueue) sort() {
	sort.Sort(pq)
}

type dijkstraGrid struct {
	Data   map[image.Point]int
	Width  int
	Height int
}

var pointUp = image.Point{Y: -1}
var pointDown = image.Point{Y: 1}
var pointRight = image.Point{X: 1}
var pointLeft = image.Point{X: -1}

func (grid *dijkstraGrid) isWall(point image.Point, atTime int) bool {
	if point.X < 0 || point.Y < 0 || point.X >= grid.Width || point.Y >= grid.Height {
		return true
	}
	fallTime, fall := grid.Data[point]

	return fall && fallTime < atTime
}

type history map[image.Point]struct{}

func (h history) has(i image.Point) bool {
	_, ok := h[i]

	return ok
}

func (grid *dijkstraGrid) findRoute(steps int) int {
	start := image.Point{}
	dest := image.Point{X: grid.Width - 1, Y: grid.Height - 1}

	if debug {
		fmt.Printf("Finding path from %v to %v\n", start, dest)
		fmt.Printf("==========================\n\n")
	}

	queue := PriorityQueue{length: 0, data: make([]*Item, 128)}
	queue.put(start, 0, dest)
	visited := history{}

	var current *Item

	limit := math.MaxInt
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
				fmt.Printf(" * %3d: %v [%d -> %d]\n", i, n.position, n.costToArrive, n.priority)
			}
		}

		for current = queue.pop(); current != nil && visited.has(current.position); current = queue.pop() {
			if debug {
				fmt.Printf(" - Skipping visited %v\n", current.position)
			}
		}
		if current == nil {
			break
		}

		visited[current.position] = struct{}{}

		if debug {
			fmt.Printf("\nExpanding node %v (cost %d)\n", current.position, current.costToArrive)
		}

		if current.position == dest {
			if debug {
				point := image.Point{X: 0, Y: 0}

				for point.Y = 0; point.Y < grid.Height; point.Y++ {
					for point.X = 0; point.X < grid.Width; point.X++ {
						if grid.isWall(point, steps) {
							fmt.Printf("#")
						} else if visited.has(point) {
							fmt.Printf("O")
						} else {
							fmt.Printf(".")
						}
					}
					fmt.Printf("\n")
				}
			}

			return current.costToArrive
		}

		var nextPoint image.Point

		nextPoint = current.position.Add(pointDown)
		if !grid.isWall(nextPoint, steps) && !visited.has(nextPoint) {
			queue.put(nextPoint, current.costToArrive+1, dest)
		}

		nextPoint = current.position.Add(pointRight)
		if !grid.isWall(nextPoint, steps) && !visited.has(nextPoint) {
			queue.put(nextPoint, current.costToArrive+1, dest)
		}

		nextPoint = current.position.Add(pointUp)
		if !grid.isWall(nextPoint, steps) && !visited.has(nextPoint) {
			queue.put(nextPoint, current.costToArrive+1, dest)
		}

		nextPoint = current.position.Add(pointLeft)
		if !grid.isWall(nextPoint, steps) && !visited.has(nextPoint) {
			queue.put(nextPoint, current.costToArrive+1, dest)
		}
	}

	return 0
}

func main() {
	defer utils.TimeTrack(time.Now(), "main")

	grid := loadData(71)

	maxSteps := len(grid.Data)
	minSteps := 1024

	for currentSteps := (maxSteps + minSteps) / 2; minSteps != maxSteps; currentSteps = (maxSteps+minSteps)/2 + 1 {
		cost := grid.findRoute(currentSteps)
		if debug {
			fmt.Printf("drops: %v, min=%d,max=%d, cost=%d\n", currentSteps, minSteps, maxSteps, cost)
		}

		if cost == 0 {
			maxSteps = currentSteps - 1
		} else {
			minSteps = currentSteps
		}
	}

	var point image.Point
	var steps int
	for point, steps = range grid.Data {
		if steps == maxSteps {
			break
		}
	}

	fmt.Printf("Max Safe Steps: %d, Next point: %v\n", maxSteps, point)
	utils.PrintMemUsage()
}

func loadData(size int) dijkstraGrid {
	defer utils.TimeTrack(time.Now(), "loadData")
	dataFile, err := os.Open("2024/input-18.txt")

	if err != nil {
		panic(err)
	}

	defer utils.CloseWithLog(dataFile)

	data := make(map[image.Point]int)

	var p image.Point
	i := 0

	for {
		c, err := fmt.Fscanf(dataFile, "%d,%d\n", &p.X, &p.Y)

		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			if errors.Is(err, io.ErrUnexpectedEOF) {
				break
			}
			panic(err)
		}
		if c != 2 {
			panic("wrong number of entries")
		}

		data[p] = i
		i++
	}

	return dijkstraGrid{
		Data:   data,
		Width:  size,
		Height: size,
	}
}
