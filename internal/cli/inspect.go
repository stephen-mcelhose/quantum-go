package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/stephen-mcelhose/quantum-go/core"
	"github.com/stephen-mcelhose/quantum-go/qasm"
)

var inspectCmd = &cobra.Command{
	Use:   "inspect [qasm-file]",
	Short: "Inspect circuit metadata",
	Long: `Inspect circuit metadata (qubits, steps, gate count).
You can inspect a QASM file or build a circuit using the --step flag.

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
	Example: `  # Inspect built-in Bell state
  strange inspect --circuit bell

  # Inspect a 3-qubit GHZ state built via flags
  strange inspect -n 3 -s "h q[0]" -s "cx q[0], q[1]" -s "cx q[1], q[2]"`,
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

		fmt.Printf("Circuit Metadata:\n")

		fmt.Printf("  Qubits: %d\n", p.NumQubits)
		fmt.Printf("  Steps:  %d\n", len(p.Steps))
		
		gateCount := 0
		gateTypes := make(map[string]int)
		for _, step := range p.Steps {
			gateCount += len(step.Gates)
			for _, g := range step.Gates {
				gateTypes[g.GetName()]++
			}
		}
		
		fmt.Printf("  Total Gates: %d\n", gateCount)
		fmt.Printf("  Gate Inventory:\n")
		for name, count := range gateTypes {
			fmt.Printf("    - %s: %d\n", name, count)
		}

		return nil
	},
}

func init() {
	inspectCmd.Flags().StringVarP(&circuitName, "circuit", "c", "", "Built-in circuit to inspect (use 'strange list circuits' for all options)")
	inspectCmd.Flags().StringToStringVarP(&circuitParams, "param", "p", nil, "Parameters for built-in circuits")
	inspectCmd.Flags().IntVarP(&numQubits, "qubits", "n", 2, "Number of qubits for the circuit")
	inspectCmd.Flags().StringArrayVarP(&circuitSteps, "step", "s", []string{}, "Add a quantum gate to the circuit (e.g. \"h q[0]\")")
	rootCmd.AddCommand(inspectCmd)
}
