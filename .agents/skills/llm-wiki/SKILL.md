---
name: llm-wiki
description: >
  Maintain the quantum-go learning wiki — a compounding knowledge base of interlinked
  Markdown files in wiki/ following the Karpathy pattern. The LLM acts as programmer;
  the wiki is the codebase. Use this skill to ingest docs, source files, and examples
  into the wiki, query it for synthesized answers, or lint it for orphans and stale
  claims. Trigger on: "/llm-wiki", "add this to the wiki", "ingest this",
  "update the wiki with", "query the wiki about", "lint the wiki", "what does my wiki
  say about", or when the user shares a file alongside any mention of the wiki.
version: "1.0.0"
allowed-tools:
  - Read
  - Write
  - Edit
  - Bash
  - Glob
  - Grep
  - WebFetch
  - WebSearch
  - AskUserQuestion
---

# llm-wiki — quantum-go Learning Wiki

A compounding personal knowledge base for the `quantum-go` simulator. Three layers:

1. **Raw sources** (immutable) — the repo's Go source files, Markdown docs, and example
   guides. The LLM reads; never writes to these.
2. **Wiki** — interlinked `.md` files in `wiki/`, synthesized and maintained by the LLM.
3. **Schema** — `wiki/AGENTS.md` defines domain conventions. **Edit this file** to
   customize how the wiki categorizes and links quantum-go concepts.

> "Obsidian is the IDE; the LLM is the programmer; the wiki is the codebase." — Karpathy

Reference: https://gist.github.com/karpathy/442a6bf555914893e9891c11519de94f

---

## Repo-Specific Configuration

| Setting       | Value                                                     |
|---------------|-----------------------------------------------------------|
| **WIKI_ROOT** | `wiki/` (relative to repo root — always use this path)   |
| **Raw sources**| `docs/`, `examples/`, `core/`, `local/`, `math/`, `qasm/`, `verification/` |
| **Domain**    | Quantum computing theory + Go simulator implementation    |

**WIKI_ROOT is fixed for this repo.** Do not ask the user where the wiki lives — it is
always `wiki/` at the repo root. Skip Step 0 of the standard flow.

---

## Step 0 — SKIP (wiki root is fixed)

The wiki always lives at `wiki/` in the repo root. Proceed directly to the requested
operation. If `wiki/index.md` does not exist, run the **Init flow** first.

---

## Init flow — bootstrapping the wiki

1. Create the `wiki/` directory.
2. Write `wiki/AGENTS.md` using the **quantum-go domain schema** below.
3. Write `wiki/index.md`:

```markdown
# quantum-go Wiki Index

| Page | Summary | Tags |
| ---- | ------- | ---- |
```

4. Write `wiki/log.md`:

```markdown
# Wiki Log

<!-- Append-only. Never edit existing entries. -->
```

5. Append the init entry to `wiki/log.md`:

```
## [YYYY-MM-DD] init | wiki initialized at wiki/
```

---

## quantum-go Domain Schema (write to wiki/AGENTS.md on init)

```markdown
# Wiki Schema — quantum-go

This wiki is maintained by an LLM using the llm-wiki skill.
Edit this file to adjust how concepts are categorized and linked.

## Domain

Quantum computing theory and the `quantum-go` simulator — a Go port of the
[Strange](https://github.com/redfx-quantum/strange) Java simulator. Topics:

- Quantum gates (fundamental, composite, parameterized)
- Quantum algorithms (Shor's, Grover's, QFT, Simon's, Deutsch-Jozsa, BB84)
- State vector simulation and optimization techniques
- Quantum error correction and thermodynamics
- Interoperability with Qiskit, IBM Quantum, OpenQASM

Intended audience: Stephen learning quantum computing via the quantum-go codebase.

## Tag Vocabulary

Use these tags consistently so related pages cluster correctly:

| Tag              | Use for                                              |
|------------------|------------------------------------------------------|
| `gates`          | Individual gate definitions and matrix forms         |
| `algorithms`     | End-to-end quantum algorithms                        |
| `simulation`     | How the simulator engine works internally            |
| `linear-algebra` | Matrix math, Kronecker products, complex numbers     |
| `entanglement`   | Bell states, GHZ, qubit correlation                  |
| `verification`   | Correctness checking against Qiskit/IBM/Quirk        |
| `architecture`   | Package structure, code design decisions             |
| `pedagogy`       | Explanatory guides bridging theory to code           |
| `thermodynamics` | Quantum energy, entropy, heat engines                |
| `cryptography`   | BB84, QKD, security protocols                        |
| `interop`        | OpenQASM, Qiskit bridge, IBM Quantum                 |

## Conventions

- **Page slugs**: kebab-case (e.g., `shors-algorithm.md`, `gate-application.md`)
- **Frontmatter**: OKF — required: `type` (default `concept`), `title`, `description`,
  `timestamp` (ISO-8601 UTC); optional: `resource`, `tags`
- **Cross-references**: `[[Page Slug]]` wikilinks — never relative paths
- **Sources section**: every page ends with `## Sources` listing its raw inputs
- **Code snippets**: include short Go examples where they illustrate a concept

## Raw Source Locations

Raw sources are immutable. The LLM reads them; never writes to them.

| Directory          | Contents                                    |
|--------------------|---------------------------------------------|
| `docs/`            | Concept guides, verification, optimization  |
| `examples/*/`      | Pedagogical algorithm walkthroughs          |
| `core/`            | Gate and circuit Go source                  |
| `local/`           | Simulation engine Go source                 |
| `math/`            | Linear algebra Go source                    |
| `qasm/`            | OpenQASM parser Go source                   |
| `verification/`    | Qiskit verification reports and scripts     |

## index.md

Structured catalog of all wiki pages. Updated on every write operation.

## log.md

Append-only chronological log. Format: `## [YYYY-MM-DD] operation | detail`
```

---

## Operation: ingest

**Trigger phrases:** "ingest this", "add this to the wiki", "add to my wiki",
"update the wiki with", user shares a file + mentions wiki.

### Steps

1. **Read the source.** If it's a file path, `Read` it. If it's a URL, `WebFetch` it.
   Never write to the source.

2. **Surface key takeaways in chat** — 3–5 ideas from this source. The user can redirect
   before any writes happen.

3. **Determine the slug.** Derive a kebab-case slug from the source topic.
   File becomes `wiki/<slug>.md`.

4. **Write the wiki page** with OKF frontmatter:

   ```markdown
   ---
   type: concept
   title: <Title>
   description: <one sentence — what this page covers and why it matters>
   resource: <source file path>
   tags: [<from tag vocabulary in AGENTS.md>]
   timestamp: <ISO-8601 UTC>
   ---

   # <Title>

   <2–4 paragraphs synthesizing the source. Cross-reference related pages with
   [[Wikilinks]]. Focus on insight over transcription. Include short Go snippets
   where they concretely illustrate a concept.>

   ## Key Points

   - …

   ## Sources

   - `<source file path>`
   ```

5. **Propagate to related pages.** Read `wiki/index.md` to identify pages whose topics
   this source touches. For each:
   - Read the page.
   - Add new information, confirm or flag contradictions, add cross-references.
   - A single ingest commonly touches 3–10 pages.

6. **Update `wiki/index.md`** — add a row for the new page:

   ```markdown
   | [[<slug>]] | <one-line summary> | <tags> |
   ```

7. **Append to `wiki/log.md`:**

   ```
   ## [YYYY-MM-DD] ingest | <title>
   ```

---

## Operation: query

**Trigger phrases:** "query the wiki", "what does my wiki say about", "ask the wiki",
"search the wiki for".

### Steps

1. Read `wiki/index.md`. Identify 3–6 pages most relevant to the question.
2. Read those pages. Follow `[[Wikilinks]]` if relevant.
3. Synthesize a cited answer. Trace every claim to a wiki page:
   > <Answer.> ([[Page Name]])
4. If synthesis produces a genuinely new insight not in the wiki, write it as
   `wiki/synthesis-<topic>-<YYYY-MM-DD>.md`.
5. Append to `wiki/log.md`:
   ```
   ## [YYYY-MM-DD] query | <question truncated to 80 chars>
   ```

---

## Operation: lint

**Trigger phrases:** "lint the wiki", "clean up the wiki", "audit the wiki".

### Steps

1. `Glob("wiki/**/*.md")` to inventory all pages. Read `wiki/index.md`.
2. Check each page for:
   - **Orphans** — no inbound `[[wikilink]]` from any other page
   - **Contradictions** — claims conflicting with another page
   - **Missing cross-references** — concept mentioned but not linked
   - **Index gaps** — page not in `wiki/index.md`
   - **OKF frontmatter gaps** — missing `type`, `title`, `description`, `timestamp`
3. Fix what you can. Mark unresolvable contradictions with:
   > ⚠️ Contradiction with [[Other Page]] — needs resolution.
4. Append to `wiki/log.md`:
   ```
   ## [YYYY-MM-DD] lint | <N pages checked, M issues found, K fixed>
   ```

---

## Wiki Conventions

| Convention                   | Rule                                                               |
|------------------------------|--------------------------------------------------------------------|
| **Page slugs**               | `kebab-case.md`                                                    |
| **Frontmatter**              | OKF: `type`, `title`, `description`, `timestamp`; optional `resource`, `tags` |
| **Cross-references**         | `[[Page Slug]]` wikilinks only — never relative file paths         |
| **Sources section**          | Every page ends with `## Sources`                                  |
| **Raw sources**              | Read-only. Never edit Go source or Markdown docs in the repo       |
| **index.md**                 | Updated on every write                                             |
| **log.md**                   | Append-only                                                        |
| **Synthesis over transcript**| Pages integrate and connect — they don't just copy a source        |
| **Go snippets welcome**      | Short code examples are encouraged to ground abstract concepts     |

---

## Rules

- **Never modify raw sources.** Repo source files are inputs, never outputs.
- **Propagate aggressively.** An ingest is a refactor — ripple changes to all affected pages.
- **Index is always current.** Update `wiki/index.md` before finishing any write operation.
- **Log is append-only.** Never rewrite or delete log entries.
- **Cite in every page.** Every claim should trace to a source.
- **Synthesize, don't transcribe.** A page that just copies a source verbatim has no value.
- **WIKI_ROOT is always `wiki/`.** Never create wiki files outside this directory.
