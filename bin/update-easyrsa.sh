#!/usr/bin/env bash
# Updates the bundled EasyRSA tarball and every reference to its version.
# Usage: bin/update-easyrsa.sh [version]   (default: latest GitHub release)
set -euo pipefail

cd "$(dirname "$0")/.."

REPO="OpenVPN/easy-rsa"
ASSETS_DIR="build/assets"
FILES=(Dockerfile-latestvpn DockerfileGitHub "$ASSETS_DIR/openvpn-install-v2.sh")

for cmd in curl jq shasum perl; do
	command -v "$cmd" >/dev/null || { echo "error: '$cmd' is required" >&2; exit 1; }
done

if [[ -n "${1:-}" ]]; then
	api_url="https://api.github.com/repos/$REPO/releases/tags/v${1#v}"
else
	api_url="https://api.github.com/repos/$REPO/releases/latest"
fi

auth=()
[[ -n "${GITHUB_TOKEN:-}" ]] && auth=(-H "Authorization: Bearer $GITHUB_TOKEN")

release=$(curl -fsSL ${auth[@]+"${auth[@]}"} "$api_url")
version=$(jq -r '.tag_name | ltrimstr("v")' <<<"$release")
file="EasyRSA-${version}.tgz"
asset=$(jq -c --arg f "$file" '.assets[] | select(.name == $f)' <<<"$release")
[[ -n "$asset" ]] || { echo "error: release v$version has no $file asset" >&2; exit 1; }

current=$(sed -nE 's/.*local version="([0-9.]+)".*/\1/p' "$ASSETS_DIR/openvpn-install-v2.sh")
echo "EasyRSA: current $current, target $version"

if [[ ! -f "$ASSETS_DIR/$file" ]]; then
	tmp=$(mktemp)
	trap 'rm -f "$tmp"' EXIT
	curl -fsSL -o "$tmp" "$(jq -r '.browser_download_url' <<<"$asset")"

	expected=$(jq -r '.digest // empty | ltrimstr("sha256:")' <<<"$asset")
	actual=$(shasum -a 256 "$tmp" | cut -d' ' -f1)
	if [[ -z "$expected" ]]; then
		echo "warning: GitHub provided no digest; sha256 is $actual" >&2
	elif [[ "$expected" != "$actual" ]]; then
		echo "error: sha256 mismatch (expected $expected, got $actual)" >&2
		exit 1
	fi
	tar tzf "$tmp" >/dev/null
	mv "$tmp" "$ASSETS_DIR/$file"
	echo "Downloaded and verified $file"
fi

# Remove older bundled versions
find "$ASSETS_DIR" -maxdepth 1 -name 'EasyRSA-*.tgz' ! -name "$file" -print -delete

perl -pi -e "s/EasyRSA-[0-9.]+\.tgz/$file/g; s/local version=\"[0-9.]+\"/local version=\"$version\"/" "${FILES[@]}"

grep -nE "EasyRSA-[0-9.]+\.tgz|local version=" "${FILES[@]}"
