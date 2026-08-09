package algorithms_test

import (
	"fmt"
	"github.com/stephen-mcelhose/quantum-go/core"
	"github.com/stephen-mcelhose/quantum-go/local"
	"github.com/stephen-mcelhose/quantum-go/math"
)

func Example_grover() {
	fmt.Println("Grover's Search Algorithm Demo (2 qubits)")
	fmt.Println("Searching for state |11> (index 3)")

	// 2 qubits can represent 4 states. Grover's finds the answer in ~1 iteration.
	p := core.NewProgram(2)

	// Step 1: Initialize superposition
	p.AddStep(core.NewStep(core.NewHadamard(0), core.NewHadamard(1)))

	// Step 2: Oracle for |11>
	// Matrix is diag(1, 1, 1, -1)
	oracleMatrix := math.NewMatrix(4, 4)
	oracleMatrix.Set(0, 0, 1)
	oracleMatrix.Set(1, 1, 1)
	oracleMatrix.Set(2, 2, 1)
	oracleMatrix.Set(3, 3, -1)
	oracle := core.NewOracle(0, oracleMatrix)
	p.AddStep(core.NewStep(oracle))

	// Step 3: Diffusion Operator
	diffMatrix := math.NewMatrix(4, 4)
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if i == j {
				diffMatrix.Set(i, j, -0.5)
			} else {
				diffMatrix.Set(i, j, 0.5)
			}
		}
	}
	diffusion := core.NewOracle(0, diffMatrix)
	p.AddStep(core.NewStep(diffusion))

	// Execute
	engine := local.NewSimpleExecutionEnvironment()
	result := engine.RunProgram(p)

	// Output
	// Should see |11> with probability 1.0
	result.PrintBinary()

	// Output:
	// Grover's Search Algorithm Demo (2 qubits)
	// Searching for state |11> (index 3)
	// Quantum Result (2 qubits):
	// |11>: 1.0000
}
