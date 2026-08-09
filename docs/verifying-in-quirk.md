# Verifying quantum-go Results in Quirk

This guide explains how to manually verify the numerical results from the `quantum-go` simulator using the visual [Quirk Quantum Simulator](https://algassert.com/quirk).

> [!WARNING]
> **Manual Verification Required**: Visual simulators like Quirk are excellent for intuition but may have different internal implementations or precision limits. Always perform a manual side-by-side comparison of the raw complex amplitudes to ensure $10^{-6}$ precision parity.

## Verification Steps

### 1. Generate quantum-go Reference Data
In your Go code, export the program to OpenQASM and print the final state vector amplitudes:

```go
// 1. Export the circuit
fmt.Println(program.ToQASM())

// 2. Print the state vector
result := engine.RunProgram(program)
for i, amp := range result.GetProbability() {
    fmt.Printf("Index %d: %v\n", i, amp)
}
```

### 2. Recreate the Circuit in Quirk
1.  Open [Quirk](https://algassert.com/quirk).
2.  Clear the workspace.
3.  Place gates on the wires matching the **QASM output** from Step 1.
    - `h q[0]` $\rightarrow$ Place a Hadamard (H) on the top wire.
    - `cx q[0], q[1]` $\rightarrow$ Place a CNOT with the control dot on the top wire and the $\oplus$ on the second wire.

### 3. Inspect raw Amplitudes
To compare the numbers precisely:
1.  Locate the **"Probes"** section in the Quirk toolbox (bottom).
2.  Drag the **"Amplitudes"** probe (blue vertical bars icon) to the end of your circuit.
3.  **Hover** your mouse over the probe. A tooltip will appear showing the **Complex Amplitudes** for every binary state.

### 4. Binary State Mapping (Endianness)
quantum-go and Quirk both use **Little-Endian** qubit ordering (LSB is Wire 0). Use this table to map indices:

| quantum-go Index | Binary State | Quirk Wire State |
| :--- | :--- | :--- |
| `0` | `|000...>` | All wires are $|0\rangle$ |
| `1` | `|001...>` | **Top wire** (0) is $|1\rangle$ |
| `2` | `|010...>` | **Second wire** (1) is $|1\rangle$ |
| `4` | `|100...>` | **Third wire** (2) is $|1\rangle$ |

### 5. Verifying Thermodynamic Gates
For specialized gates like `TimeEvolution`, Quirk may not have a native equivalent. Use a **Custom Gate**:
1.  In Quirk, click **"Make Gate"** (top menu).
2.  Select **"Matrix"**.
3.  In Go, retrieve the gate's matrix via `gate.GetMatrix()`.
4.  Paste the matrix values into Quirk's custom gate builder.
5.  Place the resulting gate into the circuit and verify the final state vector against the Go `Result`.

## Precision Check
Compare the hover-values in Quirk with the `complex128` values from Go. Both should match exactly for standard gates, or within $10^{-6}$ for calculated matrices (like QFT or Time Evolution).
