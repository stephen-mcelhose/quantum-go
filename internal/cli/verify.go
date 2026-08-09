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

var (
	verifyMode      string
	referenceFile   string
	verifyTolerance float64
)

var verifyCmd = &cobra.Command{
	Use:   "verify [qasm-file]",
	Short: "Verify a quantum circuit calculation against a reference",
	Long: `Verify the correctness of a quantum circuit calculation.
Supported modes:
  - theoretical: Compare built-in circuits (bell, ghz, qft) against math constants.
  - qiskit: Export to QASM and verify using local Qiskit Aer (requires python3 and qiskit_verify.py).
  - file: Compare against a JSON reference file.

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
	Example: `  # Verify built-in Bell state against math constants
  quantum-go verify --circuit bell --mode theoretical

  # Verify built-in 3-qubit GHZ state
  quantum-go verify --circuit ghz -n 3 --mode theoretical

  # Verify built-in 4-qubit QFT
  quantum-go verify --circuit qft -n 4 --mode theoretical

  # Verify custom circuit against Qiskit Aer
  quantum-go verify -n 2 -s "h q[0]" -s "u1(0.785398) q[0]" --mode qiskit`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var p *core.Program
		var err error
		var expected []complex128

		env := local.NewSimpleExecutionEnvironment()

		if len(circuitSteps) > 0 {
			p, err = ProgramFromSteps(circuitSteps, numQubits)
			if err != nil {
				return err
			}
		}

		switch verifyMode {
		case "theoretical":
			if p == nil {
				if circuitName == "" {
					return fmt.Errorf("--circuit is required for theoretical mode")
				}
					p, err = GetBuiltInProgram(circuitName, numQubits, circuitParams)
					if err != nil {
						return err
					}
				}
				expected, err = core.GetExpectedStateVector(circuitName, p.NumQubits)
				if err != nil {
					return err
				}

			case "qiskit":
				if p == nil {
					if len(args) > 0 {
						data, err := os.ReadFile(args[0])
						if err != nil {
							return err
						}
						p, err = qasm.Parse(string(data))
						if err != nil {
							return err
						}
					} else if circuitName != "" {
						p, err = GetBuiltInProgram(circuitName, numQubits, circuitParams)
						if err != nil {
							return err
						}

				} else {
					return fmt.Errorf("either a QASM file, --circuit, or --step must be provided")
				}
			}
			fmt.Println("Running qiskit-verify via Python bridge...")
			expected, err = runQiskitVerification(p.ToQASM())
			if err != nil {
				return err
			}

		case "file":
			if p == nil {
				if len(args) == 0 {
					return fmt.Errorf("QASM file argument or --step is required for file mode")
				}
				data, err := os.ReadFile(args[0])
				if err != nil {
					return err
				}
				p, err = qasm.Parse(string(data))
				if err != nil {
					return err
				}
			}

			if referenceFile == "" {
				return fmt.Errorf("--reference is required for file mode")
			}

			refData, err := os.ReadFile(referenceFile)
			if err != nil {
				return err
			}
			var ref core.StateVectorReference
			if err := json.Unmarshal(refData, &ref); err != nil {
				return fmt.Errorf("failed to parse reference file: %w", err)
			}
			if ref.NumQubits != p.NumQubits {
				return fmt.Errorf("qubit mismatch: circuit has %d, reference has %d", p.NumQubits, ref.NumQubits)
			}
			expected = make([]complex128, len(ref.Amplitudes))
			for i, a := range ref.Amplitudes {
				expected[i] = complex(a.Re, a.Im)
			}

		default:
			return fmt.Errorf("unknown verification mode: %s", verifyMode)
		}

		// Execute locally
		result := env.RunProgram(p)
		got := result.GetProbability()

		// Compare
		if err := core.CompareStateVectors(got, expected, verifyTolerance); err != nil {
			fmt.Printf("❌ Verification FAILED\n")
			return err
		}

		fmt.Printf("✅ Verification SUCCESSFUL (mode: %s, qubits: %d, tolerance: %g)\n", verifyMode, p.NumQubits, verifyTolerance)
		return nil
	},
}

func init() {
	verifyCmd.Flags().StringVarP(&verifyMode, "mode", "m", "theoretical", "Verification mode (theoretical, qiskit, file)")
	verifyCmd.Flags().StringVarP(&referenceFile, "reference", "r", "", "Reference JSON file (for mode=file)")
	verifyCmd.Flags().StringVarP(&circuitName, "circuit", "c", "", "Built-in circuit to verify (use 'quantum-go list circuits' for all options)")
	verifyCmd.Flags().StringToStringVarP(&circuitParams, "param", "p", nil, "Parameters for built-in circuits")
	verifyCmd.Flags().IntVarP(&numQubits, "qubits", "n", 2, "Number of qubits for the circuit")
	verifyCmd.Flags().StringArrayVarP(&circuitSteps, "step", "s", []string{}, "Add a quantum gate to the circuit (e.g. \"h q[0]\")")
	verifyCmd.Flags().Float64VarP(&verifyTolerance, "tolerance", "t", 1e-6, "Numerical tolerance for comparison")

	rootCmd.AddCommand(verifyCmd)
}
