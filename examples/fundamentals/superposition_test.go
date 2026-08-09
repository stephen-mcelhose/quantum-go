package fundamentals_test

import (
	"fmt"
	"github.com/stephen-mcelhose/quantum-go/core"
	"github.com/stephen-mcelhose/quantum-go/local"
)

func Example_superposition() {
	fmt.Println("Superposition Demo (3 qubits)")

	// Create a 3-qubit program
	p := core.NewProgram(3)

	// Step 1: Apply Hadamard to all 3 qubits
	// This creates an equal superposition of all 8 possible states (|000> to |111>)
	p.AddStep(core.NewStep(
		core.NewHadamard(0),
		core.NewHadamard(1),
		core.NewHadamard(2),
	))

	// Execute the program
	engine := local.NewSimpleExecutionEnvironment()
	result := engine.RunProgram(p)

	// Output the result in binary
	// Each of the 8 states should have a probability of 1/8 = 0.125
	result.PrintBinary()

	// Output:
	// Superposition Demo (3 qubits)
	// Quantum Result (3 qubits):
	// |000>: 0.1250
	// |001>: 0.1250
	// |010>: 0.1250
	// |011>: 0.1250
	// |100>: 0.1250
	// |101>: 0.1250
	// |110>: 0.1250
	// |111>: 0.1250
}
