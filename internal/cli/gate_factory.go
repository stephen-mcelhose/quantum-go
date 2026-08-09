package cli

import (
	"sort"
)

// GateInfo provides metadata about a quantum gate.
type GateInfo struct {
	Alias       string
	Name        string
	Description string
	Example     string
	Category    string
}

var availableGates = []GateInfo{
	// Fundamental Gates
	{Alias: "h", Name: "Hadamard", Description: "Creates an equal superposition of |0> and |1>", Example: "h q[0]", Category: "Fundamental"},
	{Alias: "x", Name: "Pauli-X", Description: "Quantum NOT gate (bit flip)", Example: "x q[0]", Category: "Fundamental"},
	{Alias: "y", Name: "Pauli-Y", Description: "Bit and phase flip", Example: "y q[0]", Category: "Fundamental"},
	{Alias: "z", Name: "Pauli-Z", Description: "Phase flip", Example: "z q[0]", Category: "Fundamental"},
	{Alias: "id", Name: "Identity", Description: "No-op gate", Example: "id q[0]", Category: "Fundamental"},

	// Multi-Qubit Gates
	{Alias: "cx", Name: "CNOT", Description: "Controlled-NOT gate", Example: "cx q[0], q[1]", Category: "Multi-Qubit"},
	{Alias: "cz", Name: "Controlled-Z", Description: "Controlled-Phase flip", Example: "cz q[0], q[1]", Category: "Multi-Qubit"},
	{Alias: "swap", Name: "SWAP", Description: "Exchanges the states of two qubits", Example: "swap q[0], q[1]", Category: "Multi-Qubit"},
	{Alias: "ccx", Name: "Toffoli", Description: "Controlled-Controlled-NOT gate", Example: "ccx q[0], q[1], q[2]", Category: "Multi-Qubit"},
	{Alias: "cswap", Name: "Fredkin", Description: "Controlled-SWAP gate", Example: "cswap q[0], q[1], q[2]", Category: "Multi-Qubit"},

	// Rotations and Phase Gates
	{Alias: "rx(theta)", Name: "Rx", Description: "Rotation around X-axis", Example: "rx(1.57) q[0]", Category: "Rotation"},
	{Alias: "ry(theta)", Name: "Ry", Description: "Rotation around Y-axis", Example: "ry(1.57) q[0]", Category: "Rotation"},
	{Alias: "rz(theta)", Name: "Rz", Description: "Rotation around Z-axis", Example: "rz(1.57) q[0]", Category: "Rotation"},
	{Alias: "u1(theta)", Name: "PhaseShift", Description: "Phase rotation on |1> state", Example: "u1(1.57) q[0]", Category: "Rotation"},
	{Alias: "cu1(theta)", Name: "Controlled-Phase", Description: "Controlled phase rotation", Example: "cu1(1.57) q[0], q[1]", Category: "Rotation"},
	{Alias: "u(theta,phi,lambda)", Name: "Universal", Description: "Universal single-qubit rotation gate", Example: "u(1.57,0,0) q[0]", Category: "Rotation"},
	{Alias: "s", Name: "S gate", Description: "Square root of Z (pi/2 phase)", Example: "s q[0]", Category: "Phase"},
	{Alias: "sdg", Name: "S-dagger", Description: "Inverse of S gate", Example: "sdg q[0]", Category: "Phase"},
	{Alias: "t", Name: "T gate", Description: "Square root of S (pi/4 phase)", Example: "t q[0]", Category: "Phase"},
	{Alias: "tdg", Name: "T-dagger", Description: "Inverse of T gate", Example: "tdg q[0]", Category: "Phase"},
	{Alias: "sx", Name: "SX gate", Description: "Square root of X (V gate)", Example: "sx q[0]", Category: "Phase"},

	// Algorithmic Gates
	{Alias: "qft", Name: "QFT", Description: "Quantum Fourier Transform", Example: "Built-in only", Category: "Algorithmic"},
	{Alias: "invqft", Name: "Inverse QFT", Description: "Inverse Quantum Fourier Transform", Example: "Built-in only", Category: "Algorithmic"},
	{Alias: "add", Name: "Adder", Description: "Quantum Addition", Example: "Built-in only", Category: "Algorithmic"},
	{Alias: "mul", Name: "Multiplier", Description: "Quantum Multiplication", Example: "Built-in only", Category: "Algorithmic"},

	// Special
	{Alias: "measure", Name: "Measurement", Description: "Collapses qubit state into classical bit", Example: "measure q[0] -> c[0]", Category: "Special"},
}

// GetAvailableGates returns a sorted list of all available gates.
func GetAvailableGates() []GateInfo {
	res := make([]GateInfo, len(availableGates))
	copy(res, availableGates)
	sort.Slice(res, func(i, j int) bool {
		return res[i].Alias < res[j].Alias
	})
	return res
}
