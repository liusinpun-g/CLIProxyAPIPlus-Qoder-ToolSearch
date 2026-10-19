#!/usr/bin/env bash
set -euo pipefail

# Build only the reviewed, clean tag; this script never publishes artifacts.
tag="${1:?Usage: package-qoder-release.sh vX.Y.Z-qoder-tools.N}"
[[ "$tag" =~ ^v[0-9]+\.[0-9]+\.[0-9]+-qoder-tools\.[1-9][0-9]*$ ]] || { echo "Invalid release tag" >&2; exit 1; }
cd "$(git rev-parse --show-toplevel)"
[[ -z "$(git status --porcelain --untracked-files=all)" ]] || { echo "Release checkout must be clean" >&2; exit 1; }
commit="$(git rev-parse HEAD)"
[[ "$(git rev-parse "refs/tags/$tag^{commit}")" = "$commit" ]] || { echo "Tag does not match HEAD" >&2; exit 1; }
shasum -a 256 -c internal/responsestools/UPSTREAM.sha256

work_dir="$(mktemp -d)"
trap 'rm -rf "$work_dir"' EXIT
package_name="CLIProxyAPIPlus-Qoder-ToolSearch_${tag}_darwin_arm64_no-plugin"
package_dir="$work_dir/$package_name"
mkdir -p "$package_dir" dist
[[ ! -e "dist/$package_name.tar.gz" ]] || { echo "Archive already exists" >&2; exit 1; }

# Install the collector for the host, then resolve the binary's target graph.
GOBIN="$work_dir/tools" go install github.com/google/go-licenses/v2@v2.0.1
GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 "$work_dir/tools/go-licenses" save ./cmd/server --save_path="$package_dir/third-party-licenses"
GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -trimpath \
  -ldflags "-s -w -X main.Version=$tag -X main.Commit=$commit -X main.BuildDate=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
  -o "$package_dir/cli-proxy-api-plus" ./cmd/server
cp LICENSE README.md README_CN.md README_UPSTREAM.md README_UPSTREAM_CN.md PUBLISHING.md THIRD_PARTY_NOTICES.md RELEASE_NOTES.md config.example.yaml "$package_dir/"
mkdir -p "$package_dir/internal/responsestools" "$package_dir/docs"
cp internal/responsestools/LICENSE internal/responsestools/UPSTREAM.md internal/responsestools/UPSTREAM.sha256 "$package_dir/internal/responsestools/"
cp docs/qoder-tool-search.md "$package_dir/docs/"
printf 'tag=%s\ncommit=%s\ngoos=darwin\ngoarch=arm64\ncgo_enabled=0\n' "$tag" "$commit" > "$package_dir/BUILD_INFO.txt"
tar -czf "dist/$package_name.tar.gz" -C "$work_dir" "$package_name"
(cd dist && shasum -a 256 "$package_name.tar.gz" > checksums.txt)
printf 'Created dist/%s.tar.gz\n' "$package_name"
