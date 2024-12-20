package main

import (
	"fmt"
	"math"
	"strings"
	utils "tea-cats.co.uk/aoc/2024"
	"time"
)

const debug = true

type bitMemory struct {
	maxBits   int
	bitsSet   uint64
	bitValues uint64
}

func (mem *bitMemory) canSet(bit int, value bool) bool {
	if bit < 0 {
		return false
	}
	if bit >= mem.maxBits {
		return !value
	}
	if mem.bitsSet&(1<<bit) == 0 {
		return true
	}
	return (value && (mem.bitValues&(1<<bit)) != 0) || (!value && (mem.bitValues&(1<<bit)) == 0)
}

func (mem *bitMemory) canSetWord(word int, value int) bool {
	return mem.canSet(word, value&1 == 1) &&
		mem.canSet(word+1, value&2 == 2) &&
		mem.canSet(word+2, value&4 == 4)
}

func (mem *bitMemory) set(bit int, value bool) {
	if !mem.canSet(bit, value) {
		panic("invalid state")
	}
	if bit >= mem.maxBits {
		return
	}
	mem.bitsSet |= 1 << bit
	if value {
		mem.bitValues |= 1 << bit
	}
}

func (mem *bitMemory) setWord(word int, value int) {
	mem.set(word, value&1 == 1)
	mem.set(word+1, value&2 == 2)
	mem.set(word+2, value&4 == 4)
}

func (mem *bitMemory) clone() bitMemory {
	return bitMemory{
		bitsSet:   mem.bitsSet,
		bitValues: mem.bitValues,
		maxBits:   mem.maxBits,
	}
}

func (mem *bitMemory) str() string {
	str := ""

	for i := uint64(1) << (mem.maxBits - 1); i > 0; i >>= 1 {
		if mem.bitsSet&i == 0 {
			str += "_"
		} else if mem.bitValues&i == 0 {
			str += "0"
		} else {
			str += "1"
		}
	}

	return str
}

func main() {
	const expectedOutput = 0o33

	if debug {
		explain()
	}

	defer utils.TimeTrack(time.Now(), "main")

	target := []int{3, 4}
	memory := bitMemory{maxBits: 3 * len(target)}

	result := explore(memory, target, 0)

	if test(result) != expectedOutput {
		panic("invalid result")
	}

	fmt.Println(result)
}

func test(a uint64) uint64 {
	b := uint64(0)
	c := uint64(0)
	out := uint64(0)

	for a > 0 {
		b = a & 7
		b = b ^ 5
		c = a >> b
		a = a >> 3
		out = out << 3
		out += (b ^ c ^ 6) & 7
	}

	return out
}

func explore(initialMemory bitMemory, targets []int, outputIndex int) uint64 {
	if outputIndex == len(targets) {
		fmt.Printf("=== REACHED A SOLUTION %v\n", initialMemory)
		return initialMemory.bitValues
	}

	targetValue := targets[outputIndex]
	var minimum uint64 = math.MaxUint64

	for lastThreeBitsOfA := 0; lastThreeBitsOfA < 8; lastThreeBitsOfA++ {
		if debug {
			fmt.Printf("%sConsidering A&7 = %d for output %d (target B = %d):", strings.Repeat("  ", outputIndex), lastThreeBitsOfA, outputIndex, targetValue)
		}
		if !initialMemory.canSetWord(outputIndex*3, lastThreeBitsOfA) {
			if debug {
				fmt.Printf(" conflict\n")
			}
			continue
		}

		memory := initialMemory.clone()
		memory.setWord(outputIndex*3, lastThreeBitsOfA)

		if debug {
			fmt.Printf(" config A:{ %s -> %s }", initialMemory.str(), memory.str())
		}

		// B(intermediate) = lastThreeBitsOfA xor 5
		// C = A >> B(intermediate)
		targetZoneForC := outputIndex*3 + (lastThreeBitsOfA ^ 5)
		// B(target) = B(intermediate) xor (C & 7) xor 6
		// B(target) = (C & 7) xor B(intermediate) xor 6
		// C & 7 = B(target) xor B(intermediate) xor 6
		// C & 7 = B(target) xor (lastThreeBitsOfA xor 5) xor 6
		requiredLastThreeBitsOfC := targetValue ^ 6 ^ (lastThreeBitsOfA ^ 5)
		if debug {
			fmt.Printf("\n%sValue of C&7 taken from bits %d->%d, needs to be %d:", strings.Repeat("  ", outputIndex+1), targetZoneForC, targetZoneForC+2, requiredLastThreeBitsOfC)
		}
		if !memory.canSetWord(targetZoneForC, requiredLastThreeBitsOfC) {
			if debug {
				fmt.Printf(" conflict\n")
			}
			continue
		}

		memory.setWord(targetZoneForC, requiredLastThreeBitsOfC)

		fmt.Printf(" config C: { -> %s }\n", memory.str())
		inside := explore(memory, targets, outputIndex+1)
		if inside < minimum {
			minimum = inside
		}
	}

	return minimum
}
