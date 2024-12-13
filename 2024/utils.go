package utils

import (
	"image"
	"log"
	"os"
	"time"
)

type Set[T comparable] map[T]struct{}

func NewSet[T comparable]() Set[T] {
	return make(Set[T])
}

func (m *Set[T]) Add(v T) {
	(*m)[v] = struct{}{}
}

func (m *Set[T]) Clear() {
	for k := range *m {
		delete(*m, k)
	}
}

func (m *Set[T]) Contains(v T) bool {
	_, ok := (*m)[v]
	return ok
}

func (m *Set[T]) Union(other Set[T]) {
	for v := range other {
		(*m)[v] = struct{}{}
	}
}

func (m *Set[T]) AddAll(v []T) {
	for _, v := range v {
		(*m)[v] = struct{}{}
	}
}

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

type Grid[T any] struct {
	Data   []T
	Width  int
	Height int
}

func (grid *Grid[T]) AtPoint(point image.Point) *T {
	return grid.At(point.X, point.Y)
}

func (grid *Grid[T]) At(x int, y int) *T {
	if x < 0 || y < 0 || x >= grid.Width || y >= grid.Height {
		return nil
	}
	return &grid.Data[y*grid.Width+x]
}
