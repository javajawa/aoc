package day13

import (
	"errors"
	"fmt"
	"image"
	"io"
	"os"
	utils "tea-cats.co.uk/aoc/2024"
	"time"
)

type Request struct {
	ButtonA image.Point
	ButtonB image.Point
	Target  image.Point
}

func LoadData() []Request {
	defer utils.TimeTrack(time.Now(), "loadData")
	dataFile, err := os.Open("2024/input-13.txt")

	if err != nil {
		panic(err)
	}

	defer utils.CloseWithLog(dataFile)

	requests := make([]Request, 0)
	var count int

	for i := 0; i < 320; i++ {
		r := Request{}

		count, err = fmt.Fscanf(dataFile, "Button A: X+%d, Y+%d\n", &r.ButtonA.X, &r.ButtonA.Y)
		if err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
				break
			}
			panic(err)
		}
		if count != 2 {
			panic("invalid input")
		}
		count, err = fmt.Fscanf(dataFile, "Button B: X+%d, Y+%d\n", &r.ButtonB.X, &r.ButtonB.Y)
		if err != nil {
			panic(err)
		}
		if count != 2 {
			panic("invalid input")
		}
		count, err = fmt.Fscanf(dataFile, "Prize: X=%d, Y=%d\n\n", &r.Target.X, &r.Target.Y)
		if err != nil {
			panic(err)
		}
		if count != 2 {
			panic("invalid input")
		}

		requests = append(requests, r)
	}

	return requests
}
