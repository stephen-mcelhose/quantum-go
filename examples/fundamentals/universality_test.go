package fundamentals_test

import (
	"fmt"
	"math"

	"github.com/stephen-mcelhose/quantum-go/core"
	"github.com/stephen-mcelhose/quantum-go/local"
)

// Example_universality demonstrates the power of the Universal Rotation gate (U).
// The U(theta, phi, lambda) gate can represent any single-qubit unitary transformation.
func Example_universality() {
	fmt.Println("Universality Demo (Universal U Gate)")

	// 1. Replicating Hadamard (H)
	// H is equivalent to U(PI/2, 0, PI)
	p1 := core.NewProgram(1)
	p1.AddStep(core.NewStep(core.NewU(math.Pi/2, 0, math.Pi, 0)))

	// 2. Replicating Pauli-X (NOT)
	// X is equivalent to U(PI, 0, PI)
	p2 := core.NewProgram(1)
	p2.AddStep(core.NewStep(core.NewU(math.Pi, 0, math.Pi, 0)))

	engine := local.NewSimpleExecutionEnvironment()

	fmt.Println("\nUniversal U(PI/2, 0, PI) - Equivalent to Hadamard:")
	res1 := engine.RunProgram(p1)
	res1.PrintBinary()

	fmt.Println("\nUniversal U(PI, 0, PI) - Equivalent to Pauli-X:")
	res2 := engine.RunProgram(p2)
	res2.PrintBinary()

	// Output:
	// Universality Demo (Universal U Gate)
	//
	// Universal U(PI/2, 0, PI) - Equivalent to Hadamard:
	// Quantum Result (1 qubits):
	// |0>: 0.5000
	// |1>: 0.5000
	//
	// Universal U(PI, 0, PI) - Equivalent to Pauli-X:
	// Quantum Result (1 qubits):
	// |1>: 1.0000
}
