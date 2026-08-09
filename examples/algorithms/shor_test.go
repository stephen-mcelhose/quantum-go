package algorithms_test

import (
	"fmt"
	"github.com/stephen-mcelhose/quantum-go/core"
	"github.com/stephen-mcelhose/quantum-go/local"
)

func Example_shor() {
	fmt.Println("Shor's Algorithm - Period Finding Demo")
	fmt.Println("Finding the period of 2^x mod 7")

	// a=2, mod=7
	// We need length = ceil(log2(7)) = 3
	// We need offset = 3 (for the precision register)
	// Total qubits = 2*length + 1 + offset = 2*3 + 1 + 3 = 10
	
	length := 3
	offset := 3
	p := core.NewProgram(2*length + 1 + offset)

	// Step 1: Initialize precision register to superposition
	for i := 0; i < offset; i++ {
		p.AddStep(core.NewStep(core.NewHadamard(i)))
	}

	// Step 2: Initialize result register to 1
	p.AddStep(core.NewStep(core.NewX(offset)))

	// Step 3: Modular Exponentiation
	a := 2
	mod := 7
	for i := length - 1; i >= 0; i-- {
		// Calculate m = a^(2^i) mod mod
		m := 1
		for j := 0; j < (1 << i); j++ {
			m = (m * a) % mod
		}
		
		fmt.Printf("Adding MulModulus for 2^(2^%d) mod 7 = %d\n", i, m)
		mul := core.NewMulModulus(offset, offset+length-1, m, mod)
		cbg := core.NewControlledBlockGate(mul, i)
		p.AddStep(core.NewStep(cbg))
	}

	// Step 4: Inverse QFT on precision register
	invQFT := core.NewFourier(offset, 0)
	invQFT.SetInverse(true)
	p.AddStep(core.NewStep(invQFT))

	// Execute
	engine := local.NewSimpleExecutionEnvironment()
	fmt.Println("Running simulation... (This may take a few seconds)")
	result := engine.RunProgram(p)

	// Output
	fmt.Println("\nMeasurement Results (Binary):")
	// Only print states with probability > 0.05 for brevity in Example
	for i, amp := range result.GetProbability() {
		prob := real(amp)*real(amp) + imag(amp)*imag(amp)
		if prob > 0.05 {
			fmt.Printf("|%0*b>: %.4f\n", result.GetNumQubits(), i, prob)
		}
	}
	
	fmt.Println("\nThe peaks in the precision register (first 3 bits) correspond to s/r where r is the period.")
	fmt.Println("For 2^x mod 7, the period r=3. Peaks should be at 0, 1/3, 2/3 of 2^offset.")

	// Output:
	// Shor's Algorithm - Period Finding Demo
	// Finding the period of 2^x mod 7
	// Adding MulModulus for 2^(2^2) mod 7 = 2
	// Adding MulModulus for 2^(2^1) mod 7 = 4
	// Adding MulModulus for 2^(2^0) mod 7 = 2
	// Running simulation... (This may take a few seconds)
	//
	// Measurement Results (Binary):
	// |0000001000>: 0.1153
	// |0000001010>: 0.0981
	// |0000001110>: 0.0731
	//
	// The peaks in the precision register (first 3 bits) correspond to s/r where r is the period.
	// For 2^x mod 7, the period r=3. Peaks should be at 0, 1/3, 2/3 of 2^offset.
}
