package models

import "testing"

func TestQoderToolSearchCapabilityRequiresOptInAndPureLocalRoutes(t *testing.T) {
	for _, enabled := range []bool{false, true} {
		entries := []map[string]any{
			{"slug": "qoder/auto", "supports_search_tool": false},
			{"slug": "mixed", "supports_search_tool": false},
			{"slug": "unknown", "supports_search_tool": false},
			{"slug": "gpt", "supports_search_tool": true},
		}
		response := map[string]any{"models": entries}
		ApplyQoderToolSearchCapability(response, func(id string) []string {
			switch id {
			case "auto":
				return []string{"qoder"}
			case "mixed":
				return []string{"qoder", "openai-compatibility"}
			case "gpt":
				return []string{"codex"}
			}
			return nil
		}, enabled)
		for i, want := range []bool{enabled, false, false, true} {
			if entries[i]["supports_search_tool"] != want {
				t.Fatalf("enabled=%v entry=%v want=%v", enabled, entries[i], want)
			}
		}
	}
}
