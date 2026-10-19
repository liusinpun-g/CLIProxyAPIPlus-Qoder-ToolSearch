# Third-party notices and provenance

This community fork retains upstream copyright and permission notices. Project links identify source and authorship; they do not imply endorsement or support by those projects.

## Project ancestry

- Original project: https://github.com/router-for-me/CLIProxyAPI
- Direct fork base: https://github.com/kaitranntt/CLIProxyAPIPlus
- Base tag: `v7.3.12-1`
- Base commit: `04048d868e96a0450abae26728fef39c734c30d0`
- Base license: MIT, preserved in the root `LICENSE` along with the copyright notice for this fork's integration changes.
- The upstream license contains an unusual `2025-2005.9` year range. It is retained verbatim rather than silently corrected.
- `README_UPSTREAM.md` and `README_UPSTREAM_CN.md` preserve the original base READMEs. References to sponsors, maintainers, services, and releases in those snapshots describe the upstream project and are not sponsorship claims by this fork.

## Imported tool-search protocol core

- Source: https://github.com/hrygo/CLIProxyAPI
- Tag: `v1.3.2-upstream8.0.23`
- Commit: `4edf1ea22abc7ef8e55f161319d3b208db66b034`
- Source directory: https://github.com/hrygo/CLIProxyAPI/tree/4edf1ea22abc7ef8e55f161319d3b208db66b034/internal/responsestools
- Local directory: `internal/responsestools/`
- Imported material: all 52 Go source and test files in that directory, unchanged at import time.
- License: MIT. The upstream license text is preserved at `internal/responsestools/LICENSE`.
- Verification manifest: `internal/responsestools/UPSTREAM.sha256`, with paths relative to the repository root.

The imported core is credited to its upstream authors and contributors. The Qoder v7 adapter, configuration switch, capability advertisement, integration tests and documentation are this fork's integration work. The protocol core is not presented as original work of this fork.

## Dependency licenses and binary distribution

Go module dependencies retain their own licenses. The root MIT license does not relicense third-party dependencies. Binary release bundles include `third-party-licenses/`, collected from the actual target-platform `./cmd/server` dependency graph by `github.com/google/go-licenses/v2@v2.0.1`. This directory retains dependency license texts and applicable notices or source files. If the collection fails, publishing must stop rather than omit notices.

Binary archives also include the root MIT license, these notices, and the protocol core's MIT license and provenance. Keep these notices when redistributing substantial portions of the source or the executable.
