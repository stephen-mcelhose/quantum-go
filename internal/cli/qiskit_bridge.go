package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
)

// runQiskitVerification executes the QASM string using the standalone qiskit_verify.py tool.
func runQiskitVerification(qasm string) ([]complex128, error) {
	// Find qiskit_verify.py - look in root of repo (one level up from where go test/build might run)
	// or assume it's in the current working directory.
	
	pythonPath, err := exec.LookPath("python3")
	if err != nil {
		return nil, fmt.Errorf("python3 not found in PATH: %w", err)
	}

	// Try common locations for the script
	scriptPaths := []string{"../qiskit_verify.py", "./qiskit_verify.py"}
	var lastErr error
	var output []byte

	for _, scriptPath := range scriptPaths {
		// gosec G204: Subprocess launched with variable.
		// These paths are internal to the tool logic and scriptPath is not user-supplied.
		cmd := exec.Command(pythonPath, scriptPath)
		cmd.Stdin = bytes.NewBufferString(qasm)
		output, err = cmd.Output()
		if err == nil {
			break
		}
		lastErr = err
	}

	if lastErr != nil && output == nil {
		if exitErr, ok := lastErr.(*exec.ExitError); ok {
			return nil, fmt.Errorf("qiskit-verify failed: %s", string(exitErr.Stderr))
		}
		return nil, fmt.Errorf("failed to run qiskit_verify.py: %w (ensure python3 and qiskit-aer are installed)", lastErr)
	}


	// 2. Parse the result
	var response struct {
		NumQubits  int `json:"num_qubits"`
		Amplitudes []struct {
			Re float64 `json:"re"`
			Im float64 `json:"im"`
		} `json:"amplitudes"`
	}
	if err := json.Unmarshal(output, &response); err != nil {
		return nil, fmt.Errorf("failed to parse qiskit-verify output: %w", err)
	}

	res := make([]complex128, len(response.Amplitudes))
	for i, a := range response.Amplitudes {
		res[i] = complex(a.Re, a.Im)
	}

	return res, nil
}
