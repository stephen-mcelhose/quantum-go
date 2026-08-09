# ADR-002 — Rename project identity from "Strange" to "quantum-go"

**Date:** 2025-08  
**Status:** Accepted  

---

## Decision

**In the context of** this project being a Go port of the Java quantum simulator
[Strange](https://github.com/redfx-quantum/strange) by Johan Vos, where the
codebase had inherited the "Strange" name wholesale — as the CLI binary name
(`strange`), the root Go package name (`package strange`), doc titles
("Strange-Go"), and throughout all internal documentation —

**facing the concern** that continuing to call ourselves "Strange" blurs the
line between our project and the upstream Java project we ported from, creates
confusion about ownership and maintenance, and makes it impossible to evolve
our identity, API, or behaviour independently without appearing to misrepresent
the original author's work,

**we decided** to rename all *internal* references from "Strange" / "Strange-Go"
to **quantum-go**, matching the Go module path already in use
(`github.com/stephen-mcelhose/quantum-go`), while explicitly preserving every
*external* attribution to the upstream project,

**to achieve** a clear, honest separation: `quantum-go` is our Go simulator that
originated as a port of Strange; Strange is Johan Vos's Java project. They are
now distinct named things that can be discussed without ambiguity.

**accepting** that this is a breaking change for anything that scripted against
the `strange` binary name, and that it requires a systematic sweep of every file
in the repository (tracked in the companion GitHub issue; verified by
`scripts/audit-strange-refs.sh`).

---

## What "rename" means in practice

| Reference type                                                                 | Action                                                  |
| :----------------------------------------------------------------------------- | :------------------------------------------------------ |
| CLI binary name / cobra `Use:` field (`strange`)                               | → `quantum-go`                                          |
| Project display name ("Strange-Go", "Strange Quantum Simulator")               | → `quantum-go`                                          |
| Root Go package declaration (`package strange`)                                | → `package quantum_go`                                  |
| `cmd/strange/` directory                                                       | → `cmd/quantum-go/`                                     |
| Build instructions (`go build -o strange …`)                                   | → `go build -o quantum-go …`                            |
| Doc titles, headings, prose referring to our tool                              | → `quantum-go`                                          |
| Variable/key names scoped to our tool (`strange_path`, `run_strange`, etc.)   | → rename (`quantum_go_path`, `run_quantum_go`, etc.)    |
| **URLs to upstream (`redfx-quantum/strange`)**                                 | **KEEP — attribution**                                  |
| **"ported from Strange" / "derived from Strange" prose**                       | **KEEP — attribution**                                  |
| **Johan Vos copyright line in `LICENSE`**                                      | **KEEP — required by BSD 3-Clause**                     |
| **Go source comments explaining upstream Strange design decisions**             | **KEEP, or clarify with "In the upstream Strange Java simulator…"** |

---

## Risk accepted

The name change is a deliberate clean break:

- **Binary rename** — any existing scripts calling `./strange` will break. This
  is acceptable because the project has not been published as a stable release.
- **Package name** — `package strange` at the repo root becomes
  `package quantum_go`. The root package is documentation-only; no importers
  reference it by package name.
- **Documentation churn** — 48 files across `docs/`, `wiki/`, `examples/`,
  `verification/`, and source contain "strange". The companion issue provides a
  per-file checklist and an acceptance-criterion verifier script.
- **Wiki re-ingest** — every wiki page synthesised from an updated source file
  must be re-ingested via the `llm-wiki` skill to reflect the new name. This is
  tracked in the issue checklist.

---

## License acknowledgement

`quantum-go` is a derivative work of Strange under BSD 3-Clause. That licence
requires that redistributions retain the original copyright notice. The `LICENSE`
file therefore carries **both** copyrights:

1. Original Strange copyright — **Johan Vos, 2018, 2023**
2. Go port copyright — **quantum-go contributors**

Published documentation must include a prominent notice that quantum-go is a Go
port of [redfx-quantum/strange](https://github.com/redfx-quantum/strange).

---

## Consequences

- All internal "strange" references are renamed to "quantum-go" per the table
  above (tracked in the companion GitHub issue).
- `cmd/strange/` is renamed to `cmd/quantum-go/`.
- `scripts/audit-strange-refs.sh` becomes a permanent, CI-ready check: any
  remaining "strange" occurrence is either a documented upstream attribution or
  a regression to fix.
- Future contributors see a consistent name throughout the codebase, with
  explicit attribution to the upstream project wherever credit is due.
