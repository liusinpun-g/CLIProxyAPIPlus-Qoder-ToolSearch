# In-process Qoder tool-search bridge

This local integration is based on CPA Plus `v7.3.12-1` (`04048d8`). It
imports the protocol core from `hrygo/CLIProxyAPI`
`v1.3.2-upstream8.0.23`; provenance and the upstream MIT license are retained
under `internal/responsestools/`.

## Enablement and scope

```yaml
qoder:
  # Keep the deployment's existing vpc-endpoint, if any.
  tool-search: true
```

The switch defaults to false. It enables in-process adaptation only for
Qoder requests in the Responses/Codex dialect that carry a client-executed
tool-search contract. Ordinary requests and other client protocols keep the
existing translation path. This does not add a listener, process, HTTP proxy
hop, tool execution engine, account, or fallback provider.

The request adapter runs before Responses-to-Chat translation. It replaces
client search with an ordinary function, keeps deferred tools out of active
declarations, and reconstructs discovered definitions from client history.
The response adapter runs after the existing Chat-to-Responses translator,
restoring search calls, discovered namespaces, item IDs and streaming events.
Attempt state is bounded, rebuilt for each retry and released on completion,
errors and cancellation. V7 event framing is adapted at this boundary, without
changing shared translators.

The local Codex model-list handler advertises search support only when the
switch is enabled and every registered provider for the model is Qoder.
Home and mixed-provider catalogs do not inherit this local capability.
Clients using a separate static model catalog still need their own
`supports_search_tool: true` entry; enable it only for routes served by this
build with the switch on. No user's current catalog is changed by this patch.

Existing custom-tool translation is retained. The integration does not
silently strip tools or relax grammar/schema constraints. The imported core
limits and request-scoped compatibility errors remain in effect.

## Validation

New regression coverage exercises the actual Qoder executor and its encoded
outbound request, for streaming and non-streaming search/discovery rounds,
namespace identity, replayed search results, lease release, cancellation,
pass-through boundaries and model capability advertisement. The complete
upstream protocol regression suite is retained.

```sh
go test ./internal/responsestools ./internal/runtime/executor/helps ./internal/runtime/executor ./internal/client/codex/models ./sdk/api/handlers/openai
go test -race ./internal/responsestools ./internal/runtime/executor/helps ./internal/runtime/executor ./internal/client/codex/models -run 'TestQoderToolSearch|TestPrepare|TestStream' -count=1
go build -o cli-proxy-api ./cmd/server
```

Isolated live validation on 2026-10-19 used a single built executable directly
against the Qoder upstream, with a temporary copy of the existing credential.
The all-tools control and both non-streaming and SSE discovery/call/result
flows passed against Qwen Flash. A real ephemeral Codex CLI invocation also
discovered and executed a read-only synthetic MCP tool and returned its
runtime-generated value exactly. Temporary processes/configuration/credential
copies were removed; production configuration and original credential hashes
were unchanged. These checks do not certify every Qoder model, large plugin
catalogs, Desktop long sessions or production deployment.

The full repository test run reported one failure:
`TestPionMediaRelayBridgesAudioAndDataChannel` in
`internal/client/codex/live`. The same test also failed on the unmodified
`04048d8` baseline with `upstream DataChannel was not created`. It is separate
from this integration; the full suite is not reported as passing.
