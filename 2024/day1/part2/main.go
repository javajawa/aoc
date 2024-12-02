package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"tea-cats.co.uk/aoc/2024"
	"time"
)

const initialLines = 1000

func main() {
	defer utils.TimeTrack(time.Now(), "main")
	dataFile, err := os.Open("2024/input-1.txt")

	if err != nil {
		log.Fatal(err)
	}

	defer utils.CloseWithLog(dataFile)

	listL, listR := processInputFile(dataFile)

	sortStart := time.Now()
	acc := 0

	for _, val := range listL {
		acc += val * listR[val]
	}

	fmt.Println(acc)
	utils.TimeTrack(sortStart, "process")
}

func processInputFile(file *os.File) ([]int, map[int]int) {
	defer utils.TimeTrack(time.Now(), "loadData")
	dataFile, err := os.Open("2024/input-1.txt")

	if err != nil {
		log.Fatal(err)
	}

	defer utils.CloseWithLog(dataFile)

	var (
		l int
		r int
	)

	listL := make([]int, initialLines)
	listR := make(map[int]int)

	for {
		n, err := fmt.Fscanln(file, &l, &r)

		if err == io.EOF {
			return listL, listR
		}

		if err != nil {
			log.Fatal(err)
		}
		if n != 2 {
			log.Fatal("Expected 2 values")
		}

		listL = append(listL, l)
		listR[r] += 1
	}
}
