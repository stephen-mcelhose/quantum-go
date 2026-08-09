package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List built-in circuits or supported gates",
	Long:  `Display a list of all built-in quantum circuits or supported quantum gates.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			fmt.Fprintf(os.Stderr, "Error showing help: %v\n", err)
		}
	},
}

var listCircuitsCmd = &cobra.Command{
	Use:   "circuits",
	Short: "List all built-in quantum circuits",
	Long:  `Display a list of all built-in quantum circuits available in the simulator with their descriptions.`,
	Run: func(cmd *cobra.Command, args []string) {
		circuits := GetBuiltInCircuits()

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tDESCRIPTION")
		for _, c := range circuits {
			fmt.Fprintf(w, "%s\t%s\n", c.Name, c.Description)
		}
		w.Flush()
	},
}

var listGatesCmd = &cobra.Command{
	Use:   "gates",
	Short: "List all supported quantum gates",
	Long:  `Display a list of all quantum gates supported by the simulator with their descriptions and examples.`,
	Run: func(cmd *cobra.Command, args []string) {
		gates := GetAvailableGates()

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ALIAS\tNAME\tCATEGORY\tDESCRIPTION\tEXAMPLE")
		for _, g := range gates {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", g.Alias, g.Name, g.Category, g.Description, g.Example)
		}
		w.Flush()
	},
}

func init() {
	listCmd.AddCommand(listCircuitsCmd)
	listCmd.AddCommand(listGatesCmd)
	rootCmd.AddCommand(listCmd)
}
