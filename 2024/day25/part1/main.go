package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"tea-cats.co.uk/aoc/2024"
	"time"
)

const debug = true

const cylinders = 5
const height = 7
const row = cylinders + 1
const totalBlobSize = row*height + 1

type key uint64
type lock uint64

func (k key) matches(l lock) bool {
	return k&key(l) == 0
}

func main() {
	defer utils.TimeTrack(time.Now(), "main")

	locks, keys := loadData()

	if debug {
		for i, lock := range locks {
			fmt.Printf("=========\nLock %d\n%s", i, toStr(lock))
		}
		for i, key := range keys {
			fmt.Printf("=========\nKey %d\n%s", i, toStr(key))
		}
	}

	defer utils.TimeTrack(time.Now(), "match")
	counter := 0
	for _, lock := range locks {
		for _, key := range keys {
			if key.matches(lock) {
				counter++
			}
		}
	}

	if debug {
		fmt.Println()
	}
	fmt.Printf("Potential Locks: %d\n", counter)
}

func loadData() ([]lock, []key) {
	defer utils.TimeTrack(time.Now(), "loadData")
	dataFile, err := os.Open("2024/input-25.txt")

	if err != nil {
		panic(err)
	}

	defer utils.CloseWithLog(dataFile)

	buffer := make([]byte, totalBlobSize)

	keys := make([]key, 0, 250)
	locks := make([]lock, 0, 250)

	for {
		_, err := dataFile.Read(buffer)

		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			panic(err)
		}

		if buffer[0] == '#' {
			locks = append(locks, bufferToObj[lock](buffer))
		} else {
			keys = append(keys, bufferToObj[key](buffer))
		}
	}

	return locks, keys
}

func bufferToObj[T key | lock](data []byte) T {
	ret := T(0)
	for i := 0; i < totalBlobSize; i++ {
		if data[i] == '#' {
			ret |= 1 << i
		}
	}
	return ret
}

func toStr[T key | lock](data T) string {
	ret := make([]byte, row*height)
	for i := 0; i < row*height; i++ {
		if i%row == cylinders {
			ret[i] = '\n'
		} else if data&(1<<i) != 0 {
			ret[i] = '#'
		} else {
			ret[i] = '.'
		}
	}
	return string(ret)
}
