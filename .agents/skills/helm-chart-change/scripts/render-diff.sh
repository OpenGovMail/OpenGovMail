#!/usr/bin/env bash
# Show how the rendered chart changes between a base ref and the working tree.
# Usage: render-diff.sh [base-ref] [extra helm template args...]
#   base-ref defaults to upstream/main. Run from the repository root.
# The base is checked out in a temporary worktree, so the working tree and the
# stash are never touched. Committed, uncommitted and untracked changes all count.
set -euo pipefail

base="${1:-upstream/main}"
shift || true

chart=charts/opengovmail
[[ -f "$chart/Chart.yaml" ]] || { echo "run from the repository root" >&2; exit 1; }
git rev-parse --verify --quiet "$base^{commit}" >/dev/null || { echo "unknown ref: $base" >&2; exit 1; }

# Placeholders for the secrets a first install requires, as in CI.
secrets=(
  --set global.ravenIdp.clientSecret=ci
  --set global.ravenIdp.directAuthSecret=ci
  --set thunderid.configuration.server.security.directAuthSecret=ci
)

tmp=$(mktemp -d)
cleanup() { git worktree remove --force "$tmp/base" >/dev/null 2>&1 || true; rm -rf "$tmp"; }
trap cleanup EXIT
git worktree add --quiet --detach "$tmp/base" "$base"

# Prefix each rendered line with a stable id: its template, which document of
# that template, and its key (or its position, for lines without a key).
annotate() {
  awk '
    /^---/              { start = NR; next }
    /^# Source: /       { src = $3; doc = ++seen[src]; start = NR }
    {
      if (match($0, /^ *-? *[^ :]+:/)) key = substr($0, 1, RLENGTH)
      else key = "line " (NR - start)
      print src "#" doc "|" key "\t" $0
    }
  '
}

# Render one tree twice with the same input. The chart generates passwords and
# keys on every offline render, so lines that differ between the two renders
# are random. Print the render, and record those lines' ids.
render() {
  local dir="$1" name="$2" out
  shift 2
  out=$(helm dependency build "$dir/$chart" 2>&1) || { echo "$out" >&2; return 1; }
  # Each side uses its own values-example.yaml, so edits to it show up.
  local args=(ci "$dir/$chart" -f "$dir/$chart/values-example.yaml" "${secrets[@]}" "$@")
  helm template "${args[@]}" | annotate > "$tmp/$name.1"
  helm template "${args[@]}" | annotate > "$tmp/$name.2"
  paste -d '\n' "$tmp/$name.1" "$tmp/$name.2" | awk -F '\t' '
    NR % 2 { id = $1; line = $0; next }
    $0 != line { print id }
  ' >> "$tmp/volatile"
}

: > "$tmp/volatile"
render "$tmp/base" before "$@"
render "." after "$@"

# Mask only the values that change between identical renders. Fixed settings,
# including fixed values in Secrets, stay visible.
mask() {
  awk -F '\t' '
    NR == FNR { volatile[$0] = 1; next }
    {
      line = substr($0, length($1) + 2)
      if ($1 in volatile) {
        if (match(line, /^ *-? *[^ :]+: /)) line = substr(line, 1, RLENGTH) "<generated>"
        else line = "<generated>"
      }
      print line
    }
  ' "$tmp/volatile" "$1"
}

mask "$tmp/before.1" > "$tmp/before.yaml"
mask "$tmp/after.1" > "$tmp/after.yaml"

if diff -u --label "before ($base)" --label "after (working tree)" "$tmp/before.yaml" "$tmp/after.yaml"; then
  echo "No change in the rendered chart."
fi
