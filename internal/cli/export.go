package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stephen-mcelhose/quantum-go/core"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export a circuit to OpenQASM 2.0",
	Long: `Export a quantum circuit to OpenQASM 2.0.
You can export built-in circuits or build your own using the --step flag.

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
	Example: `  # Export built-in Bell state
  strange export --circuit bell

  # Build custom GHZ state using manual steps
  strange export -n 3 -s "h q[0]" -s "cx q[0], q[1]" -s "cx q[1], q[2]"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var p *core.Program
		var err error

		if len(circuitSteps) > 0 {
			p, err = ProgramFromSteps(circuitSteps, numQubits)
			if err != nil {
				return err
			}
			} else {
				p, err = GetBuiltInProgram(circuitName, numQubits, circuitParams)
				if err != nil {

				return err
			}
		}

		fmt.Print(p.ToQASM())
		return nil
	},
}

func init() {
	exportCmd.Flags().StringVarP(&circuitName, "circuit", "c", "bell", "Built-in circuit to export (use 'strange list circuits' for all options)")
	exportCmd.Flags().StringToStringVarP(&circuitParams, "param", "p", nil, "Parameters for built-in circuits")
	exportCmd.Flags().IntVarP(&numQubits, "qubits", "n", 2, "Number of qubits for the circuit")
	exportCmd.Flags().StringArrayVarP(&circuitSteps, "step", "s", []string{}, "Add a quantum gate to the circuit (e.g. \"h q[0]\")")
	rootCmd.AddCommand(exportCmd)
}
