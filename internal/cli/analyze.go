package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/stephen-mcelhose/quantum-go/core"
	"github.com/stephen-mcelhose/quantum-go/local"
	"github.com/stephen-mcelhose/quantum-go/math"
	"github.com/stephen-mcelhose/quantum-go/qasm"
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze [qasm-file]",
	Short: "Analyze thermodynamic properties of a quantum circuit",
	Long: `Analyze the thermodynamic properties (Entropy, Energy) of a quantum circuit.

This command treats the quantum system as a physical ensemble and calculates:

1. Density Matrix: Converts the state vector |ψ⟩ to ρ = |ψ⟩⟨ψ|.
2. Subsystem Entropy: Calculates von Neumann entropy S(ρ_i) = -Tr(ρ_i ln ρ_i)
   for each qubit by tracing out the environment. S > 0 indicates entanglement.
3. Internal Energy: Measures the expectation value U = Tr(ρ H) for a
   standard Z-basis Hamiltonian (H = σz).

You can analyze a QASM file or build a circuit using the --step flag.

Available Gates:
| Gate      | Description            | Example                |
| :---      | :---                   | :---                   |
| h         | Hadamard               | h q[0]                 |
| x, y, z   | Pauli gates            | x q[0]                 |
| id        | Identity               | id q[0]                |
| cx        | CNOT                   | cx q[0], q[1]          |
| cz        | Controlled-Z           | cz q[0], q[1]          |
| swap      | SWAP qubits            | swap q[0], q[1]        |
| ccx       | Toffoli (CCNOT)        | ccx q[0], q[1], q[2]   |
| u1(θ)     | Phase Rotation         | u1(1.57) q[0]          |
| cu1(θ)    | Controlled Rotation    | cu1(1.57) q[0], q[1]   |
| measure   | Measurement            | measure q[0]           |`,
	Example: `  # Analyze built-in Bell state
  quantum-go analyze --circuit bell

  # Analyze a 2-qubit Bell state built via flags
  quantum-go analyze -n 2 -s "h q[0]" -s "cx q[0], q[1]"

  # Analyze a 3-qubit GHZ state
  quantum-go analyze -n 3 -s "h q[0]" -s "cx q[0], q[1]" -s "cx q[1], q[2]"`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var p *core.Program
		var err error

		if circuitName != "" {
			p, err = GetBuiltInProgram(circuitName, numQubits, circuitParams)
			if err != nil {
				return err
			}
		} else if len(circuitSteps) > 0 {
			p, err = ProgramFromSteps(circuitSteps, numQubits)
			if err != nil {
				return err
			}
		} else {
			if len(args) == 0 {
				return fmt.Errorf("qasm-file argument, --circuit, or --step is required")
			}
			data, err := os.ReadFile(args[0])
			if err != nil {
				return fmt.Errorf("failed to read file: %w", err)
			}

			p, err = qasm.Parse(string(data))
			if err != nil {
				return fmt.Errorf("failed to parse QASM: %w", err)
			}
		}

		env := local.NewSimpleExecutionEnvironment()

		result := env.RunProgram(p)

		state := result.GetProbability()
		rho := math.ToDensityMatrix(state)

		fmt.Printf("Circuit Analysis (%d qubits):\n", p.NumQubits)

		// Entropy for subsystems (if qubits > 1)
		if p.NumQubits > 1 {
			for i := 0; i < p.NumQubits; i++ {
				// Trace out all but qubit i
				traceOut := []int{}
				for j := 0; j < p.NumQubits; j++ {
					if i != j {
						traceOut = append(traceOut, j)
					}
				}
				rhoSub := math.PartialTrace(rho, p.NumQubits, traceOut)
				entropy := math.VonNeumannEntropy(rhoSub)
				fmt.Printf("  Qubit %d Entropy: %.6f\n", i, entropy)
			}
		} else {
			entropy := math.VonNeumannEntropy(rho)
			fmt.Printf("  System Entropy: %.6f\n", entropy)
		}

		// Energy (Expectation Value of Z-Hamiltonian)
		hz := math.NewMatrix(2, 2)
		hz.Data = []complex128{1, 0, 0, -1}

		var totalEnergy float64
		for i := 0; i < p.NumQubits; i++ {
			traceOut := []int{}
			for j := 0; j < p.NumQubits; j++ {
				if i != j {
					traceOut = append(traceOut, j)
				}
			}
			rhoSub := math.PartialTrace(rho, p.NumQubits, traceOut)
			ev := real(math.ExpectationValue(rhoSub, hz))
			totalEnergy += ev
		}

		fmt.Printf("  Total Energy (Z-basis): %.6f\n", totalEnergy)

		return nil
	},
}

func init() {
	analyzeCmd.Flags().StringVarP(&circuitName, "circuit", "c", "", "Built-in circuit to analyze (use 'quantum-go list circuits' for all options)")
	analyzeCmd.Flags().StringToStringVarP(&circuitParams, "param", "p", nil, "Parameters for built-in circuits")
	analyzeCmd.Flags().IntVarP(&numQubits, "qubits", "n", 2, "Number of qubits for the circuit")
	analyzeCmd.Flags().StringArrayVarP(&circuitSteps, "step", "s", []string{}, "Add a quantum gate to the circuit (e.g. \"h q[0]\")")
	rootCmd.AddCommand(analyzeCmd)
}
