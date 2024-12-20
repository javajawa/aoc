package main

import (
	"bufio"
	"fmt"
	"image"
	"log"
	"math"
	"os"
	"slices"
	"sort"
	"tea-cats.co.uk/aoc/2024"
	"time"
)

const debug = false
const quickDebug = debug || true

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
	Data   []bool
	Width  int
	Height int
	Start  image.Point
	End    image.Point
}

func (grid *dijkstraGrid) isWall(point image.Point) bool {
	if point.X < 0 || point.Y < 0 || point.X >= grid.Width || point.Y >= grid.Height {
		return true
	}
	return grid.Data[point.Y*grid.Width+point.X]
}

var pointUp = image.Point{Y: -1}
var pointDown = image.Point{Y: 1}
var pointRight = image.Point{X: 1}
var pointLeft = image.Point{X: -1}

type history map[image.Point]int

func (h history) has(i image.Point) bool {
	_, ok := h[i]

	return ok
}

func (grid *dijkstraGrid) findRoute() (int, history) {
	defer utils.TimeTrack(time.Now(), "findRoute")

	if debug {
		fmt.Printf("Finding path from %v to %v\n", grid.Start, grid.End)
		fmt.Printf("==========================\n\n")
	}

	queue := PriorityQueue{length: 0, data: make([]*Item, 128)}
	queue.put(grid.Start, 0, grid.End)
	visited := history{}
	bestCost := math.MaxInt

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

		visited[current.position] = current.costToArrive

		if debug {
			fmt.Printf("\nExpanding node %v (cost %d)\n", current.position, current.costToArrive)
		}

		if current.position == grid.End {
			if current.costToArrive < bestCost {
				bestCost = current.costToArrive
			}
			continue
		}

		var nextPoint image.Point

		nextPoint = current.position.Add(pointDown)
		if !grid.isWall(nextPoint) && !visited.has(nextPoint) {
			queue.put(nextPoint, current.costToArrive+1, grid.End)
		}

		nextPoint = current.position.Add(pointRight)
		if !grid.isWall(nextPoint) && !visited.has(nextPoint) {
			queue.put(nextPoint, current.costToArrive+1, grid.End)
		}

		nextPoint = current.position.Add(pointUp)
		if !grid.isWall(nextPoint) && !visited.has(nextPoint) {
			queue.put(nextPoint, current.costToArrive+1, grid.End)
		}

		nextPoint = current.position.Add(pointLeft)
		if !grid.isWall(nextPoint) && !visited.has(nextPoint) {
			queue.put(nextPoint, current.costToArrive+1, grid.End)
		}
	}

	return bestCost, visited
}

type cheatOptions struct {
	wallStep image.Point
	nextStep [3]image.Point
}

func newCheatOption(wallStep image.Point, turnOptionOne image.Point, turnOptionTwo image.Point) cheatOptions {
	return cheatOptions{
		wallStep: wallStep,
		nextStep: [3]image.Point{wallStep.Add(wallStep), wallStep.Add(turnOptionOne), wallStep.Add(turnOptionTwo)},
	}
}

//goland:noinspection GoBoolExpressions
func main() {
	defer utils.TimeTrack(time.Now(), "main")

	cheats := [4]cheatOptions{
		newCheatOption(pointDown, pointRight, pointLeft),
		newCheatOption(pointUp, pointRight, pointLeft),
		newCheatOption(pointRight, pointUp, pointDown),
		newCheatOption(pointLeft, pointUp, pointDown),
	}

	grid := loadData()
	cost, visited := grid.findRoute()

	fmt.Printf("Cost: %d\n", cost)

	savingsMap := make(map[int]int)
	routesWithSavings := 0
	routesWithMajorSavings := 0
	totalSavings := 0

	for visitedPoint, firstHalfCost := range visited {
		for _, cheatSteps := range cheats {
			if !grid.isWall(visitedPoint.Add(cheatSteps.wallStep)) {
				continue
			}

			lowestRoute := math.MaxInt

			for _, cheatTarget := range cheatSteps.nextStep {
				cheatTarget = visitedPoint.Add(cheatTarget)

				if grid.isWall(cheatTarget) {
					continue
				}

				newCost, ok := visited[cheatTarget]
				if ok && newCost < lowestRoute {
					lowestRoute = newCost
				}
			}

			if lowestRoute == math.MaxInt {
				continue
			}

			secondHalfCost := cost - lowestRoute
			newCost := firstHalfCost + 2 + secondHalfCost
			if newCost < cost {
				if debug {
					fmt.Printf("Cheat on %v reduces cost by %d to %d\n", cheatSteps.wallStep, cost-newCost, newCost)
				}
				routesWithSavings++
				totalSavings += cost - newCost
				savingsMap[cost-newCost]++
				if cost-newCost >= 100 {
					routesWithMajorSavings++
				}
			}
		}
	}

	if quickDebug {
		times := make([]int, 0, len(savingsMap))
		for timeSaved := range savingsMap {
			times = append(times, timeSaved)
		}
		slices.Sort(times)
		for _, timeSaved := range times {
			fmt.Printf("There are %d cheats that save %d picoseconds.\n", savingsMap[timeSaved], timeSaved)
		}
	}

	fmt.Printf("Routes with Savings: %d\n", routesWithSavings)
	fmt.Printf("Routes with >100ps Savings: %d\n", routesWithMajorSavings)

	utils.PrintMemUsage()
}

func loadData() dijkstraGrid {
	defer utils.TimeTrack(time.Now(), "loadData")
	dataFile, err := os.Open("2024/input-20.txt")

	if err != nil {
		panic(err)
	}

	defer utils.CloseWithLog(dataFile)

	scanner := bufio.NewReader(dataFile)

	var data []bool
	var x, y int
	var c byte
	var start, end image.Point

	for y = 0; true; y++ {
		line, err := scanner.ReadSlice('\n')

		if len(line) == 0 || line[0] == '\n' {
			break
		}
		if err != nil {
			panic(err)
		}

		for x, c = range line {
			switch c {
			case '\n':
				break
			case '#':
				data = append(data, true)
			case 'S':
				start = image.Point{X: x, Y: y}
				data = append(data, false)
			case 'E':
				end = image.Point{X: x, Y: y}
				data = append(data, false)
			default:
				data = append(data, false)
			}
		}
	}

	fmt.Println(len(data), x, y, start, end)

	return dijkstraGrid{
		Data:   data,
		Width:  x,
		Height: y,
		Start:  start,
		End:    end,
	}
}
