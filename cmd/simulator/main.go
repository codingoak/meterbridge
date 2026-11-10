package main

import (
	"fmt"
	"math"
)

// float32ToRegisters splits a float32 into two 16-bit Modbus registers.
// Byte order: big endian, high word first.
func float32ToRegisters(v float32) []uint16 {
	bits := math.Float32bits(v) // float32 -> uint32 (same bit pattern)
	high := uint16(bits >> 16)  // upper 16 bits
	low := uint16(bits)         // lower 16 bits
	return []uint16{high, low}
}

func main() {
	names := []string{"voltage", "current", "power", "energy"}
	values := []float32{230.0, 10.0, 2300.0, 1234.5}

	for i, v := range values {
		regs := float32ToRegisters(v)
		fmt.Printf("%-8s %8.1f -> 0x%04X 0x%04X\n", names[i], v, regs[0], regs[1])
	}
}
