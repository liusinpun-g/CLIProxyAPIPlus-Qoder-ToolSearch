package executor

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/router-for-me/CLIProxyAPI/v7/internal/config"
	cliproxyauth "github.com/router-for-me/CLIProxyAPI/v7/sdk/cliproxy/auth"
	cliproxyexecutor "github.com/router-for-me/CLIProxyAPI/v7/sdk/cliproxy/executor"
	sdktranslator "github.com/router-for-me/CLIProxyAPI/v7/sdk/translator"
)

func decodeQoderToolSearchRequest(t *testing.T, r *http.Request) map[string]any {
	t.Helper()
	raw, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatal(err)
	}
	const custom = "_doRTgHZBKcGVjlvpC,@aFSx#DPuNJme&i*MzLOEn)sUrthbf%Y^w.(kIQyXqWA!$"
	const standard = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/="
	for i, c := range raw {
		index := strings.IndexByte(custom, c)
		if index < 0 {
			t.Fatal("invalid encoded Qoder request")
		}
		raw[i] = standard[index]
	}
	a := len(raw) / 3
	raw = append(append(append([]byte{}, raw[len(raw)-a:]...), raw[a:len(raw)-a]...), raw[:a]...)
	decoded, err := base64.StdEncoding.DecodeString(string(raw))
	if err != nil {
		t.Fatal(err)
	}
	var result map[string]any
	if err := json.Unmarshal(decoded, &result); err != nil {
		t.Fatal(err)
	}
	return result
}

func qoderToolSearchResponse(t *testing.T, name, args string) string {
	t.Helper()
	var body strings.Builder
	for _, choice := range []map[string]any{
		{"index": 0, "delta": map[string]any{"tool_calls": []any{map[string]any{
			"index": 0, "id": "call_probe", "type": "function", "function": map[string]any{"name": name, "arguments": args},
		}}}, "finish_reason": nil},
		{"index": 0, "delta": map[string]any{}, "finish_reason": "tool_calls"},
	} {
		inner, _ := json.Marshal(map[string]any{"id": "chatcmpl-probe", "object": "chat.completion.chunk", "created": 1, "model": "auto", "choices": []any{choice}})
		envelope, _ := json.Marshal(map[string]any{"statusCodeValue": 200, "body": string(inner)})
		fmt.Fprintf(&body, "data: %s\n\n", envelope)
	}
	body.WriteString("data: [DONE]\n\n")
	return body.String()
}

func TestQoderToolSearchExecutorDiscoveryRoundTrip(t *testing.T) {
	for _, streaming := range []bool{false, true} {
		t.Run(fmt.Sprintf("stream=%v", streaming), func(t *testing.T) {
			cfg := &config.Config{SDKConfig: config.SDKConfig{Qoder: config.QoderConfig{ToolSearch: true}}}
			executor := NewQoderExecutor(cfg)
			auth := &cliproxyauth.Auth{Storage: testQoderStorageWithModelConfig()}
			stage := 0
			var searchCallID string
			ctx := context.WithValue(context.Background(), "cliproxy.roundtripper", qoderRoundTripFunc(func(r *http.Request) (*http.Response, error) {
				wire := decodeQoderToolSearchRequest(t, r)
				tools, _ := wire["tools"].([]any)
				want := 1
				if stage == 1 {
					want = 2
				}
				if len(tools) != want {
					t.Fatalf("stage %d: tools=%v want %d", stage, tools, want)
				}
				name := ""
				for _, raw := range tools {
					tool := raw.(map[string]any)
					fn := tool["function"].(map[string]any)
					n := fn["name"].(string)
					if stage == 0 || n != "tool_search" {
						name = n
					}
				}
				args := `{"query":"read lookup value"}`
				if stage == 1 {
					args = `{"key":"sample"}`
					messages := wire["messages"].([]any)
					found := false
					for _, raw := range messages {
						m := raw.(map[string]any)
						if m["role"] == "tool" && m["tool_call_id"] == searchCallID {
							found = true
						}
					}
					if !found {
						t.Fatal("search result history was dropped or renamed")
					}
				}
				return qoderSSEHTTPResponse(r, 200, qoderToolSearchResponse(t, name, args)), nil
			}))
			function := map[string]any{"type": "function", "name": "lookup", "defer_loading": true, "parameters": map[string]any{"type": "object", "properties": map[string]any{"key": map[string]any{"type": "string"}}, "required": []string{"key"}}}
			namespace := map[string]any{"type": "namespace", "name": "catalog", "tools": []any{function}}
			body := map[string]any{"model": "qoder/auto", "input": []any{map[string]any{"role": "user", "content": "Look up sample"}}, "tools": []any{map[string]any{"type": "tool_search", "execution": "client"}, namespace}}
			for stage = 0; stage < 2; stage++ {
				raw, _ := json.Marshal(body)
				req := cliproxyexecutor.Request{Model: "qoder/auto", Payload: raw}
				opts := cliproxyexecutor.Options{SourceFormat: sdktranslator.FormatOpenAIResponse, OriginalRequest: raw}
				var response map[string]any
				if streaming {
					result, err := executor.ExecuteStream(ctx, auth, req, opts)
					if err != nil {
						t.Fatal(err)
					}
					for chunk := range result.Chunks {
						if chunk.Err != nil {
							t.Fatal(chunk.Err)
						}
						for _, line := range bytes.Split(chunk.Payload, []byte("\n")) {
							if !bytes.HasPrefix(line, []byte("data: ")) {
								continue
							}
							var event map[string]any
							if json.Unmarshal(bytes.TrimPrefix(line, []byte("data: ")), &event) == nil && event["type"] == "response.completed" {
								response, _ = event["response"].(map[string]any)
							}
						}
					}
				} else {
					result, err := executor.Execute(ctx, auth, req, opts)
					if err != nil {
						t.Fatal(err)
					}
					if err := json.Unmarshal(result.Payload, &response); err != nil {
						t.Fatal(err)
					}
				}
				if response == nil {
					t.Fatal("no completed response")
				}
				outputs := response["output"].([]any)
				var call map[string]any
				for _, raw := range outputs {
					item := raw.(map[string]any)
					if item["type"] == "tool_search_call" || item["type"] == "function_call" {
						call = item
					}
				}
				if call == nil {
					t.Fatalf("missing tool call: %v", outputs)
				}
				if stage == 0 {
					if call["type"] != "tool_search_call" || !strings.HasPrefix(call["id"].(string), "tsc_") {
						t.Fatalf("incorrect search identity: %v", call)
					}
					searchCallID = call["call_id"].(string)
					loaded := make(map[string]any, len(function))
					for key, value := range function {
						loaded[key] = value
					}
					loaded["defer_loading"] = false
					loadedNamespace := map[string]any{"type": "namespace", "name": "catalog", "tools": []any{loaded}}
					history := append(body["input"].([]any), outputs...)
					body["input"] = append(history, map[string]any{"type": "tool_search_output", "execution": "client", "call_id": searchCallID, "tools": []any{loadedNamespace}})
				} else if call["type"] != "function_call" || call["name"] != "lookup" || call["namespace"] != "catalog" {
					t.Fatalf("incorrect discovered namespace identity: %v", call)
				}
				if n, b := executor.toolSearchLimiter.Usage(); n != 0 || b != 0 {
					t.Fatalf("leaked lease: %d/%d", n, b)
				}
			}
		})
	}
}

func TestQoderToolSearchBrokenStreamDoesNotComplete(t *testing.T) {
	for _, enabled := range []bool{false, true} {
		t.Run(fmt.Sprintf("enabled=%v", enabled), func(t *testing.T) {
			executor := NewQoderExecutor(&config.Config{SDKConfig: config.SDKConfig{Qoder: config.QoderConfig{ToolSearch: enabled}}})
			readFailure := errors.New("synthetic stream interruption")
			ctx := context.WithValue(context.Background(), "cliproxy.roundtripper", qoderRoundTripFunc(func(r *http.Request) (*http.Response, error) {
				body := strings.TrimSuffix(qoderToolSearchResponse(t, "tool_search", `{"query":"lookup"}`), "data: [DONE]\n\n")
				response := qoderSSEHTTPResponse(r, 200, "")
				response.Body = io.NopCloser(io.MultiReader(strings.NewReader(body), iotest.ErrReader(readFailure)))
				return response, nil
			}))
			req := cliproxyexecutor.Request{Model: "qoder/auto", Payload: []byte(`{"input":"lookup","tools":[{"type":"tool_search","execution":"client"}]}`)}
			result, err := executor.ExecuteStream(ctx, &cliproxyauth.Auth{Storage: testQoderStorageWithModelConfig()}, req, cliproxyexecutor.Options{SourceFormat: sdktranslator.FormatOpenAIResponse})
			if err != nil {
				t.Fatal(err)
			}
			gotError := false
			for chunk := range result.Chunks {
				if chunk.Err != nil {
					gotError = true
				}
				if bytes.Contains(chunk.Payload, []byte(`"type":"response.completed"`)) {
					t.Fatal("broken upstream stream emitted response.completed")
				}
			}
			if !gotError {
				t.Fatal("upstream read error was lost")
			}
			if n, b := executor.toolSearchLimiter.Usage(); n != 0 || b != 0 {
				t.Fatalf("leaked failed attempt: %d/%d", n, b)
			}
		})
	}
}
