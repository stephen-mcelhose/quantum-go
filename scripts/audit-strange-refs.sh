#!/usr/bin/env bash
# audit-strange-refs.sh
#
# Scans the repository for any occurrence of "strange" (case-insensitive) and
# validates that every remaining reference is a documented, intentional
# attribution to the upstream redfx-quantum/strange Java project.
#
# Exit 0  -> all occurrences are accounted-for upstream attributions.
# Exit 1  -> one or more occurrences are unrecognised and must be either renamed
#            to "quantum-go" or explicitly listed below.
#
# Usage:
#   ./scripts/audit-strange-refs.sh            # from repo root
#   ./scripts/audit-strange-refs.sh --verbose  # also print safe lines

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
VERBOSE=false
[[ "${1:-}" == "--verbose" ]] && VERBOSE=true

# ---------------------------------------------------------------------------
# Files that are entirely exempt from scanning.
# These exist solely to document or attribute the upstream Strange project.
# ---------------------------------------------------------------------------
KNOWN_UPSTREAM_FILES=(
  "LICENSE"
  "REFERENCES.md"
  "docs/adr/002-quantum-go-name-change.md"
  "scripts/audit-strange-refs.sh"
  "scripts/verify-cli.sh"   # meta-script: talks about 'strange' references by necessity
)

# ---------------------------------------------------------------------------
# Line-level safe patterns.
# A line containing any of these (case-insensitive) is a legitimate upstream
# attribution and is allowed to remain.
# ---------------------------------------------------------------------------
SAFE_LINE_PATTERNS=(
  "redfx-quantum"
  "Johan Vos"
  "original Strange Java simulator"   # upstream attribution comment in Go source
)

# ---------------------------------------------------------------------------
# Scan
# ---------------------------------------------------------------------------

UNKNOWN=0
SAFE=0

while IFS= read -r raw_line; do
  # grep -rn output: /abs/path/file:linenum:content
  rel="${raw_line#"$REPO_ROOT/"}"
  file="${rel%%:*}"
  rest="${rel#*:}"
  linenum="${rest%%:*}"
  content="${rest#*:}"

  # Skip entirely-exempt files
  exempt=false
  for ef in "${KNOWN_UPSTREAM_FILES[@]}"; do
    [[ "$file" == "$ef" ]] && { exempt=true; break; }
  done
  $exempt && continue

  # Check line-level safe patterns
  safe=false
  for pattern in "${SAFE_LINE_PATTERNS[@]}"; do
    echo "$content" | grep -qi "$pattern" && { safe=true; break; }
  done

  if $safe; then
    SAFE=$((SAFE + 1))
    $VERBOSE && printf "  OK   %s:%s\n       %s\n" "$file" "$linenum" "$content"
  else
    printf "UNKNOWN  %s:%s\n         %s\n" "$file" "$linenum" "$content"
    UNKNOWN=$((UNKNOWN + 1))
  fi

done < <(grep -rni "strange" \
  --include="*.go" \
  --include="*.md" \
  --include="*.py" \
  --include="*.sh" \
  --include="*.mod" \
  --include="*.txt" \
  --include="*.toml" \
  --include="*.yaml" \
  --include="*.yml" \
  --exclude-dir=".git" \
  --exclude-dir=".agents" \
  "$REPO_ROOT" 2>/dev/null || true)

echo
echo "---"
printf "Upstream attributions (safe):  %d\n" "$SAFE"
printf "Unrecognised references:       %d\n" "$UNKNOWN"
echo

if [[ "$UNKNOWN" -gt 0 ]]; then
  echo "FAIL -- $UNKNOWN unrecognised reference(s) must be renamed to 'quantum-go'"
  echo "or added to KNOWN_UPSTREAM_FILES / SAFE_LINE_PATTERNS with justification."
  echo "See docs/adr/002-quantum-go-name-change.md for the decision procedure."
  exit 1
fi

echo "PASS -- all 'strange' occurrences are documented upstream attributions."
exit 0
