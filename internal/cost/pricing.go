package cost

// ModelPricing contains pricing per million tokens for a model.
type ModelPricing struct {
	InputPer1M       float64 // Price per 1M input tokens
	OutputPer1M      float64 // Price per 1M output tokens
	CacheReadPer1M   float64 // Price per 1M cache read tokens
	CacheWritePer1M  float64 // Price per 1M cache write tokens
}

// Prices per million tokens for Claude models (as of 2024).
// Reference: https://www.anthropic.com/api
var modelPricing = map[string]ModelPricing{
	// Claude 4 Opus
	"claude-opus-4-5-20251101": {
		InputPer1M:      15.0,
		OutputPer1M:     75.0,
		CacheReadPer1M:  1.5,
		CacheWritePer1M: 18.75,
	},
	// Claude 3.5 Sonnet
	"claude-sonnet-4-20250514": {
		InputPer1M:      3.0,
		OutputPer1M:     15.0,
		CacheReadPer1M:  0.30,
		CacheWritePer1M: 3.75,
	},
	"claude-3-5-sonnet-20241022": {
		InputPer1M:      3.0,
		OutputPer1M:     15.0,
		CacheReadPer1M:  0.30,
		CacheWritePer1M: 3.75,
	},
	"claude-3-5-sonnet-20240620": {
		InputPer1M:      3.0,
		OutputPer1M:     15.0,
		CacheReadPer1M:  0.30,
		CacheWritePer1M: 3.75,
	},
	// Claude 3.5 Haiku
	"claude-3-5-haiku-20241022": {
		InputPer1M:      1.0,
		OutputPer1M:     5.0,
		CacheReadPer1M:  0.10,
		CacheWritePer1M: 1.25,
	},
	// Claude 3 Opus
	"claude-3-opus-20240229": {
		InputPer1M:      15.0,
		OutputPer1M:     75.0,
		CacheReadPer1M:  1.5,
		CacheWritePer1M: 18.75,
	},
	// Claude 3 Sonnet
	"claude-3-sonnet-20240229": {
		InputPer1M:      3.0,
		OutputPer1M:     15.0,
		CacheReadPer1M:  0.30,
		CacheWritePer1M: 3.75,
	},
	// Claude 3 Haiku
	"claude-3-haiku-20240307": {
		InputPer1M:      0.25,
		OutputPer1M:     1.25,
		CacheReadPer1M:  0.03,
		CacheWritePer1M: 0.30,
	},
}

// defaultPricing is used when model is not found.
// Uses Claude 3.5 Sonnet pricing as a reasonable default.
var defaultPricing = ModelPricing{
	InputPer1M:      3.0,
	OutputPer1M:     15.0,
	CacheReadPer1M:  0.30,
	CacheWritePer1M: 3.75,
}

// GetPricing returns pricing for a model.
// Falls back to default pricing if model is not found.
func GetPricing(modelID string) ModelPricing {
	if p, ok := modelPricing[modelID]; ok {
		return p
	}

	// Try prefix matching for versioned models
	for id, p := range modelPricing {
		if len(modelID) > 10 && len(id) > 10 {
			// Compare base model name (e.g., "claude-3-5-sonnet")
			if modelID[:min(len(modelID), 20)] == id[:min(len(id), 20)] {
				return p
			}
		}
	}

	return defaultPricing
}

// CalculateCost computes the total cost for a set of tokens.
func CalculateCost(modelID string, inputTokens, outputTokens, cacheRead, cacheWrite int) float64 {
	p := GetPricing(modelID)

	cost := float64(inputTokens) * p.InputPer1M / 1_000_000
	cost += float64(outputTokens) * p.OutputPer1M / 1_000_000
	cost += float64(cacheRead) * p.CacheReadPer1M / 1_000_000
	cost += float64(cacheWrite) * p.CacheWritePer1M / 1_000_000

	return cost
}

