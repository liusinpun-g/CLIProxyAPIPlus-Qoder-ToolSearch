# v7.3.12-qoder-tools.1

Experimental prerelease of CLIProxyAPI Plus with opt-in Qoder tool-search adaptation. Independently maintained; not an official upstream release.

## Changes

- Adapt Codex / Responses client tool-search contracts inside the existing Qoder executor, without a separate relay process.
- Keep deferred tool declarations out of active requests until discovered; restore tool-search calls, namespaces, call IDs and SSE events on the client side.
- Add `qoder.tool-search`, disabled by default, and scoped Qoder capability advertisement.
- Preserve ordinary request handling when the feature is disabled or the request does not use client tool search.
- Add regression coverage for discovery/replay, streaming and non-streaming calls, cancellation, interrupted streams and state cleanup.
- Include source attribution, imported-code checksums and license notices in the distribution.

## Enablement

Run this build with `qoder.tool-search: true`. For clients using a static model catalog, also enable `supports_search_tool` for the matching Qoder model. Restart the client and start a new session. Existing listener, configuration and authentication can be reused. See `README.md` or `README_CN.md`.

## Validation

Qoder Qwen Flash passed isolated non-streaming and streaming discovery/call/result checks and a real ephemeral Codex CLI invocation against a synthetic read-only MCP tool. Affected-package tests and focused race tests passed. The live checks used a locally built macOS arm64 executable; the release packaging path builds a portable CGO-disabled executable separately.

The full repository suite has one known baseline failure: `TestPionMediaRelayBridgesAudioAndDataChannel` (`upstream DataChannel was not created`), reproduced on unmodified CPA Plus `04048d8`. Other models, large real-world tool catalogs and Desktop long sessions have not been live-validated. No fixed token-saving percentage is promised.

## Assets

The initial binary target is macOS arm64, packaged with MIT and third-party notices, dependency license files, example configuration and SHA-256 checksums. Other platforms can build from source; they are not claimed as live-tested release targets. CGO-disabled builds do not support dynamic library plugins.

## Sources and licensing

- Base: [kaitranntt/CLIProxyAPIPlus](https://github.com/kaitranntt/CLIProxyAPIPlus), `v7.3.12-1` / `04048d868e96a0450abae26728fef39c734c30d0`.
- Original architecture: [router-for-me/CLIProxyAPI](https://github.com/router-for-me/CLIProxyAPI).
- Imported protocol core: [hrygo/CLIProxyAPI](https://github.com/hrygo/CLIProxyAPI), `v1.3.2-upstream8.0.23` / `4edf1ea22abc7ef8e55f161319d3b208db66b034`.
- Integration: Sinpun and contributors.
- License: MIT, preserving upstream notices. Dependencies retain their own licenses. See `THIRD_PARTY_NOTICES.md`.
