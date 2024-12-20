package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"tea-cats.co.uk/aoc/2024"
	"time"
)

func main() {
	defer utils.TimeTrack(time.Now(), "main")

	towels, requests := loadData()
	count := 0
	distinct := 0

	for _, request := range requests {
		count += countOptions(request, &towels)
		if canDistinct(request, &towels, map[string]bool{}) {
			distinct++
		}
	}

	fmt.Printf("Count: %d\n", count)
	fmt.Printf("Distinct: %d\n", distinct)
	utils.PrintMemUsage()
}

var memo = map[string]int{"": 1}

func countOptions(request string, availableTowels *[]string) int {
	count, ok := memo[request]

	if ok {
		return count
	}

	for _, towel := range *availableTowels {
		if strings.HasPrefix(request, towel) {
			count += countOptions(strings.TrimPrefix(request, towel), availableTowels)
		}
	}

	memo[request] = count

	return count
}

func canDistinct(request string, availableTowels *[]string, used map[string]bool) bool {
	if request == "" {
		return true
	}

	for _, towel := range *availableTowels {
		if !strings.HasPrefix(request, towel) {
			continue
		}
		if used[towel] {
			continue
		}

		innerUsed := make(map[string]bool)
		for t := range used {
			innerUsed[t] = true
		}
		innerUsed[towel] = true

		if canDistinct(towel, availableTowels, innerUsed) {
			return true
		}
	}

	return false
}

func loadData() ([]string, []string) {
	defer utils.TimeTrack(time.Now(), "loadData")
	dataFile, err := os.Open("2024/input-19.txt")

	if err != nil {
		panic(err)
	}

	defer utils.CloseWithLog(dataFile)

	scan := bufio.NewReader(dataFile)

	towelSpec, err := scan.ReadString('\n')
	if err != nil {
		panic(err)
	}
	towelSpec = strings.TrimSuffix(towelSpec, "\n")
	towelSpec = strings.Replace(towelSpec, ", ", " ", -1)
	towels := strings.Fields(towelSpec)

	_, err = scan.ReadString('\n')
	if err != nil {
		panic(err)
	}

	targets := make([]string, 0, 400)
	for {
		request, err := scan.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			panic(err)
		}
		if request == "\n" {
			break
		}
		targets = append(targets, strings.TrimSuffix(request, "\n"))
	}

	return towels, targets
}
