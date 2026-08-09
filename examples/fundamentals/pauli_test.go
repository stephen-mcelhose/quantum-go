package fundamentals_test

import (
	"fmt"
	"github.com/stephen-mcelhose/quantum-go/core"
	"github.com/stephen-mcelhose/quantum-go/local"
)

func Example_pauli() {
	fmt.Println("Pauli Gates Demo (X, Y, Z)")

	// Create a 3-qubit program
	p := core.NewProgram(3)

	// Step 1: Apply X to qubit 0, Y to qubit 1, Z to qubit 2
	step := core.NewStep(
		core.NewX(0),
		core.NewY(1),
		core.NewZ(2),
	)
	p.AddStep(step)

	// Execute the program
	engine := local.NewSimpleExecutionEnvironment()
	result := engine.RunProgram(p)

	// Output the result in binary
	result.PrintBinary()

	// Output:
	// Pauli Gates Demo (X, Y, Z)
	// Quantum Result (3 qubits):
	// |011>: 1.0000
}
