package day8

import (
	"fmt"
	"os"
	"strconv"
	"tea-cats.co.uk/aoc/2024"
	"time"
	"unsafe"
)

type AntennaeFrequency byte

const NoAntennaeAtLocation = '.'
const Width = 50
const Height = 50

func LoadData() [][]AntennaeFrequency {
	defer utils.TimeTrack(time.Now(), "loadData")
	dataFile, err := os.Open("2024/input-8.txt")

	if err != nil {
		panic(err)
	}

	defer utils.CloseWithLog(dataFile)
	buffer := make([][]AntennaeFrequency, Height)

	for i := range buffer {
		tempbuffer := make([]AntennaeFrequency, Width+1)

		count, err := dataFile.Read(*((*[]byte)(unsafe.Pointer(&tempbuffer))))
		if err != nil {
			panic(err)
		}
		if count != Width+1 || tempbuffer[Width] != '\n' {
			fmt.Println(tempbuffer)
			fmt.Println(string(tempbuffer))
			panic("Wrong amount of data read on line" + strconv.FormatInt(int64(i), 10))
		}

		buffer[i] = tempbuffer[:Width]
	}

	return buffer
}
