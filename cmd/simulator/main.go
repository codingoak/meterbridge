package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"os/signal"
	"time"

	"github.com/simonvetter/modbus"
)

type measurement struct {
	name  string
	unit  string
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

// meterHandler serves meter values as Modbus holding registers.
// It implements the modbus.RequestHandler interface.
type meterHandler struct {
	measurements []measurement
}

// registers builds the full register map: two registers per measurement.
func (h *meterHandler) registers() []uint16 {
	regs := make([]uint16, 0, len(h.measurements)*2)
	for _, m := range h.measurements {
		regs = append(regs, float32ToRegisters(m.value)...)
	}
	return regs
}

func (h *meterHandler) HandleHoldingRegisters(req *modbus.HoldingRegistersRequest) ([]uint16, error) {
	// Read-only: reject any write requests.
	if req.IsWrite {
		return nil, modbus.ErrIllegalFunction
	}

	regs := h.registers()
	end := int(req.Addr) + int(req.Quantity)

	// Reject requests outside the register map
	if end > len(regs) {
		return nil, modbus.ErrIllegalDataAddress
	}

	return regs[req.Addr:end], nil
}

// The meter does not support the other Modbus functions.
func (h *meterHandler) HandleCoils(req *modbus.CoilsRequest) ([]bool, error) {
	return nil, modbus.ErrIllegalFunction
}

func (h *meterHandler) HandleDiscreteInputs(req *modbus.DiscreteInputsRequest) ([]bool, error) {
	return nil, modbus.ErrIllegalFunction
}

func (h *meterHandler) HandleInputRegisters(req *modbus.InputRegistersRequest) ([]uint16, error) {
	return nil, modbus.ErrIllegalFunction
}

func main() {
	handler := &meterHandler{
		measurements: []measurement{
			{name: "voltage", unit: "V", value: 230.0},
			{name: "current", unit: "A", value: 10.0},
			{name: "power", unit: "W", value: 2300.0},
			{name: "energy", unit: "kWh", value: 1234.5},
		},
	}

	// Print the register map once at startup.
	for i, m := range handler.measurements {
		fmt.Printf("%2d %-8s %8.1f %s\n", i*2, m.name, m.value, m.unit)
	}

	server, err := modbus.NewServer(&modbus.ServerConfiguration{
		URL:        "tcp://0.0.0.0:5502",
		Timeout:    30 * time.Second,
		MaxClients: 5,
	}, handler)
	if err != nil {
		log.Fatalf("create server: %v", err)
	}

	// Start() does not block; the server runs in the background.
	if err := server.Start(); err != nil {
		log.Fatalf("start server: %v", err)
	}
	log.Println("simulator listening on :5502")

	// Block until Ctrl+C, then shut down cleanly.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig

	log.Println("shutting down...")
	server.Stop()
}
