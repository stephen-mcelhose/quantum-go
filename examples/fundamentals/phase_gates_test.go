package fundamentals_test

import (
	"fmt"

	"github.com/stephen-mcelhose/quantum-go/core"
	"github.com/stephen-mcelhose/quantum-go/local"
)

// Example_phaseGates demonstrates the Clifford-group phase gates (S, T) and the V gate (SX).
// These gates are specific rotations that are fundamental to many quantum algorithms.
func Example_phaseGates() {
	fmt.Println("Phase and Square Root Gates Demo (S, T, V)")

	// 1. S gate (Phase gate) - Square root of Z
	// S * S = Z. Applying S twice to |1> should result in Z|1> = -|1>,
	// which doesn't change measurement probability.
	pS := core.NewProgram(1)
	pS.AddStep(core.NewStep(core.NewX(0))) // Start at |1>
	pS.AddStep(core.NewStep(core.NewS(0)))
	pS.AddStep(core.NewStep(core.NewS(0)))

	// 2. T gate (PI/8 gate) - Square root of S
	// T * T = S.
	pT := core.NewProgram(1)
	pT.AddStep(core.NewStep(core.NewX(0)))
	pT.AddStep(core.NewStep(core.NewT(0)))
	pT.AddStep(core.NewStep(core.NewT(0)))
	// This should be equivalent to one S gate.

	// 3. V gate (SX gate) - Square root of X
	// V * V = X. Applying V twice to |0> should result in |1>.
	pV := core.NewProgram(1)
	pV.AddStep(core.NewStep(core.NewV(0)))
	pV.AddStep(core.NewStep(core.NewV(0)))

	engine := local.NewSimpleExecutionEnvironment()

	fmt.Println("\nS * S on |1> (Probability remains same):")
	resS := engine.RunProgram(pS)
	resS.PrintBinary()

	fmt.Println("\nV * V on |0> (Results in |1>):")
	resV := engine.RunProgram(pV)
	resV.PrintBinary()

	fmt.Println("\nT * T on |1> (Equivalent to one S gate):")
	resT := engine.RunProgram(pT)
	resT.PrintBinary()

	// Output:
	// Phase and Square Root Gates Demo (S, T, V)
	//
	// S * S on |1> (Probability remains same):
	// Quantum Result (1 qubits):
	// |1>: 1.0000
	//
	// V * V on |0> (Results in |1>):
	// Quantum Result (1 qubits):
	// |1>: 1.0000
	//
	// T * T on |1> (Equivalent to one S gate):
	// Quantum Result (1 qubits):
	// |1>: 1.0000
}
