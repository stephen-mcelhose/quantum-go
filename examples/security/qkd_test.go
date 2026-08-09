package security_test

import (
	"fmt"
	"github.com/stephen-mcelhose/quantum-go/core"
	"github.com/stephen-mcelhose/quantum-go/local"
)

func Example_qkd() {
	fmt.Println("Quantum Key Distribution (BB84) Simulation")

	// Alice wants to send a secret bit to Bob.
	// For this example, we fix the values to make the output deterministic.
	aliceBit := 1
	aliceBasis := 0 // Z basis

	fmt.Printf("Alice's Secret Bit: %d, Basis: %s\n", aliceBit, basisName(aliceBasis))

	// Step 1: Alice prepares a qubit
	p := core.NewProgram(1)
	if aliceBit == 1 {
		p.AddStep(core.NewStep(core.NewX(0)))
	}
	if aliceBasis == 1 {
		p.AddStep(core.NewStep(core.NewHadamard(0)))
	}

	// Step 2: Bob chooses a measurement basis
	bobBasis := 0 // Bob chooses same basis
	fmt.Printf("Bob chooses Basis: %s\n", basisName(bobBasis))

	if bobBasis == 1 {
		p.AddStep(core.NewStep(core.NewHadamard(0)))
	}

	// Step 3: Bob measures the qubit
	engine := local.NewSimpleExecutionEnvironment()
	result := engine.RunProgram(p)
	bobBit := result.GetQubits()[0].Measure()

	fmt.Printf("Bob measures Bit: %d\n", bobBit)

	// Step 4: Alice and Bob compare bases (Classically)
	if aliceBasis == bobBasis {
		fmt.Println("Bases Match! Shared Key Bit Found.")
		if aliceBit == bobBit {
			fmt.Println("SUCCESS: Shared bit matches.")
		} else {
			fmt.Println("FAILURE: Shared bit mismatch.")
		}
	} else {
		fmt.Println("Bases Mismatch. Discard bit.")
	}

	// Output:
	// Quantum Key Distribution (BB84) Simulation
	// Alice's Secret Bit: 1, Basis: Z (|0>, |1>)
	// Bob chooses Basis: Z (|0>, |1>)
	// Bob measures Bit: 1
	// Bases Match! Shared Key Bit Found.
	// SUCCESS: Shared bit matches.
}

func basisName(b int) string {
	if b == 0 {
		return "Z (|0>, |1>)"
	}
	return "X (+, -)"
}
