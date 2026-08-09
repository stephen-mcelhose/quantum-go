package entanglement_test

import (
	"fmt"
	"github.com/stephen-mcelhose/quantum-go/core"
	"github.com/stephen-mcelhose/quantum-go/local"
)

func Example_ghzState() {
	fmt.Println("GHZ State entanglement (|000> + |111>)")

	// Create a 3-qubit program
	p := core.NewProgram(3)

	// Step 1: Apply Hadamard to qubit 0
	p.AddStep(core.NewStep(core.NewHadamard(0)))

	// Step 2: Apply CNOT(0, 1)
	p.AddStep(core.NewStep(core.NewCnot(0, 1)))

	// Step 3: Apply CNOT(1, 2)
	p.AddStep(core.NewStep(core.NewCnot(1, 2)))

	// Execute the program
	engine := local.NewSimpleExecutionEnvironment()
	result := engine.RunProgram(p)

	// Output the result in binary
	// Should see |000> and |111> each with 0.5 probability
	result.PrintBinary()

	// Output:
	// GHZ State entanglement (|000> + |111>)
	// Quantum Result (3 qubits):
	// |000>: 0.5000
	// |111>: 0.5000
}
