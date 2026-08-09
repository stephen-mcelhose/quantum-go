package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	numQubits     int
	circuitSteps  []string
	circuitName   string
	circuitParams map[string]string
)

var rootCmd = &cobra.Command{
	Use:   "quantum-go",
	Short: "quantum-go is a quantum circuit simulator CLI",
	Long: `A high-performance quantum circuit simulator written in Go.
Go port of the original Java Strange simulator by Johan Vos (https://github.com/redfx-quantum/strange).

For a full list of supported gates, run:
  quantum-go list gates

For a list of built-in circuits, run:
  quantum-go list circuits`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			fmt.Fprintf(os.Stderr, "Error showing help: %v\n", err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Root flags can be defined here
}
