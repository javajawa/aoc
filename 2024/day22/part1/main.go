package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"tea-cats.co.uk/aoc/2024"
	"time"
)

const debug = false

func main() {
	defer utils.TimeTrack(time.Now(), "main")

	secrets := loadData()
	var total uint64 = 0

	for _, secret := range secrets {
		hashed := processSecret(secret, 2000)
		if debug {
			fmt.Printf("%d: %d\n", secret, hashed)
		}
		total += hashed
	}
	fmt.Printf("Total: %d\n", total)
	utils.PrintMemUsage()
}

func processSecret(secret uint64, rounds int) uint64 {
	for i := 0; i < rounds; i++ {
		secret ^= secret << 6
		secret &= 0xFFFFFF
		secret ^= secret >> 5
		// No prune here -- the top 8 bits can't be set due to a right shift
		secret ^= secret << 11
		secret &= 0xFFFFFF
	}
	return secret
}

func loadData() []uint64 {
	defer utils.TimeTrack(time.Now(), "loadData")
	dataFile, err := os.Open("2024/input-22.txt")

	if err != nil {
		panic(err)
	}

	defer utils.CloseWithLog(dataFile)

	var next uint64
	secrets := make([]uint64, 0, 2200)

	for {
		_, err := fmt.Fscanf(dataFile, "%d\n", &next)

		if err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
				break
			}
			panic(err)
		}
		secrets = append(secrets, next)
	}
	return secrets
}
