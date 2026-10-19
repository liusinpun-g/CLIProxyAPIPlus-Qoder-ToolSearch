package models

import "strings"

// ApplyQoderToolSearchCapability advertises the opt-in bridge only for models
// whose selectable routes all use the local Qoder executor. Mixed-provider and
// Home catalogs must not inherit a local provider's capability.
func ApplyQoderToolSearchCapability(response map[string]any, providersForModel ProvidersForModelFunc, enabled bool) {
	if !enabled || providersForModel == nil {
		return
	}
	entries, _ := response["models"].([]map[string]any)
	for _, entry := range entries {
		id, _ := entry["slug"].(string)
		providers := providersForModel(id)
		if len(providers) == 0 {
			if _, base, ok := strings.Cut(id, "/"); ok {
				providers = providersForModel(base)
			}
		}
		if len(providers) == 0 {
			continue
		}
		pureQoder := true
		for _, provider := range providers {
			if !strings.EqualFold(strings.TrimSpace(provider), "qoder") {
				pureQoder = false
				break
			}
		}
		if pureQoder {
			entry["supports_search_tool"] = true
		}
	}
}
