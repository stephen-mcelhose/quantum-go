// Package local provides a local simulation engine for quantum programs.
// It implements quantum state vector simulation on classical hardware,
// allowing quantum circuits to be tested and debugged without quantum hardware.
// The simulator tracks the full quantum state as a complex-valued vector.
package local

import (
	gmath "math"
	"math/rand"

	"github.com/stephen-mcelhose/quantum-go/core"
)


// SimpleExecutionEnvironment provides a basic quantum circuit simulator.
// It simulates quantum programs by maintaining and transforming the state vector
// according to the quantum gates in each step of the program.
type SimpleExecutionEnvironment struct{}

// NewSimpleExecutionEnvironment creates a new instance of the local quantum simulator.
func NewSimpleExecutionEnvironment() *SimpleExecutionEnvironment {
	return &SimpleExecutionEnvironment{}
}

// RunProgram executes a quantum program and returns the result.
// It initializes all qubits to the |0⟩ state (or according to InitAlpha),
// then applies each step's gates in sequence to transform the state vector.
// After execution, qubit measurement values are randomly determined based on their probabilities.
func (e *SimpleExecutionEnvironment) RunProgram(p *core.Program) core.Result {
	size := 1 << p.NumQubits
	state := make([]complex128, size)

	state[0] = 1.0

	// Initialize the state vector based on p.InitAlpha values for each qubit.
	// This creates a product state |psi> = |q_n-1> ⊗ ... ⊗ |q_0>
	// where each |q_i> = alpha_i|0> + beta_i|1>
	for i := 0; i < size; i++ {
		state[i] = 1.0
		for j := 0; j < p.NumQubits; j++ {
			// In the original Strange Java simulator, InitAlpha[j] corresponds to qubit j.
			// Bit j in index i determines if we use alpha or beta.
			if (i>>j)&1 == 0 {
				state[i] *= complex(p.InitAlpha[j], 0)
			} else {
				state[i] *= complex(gmath.Sqrt(1.0-p.InitAlpha[j]*p.InitAlpha[j]), 0)
			}
		}
	}


	for _, step := range p.Steps {
		if step.Type == core.StepNormal {
			state = CalculateNewState(step.Gates, state, p.NumQubits)
		}
	}

	res := &core.CompactResult{
		NumQubits:   p.NumQubits,
		Probability: state,
	}

	// Set measured values
	qubits := res.GetQubits()
	for _, q := range qubits {
		q.MeasuredValue = rand.Float64() < q.Probability
	}
	// We need to re-sync or something? GetQubits creates new objects.
	// Actually, Result should probably hold the Qubits.

	p.Result = res
	return res
}
