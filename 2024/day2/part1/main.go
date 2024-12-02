package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"tea-cats.co.uk/aoc/2024"
	"time"
)

func main() {
	defer utils.TimeTrack(time.Now(), "main")

	safe := loadData()
	fmt.Printf("Safe: %d\n", safe)
}

func loadData() int {
	defer utils.TimeTrack(time.Now(), "loadData")
	dataFile, err := os.Open("2024/input-2.txt")

	if err != nil {
		panic(err)
	}
	defer utils.CloseWithLog(dataFile)

	safe := 0
	scanner := bufio.NewReader(dataFile)

	for {
		line, err := scanner.ReadString('\n')

		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		strData := strings.Fields(line)
		intData := make([]int, len(strData))

		for i, d := range strData {
			intData[i], err = strconv.Atoi(d)
			if err != nil {
				panic(err)
			}
		}

		if checkSafe(intData) {
			safe += 1
		}
	}

	return safe
}

func checkSafe(readings []int) bool {
	last := 0
	var diff int

	for i := 1; i < len(readings); i++ {
		diff = readings[i] - readings[i-1]

		if diff == 0 || diff > 3 || diff < -3 || last*diff < 0 {
			return false
		}

		last = diff
	}

	return true
}
