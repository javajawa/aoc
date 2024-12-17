package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	utils "tea-cats.co.uk/aoc/2024"
	"time"
)

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
	operand int
}

func (i instruction) explain() string {
	operands := map[int]string{0: "0", 1: "1", 2: "2", 3: "3", 4: "A", 5: "B", 6: "C"}

	switch i.opcode {
	case opcodeADV:
		return "A = A / 2^" + operands[i.operand]
	case opcodeBXL:
		return "B = B xor " + strconv.Itoa(i.operand)
	case opcodeBST:
		return "B = (" + operands[i.operand] + " & 7)"
	case opcodeJNZ:
		return fmt.Sprintf("JUMP %02d", i.operand)
	case opcodeBXC:
		return "B = B xor C"
	case opcodeOUT:
		return "OUTPUT " + operands[i.operand] + " & 7"
	case opcodeBDV:
		return "B = A / 2^" + operands[i.operand]
	case opcodeCDV:
		return "C = A / 2^" + operands[i.operand]
	default:
		log.Fatalf("Unknown instruction %v\n")
	}
	return ""
}

func explain() {
	defer utils.TimeTrack(time.Now(), "explain")

	instructions := loadData()

	for i, inst := range instructions {
		fmt.Printf("%02d  %s\n", i, inst.explain())
	}
}

func loadData() []instruction {
	defer utils.TimeTrack(time.Now(), "loadData")
	dataFile, err := os.Open("2024/input-17.txt")

	if err != nil {
		panic(err)
	}

	defer utils.CloseWithLog(dataFile)

	instructions := make([]instruction, 0)
	var regA, regB, regC int
	var byteCode string

	_, _ = fmt.Fscanf(dataFile, "Register A: %d\n", &regA)
	_, _ = fmt.Fscanf(dataFile, "Register B: %d\n", &regB)
	_, _ = fmt.Fscanf(dataFile, "Register C: %d\n", &regC)
	_, _ = fmt.Fscanf(dataFile, "\n")
	_, _ = fmt.Fscanf(dataFile, "Program: %s\n", &byteCode)

	for i := 0; i < len(byteCode)-1; i += 4 {
		instructions = append(instructions, instruction{opcode: opcode(byteCode[i]), operand: int(byteCode[i+2] - '0')})
	}

	return instructions
}
