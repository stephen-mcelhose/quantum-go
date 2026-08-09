// Package local_test contains integration tests for the quantum simulator.
// It tests various quantum circuits including entangled states and gate operations.
package local_test

import (
	"math/cmplx"
	"testing"
	"github.com/stephen-mcelhose/quantum-go/core"
	"github.com/stephen-mcelhose/quantum-go/local"
)

// TestBellState creates and verifies a Bell state (|Φ+⟩ = (|00⟩ + |11⟩)/√2).
// This is a maximally entangled two-qubit state created using H and CNOT gates.
// Each qubit should have a 50% probability of being measured as |1⟩.
func TestBellState(t *testing.T) {
	// 2 qubits
	p := core.NewProgram(2)
	
	// Step 1: Hadamard on Q0
	s1 := core.NewStep(core.NewHadamard(0))
	
	// Step 2: CNOT with Q0 as control and Q1 as target
	s2 := core.NewStep(core.NewCnot(0, 1))
	
	p.AddSteps(s1, s2)
	
	e := local.NewSimpleExecutionEnvironment()
	res := e.RunProgram(p)
	
	qubits := res.GetQubits()
	if len(qubits) != 2 {
		t.Fatalf("expected 2 qubits, got %d", len(qubits))
	}
	
	// Q0 should have 0.5 probability
	if qubits[0].Probability < 0.49 || qubits[0].Probability > 0.51 {
		t.Errorf("expected Q0 probability ~0.5, got %v", qubits[0].Probability)
	}
	
	// Q1 should have 0.5 probability
	if qubits[1].Probability < 0.49 || qubits[1].Probability > 0.51 {
		t.Errorf("expected Q1 probability ~0.5, got %v", qubits[1].Probability)
	}

	// Verify SetInverse call
	s1.SetInverse(true)
}

// TestGHZState creates and verifies a three-qubit GHZ state (|000⟩ + |111⟩)/√2.
// This is a maximally entangled three-qubit state demonstrating quantum correlation.
func TestGHZState(t *testing.T) {
	// GHZ state: (|000> + |111>) / sqrt(2)
	p := core.NewProgram(3)
	
	p.AddStep(core.NewStep(core.NewHadamard(0)))
	p.AddStep(core.NewStep(core.NewCnot(0, 1)))
	p.AddStep(core.NewStep(core.NewCnot(1, 2)))
	
	e := local.NewSimpleExecutionEnvironment()
	res := e.RunProgram(p)
	
	qubits := res.GetQubits()
	
	for i, q := range qubits {
		if q.Probability < 0.49 || q.Probability > 0.51 {
			t.Errorf("Qubit %d probability expected ~0.5, got %v", i, q.Probability)
		}
	}
}

// TestOtherGates verifies the execution of various gate types including
// Pauli gates (X, Y, Z), Identity, CZ, SWAP, and Measurement.
func TestOtherGates(t *testing.T) {
	p := core.NewProgram(2)
	p.AddStep(core.NewStep(core.NewX(0), core.NewY(1)))
	p.AddStep(core.NewStep(core.NewZ(0), core.NewIdentity(1)))
	p.AddStep(core.NewStep(core.NewCz(0, 1)))
	p.AddStep(core.NewStep(core.NewSwap(0, 1)))
	p.AddStep(core.NewStep(core.NewMeasurement(0)))

	e := local.NewSimpleExecutionEnvironment()
	res := e.RunProgram(p)
	if res == nil {
		t.Fatal("expected result")
	}
	if len(res.GetProbability()) != 4 {
		t.Errorf("expected 4 states, got %d", len(res.GetProbability()))
	}
	
	_ = core.CalculateQubitStatesFromVector(res.GetProbability())
}

// TestMeasure verifies qubit measurement returns the correct deterministic value
// for states that are definitely |0⟩ or |1⟩ (not in superposition).
func TestMeasure(t *testing.T) {
	p := core.NewProgram(1)
	p.AddStep(core.NewStep(core.NewX(0)))
	e := local.NewSimpleExecutionEnvironment()
	res := e.RunProgram(p)
	qubits := res.GetQubits()
	if qubits[0].Measure() != 1 {
		t.Errorf("expected measured value 1 for X gate, got %d", qubits[0].Measure())
	}

	p2 := core.NewProgram(1)
	p2.AddStep(core.NewStep(core.NewIdentity(0)))
	res2 := e.RunProgram(p2)
	qubits2 := res2.GetQubits()
	if qubits2[0].Measure() != 0 {
		t.Errorf("expected measured value 0 for I gate, got %d", qubits2[0].Measure())
	}
}

func TestRotationGates(t *testing.T) {
	// PhaseShift gate (formerly R)
	p := core.NewProgram(1)
	p.AddStep(core.NewStep(core.NewPhaseShift(3.14159, 0)))
	e := local.NewSimpleExecutionEnvironment()
	res := e.RunProgram(p)
	if len(res.GetProbability()) != 2 {
		t.Errorf("expected 2 states")
	}

	// Cr gate
	p2 := core.NewProgram(2)
	p2.AddStep(core.NewStep(core.NewX(0)))
	p2.AddStep(core.NewStep(core.NewCr(0, 1, 3.14159)))
	res2 := e.RunProgram(p2)
	if len(res2.GetProbability()) != 4 {
		t.Errorf("expected 4 states")
	}

	// New rotation gates
	p3 := core.NewProgram(1)
	p3.AddStep(core.NewStep(core.NewRx(3.14159, 0)))
	res3 := e.RunProgram(p3)
	if len(res3.GetProbability()) != 2 {
		t.Errorf("expected 2 states")
	}
}

func TestToffoli(t *testing.T) {
	p := core.NewProgram(3)
	p.AddStep(core.NewStep(core.NewX(0), core.NewX(1)))
	p.AddStep(core.NewStep(core.NewToffoli(0, 1, 2)))
	e := local.NewSimpleExecutionEnvironment()
	res := e.RunProgram(p)
	qubits := res.GetQubits()
	if qubits[2].Measure() != 1 {
		t.Errorf("expected Q2 to be 1 after Toffoli(1,1,0), got %v", qubits[2].Measure())
	}
}

func TestFourier(t *testing.T) {
	p := core.NewProgram(2)
	p.AddStep(core.NewStep(core.NewFourier(2, 0)))
	e := local.NewSimpleExecutionEnvironment()
	res := e.RunProgram(p)
	if len(res.GetProbability()) != 4 {
		t.Errorf("expected 4 states")
	}
	// All states should have 0.25 probability
	for i, v := range res.GetProbability() {
		prob := real(v * cmplx.Conj(v))
		if prob < 0.24 || prob > 0.26 {
			t.Errorf("state %d prob expected 0.25, got %v", i, prob)
		}
	}
}

func TestAdd(t *testing.T) {
	// Add 1 + 1 = 2
	// Q0: x0, Q1: x1 (register x)
	// Q2: y0, Q3: y1 (register y)
	// We want to add y to x.
	p := core.NewProgram(4)
	p.AddStep(core.NewStep(core.NewX(0))) // x = 1 (binary 01? wait, x0 is least significant or most?)
	// In Java Add(x0, x1, y0, y1), it seems x0 is start index.
	// If x = [Q0, Q1], y = [Q2, Q3]
	// If Q0 is 1 and Q2 is 1, then x should become 2 (Q0=0, Q1=1).
	p.AddStep(core.NewStep(core.NewX(2))) // y = 1
	p.AddStep(core.NewStep(core.NewAdd(0, 1, 2, 3)))
	e := local.NewSimpleExecutionEnvironment()
	res := e.RunProgram(p)
	qubits := res.GetQubits()
	
	// x should be 2. Binary 10. Q0=0, Q1=1.
	if qubits[0].Probability > 0.01 {
		t.Errorf("expected Q0 to be 0, got %v", qubits[0].Probability)
	}
	if qubits[1].Probability < 0.99 {
		t.Errorf("expected Q1 to be 1, got %v", qubits[1].Probability)
	}
}
