package cli

import (
	"fmt"
	"strings"

	"github.com/stephen-mcelhose/quantum-go/core"
	"github.com/stephen-mcelhose/quantum-go/qasm"
)

// ProgramFromSteps builds a core.Program from a slice of QASM strings.
// If the steps do not contain a qreg declaration, one is added using numQubits.
func ProgramFromSteps(steps []string, numQubits int) (*core.Program, error) {
	hasQreg := false
	for _, s := range steps {
		if strings.Contains(s, "qreg") {
			hasQreg = true
			break
		}
	}

	var qasmInput strings.Builder
	qasmInput.WriteString("OPENQASM 2.0;\n")
	qasmInput.WriteString("include \"qelib1.inc\";\n")
	if !hasQreg {
		qasmInput.WriteString(fmt.Sprintf("qreg q[%d];\n", numQubits))
	}
	for _, s := range steps {
		s = strings.TrimSpace(s)
		qasmInput.WriteString(s)
		if !strings.HasSuffix(s, ";") {
			qasmInput.WriteString(";")
		}
		qasmInput.WriteString("\n")
	}

	return qasm.Parse(qasmInput.String())
}
