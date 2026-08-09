# Fuzzy Testing Strategy for Advanced Gates

Fuzzy testing is used to verify the mathematical correctness and stability of complex quantum gates and linear algebra operations.

## Focus Areas

### 1. Matrix Operations
- **Matrix Multiplication (`Mul`)**: Verify dimension handling and accuracy for randomized complex matrices.
- **Kronecker Product (`Tensor`)**: Ensure resulting dimensions and element-wise products are correct.

### 2. Advanced Quantum Gates
- **Toffoli (CCNOT)**: Verify CCNOT logic across all 8 base input states.
- **Quantum Fourier Transform (QFT / Fourier)**:
    - **Unitarity**: Verify that $QFT \circ IQFT$ returns the system to the initial state.
    - **Normalization**: Ensure the state vector remains normalized ($\sum |a_i|^2 \approx 1.0$).
- **Adder (`Add`)**: Verify modular addition correctness ($(x + y) \mod 2^m$) for various register sizes and initial values.

## Implementation Details

Tests are located in:
- `go/math/fuzz_test.go`: Linear algebra fuzzing.
- `go/local/fuzz_test.go`: Simulation and gate logic fuzzing.

## Execution

Run fuzz tests using:
```bash
go test -v -fuzz=FuzzToffoli -fuzztime=10s ./local
go test -v -fuzz=FuzzFourier -fuzztime=10s ./local
go test -v -fuzz=FuzzAdd -fuzztime=10s ./local
```
