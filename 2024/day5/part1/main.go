package main

import (
	"fmt"
	"io"
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

nextRun:
	for _, printRun := range printRuns {
		seenPages := make([]uint8, 5)
		for _, page := range printRun {
			for _, pageICantFollow := range rules[page] {
				if slices.Contains(seenPages, pageICantFollow) {
					//fmt.Printf("Order %d fails validation: %d can not follow %d  (%v)\n", order, pageICantFollow, page, printRun)
					continue nextRun
				}
			}
			seenPages = append(seenPages, page)
		}
		//fmt.Printf("Order %d validates OK!  (%v)\n", order, printRun)
		counter += uint32(printRun[len(printRun)>>1])
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
