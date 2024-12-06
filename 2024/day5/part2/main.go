package main

import (
	"fmt"
	"io"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
	"tea-cats.co.uk/aoc/2024"
	"time"
)

func main() {
	defer utils.TimeTrack(time.Now(), "main")

	rules, printRuns := loadData()

	defer utils.TimeTrack(time.Now(), "process")

	counter := uint32(0)

	for order, printRun := range printRuns {
		revised := false
		order += 0

		for currentIndex, page := range printRun {
			cutIndex := math.MaxInt

			for _, pageICantFollow := range rules[page] {
				idx := slices.Index(printRun[:currentIndex], pageICantFollow)

				if idx > -1 && idx < cutIndex {
					cutIndex = idx
				}
			}

			if cutIndex == math.MaxInt {
				continue
			}

			//fmt.Printf("Revising Order %d by moving %d (pos %d) before %d (pos %d)\n", order, page, currentIndex, printRun[cutIndex], cutIndex)

			revisedRun := make([]uint8, 0)
			revisedRun = append(revisedRun, printRun[:cutIndex]...)
			revisedRun = append(revisedRun, page)
			revisedRun = append(revisedRun, printRun[cutIndex:currentIndex]...)
			revisedRun = append(revisedRun, printRun[currentIndex+1:]...)

			copy(printRun, revisedRun)
			revised = true
		}

		if revised {
			counter += uint32(printRun[len(printRun)>>1])
		}
	}

	fmt.Printf("Counter: %d\n", counter)
}

func loadData() (map[uint8][]uint8, [][]uint8) {
	defer utils.TimeTrack(time.Now(), "loadData")
	dataFile, err := os.Open("2024/input-5.txt")

	if err != nil {
		panic(err)
	}

	defer utils.CloseWithLog(dataFile)

	rules := make(map[uint8][]uint8)

	for {
		var pageMustComeBefore uint8
		var pageMustComeLater uint8

		fields, err := fmt.Fscanf(dataFile, "%d|%d\n", &pageMustComeBefore, &pageMustComeLater)
		if err != nil {
			if err.Error() == "unexpected newline" {
				break
			}
			panic(err)
		}
		if fields != 2 {
			panic("unexpected fields")
		}

		rules[pageMustComeBefore] = append(rules[pageMustComeBefore], pageMustComeLater)
	}

	var line string
	printRuns := make([][]uint8, 0)

	for run := 0; true; run++ {
		_, err := fmt.Fscanln(dataFile, &line)

		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		strData := strings.Split(line, ",")
		intData := make([]uint8, len(strData))

		for printOrder, pageNumber := range strData {
			page, err := strconv.Atoi(pageNumber)
			if err != nil {
				panic(err)
			}
			intData[printOrder] = uint8(page)
		}

		printRuns = append(printRuns, intData)
	}

	//fmt.Println(rules)
	//fmt.Println(printRuns)
	return rules, printRuns
}
