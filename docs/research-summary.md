# Strange Quantum Simulator: Go Port Research Summary

This document establishes the requirements for porting the 'Strange' quantum simulator to Go.

## Core Findings

- **State Representation**: The simulator uses a dense-vector state representation.
- **Math Engine**: Relies on a custom complex number and matrix math engine.
- **Porting Strategy**:
    - Leverage Go's native `complex128` type for scalar math.
    - Implement a custom `math` package for matrix operations like Kronecker (tensor) product and matrix multiplication.
    - Preserve optimized gate application paths from the original `Computations.java`.
    - Map Java's `Program`/`Step`/`Gate` hierarchy to Go structs and interfaces.
    - Replace Java's thread-based parallelism with Go's goroutines.

## Key Files (Java)

- `Complex.java`: Scalar and matrix complex math foundation.
- `Program.java`: Core model for a quantum circuit.
- `SimpleQuantumExecutionEnvironment.java`: Main simulation engine for local execution.
- `Computations.java`: Low-level quantum state evolution logic.

## Architecture Patterns

- Dense vector state representation.
- Hierarchical structure: `Program` -> `Step` -> `Gate`.
- Recursive and bit-manipulation based state evolution.
