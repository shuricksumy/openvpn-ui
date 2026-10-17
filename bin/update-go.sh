#!/usr/bin/env bash
# Bumps the minimum Go version in go.mod and updates all Go dependencies.
# The golang Docker images are pinned by digest and updated by Renovate.
# Usage: bin/update-go.sh [version]   (default: latest stable minor from go.dev)
set -euo pipefail

cd "$(dirname "$0")/.."

for cmd in curl go; do
	command -v "$cmd" >/dev/null || { echo "error: '$cmd' is required" >&2; exit 1; }
done

latest=$(curl -fsSL 'https://go.dev/VERSION?m=text' | head -n1)
version="${1:-${latest#go}}"
version="${version#go}"
# Require only the minor release (e.g. 1.27.0) so any 1.27.x builder image works
version="$(cut -d. -f1-2 <<<"$version").0"

current=$(go mod edit -json | sed -nE 's/^[[:space:]]*"Go": "([^"]+)".*/\1/p')
echo "Go: current $current, target $version"

go mod edit -go="$version"

# Lets an older local Go download the required toolchain
export GOTOOLCHAIN=auto

echo "Updating dependencies..."
go get -u ./...
go mod tidy
go build ./...

git --no-pager diff --stat -- go.mod go.sum || true
