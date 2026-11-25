#!/usr/bin/env bash
# Set the release version in every file that carries it.
# Usage: bump-version.sh <version>    (for example: bump-version.sh 0.7.0)
# Run from the repository root.
set -euo pipefail

version="${1:-}"
if [[ ! "$version" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  echo "usage: $0 <major.minor.patch>  (no leading v)" >&2
  exit 1
fi

chart=charts/opengovmail/Chart.yaml
images=(
  charts/opengovmail/charts/postfix/values.yaml
  charts/opengovmail/charts/opendkim/values.yaml
)

for f in "$chart" "${images[@]}"; do
  [[ -f "$f" ]] || { echo "missing $f: run from the repository root" >&2; exit 1; }
done

# perl -pi works the same with BSD (macOS) and GNU tools, unlike sed -i.
perl -pi -e "s/^version: .*/version: $version/; s/^appVersion: .*/appVersion: \"$version\"/" "$chart"

# Only the first `tag:` in each file, which is image.tag.
for f in "${images[@]}"; do
  perl -pi -e "if (!\$done && /^(\s+)tag: /) { s/^(\s+)tag: .*/\$1tag: \"$version\"/; \$done = 1 }" "$f"
done

grep -nE '^(version|appVersion):' "$chart"
for f in "${images[@]}"; do
  echo "$f: $(grep -m1 -E '^\s+tag:' "$f" | sed 's/^ *//')"
done
