#!/usr/bin/env python3
"""
Verification script to compare quantum-go quantum operations against Qiskit.
This script runs each quantum operation from the learning transcript in both
quantum-go and Qiskit, then compares the results.
"""

import subprocess
import json
import numpy as np
from qiskit import QuantumCircuit
from qiskit_aer import AerSimulator
from typing import Dict, List, Tuple
import sys
import os
import argparse

class QuantumVerifier:
    """Compares quantum operations between quantum-go and Qiskit."""
    
    def __init__(self, quantum_go_path: str = "./go/quantum-go"):
        self.quantum_go_path = quantum_go_path
        self.simulator = AerSimulator(method='statevector')
        self.tolerance = 1e-4  # Tolerance for probability comparison
        
        # Check if quantum-go binary exists
        if not os.path.exists(self.quantum_go_path):
            print("\n" + "="*70)
            print(" ERROR: quantum-go binary not found")
            print("="*70)
            print(f"\nCannot find quantum-go binary at: {self.quantum_go_path}")
            print("\nTo fix this:")
            print("  1. Build quantum-go: cd go && go build -o quantum-go ./cmd/quantum-go")
            print("  2. Run from repo root: python3 go/verification/verify_against_qiskit.py")
            print("  3. Or specify path: python3 verify_against_qiskit.py /path/to/quantum-go")
            print("="*70 + "\n")
            sys.exit(1)
        
    def run_quantum_go(self, num_qubits: int, gates: str = None, 
                    circuit: str = None, params: Dict[str, str] = None) -> Dict[str, float]:
        """
        Run a quantum circuit in quantum-go and return probability distribution.
        
        Args:
            num_qubits: Number of qubits
            gates: Gate sequence (e.g., "h q[0]")
            circuit: Built-in circuit name
            params: Parameters for built-in circuit
            
        Returns:
            Dictionary mapping states to probabilities
        """
        cmd = [self.quantum_go_path, "run"]
        if circuit:
            cmd.extend(["--circuit", circuit])
            if params:
                param_str = ",".join([f"{k}={v}" for k, v in params.items()])
                cmd.extend(["-p", param_str])
        else:
            cmd.extend(["-n", str(num_qubits), "-s", gates])
            
        try:
            result = subprocess.run(cmd, capture_output=True, text=True, check=True)
            return self._parse_quantum_go_output(result.stdout)
        except FileNotFoundError:
            print(f"\nERROR: quantum-go binary not found at: {self.quantum_go_path}")
            print("Build it with: cd go && go build -o quantum-go ./cmd/quantum-go")
            return {}
        except subprocess.CalledProcessError as e:
            print(f"Error running quantum-go: {e.stderr}")
            return {}
    
    def _parse_quantum_go_output(self, output: str) -> Dict[str, float]:
        """Parse quantum-go CLI output to extract probability distribution."""
        probs = {}
        for line in output.strip().split('\n'):
            if line.startswith('|') and ':' in line:
                # Parse format: |00>: 0.5000
                state = line.split(':')[0].strip()
                prob = float(line.split(':')[1].strip())
                probs[state] = prob
        return probs
    
    def run_qiskit(self, num_qubits: int, gates: List[Tuple[str, List[int], float]]) -> Dict[str, float]:
        """
        Run a quantum circuit in Qiskit and return probability distribution.
        
        Args:
            num_qubits: Number of qubits
            gates: List of (gate_name, qubit_indices, parameter) tuples
            
        Returns:
            Dictionary mapping states to probabilities
        """
        qc = QuantumCircuit(num_qubits)
        
        for gate_info in gates:
            gate_name = gate_info[0]
            qubits = gate_info[1]
            param = gate_info[2] if len(gate_info) > 2 else None
            
            if gate_name == 'h':
                qc.h(qubits[0])
            elif gate_name == 'x':
                qc.x(qubits[0])
            elif gate_name == 'y':
                qc.y(qubits[0])
            elif gate_name == 'z':
                qc.z(qubits[0])
            elif gate_name == 'id':
                qc.id(qubits[0])
            elif gate_name == 'cx' or gate_name == 'cnot':
                qc.cx(qubits[0], qubits[1])
            elif gate_name == 'cz':
                qc.cz(qubits[0], qubits[1])
            elif gate_name == 'swap':
                qc.swap(qubits[0], qubits[1])
            elif gate_name == 'cswap' or gate_name == 'fredkin':
                qc.cswap(qubits[0], qubits[1], qubits[2])
            elif gate_name == 'ccx' or gate_name == 'toffoli':
                qc.ccx(qubits[0], qubits[1], qubits[2])
            elif gate_name == 'u1':
                qc.p(param, qubits[0])  # u1 is now p (phase gate) in Qiskit
            elif gate_name == 'cu1' or gate_name == 'cp':
                qc.cp(param, qubits[0], qubits[1])  # cu1 is now cp in Qiskit
            elif gate_name == 'rx':
                qc.rx(param, qubits[0])
            elif gate_name == 'ry':
                qc.ry(param, qubits[0])
            elif gate_name == 'rz':
                qc.rz(param, qubits[0])
            elif gate_name == 's':
                qc.s(qubits[0])
            elif gate_name == 't':
                qc.t(qubits[0])
            elif gate_name == 'sx' or gate_name == 'v':
                qc.sx(qubits[0])
            elif gate_name == 'u' or gate_name == 'u3':
                # param for 'u' is (theta, phi, lambda)
                if isinstance(param, (list, tuple)):
                    qc.u(param[0], param[1], param[2], qubits[0])
                else:
                    qc.u(param, 0, 0, qubits[0])
            else:
                print(f"Warning: Unknown gate {gate_name}")
        
        qc.save_statevector()
        result = self.simulator.run(qc).result()
        statevector = np.asarray(result.get_statevector())
        
        # Convert statevector to probability distribution
        probs = {}
        for i, amplitude in enumerate(statevector):
            prob = abs(amplitude) ** 2
            if prob > 1e-10:  # Only include non-zero probabilities
                state = f"|{format(i, f'0{num_qubits}b')}>"
                probs[state] = float(prob)
        
        return probs
    
    def compare_results(self, quantum_go_probs: Dict[str, float], 
                       qiskit_probs: Dict[str, float], 
                       test_name: str) -> bool:
        """
        Compare probability distributions from quantum-go and Qiskit.
        
        This method performs a detailed comparison of quantum measurement probabilities
        from both simulators to verify that quantum-go produces correct results.
        
        How it works:
        1. Collects all quantum states that appear in either result (union of states)
        2. For each state, compares the probability values:
           - quantum-go probability: P_s(|ψ⟩) = |⟨ψ|φ_s⟩|²
           - Qiskit probability: P_q(|ψ⟩) = |⟨ψ|φ_q⟩|²
        3. Calculates absolute difference: |P_s - P_q|
        4. Checks if difference is within tolerance (default: 0.0001 = 0.01%)
        5. Marks test as PASS only if ALL states match within tolerance
        
        Example comparison for Bell state |00⟩ + |11⟩:
            State    quantum-go    Qiskit     Diff        Match
            |00⟩     0.5000     0.5000     0.000000    ✓
            |11⟩     0.5000     0.5000     0.000000    ✓
        
        Why this matters:
        - Quantum states must match EXACTLY (within floating-point precision)
        - Even small differences indicate implementation bugs
        - Probabilities must sum to 1.0 (verified implicitly)
        
        Args:
            quantum_go_probs: Dict mapping states like "|00⟩" to probabilities [0,1]
            qiskit_probs: Dict mapping states like "|00⟩" to probabilities [0,1]
            test_name: Human-readable name for the test being run
            
        Returns:
            True if ALL state probabilities match within self.tolerance
            False if ANY state probability differs beyond tolerance
            
        Side effects:
            Prints detailed comparison table showing:
            - Each quantum state (e.g., |00⟩, |01⟩, |10⟩, |11⟩)
            - Probability from quantum-go simulator
            - Probability from Qiskit simulator
            - Absolute difference between the two
            - Pass/fail indicator (✓ or ✗) for each state
            - Overall test result (PASS ✓ or FAIL ✗)
        """
        print(f"\n{'='*60}")
        print(f"Test: {test_name}")
        print(f"{'='*60}")
        
        # Collect all states from both results (union)
        # 
        # WHY UNION? This catches discrepancies where one simulator produces states
        # that the other doesn't:
        #
        # Example Bug Scenario:
        #   quantum-go:  {|00⟩: 0.5, |11⟩: 0.5}        ← Correct
        #   Qiskit:   {|00⟩: 0.5, |01⟩: 0.3, |11⟩: 0.2}  ← Bug: extra state!
        #
        # If we only checked quantum-go's states, we'd miss that Qiskit produced |01⟩
        # If we only checked Qiskit's states, we'd miss if quantum-go had extra states
        #
        # By using UNION (|), we check ALL states from BOTH simulators:
        #   - States in both: Compare probabilities directly
        #   - States in only one: Other gets 0.0, causing large diff → FAIL
        #
        # This ensures:
        #   1. No spurious states (would show up with prob > 0 vs 0)
        #   2. No missing states (would show up as 0 vs prob > 0)
        #   3. All probabilities match exactly
        all_states = set(quantum_go_probs.keys()) | set(qiskit_probs.keys())
        
        match = True  # Assume success until we find a mismatch
        print(f"\n{'State':<10} {'quantum-go':<14} {'Qiskit':<12} {'Diff':<12} {'Match'}")
        print("-" * 60)
        
        for state in sorted(all_states):
            # Get probability for this state (default to 0.0 if missing)
            quantum_go_prob = quantum_go_probs.get(state, 0.0)
            qiskit_prob = qiskit_probs.get(state, 0.0)
            
            # Calculate absolute difference
            diff = abs(quantum_go_prob - qiskit_prob)
            
            # Check if within tolerance
            matches = diff < self.tolerance
            
            # If ANY state fails, mark entire test as failed
            if not matches:
                match = False
            
            # Display checkmark or X for this state
            status = "✓" if matches else "✗"
            
            # Print formatted row: |00⟩  0.5000  0.5000  0.000000  ✓
            print(f"{state:<10} {quantum_go_prob:>10.4f}   {qiskit_prob:>10.4f}   "
                  f"{diff:>10.6f}   {status}")
        
        # Print overall test result
        print(f"\n{'Overall Result:':<30} {'PASS ✓' if match else 'FAIL ✗'}")
        return match
    
    def verify_test(self, test_name: str, num_qubits: int, 
                    quantum_go_gates: str = None, qiskit_gates: List[Tuple] = None,
                    circuit: str = None, params: Dict[str, str] = None) -> bool:
        """Run a single verification test."""
        quantum_go_probs = self.run_quantum_go(num_qubits, quantum_go_gates, circuit, params)
        qiskit_probs = self.run_qiskit(num_qubits, qiskit_gates)
        return self.compare_results(quantum_go_probs, qiskit_probs, test_name)


def main():
    """Run all verification tests from the learning transcript."""
    parser = argparse.ArgumentParser(
        description='Verify quantum-go simulator against IBM Qiskit',
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  python3 verify_against_qiskit.py                    # Use default path ./go/quantum-go
  python3 verify_against_qiskit.py ../quantum-go         # Use relative path
  python3 verify_against_qiskit.py /usr/local/b../quantum-go  # Use absolute path
        """
    )
    parser.add_argument(
        'quantum_go_path',
        nargs='?',
        default='./go/quantum-go',
        help='Path to the quantum-go binary (default: ./go/quantum-go)'
    )
    
    args = parser.parse_args()
    
    verifier = QuantumVerifier(args.quantum_go_path)
    
    tests = [
        # Objective 1: Understanding Superposition
        {
            "name": "Single Qubit Superposition (H on q[0])",
            "num_qubits": 1,
            "quantum_go_gates": "h q[0]",
            "qiskit_gates": [("h", [0])]
        },
        {
            "name": "Two Qubit Superposition (H on both)",
            "num_qubits": 2,
            "quantum_go_gates": "h q[0]; h q[1]",
            "qiskit_gates": [("h", [0]), ("h", [1])]
        },
        {
            "name": "Three Qubit Superposition (H on all)",
            "num_qubits": 3,
            "quantum_go_gates": "h q[0]; h q[1]; h q[2]",
            "qiskit_gates": [("h", [0]), ("h", [1]), ("h", [2])]
        },
        
        # Objective 2: Selective Superposition
        {
            "name": "Selective - Only q[0] in superposition",
            "num_qubits": 3,
            "quantum_go_gates": "h q[0]",
            "qiskit_gates": [("h", [0])]
        },
        {
            "name": "Selective - q[0] and q[2] in superposition",
            "num_qubits": 3,
            "quantum_go_gates": "h q[0]; h q[2]",
            "qiskit_gates": [("h", [0]), ("h", [2])]
        },
        {
            "name": "Selective - Only q[1] in superposition",
            "num_qubits": 3,
            "quantum_go_gates": "h q[1]",
            "qiskit_gates": [("h", [1])]
        },
        
        # Objective 3: Gate Reversibility
        {
            "name": "Double Hadamard (H·H = I)",
            "num_qubits": 1,
            "quantum_go_gates": "h q[0]; h q[0]",
            "qiskit_gates": [("h", [0]), ("h", [0])]
        },
        {
            "name": "X Gate (NOT)",
            "num_qubits": 1,
            "quantum_go_gates": "x q[0]",
            "qiskit_gates": [("x", [0])]
        },
        {
            "name": "Double X Gate (X·X = I)",
            "num_qubits": 1,
            "quantum_go_gates": "x q[0]; x q[0]",
            "qiskit_gates": [("x", [0]), ("x", [0])]
        },
        
        # Objective 4: Z Gate Investigation
        {
            "name": "Z Gate on |0⟩ (no visible change)",
            "num_qubits": 1,
            "quantum_go_gates": "z q[0]",
            "qiskit_gates": [("z", [0])]
        },
        {
            "name": "Z Gate phase effect (H·Z·H)",
            "num_qubits": 1,
            "quantum_go_gates": "h q[0]; z q[0]; h q[0]",
            "qiskit_gates": [("h", [0]), ("z", [0]), ("h", [0])]
        },
        
        # Objective 5: Y Gate
        {
            "name": "Y Gate (flips like X)",
            "num_qubits": 1,
            "quantum_go_gates": "y q[0]",
            "qiskit_gates": [("y", [0])]
        },
        {
            "name": "Double Y Gate (Y·Y = I)",
            "num_qubits": 1,
            "quantum_go_gates": "y q[0]; y q[0]",
            "qiskit_gates": [("y", [0]), ("y", [0])]
        },
        
        # Objective 6: Gate Compositions
        {
            "name": "X then Z",
            "num_qubits": 1,
            "quantum_go_gates": "x q[0]; z q[0]",
            "qiskit_gates": [("x", [0]), ("z", [0])]
        },
        {
            "name": "Z then X",
            "num_qubits": 1,
            "quantum_go_gates": "z q[0]; x q[0]",
            "qiskit_gates": [("z", [0]), ("x", [0])]
        },
        {
            "name": "H-sandwich with X (H·X·H)",
            "num_qubits": 1,
            "quantum_go_gates": "h q[0]; x q[0]; h q[0]",
            "qiskit_gates": [("h", [0]), ("x", [0]), ("h", [0])]
        },
        {
            "name": "H-sandwich with Z (H·Z·H)",
            "num_qubits": 1,
            "quantum_go_gates": "h q[0]; z q[0]; h q[0]",
            "qiskit_gates": [("h", [0]), ("z", [0]), ("h", [0])]
        },
        {
            "name": "Triple combo (X·Y·Z)",
            "num_qubits": 1,
            "quantum_go_gates": "x q[0]; y q[0]; z q[0]",
            "qiskit_gates": [("x", [0]), ("y", [0]), ("z", [0])]
        },
        
        # Examples from CLI help: Basic Gates
        {
            "name": "CLI Example: Identity gate",
            "num_qubits": 1,
            "quantum_go_gates": "id q[0]",
            "qiskit_gates": [("id", [0])]
        },
        {
            "name": "CLI Example: CNOT gate",
            "num_qubits": 2,
            "quantum_go_gates": "cx q[0], q[1]",
            "qiskit_gates": [("cx", [0, 1])]
        },
        {
            "name": "CLI Example: CZ gate",
            "num_qubits": 2,
            "quantum_go_gates": "cz q[0], q[1]",
            "qiskit_gates": [("cz", [0, 1])]
        },
        {
            "name": "CLI Example: SWAP gate",
            "num_qubits": 2,
            "quantum_go_gates": "swap q[0], q[1]",
            "qiskit_gates": [("swap", [0, 1])]
        },
        {
            "name": "CLI Example: Toffoli (CCNOT) gate",
            "num_qubits": 3,
            "quantum_go_gates": "ccx q[0], q[1], q[2]",
            "qiskit_gates": [("ccx", [0, 1, 2])]
        },
        {
            "name": "CLI Example: U1 Phase Rotation (π/2)",
            "num_qubits": 1,
            "quantum_go_gates": "u1(1.57) q[0]",
            "qiskit_gates": [("u1", [0], 1.57)]
        },
        {
            "name": "CLI Example: Controlled U1 Rotation",
            "num_qubits": 2,
            "quantum_go_gates": "cu1(1.57) q[0], q[1]",
            "qiskit_gates": [("cu1", [0, 1], 1.57)]
        },
        
        # Examples from 'quantum-go run' help
        {
            "name": "Run Example: Bell state (H + CNOT)",
            "num_qubits": 2,
            "quantum_go_gates": "h q[0]; cx q[0], q[1]",
            "qiskit_gates": [("h", [0]), ("cx", [0, 1])]
        },
        
        # Examples from 'quantum-go inspect' help
        {
            "name": "Inspect Example: 3-qubit GHZ state",
            "num_qubits": 3,
            "quantum_go_gates": "h q[0]; cx q[0], q[1]; cx q[1], q[2]",
            "qiskit_gates": [("h", [0]), ("cx", [0, 1]), ("cx", [1, 2])]
        },
        
        # Examples from 'quantum-go export' help
        {
            "name": "Export Example: Custom 3-qubit GHZ",
            "num_qubits": 3,
            "quantum_go_gates": "h q[0]; cx q[0], q[1]; cx q[1], q[2]",
            "qiskit_gates": [("h", [0]), ("cx", [0, 1]), ("cx", [1, 2])]
        },
        
        # Examples from 'quantum-go verify' help
        {
            "name": "Verify Example: H + Phase rotation",
            "num_qubits": 2,
            "quantum_go_gates": "h q[0]; u1(0.785398) q[0]",
            "qiskit_gates": [("h", [0]), ("u1", [0], 0.785398)]
        },
        
        # Examples from 'quantum-go analyze' help
        {
            "name": "Analyze Example: Bell state",
            "num_qubits": 2,
            "quantum_go_gates": "h q[0]; cx q[0], q[1]",
            "qiskit_gates": [("h", [0]), ("cx", [0, 1])]
        },
        {
            "name": "Analyze Example: 3-qubit GHZ",
            "num_qubits": 3,
            "quantum_go_gates": "h q[0]; cx q[0], q[1]; cx q[1], q[2]",
            "qiskit_gates": [("h", [0]), ("cx", [0, 1]), ("cx", [1, 2])]
        },
        
        # Additional multi-qubit tests
        {
            "name": "CNOT with swapped qubits",
            "num_qubits": 2,
            "quantum_go_gates": "x q[1]; cx q[1], q[0]",
            "qiskit_gates": [("x", [1]), ("cx", [1, 0])]
        },
        {
            "name": "Multiple CNOTs in sequence",
            "num_qubits": 2,
            "quantum_go_gates": "h q[0]; cx q[0], q[1]; cx q[0], q[1]",
            "qiskit_gates": [("h", [0]), ("cx", [0, 1]), ("cx", [0, 1])]
        },
        {
            "name": "Toffoli with X on controls",
            "num_qubits": 3,
            "quantum_go_gates": "x q[0]; x q[1]; ccx q[0], q[1], q[2]",
            "qiskit_gates": [("x", [0]), ("x", [1]), ("ccx", [0, 1, 2])]
        },
        {
            "name": "Phase rotation combinations",
            "num_qubits": 1,
            "quantum_go_gates": "h q[0]; u1(1.57) q[0]; h q[0]",
            "qiskit_gates": [("h", [0]), ("u1", [0], 1.57), ("h", [0])]
        },
        {
            "name": "SWAP effect verification",
            "num_qubits": 2,
            "quantum_go_gates": "x q[0]; swap q[0], q[1]",
            "qiskit_gates": [("x", [0]), ("swap", [0, 1])]
        },
        {
            "name": "CZ symmetry test",
            "num_qubits": 2,
            "quantum_go_gates": "h q[0]; h q[1]; cz q[0], q[1]",
            "qiskit_gates": [("h", [0]), ("h", [1]), ("cz", [0, 1])]
        },
        {
            "name": "Controlled phase with superposition",
            "num_qubits": 2,
            "quantum_go_gates": "h q[0]; cu1(0.785398) q[0], q[1]",
            "qiskit_gates": [("h", [0]), ("cu1", [0, 1], 0.785398)]
        },
        # New Standard Gates Verification
        {
            "name": "Rotation X (Rx(π))",
            "num_qubits": 1,
            "quantum_go_gates": "rx(3.141593) q[0]",
            "qiskit_gates": [("rx", [0], 3.141593)]
        },
        {
            "name": "Rotation Y (Ry(π))",
            "num_qubits": 1,
            "quantum_go_gates": "ry(3.141593) q[0]",
            "qiskit_gates": [("ry", [0], 3.141593)]
        },
        {
            "name": "Rotation Z (Rz(π))",
            "num_qubits": 1,
            "quantum_go_gates": "rz(3.141593) q[0]",
            "qiskit_gates": [("rz", [0], 3.141593)]
        },
        {
            "name": "Phase Gate (S)",
            "num_qubits": 1,
            "quantum_go_gates": "s q[0]",
            "qiskit_gates": [("s", [0])]
        },
        {
            "name": "T Gate",
            "num_qubits": 1,
            "quantum_go_gates": "t q[0]",
            "qiskit_gates": [("t", [0])]
        },
        {
            "name": "Square Root of X (V/SX)",
            "num_qubits": 1,
            "quantum_go_gates": "sx q[0]",
            "qiskit_gates": [("sx", [0])]
        },
        {
            "name": "Universal Gate (U(π/2, 0, π))",
            "num_qubits": 1,
            "quantum_go_gates": "u3(1.570796, 0, 3.141593) q[0]",
            "qiskit_gates": [("u3", [0], (1.570796, 0, 3.141593))]
        },
        {
            "name": "Universal Gate (U(0, 0, 0) - Identity)",
            "num_qubits": 1,
            "quantum_go_gates": "u3(0, 0, 0) q[0]",
            "qiskit_gates": [("u3", [0], (0, 0, 0))]
        },
        # 2-Qubit Standard Gates Verification
        {
            "name": "2-Qubit Rx (Rx(π) on both)",
            "num_qubits": 2,
            "quantum_go_gates": "rx(3.141593) q[0]; rx(3.141593) q[1]",
            "qiskit_gates": [("rx", [0], 3.141593), ("rx", [1], 3.141593)]
        },
        {
            "name": "2-Qubit Ry (Ry(π) on both)",
            "num_qubits": 2,
            "quantum_go_gates": "ry(3.141593) q[0]; ry(3.141593) q[1]",
            "qiskit_gates": [("ry", [0], 3.141593), ("ry", [1], 3.141593)]
        },
        {
            "name": "2-Qubit Rz (Rz(π) on both)",
            "num_qubits": 2,
            "quantum_go_gates": "rz(3.141593) q[0]; rz(3.141593) q[1]",
            "qiskit_gates": [("rz", [0], 3.141593), ("rz", [1], 3.141593)]
        },
        {
            "name": "2-Qubit S Gate (S on both)",
            "num_qubits": 2,
            "quantum_go_gates": "s q[0]; s q[1]",
            "qiskit_gates": [("s", [0]), ("s", [1])]
        },
        {
            "name": "2-Qubit T Gate (T on both)",
            "num_qubits": 2,
            "quantum_go_gates": "t q[0]; t q[1]",
            "qiskit_gates": [("t", [0]), ("t", [1])]
        },
        {
            "name": "2-Qubit SX Gate (SX on both)",
            "num_qubits": 2,
            "quantum_go_gates": "sx q[0]; sx q[1]",
            "qiskit_gates": [("sx", [0]), ("sx", [1])]
        },
        {
            "name": "2-Qubit Mixed: Rx and Ry",
            "num_qubits": 2,
            "quantum_go_gates": "rx(1.570796) q[0]; ry(1.570796) q[1]",
            "qiskit_gates": [("rx", [0], 1.570796), ("ry", [1], 1.570796)]
        },
        {
            "name": "2-Qubit Mixed: S and T",
            "num_qubits": 2,
            "quantum_go_gates": "s q[0]; t q[1]",
            "qiskit_gates": [("s", [0]), ("t", [1])]
        },
        {
            "name": "2-Qubit Mixed: U and H",
            "num_qubits": 2,
            "quantum_go_gates": "u3(1.570796, 0, 3.141593) q[0]; h q[1]",
            "qiskit_gates": [("u3", [0], (1.570796, 0, 3.141593)), ("h", [1])]
        },
        {
            "name": "Fredkin Gate (CSWAP)",
            "num_qubits": 3,
            "quantum_go_gates": "x q[0]; swap q[1], q[2] controlled by q[0]",
            "qiskit_gates": [("x", [0]), ("cswap", [0, 1, 2])]
        },
        {
            "name": "Built-in: Superdense Coding (11)",
            "num_qubits": 2,
            "circuit": "superdense",
            "qiskit_gates": [("h", [0]), ("cx", [0, 1]), ("x", [0]), ("z", [0]), ("cx", [0, 1]), ("h", [0])]
        },
        {
            "name": "Built-in: Bernstein-Vazirani (s=101)",
            "num_qubits": 4,
            "circuit": "bernstein-vazirani",
            "params": {"s": "101"},
            "qiskit_gates": [("x", [3]), ("h", [0]), ("h", [1]), ("h", [2]), ("h", [3]), ("cx", [0, 3]), ("cx", [2, 3]), ("h", [0]), ("h", [1]), ("h", [2])]
        },
        {
            "name": "Built-in: Deutsch-Jozsa (Balanced)",
            "num_qubits": 3,
            "circuit": "deutsch-jozsa",
            "params": {"n": "2", "balanced": "true"},
            "qiskit_gates": [("x", [2]), ("h", [0]), ("h", [1]), ("h", [2]), ("cx", [0, 2]), ("h", [0]), ("h", [1])]
        },
        {
            "name": "Built-in: Simon's Algorithm (s=11)",
            "num_qubits": 4,
            "circuit": "simon",
            "params": {"s": "11"},
            "qiskit_gates": [("h", [0]), ("h", [1]), ("cx", [0, 2]), ("cx", [1, 2]), ("h", [0]), ("h", [1])]
        },
        {
            "name": "Built-in: Error Correction (bit 1)",
            "num_qubits": 3,
            "circuit": "error-correction",
            "params": {"bit": "1"},
            "qiskit_gates": [("x", [0]), ("cx", [0, 1]), ("cx", [0, 2]), ("x", [1]), ("cx", [0, 1]), ("cx", [0, 2]), ("ccx", [1, 2, 0])]
        },
        {
            "name": "Built-in: Toffoli Gate",
            "num_qubits": 3,
            "circuit": "toffoli",
            "qiskit_gates": [("x", [0]), ("x", [1]), ("ccx", [0, 1, 2])]
        },
        {
            "name": "Built-in: Deutsch-Jozsa (Constant)",
            "num_qubits": 3,
            "circuit": "deutsch-jozsa",
            "params": {"n": "2", "balanced": "false"},
            "qiskit_gates": [("x", [2]), ("h", [0]), ("h", [1]), ("h", [2]), ("h", [0]), ("h", [1])]
        },
        {
            "name": "Built-in: Error Correction (bit 0)",
            "num_qubits": 3,
            "circuit": "error-correction",
            "params": {"bit": "0"},
            "qiskit_gates": [("cx", [0, 1]), ("cx", [0, 2]), ("x", [1]), ("cx", [0, 1]), ("cx", [0, 2]), ("ccx", [1, 2, 0])]
        },
    ]
    
    print("\n" + "="*70)
    print(" quantum-go vs QISKIT VERIFICATION SUITE")
    print(" Includes: Learning Transcript + CLI Help Examples")
    print("="*70)
    
    passed = 0
    failed = 0
    
    for test in tests:
        try:
            result = verifier.verify_test(
                test["name"],
                test["num_qubits"],
                test.get("quantum_go_gates"),
                test.get("qiskit_gates"),
                test.get("circuit"),
                test.get("params")
            )
            if result:
                passed += 1
            else:
                failed += 1
        except Exception as e:
            print(f"\nError running test '{test['name']}': {e}")
            failed += 1
    
    # Summary
    print("\n" + "="*70)
    print(" SUMMARY")
    print("="*70)
    print(f"Total Tests: {passed + failed}")
    print(f"Passed:      {passed} ✓")
    print(f"Failed:      {failed} ✗")
    print(f"Success Rate: {100 * passed / (passed + failed):.1f}%")
    print("="*70 + "\n")
    
    return 0 if failed == 0 else 1


if __name__ == "__main__":
    sys.exit(main())
