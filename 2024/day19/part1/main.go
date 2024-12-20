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

	for _, request := range requests {
		if hasOptions(request, &towels) {
			count++
		}
	}

	fmt.Printf("Count: %d\n", count)
	utils.PrintMemUsage()
}

var memo = map[string]bool{"": true}

func hasOptions(request string, availableTowels *[]string) bool {
	count, ok := memo[request]

	if ok {
		return count
	}

	for _, towel := range *availableTowels {
		if strings.HasPrefix(request, towel) && hasOptions(request[len(towel):], availableTowels) {
			memo[request] = true
			return true
		}
	}

	memo[request] = false
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
