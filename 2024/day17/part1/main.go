package main

import (
	"fmt"
	"log"
	"os"
	utils "tea-cats.co.uk/aoc/2024"
	"time"
)

const debug = true

type opcode byte

const (
	opcodeADV opcode = '0'
	opcodeBXL opcode = '1'
	opcodeBST opcode = '2'
	opcodeJNZ opcode = '3'
	opcodeBXC opcode = '4'
	opcodeOUT opcode = '5'
	opcodeBDV opcode = '6'
	opcodeCDV opcode = '7'
)

type instruction struct {
	opcode  opcode
	operand uint64
}

var literal0 uint64 = 0
var literal1 uint64 = 1
var literal2 uint64 = 2
var literal3 uint64 = 3

type machineState struct {
	instructions                    []instruction
	instructionPointer              int
	registerA, registerB, registerC *uint64
	values                          [7]*uint64
	output                          []byte
}

func newMachine(instructions []instruction, regA uint64, regB uint64, regC uint64) machineState {
	machine := machineState{
		instructions:       instructions,
		instructionPointer: 0,
		registerA:          &regA,
		registerB:          &regB,
		registerC:          &regC,
		values:             [7]*uint64{&literal0, &literal1, &literal2, &literal3, &regA, &regB, &regC},
		output:             make([]byte, 0),
	}

	return machine
}

func (m *machineState) step() bool {
	if m.instructionPointer < 0 || m.instructionPointer >= len(m.instructions) {
		return false
	}

	step := m.instructions[m.instructionPointer]

	if debug {
		fmt.Printf("pc=%d, instruction=%v", m.instructionPointer, step)
	}

	switch step.opcode {
	case opcodeADV:
		if debug {
			fmt.Printf(" ADV %d div 2^%d", *m.registerA, *m.values[step.operand])
		}
		*m.registerA = *m.registerA / (1 << *m.values[step.operand])
	case opcodeBXL:
		if debug {
			fmt.Printf(" BXL %d xor %d", *m.registerB, step.operand)
		}
		*m.registerB = *m.registerB ^ step.operand
	case opcodeBST:
		if debug {
			fmt.Printf(" BST %d", *m.values[step.operand]&7)
		}
		*m.registerB = *m.values[step.operand] & 7
	case opcodeJNZ:
		if *m.registerA != 0 {
			if debug {
				fmt.Printf(" JMP %d\n", step.operand)
			}
			m.instructionPointer = int(step.operand)
			return true
		}
		if debug {
			fmt.Printf(" No-JMP")
		}
	case opcodeBXC:
		if debug {
			fmt.Printf(" BXC %d xor %d", *m.registerB, *m.registerC)
		}
		*m.registerB ^= *m.registerC
	case opcodeOUT:
		if debug {
			fmt.Printf(" OUT %d (%d)", *m.values[step.operand]&7, *m.values[step.operand])
		}
		m.output = append(m.output, []byte{'0' + byte(*m.values[step.operand]&7), ','}...)
	case opcodeBDV:
		if debug {
			fmt.Printf(" BDV %d div 2^%d", *m.registerA, *m.values[step.operand])
		}
		*m.registerB = *m.registerA / (1 << *m.values[step.operand])
	case opcodeCDV:
		if debug {
			fmt.Printf(" CDV %d div 2^%d", *m.registerA, *m.values[step.operand])
		}
		*m.registerC = *m.registerA / (1 << *m.values[step.operand])
	default:
		log.Fatalf("Unknown instruction %v\n", step)
	}

	if debug {
		fmt.Printf(" => a=%d,b=%d,c=%d\n", *m.registerA, *m.registerB, *m.registerC)
	}

	m.instructionPointer++
	return true
}

func main() {
	defer utils.TimeTrack(time.Now(), "main")

	var machine machineState

	for machine = loadData(); machine.step(); {
	}

	fmt.Printf("Result: %s\n", string(machine.output))
}

func loadData() machineState {
	defer utils.TimeTrack(time.Now(), "loadData")
	dataFile, err := os.Open("2024/input-17.txt")

	if err != nil {
		panic(err)
	}

	defer utils.CloseWithLog(dataFile)

	instructions := make([]instruction, 0)
	var regA, regB, regC uint64
	var byteCode string

	_, _ = fmt.Fscanf(dataFile, "Register A: %d\n", &regA)
	_, _ = fmt.Fscanf(dataFile, "Register B: %d\n", &regB)
	_, _ = fmt.Fscanf(dataFile, "Register C: %d\n", &regC)
	_, _ = fmt.Fscanf(dataFile, "\n")
	_, _ = fmt.Fscanf(dataFile, "Program: %s\n", &byteCode)

	for i := 0; i < len(byteCode)-1; i += 4 {
		instructions = append(instructions, instruction{opcode: opcode(byteCode[i]), operand: uint64(byteCode[i+2] - '0')})
	}

	return newMachine(instructions, regA, regB, regC)
}
