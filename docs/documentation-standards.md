# Strange-Go Pedagogical Documentation Standards

This document codifies the rubric and standards for creating pedagogical learning aids (Markdown guides) for quantum algorithms and fundamental concepts in the `quantum-go` simulator.

## Goal
The goal of these guides is to bridge the gap between abstract quantum mechanics theory and practical Go implementation, providing developers with a clear "why" and "how" for each example.

## Rubric for Algorithm Guides

Each learning aid (e.g., `shor.md`, `grover.md`) MUST follow this structure:

### 1. Title
A clear, action-oriented title.
*Example: Learning Grover's Algorithm with Strange-Go*

### 2. Introduction (The Challenge)
*   **Metaphor**: Use a real-world analogy (e.g., "needle in a haystack" for Grover's).
*   **Quantum Advantage**: Explain why a quantum computer is better for this task (e.g., $O(\sqrt{N})$ vs $O(N)$).
*   **Preliminary Link**: Include at least one verified external link for high-level theory (e.g., Wikipedia).

### 3. Core Concepts
Define the foundational physics and math concepts used in the implementation.
*   **Definitions**: Provide clear, developer-focused definitions for terms like Superposition, Entanglement, Oracle, Phase Flip, etc.
*   **Geometric/Mathematical Intuition**: Use LaTeX notation for matrices and state vectors to provide rigorous context.
*   **Hilbert Space Context**: Explicitly link concepts to rotations or transformations in Hilbert Space where appropriate.

### 4. Walkthrough (Implementation Mapping)
A step-by-step breakdown of the corresponding Go test/example code.
*   **Code Snippets**: Use fenced code blocks with the `go` language identifier.
*   **Concept Mapping**: Explicitly explain which lines of code correspond to which part of the theoretical algorithm.
*   **Gate Usage**: Explain the role of specific gates (e.g., Hadamard for superposition, Oracle for marking).

### 5. Interpreting Results
Explain how to read the output of the simulation.
*   **Probability Peaks**: Explain why certain binary states have non-zero probabilities.
*   **Physical Meaning**: Link the measured bitstrings back to the original problem (e.g., "The peak at 3 represents the correct index").
*   **Risks/Pitfalls**: Document implementation risks like "Overcooking" (over-rotation) or insufficient precision.

### 6. Self-Assessment
3-5 multiple-choice or short-answer questions to help the learner verify their conceptual understanding.

### 7. Hands-on Exploration
Direct the user to run the example and provide CLI tool calls using the `strange` utility.
*   **Run the Example**: Specific command to run the Go test or example (e.g., `go test -v ./go/examples/...`).
*   **CLI Commands**: Show how to use the `strange` CLI tool to `export`, `run`, or `analyze` the circuit described in the guide.

### 8. References & Further Reading
A list of **verified** web-accessible resources. 
*   **Standards**: Wikipedia, Qiskit Textbook, Brilliant.org, ArXiv (Original Papers).
*   **Verification**: All links MUST be verified using `curl` or manual inspection before inclusion.

## Formatting Standards

*   **LaTeX**: Use single `$` for inline math (e.g., $| \psi \rangle$) and double `$$` for block math (e.g., matrices).
*   **Admonitions**: Use blockquotes or callouts for critical definitions (e.g., Hilbert Space).
*   **Links**: Use descriptive link text rather than raw URLs.

## Verification Checklist
- [ ] Are all core physics terms defined?
- [ ] Does the walkthrough match the latest code in `go/examples/`?
- [ ] Are matrices represented in LaTeX?
- [ ] Have all links been checked for 404s?
- [ ] Is there a "Hands-on Exploration" section with CLI commands?
- [ ] Is there a "Next Step" suggestion for the user to modify the code?
