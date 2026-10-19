# Upstream protocol core

These Go sources and tests are imported without modification from
[`hrygo/CLIProxyAPI`, tag `v1.3.2-upstream8.0.23`](https://github.com/hrygo/CLIProxyAPI/tree/4edf1ea22abc7ef8e55f161319d3b208db66b034/internal/responsestools),
commit `4edf1ea22abc7ef8e55f161319d3b208db66b034`.
The upstream MIT license is retained in this directory. `UPSTREAM.sha256`
records all 52 imported Go files, with paths relative to the repository root.

The CPA Plus v7 integration lives in
`internal/runtime/executor/helps/qoder_tool_search.go`. Only opt-in Qoder
Responses requests carrying an explicit client-search contract enter the core.
No proxy listener, process, provider credential, or tool execution is owned by
this package. Future imports should preserve the upstream regression suite,
record the new upstream commit and refresh the checksum manifest.
