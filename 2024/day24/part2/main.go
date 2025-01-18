package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"tea-cats.co.uk/aoc/2024"
	"time"
)

const fieldSize = 45

type operand byte

const (
	and operand = iota
	or  operand = iota
	xor operand = iota
)

type signal string

type gate struct {
	inputLeft  signal
	inputRight signal
	op         operand
}

type adder struct {
	signals map[signal]bool
	setters map[signal]gate
}

func (a adder) findBugs() []signal {
	// Full adder:
	//
	// int  = A   XOR B
	// S    = int XOR Cin
	// Cout = (A AND B) OR (int AND Cin)

	carryNames := make([]signal, fieldSize)
	adderNames := make([]signal, fieldSize)

	for i := fieldSize - 1; i > 0; i-- {
		gateS, _ := a.setters[signal(fmt.Sprintf("z%02d", i))]
		if gateS.op != xor {
			fmt.Printf("Wrong gate in adder z%02d = %v\n", i, gateS)
		}

		gateLeft, exists := a.setters[gateS.inputLeft]
		if !exists {
			fmt.Printf("Second half adder for z%02d (%v) uses raw input %s\n", i, gateS, gateS.inputLeft)
		}
		gateRight, exists := a.setters[gateS.inputRight]
		if !exists {
			fmt.Printf("Second half adder for z%02d (%v) uses raw input %s\n", i, gateS, gateS.inputRight)
		}

		if gateLeft.op == or && gateRight.op == xor {
			carryNames[i] = gateS.inputLeft
			adderNames[i] = gateS.inputRight
		} else if gateLeft.op == xor && gateRight.op == or {
			adderNames[i] = gateS.inputLeft
			carryNames[i] = gateS.inputRight
		} else {
			fmt.Printf("Input gates to geneate z%02d have wrong operations: %s = %v, %s = %v\n", i, gateS.inputLeft, gateLeft, gateS.inputRight, gateRight)
			continue
		}

		carryGate := a.setters[carryNames[i]]
		adderGate := a.setters[adderNames[i]]
	}
}

func (a adder) resolve() uint64 {
	output := uint64(0)
	for i := 0; i <= fieldSize; i++ {
		if a.find(signal(fmt.Sprintf("z%02d", i))) {
			output |= 1 << i
		}
	}

	return output
}

func (a adder) find(s signal) bool {
	previous, exists := a.signals[s]
	if exists {
		return previous
	}
	gate, exists := a.setters[s]
	if !exists {
		panic("Unknown gate " + s)
	}
	leftOperand := a.find(gate.inputLeft)
	rightOperand := a.find(gate.inputRight)

	val := false
	switch gate.op {
	case and:
		val = leftOperand && rightOperand
	case or:
		val = leftOperand || rightOperand
	case xor:
		val = leftOperand != rightOperand
	}

	a.signals[s] = val
	return val
}

func main() {
	defer utils.TimeTrack(time.Now(), "main")

	adder := loadData()
	fmt.Println(adder.resolve())
	fmt.Println(len(adder.setters))
	fmt.Println(len(adder.signals))
}

func loadData() adder {
	defer utils.TimeTrack(time.Now(), "loadData")
	dataFile, err := os.Open("2024/input-24.txt")

	if err != nil {
		panic(err)
	}

	defer utils.CloseWithLog(dataFile)

	word := make([]byte, 3)
	byte := make([]byte, 1)

	output := adder{
		signals: map[signal]bool{},
		setters: map[signal]gate{},
	}

	// x00-xnn
	for i := 0; i < fieldSize; i++ {
		_, _ = dataFile.Read(word)
		_, _ = dataFile.Read(byte) // :
		_, _ = dataFile.Read(byte) // _
		_, _ = dataFile.Read(byte) // 1 or 0
		output.signals[signal(word)] = byte[0] == '1'
		_, _ = dataFile.Read(byte) // \n
	}
	// y00-ynn
	for i := 0; i < fieldSize; i++ {
		_, _ = dataFile.Read(word)
		_, _ = dataFile.Read(byte) // :
		_, _ = dataFile.Read(byte) // _
		_, _ = dataFile.Read(byte) // 1 or 0
		output.signals[signal(word)] = byte[0] == '1'
		_, _ = dataFile.Read(byte) // \n
	}
	_, _ = dataFile.Read(byte) // \n
	for {
		g := gate{}
		_, err = dataFile.Read(word)

		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			panic(err)
		}

		_, _ = dataFile.Read(byte) // _
		g.inputLeft = signal(word)

		_, _ = dataFile.Read(word)
		switch string(word) {
		case "AND":
			g.op = and
			_, _ = dataFile.Read(byte) // _
		case "OR ":
			g.op = or
		case "XOR":
			g.op = xor
			_, _ = dataFile.Read(byte) // _
		}

		_, _ = dataFile.Read(word)
		_, _ = dataFile.Read(byte) // _
		g.inputRight = signal(word)

		_, _ = dataFile.Read(word) // ->
		_, _ = dataFile.Read(word)
		output.setters[signal(word)] = g

		_, _ = dataFile.Read(byte) // \n
	}

	return output
}
