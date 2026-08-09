// Package local provides quantum state computation functions.
// This file contains the core state vector transformation logic that applies
// quantum gates to state vectors efficiently.
package local

import (
	"github.com/stephen-mcelhose/quantum-go/core"
)

func init() {
	core.GlobalStepExecutor = CalculateNewState
}

// CalculateNewState calculates the next state vector given a list of gates.
// It applies each gate to the state vector in sequence, using optimized implementations
// when available. Identity gates are skipped for efficiency.
func CalculateNewState(gates []core.Gate, vector []complex128, numQubits int) []complex128 {
	w := vector
	for _, gate := range gates {
		w = applyGate(gate, w, nil)
	}
	return w
}

func applyGate(gate core.Gate, v []complex128, controls []int) []complex128 {
	if _, ok := gate.(*core.Identity); ok {
		return v
	}

	// Handle ControlledGate by merging controls
	if cg, ok := gate.(*core.ControlledGate); ok {
		newControls := make([]int, 0, len(controls)+len(cg.ControlIndexes))
		newControls = append(newControls, controls...)
		newControls = append(newControls, cg.ControlIndexes...)
		return applyGate(cg.RootGate, v, newControls)
	}

	// Handle ControlledBlockGate
	if cbg, ok := gate.(*core.ControlledBlockGate); ok {
		newControls := append(controls, cbg.ControlIndex)
		return applyBlock(cbg.Block, v, newControls, cbg.IsInverse())
	}

	// Handle BlockGate (and its derivatives like Fourier, Add)
	if bgi, ok := gate.(core.BlockGateInterface); ok {
		return applyBlock(bgi.GetBlock(), v, controls, bgi.IsInverse())
	}

	// If no controls and has optimized path, use it
	if len(controls) == 0 && gate.HasOptimization() {
		return gate.ApplyOptimize(v)
	}

	// Fallback to general gate application with controls
	return processGate(gate, v, controls)
}

func applyBlock(block *core.Block, v []complex128, controls []int, inverse bool) []complex128 {
	w := v
	if inverse {
		for i := len(block.Steps) - 1; i >= 0; i-- {
			step := block.Steps[i]
			// We need to apply gates in reverse order and inverted
			for j := len(step.Gates) - 1; j >= 0; j-- {
				g := step.Gates[j]
				g.SetInverse(!g.IsInverse())
				w = applyGate(g, w, controls)
				g.SetInverse(!g.IsInverse()) // Restore
			}
		}
	} else {
		for _, step := range block.Steps {
			for _, g := range step.Gates {
				w = applyGate(g, w, controls)
			}
		}
	}
	return w
}

func checkControls(i int, controls []int) bool {
	for _, c := range controls {
		if (i>>c)&1 == 0 {
			return false
		}
	}
	return true
}

// processGate applies a single gate to the state vector, respecting control qubits.
func processGate(gate core.Gate, v []complex128, controls []int) []complex128 {
	size := len(v)
	answer := make([]complex128, size)
	copy(answer, v)

	affected := gate.GetAffectedQubitIndexes()
	
	// Single qubit gates
	if len(affected) == 1 {
		idx := affected[0]
		matrix := gate.GetMatrix()
		qdelta := 1 << idx
		ngroups := size / (2 * qdelta)

		for group := 0; group < ngroups; group++ {
			for j := 2 * group * qdelta; j < (2*group+1)*qdelta; j++ {
				// Both j and j+qdelta have the same bits except at idx.
				// Since idx is not a control bit, we only need to check one of them.
				if checkControls(j, controls) {
					v0 := v[j]
					v1 := v[j+qdelta]
					answer[j] = matrix.Get(0, 0)*v0 + matrix.Get(0, 1)*v1
					answer[j+qdelta] = matrix.Get(1, 0)*v0 + matrix.Get(1, 1)*v1
				}
			}
		}
		return answer
	}

	// 2-qubit gates (optimized)
	if len(affected) == 2 {
		control := affected[0]
		target := affected[1]

		if gate.GetCaption() == "CNOT" {
			for i := 0; i < size; i++ {
				if (i>>control)&1 == 1 && checkControls(i, controls) {
					answer[i] = v[i^(1<<target)]
				}
			}
			return answer
		}

		if gate.GetCaption() == "CZ" {
			for i := 0; i < size; i++ {
				if (i>>control)&1 == 1 && (i>>target)&1 == 1 && checkControls(i, controls) {
					answer[i] = -v[i]
				}
			}
			return answer
		}

		if gate.GetCaption() == "SWAP" {
			for i := 0; i < size; i++ {
				if checkControls(i, controls) {
					bControl := (i >> control) & 1
					bTarget := (i >> target) & 1
					if bControl != bTarget {
						j := i ^ (1 << control) ^ (1 << target)
						answer[i] = v[j]
					}
				}
			}
			return answer
		}

		if gate.GetCaption() == "CR" {
			matrix := gate.GetMatrix()
			rot := matrix.Get(3, 3)
			for i := 0; i < size; i++ {
				if (i>>control)&1 == 1 && (i>>target)&1 == 1 && checkControls(i, controls) {
					answer[i] = v[i] * rot
				}
			}
			return answer
		}
	}

	// 3-qubit gates (optimized)
	if len(affected) == 3 {
		a := affected[0]
		b := affected[1]
		c := affected[2]
		if gate.GetCaption() == "CCNOT" {
			for i := 0; i < size; i++ {
				if (i>>a)&1 == 1 && (i>>b)&1 == 1 && checkControls(i, controls) {
					answer[i] = v[i^(1<<c)]
				}
			}
			return answer
		}
	}

	// Fallback to full matrix multiplication
	matrix := gate.GetMatrix()
	if matrix.Rows == size {
		for i := 0; i < size; i++ {
			if checkControls(i, controls) {
				var sum complex128
				for j := 0; j < size; j++ {
					sum += matrix.Get(i, j) * v[j]
				}
				answer[i] = sum
			}
		}
		return answer
	}

	return answer
}
