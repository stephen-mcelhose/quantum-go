package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/stephen-mcelhose/quantum-go/core"
	"github.com/stephen-mcelhose/quantum-go/local"
	"github.com/stephen-mcelhose/quantum-go/qasm"
)

var outputJSON bool

var runCmd = &cobra.Command{
	Use:   "run [qasm-file]",
	Short: "Run a quantum circuit",
	Long: `Run a quantum circuit from a QASM file or build it using the --step flag.

Use 'quantum-go list circuits' to see all available built-in circuits.
Use 'quantum-go list gates' to see all supported gates.`,
		Example: `  # List all built-in circuits
  quantum-go list circuits

  # List all supported gates
  quantum-go list gates

  # Run built-in Bell state
  quantum-go run --circuit bell

  # Run built-in Bernstein-Vazirani with custom hidden string
  quantum-go run --circuit bernstein-vazirani -p s=101

  # Run built-in Shor's algorithm for 15
  quantum-go run --circuit shor -p mod=15,a=7,precision=8`,
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

		if outputJSON {
			ref := result.GetStateVectorReference()
			jsonData, err := json.MarshalIndent(ref, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal result to JSON: %w", err)
			}
			fmt.Println(string(jsonData))
		} else {
			result.PrintBinary()
		}
		return nil
	},
}

func init() {
	runCmd.Flags().BoolVar(&outputJSON, "json", false, "Output result as JSON")
	runCmd.Flags().StringVarP(&circuitName, "circuit", "c", "", "Built-in circuit to run (use 'quantum-go list circuits' for all options)")
	runCmd.Flags().StringToStringVarP(&circuitParams, "param", "p", nil, "Parameters for built-in circuits (e.g. s=101)")
	runCmd.Flags().IntVarP(&numQubits, "qubits", "n", 2, "Number of qubits for the circuit")
	runCmd.Flags().StringArrayVarP(&circuitSteps, "step", "s", []string{}, "Add a quantum gate to the circuit (e.g. \"h q[0]\")")
	rootCmd.AddCommand(runCmd)
}
