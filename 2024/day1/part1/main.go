package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"tea-cats.co.uk/aoc/2024"
	"time"
)

const initialLines = 1000

func main() {
	defer utils.TimeTrack(time.Now(), "main")

	listL, listR := loadData()

	sortStart := time.Now()
	sort.Ints(listL)
	sort.Ints(listR)
	utils.TimeTrack(sortStart, "sort")

	sortStart = time.Now()
	acc := 0

	for i, l := range listL {
		acc += utils.Abs(l - listR[i])
	}

	fmt.Println(acc)
	utils.TimeTrack(sortStart, "process")
}

func loadData() ([]int, []int) {
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
	listR := make([]int, initialLines)

	for {
		n, err := fmt.Fscanln(dataFile, &l, &r)

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
		listR = append(listR, r)
	}
}
