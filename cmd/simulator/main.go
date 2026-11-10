package main

import (
	"fmt"
	"math"
)

type measurement struct {
	name string
	unit string
	value float32
}

// float32ToRegisters splits a float32 into two 16-bit Modbus registers.
// Byte order: big endian, high word first.
func float32ToRegisters(v float32) []uint16 {
	bits := math.Float32bits(v) // float32 -> uint32 (same bit pattern)
	high := uint16(bits >> 16)  // upper 16 bits
	low := uint16(bits)         // lower 16 bits
	return []uint16{high, low}
}

func main() {
	// Static test values; order defines the register map (2 registers each).
	measurements := []measurement{
		{name: "voltage", unit: "V", value: 230.0},
		{name: "current", unit: "A", value: 10.0},
		{name: "power", unit: "W", value: 2300.0},
		{name: "energy", unit: "kWh", value: 1234.5},
	}

	for i, m := range measurements {
		regs := float32ToRegisters(m.value)
		addr := i * 2 // each float32 occupies two registers
		fmt.Printf("%2d %-8s %8.1f %-3s -> 0x%04X 0x%04X\n",
			addr, m.name, m.value, m.unit, regs[0], regs[1])
	}
}
