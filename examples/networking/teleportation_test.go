package networking_test

import (
	"fmt"
	"github.com/stephen-mcelhose/quantum-go/core"
	"github.com/stephen-mcelhose/quantum-go/local"
)

func Example_teleportation() {
	fmt.Println("Quantum Teleportation Demo")

	// Qubits:
	// q0: State to teleport (|psi>)
	// q1: Alice's half of entangled pair
	// q2: Bob's half of entangled pair (Destination)
	p := core.NewProgram(3)

	// Step 1: Prepare the state to teleport (|psi> = H|0>)
	p.AddStep(core.NewStep(core.NewHadamard(0)))

	// Step 2: Create entanglement between q1 and q2 (Bell pair)
	p.AddStep(core.NewStep(core.NewHadamard(1)))
	p.AddStep(core.NewStep(core.NewCnot(1, 2)))

	// Step 3: Alice performs Bell measurement on q0 and q1
	p.AddStep(core.NewStep(core.NewCnot(0, 1)))
	p.AddStep(core.NewStep(core.NewHadamard(0)))

	// Step 4: Measurements (Alice measures q0 and q1)
	// In a real experiment, Alice sends these 2 classical bits to Bob.
	// Bob then applies corrections based on these bits.
	// In the simulator, we can use controlled gates to simulate the correction.
	
	// If q1 is 1, Bob applies X to q2
	p.AddStep(core.NewStep(core.NewCnot(1, 2)))
	
	// If q0 is 1, Bob applies Z to q2
	p.AddStep(core.NewStep(core.NewCz(0, 2)))

	// Execute the program
	engine := local.NewSimpleExecutionEnvironment()
	result := engine.RunProgram(p)

	// Verify the result
	// The state of q2 should now be the original state of q0 (H|0>).
	// So if we measure q2, it should be 0 or 1 with 50% probability.
	// However, the qubits q0 and q1 are also measured.
	// We use PrintBinary to see the full state.
	result.PrintBinary()

	// Output:
	// Quantum Teleportation Demo
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
