package arithmetic_test

import (
	"fmt"
	"github.com/stephen-mcelhose/quantum-go/core"
	"github.com/stephen-mcelhose/quantum-go/local"
)

func Example_adder() {
	fmt.Println("Quantum Adder Demo (|x>|y> -> |x>|x+y>)")

	// We want to add 2 + 3 = 5
	// Register x: 2 (binary 10) -> 2 qubits
	// Register y: 3 (binary 11) -> 2 qubits
	
	p := core.NewProgram(5)

	// Step 1: Initialize x = 2 (|10>)
	// x0 is bit 0, x1 is bit 1. So 2 is q1=1, q0=0.
	p.AddStep(core.NewStep(core.NewX(1)))

	// Step 2: Initialize y = 3 (|011>)
	// y0 is q2, y1 is q3, y2 is q4. 3 is q3=1, q2=1, q4=0.
	p.AddStep(core.NewStep(core.NewX(2), core.NewX(3)))

	// Step 3: Apply Add gate
	// x0=0, x1=1, y0=2, y1=4
	// In Strange, the Add gate result is stored in the first register (x).
	// So x will become (x+y) mod 2^m.
	// Since x has 2 bits (m=2), it will be (2+3) mod 4 = 1.
	// Register y (3 qubits) remains unchanged (3).
	adder := core.NewAdd(0, 1, 2, 4)
	p.AddStep(core.NewStep(adder))

	// Execute
	engine := local.NewSimpleExecutionEnvironment()
	result := engine.RunProgram(p)

	// Output
	// Expected x=1 (q1=0, q0=1) and y=3 (q4=0, q3=1, q2=1)
	// |01101> (q4 q3 q2 q1 q0)
	result.PrintBinary()

	// Output:
	// Quantum Adder Demo (|x>|y> -> |x>|x+y>)
	// Quantum Result (5 qubits):
	// |01101>: 1.0000
}
