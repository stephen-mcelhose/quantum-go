package core

import (
	"fmt"
	"strings"
)

// ToQASM converts the quantum program to an OpenQASM 2.0 string representation.
// This is useful for circuit visualization and verification against other simulators.
func (p *Program) ToQASM() string {
	var sb strings.Builder
	sb.WriteString("OPENQASM 2.0;\n")
	sb.WriteString("include \"qelib1.inc\";\n")
	sb.WriteString(fmt.Sprintf("qreg q[%d];\n", p.NumQubits))
	sb.WriteString(fmt.Sprintf("creg c[%d];\n", p.NumQubits))

	for _, step := range p.Steps {
		if step.Type == StepNormal {
			for _, gate := range step.Gates {
				sb.WriteString(gateToQASM(gate))
			}
		}
	}

	return sb.String()
}

func gateToQASM(g Gate) string {
	indices := g.GetAffectedQubitIndexes()
	
	// Handle BlockGates by decomposing them
	if bg, ok := g.(BlockGateInterface); ok {
		var sb strings.Builder
		block := bg.GetBlock()
		for _, step := range block.Steps {
			for _, subGate := range step.Gates {
				sb.WriteString(gateToQASM(subGate))
			}
		}
		return sb.String()
	}

	switch gate := g.(type) {
		case *Hadamard:
			return fmt.Sprintf("h q[%d];\n", indices[0])
		case *X:
			return fmt.Sprintf("x q[%d];\n", indices[0])
		case *Y:
			return fmt.Sprintf("y q[%d];\n", indices[0])
		case *Z:
			return fmt.Sprintf("z q[%d];\n", indices[0])
		case *S:
			if gate.IsInverse() {
				return fmt.Sprintf("sdg q[%d];\n", indices[0])
			}
			return fmt.Sprintf("s q[%d];\n", indices[0])
		case *T:
			if gate.IsInverse() {
				return fmt.Sprintf("tdg q[%d];\n", indices[0])
			}
			return fmt.Sprintf("t q[%d];\n", indices[0])
		case *V:
			return fmt.Sprintf("sx q[%d];\n", indices[0])
		case *Rx:
			return fmt.Sprintf("rx(%.6f) q[%d];\n", gate.Theta, indices[0])
		case *Ry:
			return fmt.Sprintf("ry(%.6f) q[%d];\n", gate.Theta, indices[0])
		case *Rz:
			return fmt.Sprintf("rz(%.6f) q[%d];\n", gate.Theta, indices[0])
		case *U:
			return fmt.Sprintf("u3(%.6f,%.6f,%.6f) q[%d];\n", gate.Theta, gate.Phi, gate.Lambda, indices[0])
		case *Cnot:
			return fmt.Sprintf("cx q[%d], q[%d];\n", indices[0], indices[1])
		case *Cz:
			return fmt.Sprintf("cz q[%d], q[%d];\n", indices[0], indices[1])
		case *Swap:
			return fmt.Sprintf("swap q[%d], q[%d];\n", indices[0], indices[1])
		case *Toffoli:
			return fmt.Sprintf("ccx q[%d], q[%d], q[%d];\n", indices[0], indices[1], indices[2])
		case *PhaseShift:
			return fmt.Sprintf("u1(%.6f) q[%d];\n", gate.Theta, indices[0])
		case *Cr:
			return fmt.Sprintf("cu1(%.6f) q[%d], q[%d];\n", gate.Theta, indices[0], indices[1])

	case *Identity:
		return fmt.Sprintf("id q[%d];\n", indices[0])
	case *Measurement:
		return fmt.Sprintf("measure q[%d] -> c[%d];\n", indices[0], indices[0])
	default:
		// Fallback for unknown gates: comment them out or use caption
		return fmt.Sprintf("// gate %s q%v\n", g.GetName(), indices)
	}
}
