package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"tea-cats.co.uk/aoc/2024"
	"time"
)

const debug = false

type node uint16

func (n node) String() string {
	return string([]byte{byte(n >> 8), byte(n & 0xff)})
}

type network uint64

func normaliseIdent(node1 node, node2 node, node3 node) network {
	slice := []node{node1, node2, node3}
	slices.Sort(slice)

	return network(slice[0]) | network(slice[1])<<16 | network(slice[2])<<32
}

func (n network) String() string {
	result := []byte{'*', '*', '-', '*', '*', '-', '*', '*'}
	for i := 0; i < 9; i += 3 {
		result[i+1] = byte(n & 0xFF)
		n >>= 8
		result[i] = byte(n & 0xFF)
		n >>= 8
	}
	return string(result)
}

const TMask node = 0xFF00
const TValue node = 't' << 8

type neighbours map[node]map[node]struct{}

func (c neighbours) addEdge(start []byte, end []byte) {
	if len(start) != 2 || len(end) != 2 {
		panic("Wrong data size")
	}

	startNode := (node(start[0]) << 8) + node(start[1])
	endNode := (node(end[0]) << 8) + node(end[1])

	_, ok := c[startNode]
	if !ok {
		c[startNode] = make(map[node]struct{})
	}
	_, ok = c[endNode]
	if !ok {
		c[endNode] = make(map[node]struct{})
	}

	c[startNode][endNode] = struct{}{}
	c[endNode][startNode] = struct{}{}
}

func (c neighbours) uniquesStartingWithT() map[network]struct{} {
	networks := make(map[network]struct{})

	for anchor, firstConnections := range c {
		if anchor&TMask != TValue {
			continue
		}

		for step1 := range firstConnections {
			for step2 := range c[step1] {
				_, ok := c[step2][anchor]

				if ok {
					networks[normaliseIdent(step1, step2, anchor)] = struct{}{}
				}
			}
		}
	}
	return networks
}

func (c neighbours) String() string {
	result := ""
	for left, connections := range c {
		for right := range connections {
			result += fmt.Sprintf("%s-%s\n", left.String(), right.String())
		}
	}
	return result
}

func main() {
	defer utils.TimeTrack(time.Now(), "main")

	c := loadData()
	if debug {
		fmt.Println(c.String())
	}

	networks := c.uniquesStartingWithT()
	if debug {
		for net := range networks {
			fmt.Println(net.String())
		}
	}
	fmt.Printf("Networks found: %d\n", len(networks))
}

func loadData() neighbours {
	defer utils.TimeTrack(time.Now(), "loadData")
	dataFile, err := os.Open("2024/input-23.txt")

	if err != nil {
		panic(err)
	}

	defer utils.CloseWithLog(dataFile)

	start := make([]byte, 2)
	end := make([]byte, 2)
	b := make([]byte, 1)
	c := make(neighbours)

	for {
		_, err = dataFile.Read(start)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			if errors.Is(err, io.ErrUnexpectedEOF) {
				break
			}
			panic(err)
		}
		_, err = dataFile.Read(b)
		if err != nil || b[0] != '-' {
			panic(err)
		}
		_, err = dataFile.Read(end)
		if err != nil {
			panic(err)
		}
		_, err = dataFile.Read(b)
		if err != nil || b[0] != '\n' {
			panic(err)
		}

		c.addEdge(start, end)
	}

	return c
}
