# ADR-001 — Use `vectorize` CLI as the embedding layer for semantic wiki search

**Date:** 2026-08  
**Status:** Accepted  
**Issue:** [#6](https://github.com/stephen-mcelhose/quantum-go/issues/6), closes [#7](https://github.com/stephen-mcelhose/quantum-go/issues/7)

---

## Decision

**In the context of** adding semantic goal resolution to `wikigraph goal` (issue #6), where we need to embed wiki pages into a vector space so a natural-language query can be resolved to relevant pages,

**facing the concern** that building a bespoke `wikigraph vectorize` subcommand would require choosing and integrating an embedding runtime (Ollama, ONNX, etc.), managing a vector store format, and maintaining that code long-term,

**we decided** to use the existing `vectorize` CLI (part of the csgdaa-code toolchain, installed at `/opt/homebrew/bin/vectorize`) as the embedding and search layer,

**to achieve** semantic nearest-neighbour lookup over wiki pages without writing or maintaining any embedding infrastructure,

**accepting** that `wikigraph goal --semantic` will have a hard dependency on `vectorize` being installed and `vectorize local-index` having been run against the wiki repo.

---

## Context

Issue #7 proposed a `wikigraph vectorize` subcommand that would:
- Call an Ollama HTTP endpoint to embed each page
- Write a bespoke `.vectors/index.json` file
- Require users to run `ollama serve` and pull a model

During review we discovered that `vectorize` (v0.2.21, Rust binary, shipped with csgdaa-code) already does all of this:

| Capability | `vectorize` |
| --- | --- |
| Chunk + embed files from a Git repo | `vectorize local-index` |
| Local SQLite vector store | `~/.local/share/csgdaa-code/vectorize/vectorize.db` |
| Semantic search with repo/workspace scope | `vectorize local-search "query" --repo <name>` |
| No external API key required | ✅ local-only path |
| Incremental re-index | `vectorize local-index` (skips unchanged) |

The wiki already lives in a Git repo (`quantum-go`), so it can be indexed directly.

---

## Consequences

- **`wikigraph goal --semantic`** will shell out to `vectorize local-search` to resolve a natural-language query to a ranked list of slugs, then feed those slugs as the MFPT target set.
- **Setup step** becomes `vectorize local-index --repo quantum-go` (or scoped to `wiki/`), documented in the runbook.
- **No bespoke embedding code** in this repo; model updates and storage improvements come for free via csgdaa-code releases.
- **Portability trade-off:** the semantic goal feature only works in environments where `vectorize` is installed. The plain slug-based `goal` subcommand remains fully self-contained.
