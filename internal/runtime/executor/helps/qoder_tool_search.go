package helps

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/router-for-me/CLIProxyAPI/v7/internal/responsestools"
	translatorcommon "github.com/router-for-me/CLIProxyAPI/v7/internal/translator/common"
	cliproxyexecutor "github.com/router-for-me/CLIProxyAPI/v7/sdk/cliproxy/executor"
	sdktranslator "github.com/router-for-me/CLIProxyAPI/v7/sdk/translator"
)

// PrepareQoderToolSearch adapts only explicit client-search contracts. Ordinary
// requests and other client protocols retain the existing Qoder translation.
// State belongs to this attempt and is reconstructed from history on each retry.
func PrepareQoderToolSearch(req cliproxyexecutor.Request, opts cliproxyexecutor.Options, enabled bool, limiter *responsestools.Limiter) (cliproxyexecutor.Request, cliproxyexecutor.Options, *responsestools.Attempt, error) {
	if !enabled || (opts.SourceFormat != sdktranslator.FormatOpenAIResponse && opts.SourceFormat != sdktranslator.FormatCodex) {
		return req, opts, nil, nil
	}
	contract := responsestools.ParseContract(req.Payload)
	if contract == nil || !contract.ClientSearch {
		return req, opts, nil, nil
	}
	policy := responsestools.ConventionPolicy(responsestools.Route{
		Provider: "qoder", AuthKind: "oauth", UpstreamModel: req.Model, UpstreamFormat: "openai",
	})
	prepared, err := responsestools.Prepare(req.Payload, policy, responsestools.DefaultLimits(), limiter)
	if err != nil {
		return req, opts, nil, err
	}
	req.Payload = prepared.Body
	// Existing response translators must resolve names against the prepared
	// wire declarations, before the core restores the original client identity.
	opts.OriginalRequest = prepared.Body
	return req, opts, prepared.Attempt, nil
}

// RestoreQoderToolSearchStream adapts Responses frames in process. There is no
// HTTP hop, secondary service, shared attempt history, or tool execution here.
func RestoreQoderToolSearchStream(ctx context.Context, cancel context.CancelFunc, upstream *cliproxyexecutor.StreamResult, attempt *responsestools.Attempt) *cliproxyexecutor.StreamResult {
	out := make(chan cliproxyexecutor.StreamChunk)
	go func() {
		defer close(out)
		defer attempt.Close()
		defer func() {
			cancel()
			// The underlying executor can be finishing an error send when the
			// adapter rejects a frame. Drain it after cancelling its HTTP request.
			for range upstream.Chunks {
			}
		}()
		send := func(chunk cliproxyexecutor.StreamChunk) bool {
			select {
			case out <- chunk:
				return true
			case <-ctx.Done():
				return false
			}
		}
		for {
			select {
			case <-ctx.Done():
				return
			case chunk, ok := <-upstream.Chunks:
				if !ok {
					frames, err := attempt.Finish()
					if err != nil {
						send(cliproxyexecutor.StreamChunk{Err: err})
						return
					}
					for _, frame := range frames {
						if !send(cliproxyexecutor.StreamChunk{Payload: qoderToolSearchEvent(frame)}) {
							return
						}
					}
					return
				}
				if chunk.Err != nil {
					send(chunk)
					return
				}
				// V7 translators emit a complete event without the SSE delimiter;
				// the HTTP handler normally adds it. The core consumes framed SSE.
				frame := chunk.Payload
				if bytes.HasPrefix(frame, []byte("event:")) || bytes.HasPrefix(frame, []byte("data:")) {
					frame = append(bytes.Clone(bytes.TrimRight(frame, "\r\n")), '\n', '\n')
				}
				frames, err := attempt.Feed(frame)
				if err != nil {
					send(cliproxyexecutor.StreamChunk{Err: err})
					return
				}
				for _, frame := range frames {
					if !send(cliproxyexecutor.StreamChunk{Payload: qoderToolSearchEvent(frame)}) {
						return
					}
				}
			}
		}
	}()
	return &cliproxyexecutor.StreamResult{Headers: upstream.Headers, Chunks: out}
}

// The core returns JSON event payloads; V7 handlers expect the event/data
// envelope and add the final delimiter themselves.
func qoderToolSearchEvent(frame []byte) []byte {
	var event struct {
		Type string `json:"type"`
	}
	if json.Unmarshal(frame, &event) == nil && event.Type != "" {
		return translatorcommon.SSEEventData(event.Type, frame)
	}
	return bytes.TrimRight(frame, "\r\n")
}
