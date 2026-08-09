// Package local_test contains fuzz tests for quantum gate operations.
// Fuzz testing helps discover edge cases in quantum circuit simulation.
package local_test

import (
	"math"
	"math/cmplx"
	"testing"

	"github.com/stephen-mcelhose/quantum-go/core"
	"github.com/stephen-mcelhose/quantum-go/local"
)

// FuzzToffoli uses fuzzing to test the Toffoli (CCNOT) gate implementation
// across all possible three-qubit initial states.
// The Toffoli gate should flip the target qubit if and only if both control qubits are 1.
func FuzzToffoli(f *testing.F) {
	// Seed corpus with all 8 possible 3-qubit states
	for i := 0; i < 8; i++ {
		f.Add(i)
	}

	f.Fuzz(func(t *testing.T, initialState int) {
		if initialState < 0 || initialState > 7 {
			return
		}

		p := core.NewProgram(3)
		// Set initial state
		// We can use X gates to set bits
		if (initialState >> 0) & 1 == 1 {
			p.AddStep(core.NewStep(core.NewX(0)))
		}
		if (initialState >> 1) & 1 == 1 {
			p.AddStep(core.NewStep(core.NewX(1)))
		}
		if (initialState >> 2) & 1 == 1 {
			p.AddStep(core.NewStep(core.NewX(2)))
		}

		// Apply Toffoli (control1=0, control2=1, target=2)
		p.AddStep(core.NewStep(core.NewToffoli(0, 1, 2)))

		e := local.NewSimpleExecutionEnvironment()
		res := e.RunProgram(p)

		// Expected state index
		expectedState := initialState
		if (initialState & 3) == 3 { // Both Q0 and Q1 are 1
			expectedState ^= 4 // Flip Q2
		}

			// Verify result
			for i, amp := range res.GetProbability() {
				prob := real(amp * cmplx.Conj(amp))

			if i == expectedState {
				if prob < 0.99 {
					t.Errorf("expected state %d to have prob 1.0, got %v", i, prob)
				}
			} else {
				if prob > 0.01 {
					t.Errorf("expected state %d to have prob 0.0, got %v", i, prob)
				}
			}
		}
	})
}

func FuzzFourier(f *testing.F) {
	f.Add(2)
	f.Add(3)
	f.Add(4)

	f.Fuzz(func(t *testing.T, numQubits int) {
		if numQubits < 2 || numQubits > 6 {
			return
		}

		p := core.NewProgram(numQubits)
		p.AddStep(core.NewStep(core.NewFourier(numQubits, 0)))
		
		// Unitarity check: QFT followed by IQFT should be identity
		invFourier := core.NewFourier(numQubits, 0)
		invFourier.SetInverse(true)
		p.AddStep(core.NewStep(invFourier))

			e := local.NewSimpleExecutionEnvironment()
			res := e.RunProgram(p)

			// Should be back to |0...0>
			prob0 := real(res.GetProbability()[0] * cmplx.Conj(res.GetProbability()[0]))
			if prob0 < 0.99 {
				t.Errorf("Fourier/InverseFourier identity failed for %d qubits, prob[0]=%v", numQubits, prob0)
			}
			
			// Check normalization
			var totalProb float64
			for _, amp := range res.GetProbability() {
				totalProb += real(amp * cmplx.Conj(amp))
			}

		if math.Abs(totalProb - 1.0) > 1e-9 {
			t.Errorf("Normalization failed: total prob = %v", totalProb)
		}
	})
}

func FuzzAdd(f *testing.F) {
	f.Add(2, 1, 1) // 1+1=2 (mod 4)
	f.Add(3, 3, 1) // 3+1=4 (mod 8)

	f.Fuzz(func(t *testing.T, m, x, y int) {
		if m < 1 || m > 4 {
			return
		}
		maxVal := 1 << m
		x %= maxVal
		y %= maxVal
		if x < 0 { x += maxVal }
		if y < 0 { y += maxVal }

		p := core.NewProgram(2 * m)
		// Initial x in Q0...Qm-1
		for i := 0; i < m; i++ {
			if (x >> i) & 1 == 1 {
				p.AddStep(core.NewStep(core.NewX(i)))
			}
		}
		// Initial y in Qm...Q2m-1
		for i := 0; i < m; i++ {
			if (y >> i) & 1 == 1 {
				p.AddStep(core.NewStep(core.NewX(m + i)))
			}
		}

		// Add y to x
		p.AddStep(core.NewStep(core.NewAdd(0, m-1, m, 2*m-1)))

		e := local.NewSimpleExecutionEnvironment()
		res := e.RunProgram(p)
		qubits := res.GetQubits()

		expectedX := (x + y) % maxVal
		for i := 0; i < m; i++ {
			expectedBit := (expectedX >> i) & 1
			if expectedBit == 1 {
				if qubits[i].Probability < 0.99 {
					t.Errorf("m=%d, x=%d, y=%d: expected Q%d to be 1, prob=%v", m, x, y, i, qubits[i].Probability)
				}
			} else {
				if qubits[i].Probability > 0.01 {
					t.Errorf("m=%d, x=%d, y=%d: expected Q%d to be 0, prob=%v", m, x, y, i, qubits[i].Probability)
				}
			}
		}
	})
}
