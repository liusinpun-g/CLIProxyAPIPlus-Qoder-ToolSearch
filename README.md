# CLIProxyAPI Plus — Qoder Tool Search

English | [中文](README_CN.md)

A community fork of [kaitranntt/CLIProxyAPIPlus](https://github.com/kaitranntt/CLIProxyAPIPlus), based on `v7.3.12-1`, adding opt-in, in-process tool-search adaptation for Qoder models used through the Codex / OpenAI Responses protocol.

This fork is independently maintained. It is not an official release of CLIProxyAPI, CLIProxyAPI Plus, Qoder, or OpenAI. The feature is currently experimental.

## What this fork adds

Codex can defer tool definitions and discover them when needed instead of declaring the entire tool catalog at the beginning of a session. The adaptation runs inside the existing CPA Plus executable and uses its existing listener and Qoder authentication. No additional relay process is needed.

This can reduce tool-definition overhead. It does not increase a model's context window, eliminate the base system prompt, or guarantee a particular token reduction. Providers other than Qoder are outside the scope of this patch.

## Enable it

Build with Go 1.26 or later:

```sh
go build -o cli-proxy-api-plus ./cmd/server
```

Start from `config.example.yaml`, configure Qoder using the existing CPA Plus login/configuration flow, and add this key under the existing `qoder` block:

```yaml
qoder:
  tool-search: true
```

Preserve any existing endpoint and authentication settings. Start the service with your configuration:

```sh
./cli-proxy-api-plus -config config.yaml
```

The switch defaults to `false`. The Codex client must also support tool search. If the client uses a static model catalog, set `supports_search_tool: true` only for a Qoder route served by this build with the switch enabled. Restart the client and use a new session. A catalog manager may overwrite manual catalog edits when it regenerates the file.

See [integration details](docs/qoder-tool-search.md) and the [Chinese guide](README_CN.md). For inherited providers and configuration, see the [upstream README snapshot](README_UPSTREAM.md).

## Validation and limits

Validated with Qwen Flash through Qoder: non-streaming and streaming search/discovery/tool-result flows, a real ephemeral Codex CLI session with a synthetic read-only MCP tool, affected-package tests, and focused race tests.

The full test suite is **not** reported as passing: `TestPionMediaRelayBridgesAudioAndDataChannel` failed with `upstream DataChannel was not created`; the same failure was reproduced on the unmodified CPA Plus base. Other Qoder models, large real-world plugin catalogs, Desktop long sessions, and other operating-system binaries have not been live-validated for this feature.

## Credits and license

| Component | Source | Version / commit | Role |
| --- | --- | --- | --- |
| Original project | [router-for-me/CLIProxyAPI](https://github.com/router-for-me/CLIProxyAPI) | Inherited through CPA Plus | Proxy architecture and provider foundations |
| Direct base | [kaitranntt/CLIProxyAPIPlus](https://github.com/kaitranntt/CLIProxyAPIPlus) | `v7.3.12-1`, `04048d868e96a0450abae26728fef39c734c30d0` | CPA Plus v7 codebase and Qoder provider |
| Imported protocol core | [hrygo/CLIProxyAPI](https://github.com/hrygo/CLIProxyAPI) | `v1.3.2-upstream8.0.23`, `4edf1ea22abc7ef8e55f161319d3b208db66b034` | Unmodified `internal/responsestools/*.go` sources and tests |
| This fork | Sinpun and contributors | Qoder tool-search integration | Opt-in configuration, executor adapters, capability advertisement, tests and documentation |

Released under [MIT](LICENSE), retaining the upstream copyright and permission notices. Commercial use, modification, and redistribution are permitted under that license. Imported code remains attributed to its upstream authors; this fork does not claim authorship of the protocol core. See [third-party notices](THIRD_PARTY_NOTICES.md) and [import provenance](internal/responsestools/UPSTREAM.md). Go dependencies retain their own licenses.

## Releases

The first proposed release is `v7.3.12-qoder-tools.1`, an experimental prerelease based on CPA Plus `v7.3.12-1`. See the [release notes](RELEASE_NOTES.md) and [publishing guide](PUBLISHING.md). Upstream automation is preserved under `.github/upstream-workflows/` as inactive reference material; this fork does not automatically merge upstream or publish Docker images.
