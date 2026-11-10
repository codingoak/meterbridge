package main

import "testing"

func TestFloat32ToRegisters(t *testing.T) {
	// Each case: input value and the expected high/low register.
	tests := []struct {
		name string
		in   float32
		high uint16
		low  uint16
	}{
		{"voltage", 230.0, 0x4366, 0x0000},
		{"current", 10.0, 0x4120, 0x0000},
		{"power", 2300.0, 0x450F, 0xC000},
		{"energy", 1234.5, 0x449A, 0x5000},
		{"zero", 0.0, 0x0000, 0x0000},
		{"negative", -1.0, 0xBF80, 0x0000},
	}

	for _, tt := range tests {
		// t.Run creates a named subtest, shown in the test output.
		t.Run(tt.name, func(t *testing.T) {
			got := float32ToRegisters(tt.in)

			if len(got) != 2 {
				t.Fatalf("expected 2 registers, got %d", len(got))
			}
			if got[0] != tt.high || got[1] != tt.low {
				t.Errorf("float32ToRegisters(%v) = 0x%04X 0x%04X, want 0x%04X 0x%04X",
					tt.in, got[0], got[1], tt.high, tt.low)
			}
		})
	}
}
