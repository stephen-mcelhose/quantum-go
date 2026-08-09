package qasm

import (
	"testing"
)

func TestParse(t *testing.T) {
	input := `
OPENQASM 2.0;
include "qelib1.inc";
qreg q[2];
creg c[2];
h q[0];
cx q[0], q[1];
measure q[0] -> c[0];
`
	p, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if p.NumQubits != 2 {
		t.Errorf("expected 2 qubits, got %d", p.NumQubits)
	}

	// Steps:
	// 1. h
	// 2. cx
	// 3. measure
	if len(p.Steps) != 3 {
		t.Errorf("expected 3 steps, got %d", len(p.Steps))
	}
}

func TestParseAdvanced(t *testing.T) {
	input := `
OPENQASM 2.0;
include "qelib1.inc";
qreg q[2];
rx(1.57) q[0];
ry(1.57) q[1];
rz(1.57) q[0];
s q[1];
sdg q[0];
t q[1];
tdg q[0];
sx q[1];
u3(1.57, 0, 3.14) q[0];
barrier q[0], q[1];
measure q[0] -> c[0];
`
	p, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Rx, Ry, Rz, S, sdg, T, tdg, sx, u3, (barrier ignored), measure
	// Total steps should be 10 (excluding barrier)
	if len(p.Steps) != 10 {
		t.Errorf("expected 10 steps, got %d", len(p.Steps))
	}
}

func TestParseU1(t *testing.T) {
	input := `
qreg q[1];
u1(1.570796) q[0];
`
	p, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(p.Steps) != 1 {
		t.Errorf("expected 1 step, got %d", len(p.Steps))
	}
}
