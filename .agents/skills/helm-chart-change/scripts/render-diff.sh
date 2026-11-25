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

# The chart generates passwords and keys on every offline render, and checksum
# annotations follow them. Mask both so only real changes show.
mask() {
  awk '
    /^---/                 { secret = 0; section = "" }
    /^kind: Secret[ \t]*$/ { secret = 1 }
    /^[A-Za-z]/            { section = $1 }
    secret && (section == "data:" || section == "stringData:") && /^  [^ ].*: / { sub(/: .*/, ": <secret>") }
    /checksum\/[A-Za-z0-9_.-]+: / { sub(/: .*/, ": <checksum>") }
    { print }
  '
}

render() {
  local dir="$1" out
  shift
  out=$(helm dependency build "$dir/$chart" 2>&1) || { echo "$out" >&2; return 1; }
  # Both sides use the working tree's values-example.yaml, so a change to it
  # shows as its effect on the rendered chart.
  helm template ci "$dir/$chart" -f "$chart/values-example.yaml" "${secrets[@]}" "$@" | mask
}

render "$tmp/base" "$@" > "$tmp/before.yaml"
render "." "$@" > "$tmp/after.yaml"

if diff -u --label "before ($base)" --label "after (working tree)" "$tmp/before.yaml" "$tmp/after.yaml"; then
  echo "No change in the rendered chart."
fi
