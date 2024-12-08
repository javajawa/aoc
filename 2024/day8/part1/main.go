package main

import (
	"fmt"
	"image"
	"tea-cats.co.uk/aoc/2024"
	"tea-cats.co.uk/aoc/2024/day8"
	"time"
)

const debug = false

func main() {
	defer utils.TimeTrack(time.Now(), "main")

	buffer := day8.LoadData()
	countOfNodes := process(&buffer)

	fmt.Printf("Answer: %d\n", countOfNodes)
}

// An antinode occurs at any point that is perfectly in line with two antennas
// of the same frequency - but only when one of the antennas is twice as far
// away as the other.
// This means that for any pair of antennas with the same frequency,
// there are two antinodes, one on either side of them.
func process(data *[][]day8.AntennaeFrequency) int {
	defer utils.TimeTrack(time.Now(), "process")

	mapBoundingBox := image.Rect(0, 0, day8.Width, day8.Height)
	knownAntennae := make(map[day8.AntennaeFrequency][]image.Point)
	knownAntiNodes := make(utils.Set[image.Point])

	for row, rowData := range *data {
		for column, character := range rowData {
			if character == day8.NoAntennaeAtLocation {
				continue
			}

			point := image.Point{X: column, Y: row}

			if debug {
				fmt.Printf("Found antenae of type %s at %v\n", string(character), point)
			}

			existingTowersOfType := knownAntennae[character]
			for _, previousTower := range existingTowersOfType {
				//  *...🗼...🗼...*
				// Antinodes are mirrors of the tower as seen from each other

				vector := previousTower.Sub(point)
				node1 := point.Sub(vector)
				node2 := previousTower.Add(vector)

				if debug {
					fmt.Printf("  Nodes for pair %v/%v are at %v and %v\n", previousTower, vector, node1, node2)
				}

				if node1.In(mapBoundingBox) {
					knownAntiNodes[node1] = struct{}{}
				}
				if node2.In(mapBoundingBox) {
					knownAntiNodes[node2] = struct{}{}
				}
			}

			knownAntennae[character] = append(existingTowersOfType, point)
		}
	}

	return len(knownAntiNodes)
}
