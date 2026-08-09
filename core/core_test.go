// Package core_test contains unit tests for the core quantum circuit components.
// It tests the basic functionality of gates, steps, programs, and related structures.
package core_test

import (
	"testing"
	"github.com/stephen-mcelhose/quantum-go/core"
)

// TestBaseGate verifies the basic gate interface implementation.
// It tests getter methods and inverse flag handling on a Hadamard gate.
func TestBaseGate(t *testing.T) {
	g := core.NewHadamard(2)
	if g.GetHighestAffectedQubitIndex() != 2 {
		t.Errorf("expected 2, got %d", g.GetHighestAffectedQubitIndex())
	}
	if g.GetName() != "Hadamard" {
		t.Errorf("expected Hadamard, got %s", g.GetName())
	}
	if g.GetSize() != 1 {
		t.Errorf("expected 1, got %d", g.GetSize())
	}
	if g.GetCaption() != "H" {
		t.Errorf("expected H, got %s", g.GetCaption())
	}
	if g.GetGroup() != "" {
		t.Errorf("expected empty group, got %s", g.GetGroup())
	}
	g.SetInverse(true)
}

// TestStep verifies that steps can be created and have their inverse flag set.
func TestStep(t *testing.T) {
	s := core.NewStep(core.NewHadamard(0))
	s.SetInverse(true)
}

// TestStepUnique verifies that adding two gates affecting the same qubit panics.
// This ensures parallel gates in a step operate on disjoint qubit sets.
func TestStepUnique(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic when adding overlapping gates")
		}
	}()
	s := core.NewStep(core.NewHadamard(0), core.NewX(0))
	_ = s
}

// TestGatesMatrix verifies that various gates return non-empty matrices.
func TestGatesMatrix(t *testing.T) {
	gates := []core.Gate{
			core.NewX(0),
			core.NewY(0),
			core.NewZ(0),
			core.NewHadamard(0),
			core.NewIdentity(0),
			core.NewMeasurement(0),
			core.NewS(0),
			core.NewT(0),
			core.NewV(0),
			core.NewRx(0.5, 0),
			core.NewRy(0.5, 0),
			core.NewRz(0.5, 0),
			core.NewPhaseShift(0.5, 0),
			core.NewU(0.5, 0.5, 0.5, 0),
			core.NewCz(0, 1),
			core.NewSwap(0, 1),
			core.NewCnot(0, 1),
			core.NewCr(0, 1, 0.5),
			core.NewToffoli(0, 1, 2),
		}

	for _, g := range gates {
		m := g.GetMatrix()
		if m.Rows == 0 || m.Cols == 0 {
			t.Errorf("gate %s has empty matrix", g.GetCaption())
		}
	}
}

// TestProgram verifies basic program creation and step addition.
func TestProgram(t *testing.T) {
	p := core.NewProgram(2)
	p.AddStep(core.NewStep(core.NewHadamard(0)))
	if p.NumQubits != 2 {
		t.Errorf("expected 2 qubits, got %d", p.NumQubits)
	}
}
