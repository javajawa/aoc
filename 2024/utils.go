package utils

import (
	"log"
	"os"
	"time"
)

type Set[T comparable] map[T]struct{}

func CloseWithLog(file *os.File) {
	err := file.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%07.3fms: %s", float64(elapsed.Microseconds())/1000.0, name)
}
