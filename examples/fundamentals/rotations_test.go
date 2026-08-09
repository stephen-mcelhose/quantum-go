package fundamentals_test

import (
	"fmt"
	"math"

	"github.com/stephen-mcelhose/quantum-go/core"
	"github.com/stephen-mcelhose/quantum-go/local"
)

// Example_rotations demonstrates the use of parameterized rotation gates (Rx, Ry, Rz).
// These gates allow for arbitrary rotations around the X, Y, and Z axes of the Bloch Sphere.
func Example_rotations() {
	fmt.Println("Rotation Gates Demo (Rx, Ry, Rz)")

	// 1. Rx(PI) - Equivalent to a bit flip (Pauli-X) up to a global phase
	p1 := core.NewProgram(1)
	p1.AddStep(core.NewStep(core.NewRx(math.Pi, 0)))
	
	// 2. Ry(PI/2) - Creates a superposition state (|0> + |1>)/sqrt(2)
	// Similar to Hadamard, but with a different phase alignment on the Bloch Sphere equator.
	p2 := core.NewProgram(1)
	p2.AddStep(core.NewStep(core.NewRy(math.Pi/2, 0)))

	// 3. Rz(PI/4) - Applies a phase rotation to the |1> state.
	// We first put the qubit in superposition so the phase effect is visible.
	p3 := core.NewProgram(1)
	p3.AddStep(core.NewStep(core.NewHadamard(0)))
	p3.AddStep(core.NewStep(core.NewRz(math.Pi/4, 0)))

	engine := local.NewSimpleExecutionEnvironment()

	fmt.Println("\nRx(PI) Result (Qubit starts at |0>):")
	res1 := engine.RunProgram(p1)
	res1.PrintBinary()

	fmt.Println("\nRy(PI/2) Result (Qubit starts at |0>):")
	res2 := engine.RunProgram(p2)
	res2.PrintBinary()

	fmt.Println("\nH + Rz(PI/4) Result:")
	res3 := engine.RunProgram(p3)
	// Phase doesn't affect the measurement probability, so we'll see 50/50.
	res3.PrintBinary()

	// Output:
	// Rotation Gates Demo (Rx, Ry, Rz)
	//
	// Rx(PI) Result (Qubit starts at |0>):
	// Quantum Result (1 qubits):
	// |1>: 1.0000
	//
	// Ry(PI/2) Result (Qubit starts at |0>):
	// Quantum Result (1 qubits):
	// |0>: 0.5000
	// |1>: 0.5000
	//
	// H + Rz(PI/4) Result:
	// Quantum Result (1 qubits):
	// |0>: 0.5000
	// |1>: 0.5000
}
