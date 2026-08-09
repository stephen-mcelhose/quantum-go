# Wiki Schema

This wiki is maintained by an LLM using the llm-wiki skill
(https://gist.github.com/karpathy/442a6bf555914893e9891c11519de94f).

## Domain

Quantum computing theory, gate-level simulation in Go, algorithm implementations, and the quantum-go codebase. Intended as a personal learning reference for Stephen.

Topics covered:
- Quantum gates (single-qubit, multi-qubit, parameterized, composite)
- Quantum algorithms (Grover's, Shor's, QFT, BB84, teleportation, error correction)
- State vector simulation internals (bitwise optimizations, Kronecker products)
- The quantum-go Go codebase — packages, DSL, engine, arithmetic, oracles
- Verification methodology (Qiskit, IBM Quantum, Quirk)
- Quantum information theory concepts (entanglement, superposition, measurement)

## Conventions

- **Page slugs**: kebab-case (e.g., `grovers-algorithm.md`)
- **Frontmatter**: OKF — `type` (default `concept`), `title`, `description`, `timestamp` (ISO-8601 UTC); optional `resource`, `tags`
- **Cross-references**: `[[Page Slug]]` wikilinks
- **Sources section**: every page ends with `## Sources` listing its raw inputs

## Operations

Run these via the `llm-wiki` skill:

- `ingest <source>` — read a new source, write a summary page, propagate to related pages
- `query <question>` — synthesize an answer from wiki pages, optionally write back
- `lint` — audit for orphans, contradictions, stale claims, missing links

## Raw Sources

Raw source files live in the quantum-go repository. They are immutable — the LLM reads them but never writes to them.

External sources (arXiv, GitHub, specs) are fetched with curl and read directly — no local raw copy is stored. The URL in `## Sources` is the canonical raw pointer. There is no `wiki/raw/` directory; external URLs serve that role.

## index.md

Structured catalog of all wiki pages. Updated on every write operation.

## log.md

Append-only chronological log. Format: `## [YYYY-MM-DD] operation | detail`
