package fundamentals_test

import (
	"fmt"
	"github.com/stephen-mcelhose/quantum-go/core"
	"github.com/stephen-mcelhose/quantum-go/local"
)

func Example_hadamard() {
	fmt.Println("Hadamard Gate Demo (Superposition)")

	// Create a 1-qubit program
	p := core.NewProgram(1)

	// Step 1: Apply Hadamard to qubit 0
	// This puts the qubit into a state where it is 50% |0> and 50% |1>
	p.AddStep(core.NewStep(core.NewHadamard(0)))

	// Execute the program
	engine := local.NewSimpleExecutionEnvironment()
	result := engine.RunProgram(p)

	// Output the result in binary
	result.PrintBinary()

	// Output:
	// Hadamard Gate Demo (Superposition)
	// Quantum Result (1 qubits):
	// |0>: 0.5000
	// |1>: 0.5000
}
