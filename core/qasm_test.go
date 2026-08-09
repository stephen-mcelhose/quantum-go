package core

import (
	"strings"
	"testing"
)

func TestToQASM(t *testing.T) {
	p := NewProgram(2)
	p.AddStep(NewStep(NewHadamard(0)))
	p.AddStep(NewStep(NewCnot(0, 1)))
	p.AddStep(NewStep(NewS(0)))
	p.AddStep(NewStep(NewT(1)))
	sInv := NewS(0)
	sInv.SetInverse(true)
	p.AddStep(NewStep(sInv))

	qasm := p.ToQASM()

	expectedLines := []string{
		"OPENQASM 2.0;",
		"include \"qelib1.inc\";",
		"qreg q[2];",
		"creg c[2];",
		"h q[0];",
		"cx q[0], q[1];",
		"s q[0];",
		"t q[1];",
		"sdg q[0];",
	}

	for _, line := range expectedLines {
		if !strings.Contains(qasm, line) {
			t.Errorf("expected QASM to contain %q, but it did not.\nFull QASM:\n%s", line, qasm)
		}
	}
}

func TestToQASMDecomposition(t *testing.T) {
	// Test that Fourier block is decomposed
	p := NewProgram(2)
	p.AddStep(NewStep(NewFourier(2, 0)))

	qasm := p.ToQASM()

	// QFT for 2 qubits involves H and CR gates
	if !strings.Contains(qasm, "h q[0]") {
		t.Errorf("expected QASM to contain h q[0] from decomposed QFT")
	}
	if !strings.Contains(qasm, "cu1") {
		t.Errorf("expected QASM to contain cu1 from decomposed QFT")
	}
	if !strings.Contains(qasm, "swap") {
		t.Errorf("expected QASM to contain swap from decomposed QFT")
	}
}
