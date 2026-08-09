package entanglement_test

import (
	"fmt"
	"github.com/stephen-mcelhose/quantum-go/core"
	"github.com/stephen-mcelhose/quantum-go/local"
)

func Example_bellState() {
	fmt.Println("Bell State entanglement (|00> + |11>)")

	// Create a 2-qubit program
	p := core.NewProgram(2)

	// Step 1: Apply Hadamard to qubit 0
	p.AddStep(core.NewStep(core.NewHadamard(0)))

	// Step 2: Apply CNOT with qubit 0 as control and qubit 1 as target
	p.AddStep(core.NewStep(core.NewCnot(0, 1)))

	// Execute the program
	engine := local.NewSimpleExecutionEnvironment()
	result := engine.RunProgram(p)

	// Output the result in binary
	// Should see |00> and |11> each with 0.5 probability
	result.PrintBinary()

	// Output:
	// Bell State entanglement (|00> + |11>)
	// Quantum Result (2 qubits):
	// |00>: 0.5000
	// |11>: 0.5000
}
