package helps

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/router-for-me/CLIProxyAPI/v7/internal/responsestools"
	cliproxyexecutor "github.com/router-for-me/CLIProxyAPI/v7/sdk/cliproxy/executor"
	sdktranslator "github.com/router-for-me/CLIProxyAPI/v7/sdk/translator"
)

func TestQoderToolSearchPassThroughBoundaries(t *testing.T) {
	search := []byte(`{"input":[],"tools":[{"type":"tool_search","execution":"client"}]}`)
	ordinary := []byte(`{"input":"hello","tools":[{"type":"function","name":"read","parameters":{"type":"object"}}]}`)
	for _, tc := range []struct {
		name    string
		body    []byte
		enabled bool
		format  sdktranslator.Format
	}{
		{"disabled", search, false, sdktranslator.FormatOpenAIResponse},
		{"chat client", search, true, sdktranslator.FormatOpenAI},
		{"ordinary Responses", ordinary, true, sdktranslator.FormatOpenAIResponse},
	} {
		t.Run(tc.name, func(t *testing.T) {
			opts := cliproxyexecutor.Options{SourceFormat: tc.format, OriginalRequest: []byte("original")}
			req, got, attempt, err := PrepareQoderToolSearch(cliproxyexecutor.Request{Payload: tc.body}, opts, tc.enabled, nil)
			if err != nil || attempt != nil || !bytes.Equal(req.Payload, tc.body) || !bytes.Equal(got.OriginalRequest, opts.OriginalRequest) {
				t.Fatalf("pass-through changed: err=%v attempt=%v", err, attempt)
			}
		})
	}
}

func TestQoderToolSearchCancellationReleasesState(t *testing.T) {
	limits := responsestools.DefaultLimits()
	limiter := responsestools.NewLimiter(limits)
	req := cliproxyexecutor.Request{Payload: []byte(`{"input":[],"tools":[{"type":"tool_search","execution":"client"}]}`)}
	_, _, attempt, err := PrepareQoderToolSearch(req, cliproxyexecutor.Options{SourceFormat: sdktranslator.FormatOpenAIResponse}, true, limiter)
	if err != nil || attempt == nil {
		t.Fatalf("prepare: attempt=%v err=%v", attempt, err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	upstream := make(chan cliproxyexecutor.StreamChunk)
	go func() { <-ctx.Done(); close(upstream) }()
	stream := RestoreQoderToolSearchStream(ctx, cancel, &cliproxyexecutor.StreamResult{Chunks: upstream}, attempt)
	cancel()
	select {
	case _, ok := <-stream.Chunks:
		if ok {
			t.Fatal("cancelled stream emitted a chunk")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("cancelled stream did not close")
	}
	if n, b := limiter.Usage(); n != 0 || b != 0 {
		t.Fatalf("retained protocol state: %d attempts / %d bytes", n, b)
	}
}
