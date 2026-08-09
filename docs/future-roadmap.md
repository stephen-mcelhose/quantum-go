# Research Report: quantum-go Port Roadmap & Future Enhancements

## Summary
The current state of the `quantum-go` Go port is stable and functionally rich, covering fundamental quantum gates, complex algorithmic blocks (QFT, Adder), and several built-in quantum algorithms. This report outlines a roadmap for future enhancements to transition the simulator from a basic port to a feature-complete quantum development tool.

## 1. Algorithmic Roadmap
Several high-value algorithms from the original Strange Java simulator ecosystem and broader quantum computing literature are recommended for implementation:

- **Simon's Algorithm**: A key precursor to Shor's algorithm that demonstrates an exponential speedup.
- **Quantum Phase Estimation (QPE)**: A fundamental building block for Shor's and many other algorithms.
- **Variational Framework**: Basic support for VQE (Variational Quantum Eigensolver) by exposing an interface for classical optimizers to update gate parameters.

## 2. CLI & User Experience
The discovery commands (`list circuits`, `list gates`) have significantly improved the UX. Further enhancements include:

- **ASCII Circuit Drawing**: Implementing a `quantum-go draw` command to visualize circuits. The `core.Program` structure (ordered steps of disjoint gates) is perfectly suited for a grid-based ASCII renderer.
- **Enhanced QASM Support**: Extending the `qasm` parser to support the algorithmic gates already present in `core` (e.g., `qft q[0], q[1], q[2]`).
- **Intermediate State Inspection**: Adding an `--inspect-steps` flag to `quantum-go run` to display state vector transitions between steps.

## 3. Advanced Physics & Thermodynamics
Building on the existing thermodynamic capabilities:

- **Noise Models**: Implementing open quantum system simulation (decoherence) via Kraus operators or Lindblad master equations.
- **Advanced Thermodynamic Cycles**: Replicating more complex benchmarks from the "Quantum Heat Engine" research (sciadv.adw8462).

## 4. Maintenance & Reliability
- **Fuzzing Expansion**: Adding fuzz tests for the `Fredkin` gate and `MulModulus` to ensure correctness across edge cases.
- **Algorithm Verification**: Formalizing the verification of built-in algorithms against their theoretical expected state vectors (similar to `ExpectedBellState`).

## 5. Identified Potential Issues
- **Multiplier Registry Alias**: The alias `mul` for `Multiplier` in the gate registry might be ambiguous; `mulm` (Modular Multiplier) may be more accurate as it reflects the current `core` implementation (`MulModulus`).
- **Fredkin Gate Target Specificity**: The current `Fredkin` matrix implementation assumes a hardcoded qubit order (0=control, 1,2=targets). A more general approach using qubit permutations or bit-manipulation is needed for full flexibility.
