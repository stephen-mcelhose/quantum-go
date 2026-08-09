#!/usr/bin/env bash
# verify-cli.sh
#
# End-to-end verification of the quantum-go CLI.
# Builds the binary then exercises every subcommand against every built-in
# circuit, plus inline QASM, QASM round-trip, binary identity, and the
# Qiskit comparison verifier (64 test cases, opt-in with --qiskit).
#
# Exit 0  -> all checks passed.
# Exit 1  -> one or more checks failed (failures printed inline).
#
# Usage:
#   ./scripts/verify-cli.sh              # builds from source, runs all checks
#   ./scripts/verify-cli.sh --no-build   # skip build (use existing binary)
#   ./scripts/verify-cli.sh --qiskit     # also run the Qiskit verification
#                                        # (requires venv/ with qiskit-aer)
#
# Prerequisites:
#   go 1.24+  (for build)
#   python3 + qiskit-aer in venv/  (for --qiskit)

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
BIN="$REPO_ROOT/quantum-go"
BUILD=true
RUN_QISKIT=false

for arg in "$@"; do
  case "$arg" in
    --no-build)  BUILD=false ;;
    --qiskit)    RUN_QISKIT=true ;;
  esac
done

PASS=0
FAIL=0

# ---------------------------------------------------------------------------
run() {
  local desc="$1"; shift
  if "$@" > /dev/null 2>&1; then
    printf "  PASS  %s\n" "$desc"
    PASS=$((PASS + 1))
  else
    printf "  FAIL  %s\n" "  -> $desc"
    printf "        cmd: %s\n" "$*"
    FAIL=$((FAIL + 1))
  fi
}

run_output() {
  # Like run, but also checks that stdout is non-empty
  local desc="$1"; shift
  local out
  if out=$("$@" 2>/dev/null) && [[ -n "$out" ]]; then
    printf "  PASS  %s\n" "$desc"
    PASS=$((PASS + 1))
  else
    printf "  FAIL  %s\n" "$desc"
    printf "        cmd: %s\n" "$*"
    FAIL=$((FAIL + 1))
  fi
}

run_contains() {
  # Passes if output contains the expected string
  local desc="$1"; local expected="$2"; shift 2
  local out
  out=$("$@" 2>&1) || true
  if echo "$out" | grep -q "$expected"; then
    printf "  PASS  %s\n" "$desc"
    PASS=$((PASS + 1))
  else
    printf "  FAIL  %s (expected '%s' in output)\n" "$desc" "$expected"
    FAIL=$((FAIL + 1))
  fi
}

run_excludes() {
  # Passes if output does NOT contain the excluded string (case-insensitive)
  local desc="$1"; local excluded="$2"; shift 2
  local out
  out=$("$@" 2>&1) || true
  if echo "$out" | grep -qi "$excluded"; then
    printf "  FAIL  %s (found '%s' in output)\n" "$desc" "$excluded"
    FAIL=$((FAIL + 1))
  else
    printf "  PASS  %s\n" "$desc"
    PASS=$((PASS + 1))
  fi
}
# ---------------------------------------------------------------------------

cd "$REPO_ROOT"

# All 17 built-in circuits
ALL_CIRCUITS="adder bell bernstein-vazirani deutsch-jozsa engine error-correction fredkin ghz grover qft qkd shor simon superdense superposition teleportation toffoli"

# ---------------------------------------------------------------------------
echo "=== Build ==="
if $BUILD; then
  if go build -o quantum-go ./cmd/quantum-go/main.go; then
    echo "  PASS  go build -o quantum-go ./cmd/quantum-go/main.go"
    PASS=$((PASS + 1))
  else
    echo "  FAIL  build failed — cannot continue"
    exit 1
  fi
else
  echo "  SKIP  --no-build"
  if [[ ! -x "$BIN" ]]; then
    echo "  ERROR binary not found at $BIN — build first or remove --no-build"
    exit 1
  fi
fi

# ---------------------------------------------------------------------------
echo
echo "=== Go unit tests ==="
if go test ./... > /tmp/qgo_test.out 2>&1; then
  echo "  PASS  go test ./..."
  PASS=$((PASS + 1))
else
  echo "  FAIL  go test ./..."
  cat /tmp/qgo_test.out
  FAIL=$((FAIL + 1))
fi

# ---------------------------------------------------------------------------
echo
echo "=== go vet ==="
run "go vet ./..." go vet ./...

# ---------------------------------------------------------------------------
echo
echo "=== Binary identity ==="
# Verify the binary identifies itself as quantum-go, never as strange
run_contains  "root -h mentions quantum-go"  "quantum-go"  "$BIN" -h
run_excludes  "root -h has no 'strange'"     "^strange"    "$BIN" -h
run_contains  "run -h mentions quantum-go"   "quantum-go"  "$BIN" run -h
run_contains  "module path is quantum-go"    "quantum-go"  go list -m

# ---------------------------------------------------------------------------
echo
echo "=== Help / discovery ==="
run_output "quantum-go -h"                  "$BIN" -h
run_output "quantum-go list circuits"       "$BIN" list circuits
run_output "quantum-go list gates"          "$BIN" list gates
run_output "quantum-go run -h"             "$BIN" run -h
run_output "quantum-go export -h"          "$BIN" export -h
run_output "quantum-go inspect -h"         "$BIN" inspect -h
run_output "quantum-go analyze -h"         "$BIN" analyze -h
run_output "quantum-go verify -h"          "$BIN" verify -h

# ---------------------------------------------------------------------------
echo
echo "=== run: all built-in circuits ==="
run_output "run bell"                        "$BIN" run --circuit bell
run_output "run ghz (3 qubits)"             "$BIN" run --circuit ghz -n 3
run_output "run ghz (5 qubits)"             "$BIN" run --circuit ghz -n 5
run_output "run teleportation"              "$BIN" run --circuit teleportation
run_output "run superdense"                 "$BIN" run --circuit superdense
run_output "run toffoli"                    "$BIN" run --circuit toffoli
run_output "run fredkin"                    "$BIN" run --circuit fredkin
run_output "run superposition (4 qubits)"   "$BIN" run --circuit superposition -n 4
run_output "run error-correction (bit=0)"   "$BIN" run --circuit error-correction -p bit=0
run_output "run error-correction (bit=1)"   "$BIN" run --circuit error-correction -p bit=1
run_output "run qkd (bit=0,basis=0)"        "$BIN" run --circuit qkd -p bit=0,basis=0
run_output "run qkd (bit=1,basis=1)"        "$BIN" run --circuit qkd -p bit=1,basis=1
run_output "run deutsch-jozsa (balanced)"   "$BIN" run --circuit deutsch-jozsa -p balanced=true
run_output "run deutsch-jozsa (constant)"   "$BIN" run --circuit deutsch-jozsa -p balanced=false
run_output "run bernstein-vazirani (s=101)" "$BIN" run --circuit bernstein-vazirani -p s=101
run_output "run bernstein-vazirani (s=1011)""$BIN" run --circuit bernstein-vazirani -p s=1011
run_output "run simon (s=11)"               "$BIN" run --circuit simon -p s=11
run_output "run adder (x=2,y=3)"            "$BIN" run --circuit adder -p x=2,y=3
run_output "run shor (default)"             "$BIN" run --circuit shor
run_output "run shor (mod=15,a=7)"          "$BIN" run --circuit shor -p mod=15,a=7,precision=8
run_output "run engine"                     "$BIN" run --circuit engine
run_output "run grover"                     "$BIN" run --circuit grover
run_output "run qft (3 qubits)"             "$BIN" run --circuit qft -n 3

# ---------------------------------------------------------------------------
echo
echo "=== run: inline QASM steps ==="
run_output "run inline: H gate"             "$BIN" run -n 1 -s "h q[0]"
run_output "run inline: X gate"             "$BIN" run -n 1 -s "x q[0]"
run_output "run inline: bell (steps)"       "$BIN" run -n 2 -s "h q[0]" -s "cx q[0], q[1]"
run_output "run inline: GHZ (steps)"        "$BIN" run -n 3 -s "h q[0]" -s "cx q[0], q[1]" -s "cx q[1], q[2]"
run_output "run inline: phase rotation"     "$BIN" run -n 1 -s "u1(0.785398) q[0]"
run_output "run inline: RX gate"            "$BIN" run -n 1 -s "rx(3.141593) q[0]"
run_output "run inline: multi-gate step"    "$BIN" run -n 2 -s "h q[0]" -s "u1(1.5708) q[0]"

# ---------------------------------------------------------------------------
echo
echo "=== run: JSON output ==="
run_output "run bell --json"                "$BIN" run --circuit bell --json
run_output "run ghz --json"                 "$BIN" run --circuit ghz -n 3 --json
run_output "run inline bell --json"         "$BIN" run -n 2 -s "h q[0]" -s "cx q[0], q[1]" --json

# ---------------------------------------------------------------------------
echo
echo "=== QASM round-trip (bell, ghz, toffoli) ==="
for rt_circuit in bell ghz toffoli; do
  TMPQASM="$(mktemp /tmp/qgo_XXXXXX.qasm)"
  if "$BIN" export --circuit "$rt_circuit" > "$TMPQASM" 2>/dev/null; then
    run_output "export $rt_circuit → QASM file"  cat "$TMPQASM"
    run_output "run from QASM ($rt_circuit)"     "$BIN" run "$TMPQASM"
    run_output "inspect from QASM ($rt_circuit)" "$BIN" inspect "$TMPQASM"
  else
    echo "  FAIL  export $rt_circuit (could not create QASM)"
    FAIL=$((FAIL + 1))
  fi
  rm -f "$TMPQASM"
done

# ---------------------------------------------------------------------------
echo
echo "=== export: all 17 built-in circuits ==="
for circuit in $ALL_CIRCUITS; do
  run_output "export $circuit" "$BIN" export --circuit "$circuit"
done

# ---------------------------------------------------------------------------
echo
echo "=== inspect: all 17 built-in circuits ==="
for circuit in $ALL_CIRCUITS; do
  run_output "inspect $circuit" "$BIN" inspect --circuit "$circuit"
done
run_output "inspect inline GHZ"  "$BIN" inspect -n 3 -s "h q[0]" -s "cx q[0], q[1]" -s "cx q[1], q[2]"

# ---------------------------------------------------------------------------
echo
echo "=== analyze: all 17 built-in circuits ==="
for circuit in $ALL_CIRCUITS; do
  run_output "analyze $circuit" "$BIN" analyze --circuit "$circuit"
done
run_output "analyze inline bell"     "$BIN" analyze -n 2 -s "h q[0]" -s "cx q[0], q[1]"
run_output "analyze inline GHZ"      "$BIN" analyze -n 3 -s "h q[0]" -s "cx q[0], q[1]" -s "cx q[1], q[2]"

# ---------------------------------------------------------------------------
echo
echo "=== verify: theoretical mode ==="
# Only circuits that have a built-in theoretical reference (bell, ghz, qft)
run "verify bell (theoretical)"      "$BIN" verify --circuit bell --mode theoretical
run "verify ghz 3 (theoretical)"     "$BIN" verify --circuit ghz -n 3 --mode theoretical
run "verify ghz 4 (theoretical)"     "$BIN" verify --circuit ghz -n 4 --mode theoretical
run "verify qft 4 (theoretical)"     "$BIN" verify --circuit qft -n 4 --mode theoretical

# ---------------------------------------------------------------------------
echo
echo "=== verify: file mode (JSON round-trip) ==="
# file mode requires inline steps (--step), not --circuit; generate reference then re-verify
TMPJSON="$(mktemp /tmp/qgo_XXXXXX.json)"
if "$BIN" run -n 2 -s "h q[0]" -s "cx q[0], q[1]" --json > "$TMPJSON" 2>/dev/null; then
  run "verify inline bell (file mode)" \
    "$BIN" verify -n 2 -s "h q[0]" -s "cx q[0], q[1]" --mode file --reference "$TMPJSON"
else
  echo "  SKIP  verify file mode (inline bell JSON failed)"
fi
rm -f "$TMPJSON"

TMPJSON="$(mktemp /tmp/qgo_XXXXXX.json)"
if "$BIN" run -n 3 -s "h q[0]" -s "cx q[0], q[1]" -s "cx q[1], q[2]" --json > "$TMPJSON" 2>/dev/null; then
  run "verify inline GHZ (file mode)" \
    "$BIN" verify -n 3 -s "h q[0]" -s "cx q[0], q[1]" -s "cx q[1], q[2]" --mode file --reference "$TMPJSON"
else
  echo "  SKIP  verify file mode (inline GHZ JSON failed)"
fi
rm -f "$TMPJSON"

# ---------------------------------------------------------------------------
echo
echo "=== Qiskit verification (64 test cases) ==="
if $RUN_QISKIT; then
  if [[ -f "$REPO_ROOT/venv/bin/activate" ]]; then
    # shellcheck disable=SC1091
    source "$REPO_ROOT/venv/bin/activate"
    if python3 verification/verify_against_qiskit.py "$BIN" > /tmp/qgo_qiskit.out 2>&1; then
      echo "  PASS  verify_against_qiskit.py (64 tests)"
      PASS=$((PASS + 1))
    else
      echo "  FAIL  verify_against_qiskit.py"
      tail -30 /tmp/qgo_qiskit.out
      FAIL=$((FAIL + 1))
    fi
  else
    echo "  SKIP  venv/ not found"
    echo "        Setup: python3 -m venv venv && source venv/bin/activate && pip install qiskit qiskit-aer"
    echo "        Then:  bash scripts/verify-cli.sh --qiskit"
  fi
else
  echo "  SKIP  (pass --qiskit to enable — requires venv/ with qiskit-aer)"
fi

# ---------------------------------------------------------------------------
echo
echo "=== Audit: no stray 'strange' references ==="
if bash "$REPO_ROOT/scripts/audit-strange-refs.sh" > /tmp/qgo_audit.out 2>&1; then
  echo "  PASS  audit-strange-refs.sh (0 unrecognised refs)"
  PASS=$((PASS + 1))
else
  echo "  FAIL  audit-strange-refs.sh"
  grep "^UNKNOWN" /tmp/qgo_audit.out | head -20
  FAIL=$((FAIL + 1))
fi

# ---------------------------------------------------------------------------
echo
echo "==================================================================="
printf "PASS: %d   FAIL: %d\n" "$PASS" "$FAIL"
echo "==================================================================="

[[ "$FAIL" -eq 0 ]]
