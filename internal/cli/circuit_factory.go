package cli

import (
	"fmt"
	"sort"
	"strconv"
	"github.com/stephen-mcelhose/quantum-go/core"
)

// CircuitInfo provides metadata about a built-in circuit.
type CircuitInfo struct {
	Name        string
	Description string
	Factory     func(qubits int, params map[string]string) (*core.Program, error)
}

var builtInCircuits = map[string]CircuitInfo{
	"bell": {
		Name:        "Bell State",
		Description: "2-qubit entanglement (|00> + |11>) / sqrt(2)",
		Factory: func(qubits int, params map[string]string) (*core.Program, error) {
			return core.NewBellProgram(), nil
		},
	},
	"ghz": {
		Name:        "GHZ State",
		Description: "n-qubit entanglement (|0...0> + |1...1>) / sqrt(2)",
		Factory: func(qubits int, params map[string]string) (*core.Program, error) {
			return core.NewGHZProgram(qubits), nil
		},
	},
	"qft": {
		Name:        "Quantum Fourier Transform",
		Description: "Performs QFT on n qubits",
		Factory: func(qubits int, params map[string]string) (*core.Program, error) {
			return core.NewQFTProgram(qubits), nil
		},
	},
	"grover": {
		Name:        "Grover's Search",
		Description: "2-qubit search for state |11>",
		Factory: func(qubits int, params map[string]string) (*core.Program, error) {
			return core.NewGroverProgram(), nil
		},
	},
	"teleportation": {
		Name:        "Quantum Teleportation",
		Description: "3-qubit teleportation protocol",
		Factory: func(qubits int, params map[string]string) (*core.Program, error) {
			return core.NewTeleportationProgram(), nil
		},
	},
	"superdense": {
		Name:        "Superdense Coding",
		Description: "2-qubit protocol to send two classical bits via one qubit",
		Factory: func(qubits int, params map[string]string) (*core.Program, error) {
			return core.NewSuperdenseCodingProgram(), nil
		},
	},
	"toffoli": {
		Name:        "Toffoli Gate",
		Description: "3-qubit CCNOT gate, universal for reversible computing",
		Factory: func(qubits int, params map[string]string) (*core.Program, error) {
			return core.NewToffoliProgram(), nil
		},
	},
	"fredkin": {
		Name:        "Fredkin Gate",
		Description: "3-qubit CSWAP gate, swaps two qubits based on a control",
		Factory: func(qubits int, params map[string]string) (*core.Program, error) {
			return core.NewFredkinProgram(), nil
		},
	},
	"shor": {
		Name:        "Shor's Algorithm",
		Description: "Factoring algorithm. Params: a (default 2), mod (default 7), precision (default 3)",
		Factory: func(qubits int, params map[string]string) (*core.Program, error) {
			a := 2
			mod := 7
			prec := 3
			if v, ok := params["a"]; ok {
				a, _ = strconv.Atoi(v)
			}
			if v, ok := params["mod"]; ok {
				mod, _ = strconv.Atoi(v)
			}
			if v, ok := params["precision"]; ok {
				prec, _ = strconv.Atoi(v)
			}
			return core.NewShorProgram(a, mod, prec), nil
		},
	},
	"adder": {
		Name:        "Quantum Adder",
		Description: "Adds two numbers using Draper adder. Params: x (default 1), y (default 1)",
		Factory: func(qubits int, params map[string]string) (*core.Program, error) {
			x := 1
			y := 1
			if v, ok := params["x"]; ok {
				x, _ = strconv.Atoi(v)
			}
			if v, ok := params["y"]; ok {
				y, _ = strconv.Atoi(v)
			}
			x0, x1 := x&1, (x>>1)&1
			y0, y1, y2 := y&1, (y>>1)&1, (y>>2)&1
			return core.NewAdderProgram(x0, x1, y0, y1, y2), nil
		},
	},
	"qkd": {
		Name:        "Quantum Key Distribution",
		Description: "Single bit preparation for BB84 protocol. Params: bit (0/1), basis (0/1)",
		Factory: func(qubits int, params map[string]string) (*core.Program, error) {
			bit := 1
			basis := 1
			if v, ok := params["bit"]; ok {
				bit, _ = strconv.Atoi(v)
			}
			if v, ok := params["basis"]; ok {
				basis, _ = strconv.Atoi(v)
			}
			return core.NewQKDProgram(bit, basis), nil
		},
	},
	"deutsch-jozsa": {
		Name:        "Deutsch-Jozsa",
		Description: "Classic algorithm to determine if a function is constant or balanced. Params: n (default 2), balanced (true/false)",
		Factory: func(qubits int, params map[string]string) (*core.Program, error) {
			n := 2
			balanced := true
			if v, ok := params["n"]; ok {
				n, _ = strconv.Atoi(v)
			}
			if v, ok := params["balanced"]; ok {
				balanced, _ = strconv.ParseBool(v)
			}
			return core.NewDeutschJozsaProgram(n, balanced), nil
		},
	},
	"bernstein-vazirani": {
		Name:        "Bernstein-Vazirani",
		Description: "Algorithm to find a hidden bitstring. Params: s (default 11)",
		Factory: func(qubits int, params map[string]string) (*core.Program, error) {
			s := "11"
			if v, ok := params["s"]; ok {
				s = v
			}
			return core.NewBernsteinVaziraniProgram(s), nil
		},
	},
	"simon": {
		Name:        "Simon's Algorithm",
		Description: "Algorithm to find a hidden period in a function. Params: s (default 11)",
		Factory: func(qubits int, params map[string]string) (*core.Program, error) {
			s := "11"
			if v, ok := params["s"]; ok {
				s = v
			}
			return core.NewSimonsProgram(s), nil
		},
	},
	"error-correction": {
		Name:        "Error Correction",
		Description: "3-qubit bit-flip code demonstration. Params: bit (0/1, default 1)",
		Factory: func(qubits int, params map[string]string) (*core.Program, error) {
			bit := 1
			if v, ok := params["bit"]; ok {
				bit, _ = strconv.Atoi(v)
			}
			return core.NewErrorCorrectionProgram(bit), nil
		},
	},
	"superposition": {
		Name:        "Superposition",
		Description: "Equal superposition of all states for n qubits",
		Factory: func(qubits int, params map[string]string) (*core.Program, error) {
			return core.NewSuperpositionProgram(qubits), nil
		},
	},
	"engine": {
		Name:        "Quantum Engine",
		Description: "Thermodynamic cycle simulation (requires 'analyze' command)",
		Factory: func(qubits int, params map[string]string) (*core.Program, error) {
			return core.NewEngineProgram(), nil
		},
	},
}

// GetBuiltInProgram returns a core.Program for a recognized circuit name.
func GetBuiltInProgram(name string, qubits int, params map[string]string) (*core.Program, error) {
	if info, ok := builtInCircuits[name]; ok {
		return info.Factory(qubits, params)
	}
	return nil, fmt.Errorf("unknown circuit: %s", name)
}

// GetBuiltInCircuitNames returns a sorted list of all built-in circuit names.
func GetBuiltInCircuitNames() []string {
	names := make([]string, 0, len(builtInCircuits))
	for name := range builtInCircuits {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// GetBuiltInCircuits returns a sorted list of all built-in circuit info.
func GetBuiltInCircuits() []CircuitInfo {
	names := GetBuiltInCircuitNames()
	res := make([]CircuitInfo, len(names))
	for i, name := range names {
		res[i] = builtInCircuits[name]
		res[i].Name = name // Use the key as the identifier
	}
	return res
}
